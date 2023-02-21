# HTTP/2.0 请求处理

服务器启动后会处于监听中，通过浏览器访问本站时，首先建立 TCP 连接，然后进行 TLS 握手，通过 ALPN 协商 HTTP 版本号为 h2（http/2.0）。完成 TLS 握手后，后续交换的数据都是加密的。而且后续的所有 http 请求，浏览器会优先采用 2.0 协议向服务器发送请求。如果有的请求没办法用 2.0 协议，则自动降级为 1.1 协议，如 websocket。关于 websocket 请求的处理后面会在展开介绍。本节内容是对 http/2.0 请求的转发处理，是反向代理实现的核心代码之一。

代码位于 `serve/stream.js` 文件

```javascript
const http2 = require('http2')

function _serveStream(stream, headers, next, Client, config) {
  let reqPath = headers[HTTP2_HEADER_PATH],
      proxy = Client.get(reqPath),
      chunks = []

  if(!proxy) {
    console.debug('server.onStream._serveStream', `no proxy for ${reqPath}`, stream.id, stream._context.sid, stream._context.reqid)
    // 没有代理能处理此请求，返回 404
    stream.safe_respond({ ':status': HTTP_STATUS_NOT_FOUND })
    stream.end()
    next(stream, headers)
    return
  }

  stream
  .on('data', (chunk) => {
    console.debug('server.onStream.onData', `got ${chunk.length} bytes`, stream.id, stream._context.sid, stream._context.reqid)
    if(chunk.length > 0) {
      chunks.push(chunk)
    }
  })
  .on('end', () => {
    console.debug('server.onStream.onEnd', `got ${chunks.length} chunks`, stream.id, stream._context.sid, stream._context.reqid)
    headers['x-real-ip'] = stream.session.socket.remoteAddress
    headers['x-forwarded-for'] = stream.session.socket.remoteAddress
    _stream_1_1(stream, headers, next, proxy, chunks)
  })
}

module.exports = (Client, config) => {
  let serveStream = (stream, headers, next) => _serveStream(stream, headers, next, Client, config)
  return serveStream
}
```

如上述代码，stream.js 模块导出了 serveStream 组件。该模块参数是前面已经分别介绍过的远程连接管理组件 Client（client.js） 和配置组件 config（config.js）。serveStream 组件会在服务器接受 2.0 请求时被调用，此时服务器会触发 stream 事件。serveStream 会得到传入的 stream 与 headers 对象，代表了与浏览器之间的双向流和请求头部集合。

通过 `proxy = Client.get(reqPath)` 获取转发代理 proxyClient 对象。通过在 stream 对象上注册 `data` 及 `end` 事件处理函数，可得到请求 body 数据 `chunks`。因为多了一层 hawkey 转发代理，在后端博客程序中无法获取客户端远程 IP 地址，此时需要增加 `x-real-ip` 及 `x-forwarded-for` 头部，方便后端程序读取。真正的转发处理逻辑位于 `_stream_1_1` 方法中，会在接收完请求体所有数据后调用。

```javascript
function _stream_1_1(stream, headers, next, proxy, chunks) {
  // part-2: 代理转发请求后，从后端程序接收响应
  let onres = (res) => {
    console.debug('server.onStream.onEnd.onres', 'http/1.1', stream.id, stream._context.sid, stream._context.reqid, res.headers)
    let resChunks = []
    res
    .on('data', (chunk) => {
      console.debug('server.onStream.onEnd.onres.onData', 'http/1.1', stream.id, stream._context.sid, stream._context.reqid, `got ${chunk.length} bytes`)
      if(chunk.length > 0) {
        resChunks.push(chunk)
      }
    })
    .on('end', () => {
      console.debug('server.onStream.onEnd.onres.onEnd', 'http/1.1', stream.id, stream._context.sid, stream._context.reqid, `got ${resChunks.length} chunks`)
      let endStream = (resChunks.length == 0)

      // h2 <- http/1.1 header 兼容处理
      // HTTP/1 Connection specific headers are forbidden: "connection, keep-alive, transfer-encoding"
      delete res.headers[HTTP2_HEADER_CONNECTION]
      delete res.headers[HTTP2_HEADER_KEEP_ALIVE]
      delete res.headers[HTTP2_HEADER_TRANSFER_ENCODING]

      res.headers[HTTP2_HEADER_STATUS] = res.statusCode
      stream.safe_respond(res.headers, { endStream: endStream })
      if(!endStream) {
        stream.end(Buffer.concat(resChunks))
      }
      next(stream, headers)
    })
    .on('abort', () => {
      stream.safe_respond({ ':status': HTTP_STATUS_INTERNAL_SERVER_ERROR })
      stream.end()
      next(stream, headers)
    })
  }

  // part-1: 代理转发请求
  // h2 -> http/1.1 header 兼容处理
  headers['method'] = headers[HTTP2_HEADER_METHOD]
  headers['path'] = headers[HTTP2_HEADER_PATH]
  delete headers[HTTP2_HEADER_AUTHORITY]
  delete headers[HTTP2_HEADER_METHOD]
  delete headers[HTTP2_HEADER_PATH]
  delete headers[HTTP2_HEADER_SCHEME]

  let client = proxy.get()
  let endStream = (chunks.length == 0)
  const reqStream = client.request(headers, { endStream, onres })
  if(!endStream) {
    reqStream.end(Buffer.concat(chunks))
  }
  reqStream
  .on('error', (err) => {
    console.error('server.onStream.onEnd.error', 'http/1.1', stream.id, stream._context.sid, stream._context.reqid, err)
    stream.safe_respond({ ':status': HTTP_STATUS_INTERNAL_SERVER_ERROR })
    stream.end()
    next(stream, headers, err)
  })
  .setTimeout(5000, () => {
    console.error('server.onStream.onEnd.onTimeout', stream._context.sid, stream._context.reqid)
    stream.safe_respond({ ':status': HTTP_STATUS_GATEWAY_TIMEOUT })
    stream.end()
    next(stream, headers)
  })
}
```

上述代码为 `_stream_1_1` 方法，主要分为 2 部分（参见代码中注释）：

* part-1: 代理转发请求
  
  后端程序最高支持 http/1.1 协议，所以不能把浏览器发送请求的头部对象 headers 直接转发，需进行 h2 -> http/1.1 header 兼容处理，然后 `let client = proxy.get()`，最后代理转发请求 `client.request(headers, { endStream, onres })`，此处会注册 `onres` 回调函数，该函数于 part-2 部分定义

* part-2: 代理转发请求后，从后端程序接收响应
  
  这部分代码主要是定义 `onres` 回调，通过在后端程序（转发目标）响应对象 `res` 上注册 `data` 及 `end` 事件处理函数，可获得响应 body 完整数据。在读取 res 头部及响应数据后，向浏览器转发响应前，还需进行 h2 <- http/1.1 header 兼容处理。然后向浏览器发送头部信息 `stream.safe_respond(res.headers, { endStream: endStream })`，如果还有响应 body 数据，则继续发送 `if(!endStream) { reqStream.end(Buffer.concat(chunks)) }`