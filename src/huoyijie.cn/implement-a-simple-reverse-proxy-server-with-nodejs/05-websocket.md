# Websocket 请求处理

本节内容是对 http/1.1 及以下协议的请求转发处理介绍，是反向代理实现的核心代码之一，会着重介绍 websocket 代理转发。代码位于 `serve/request.js` 文件中。

```javascript
module.exports = (Client) => {
    let serveRequest = (req, res, next) => _serveRequest(req, res, next, Client)
    return serveRequest
}
```

`request.js` 模块导出了 serveRequest 组件，该模块的参数是前面介绍过的远程连接管理组件 Client。serveRequest 组件会在服务器接受 1.1 及以下协议请求时被调用，此时服务器会触发 request 事件。serveRequest 会得到传入的 req 与 res 对象，分别代表了浏览器与 hawkey 之间的请求与响应对象。serveRequest 组件的真正实现为 _serveRequest 方法。

```javascript
function _serveRequest(req, res, next, Client) {
  let proxy = Client.get(req.url)
  if(!proxy) {
    console.debug('server.onRequest._serveRequest',
      `no proxy for ${req.url}`,
      req._context.sid,
      req._context.reqid)
    // 没有代理能处理此请求，返回 404
    res.writeHead(HTTP_STATUS_NOT_FOUND, { 'content-type': 'text/html; charset=utf-8' })
    res.end()
    next(req, res)
    return
  }

  if(req.httpVersion === '2.0') {
    console.debug(
      'server.onRequest._serveRequest',
      'doNothing',
      req._context.sid,
      req._context.reqid)
  } else if(req.httpVersion === '1.1') {
    req.headers['x-real-ip'] = req.socket.remoteAddress
    req.headers['x-forwarded-for'] = req.socket.remoteAddress
    req.headers.method = req.method
    req.headers.path = req.url

    if(_upgradeWebsocket(req, res, next, proxy)) {
      console.debug('server.onRequest._serveRequest',
      'doUpgrade',
      req._context.sid,
      req._context.reqid)
    } else {
      _request_1_1(req, res, next, proxy)
    }
  } else {
    console.debug(
      'server.onRequest._serveRequest',
      `http ${req.httpVersion} Not Supported`,
      req._context.sid,
      req._context.reqid)
    res.writeHead(HTTP_STATUS_INTERNAL_SERVER_ERROR, { 'content-type': 'text/html; charset=utf-8' })
    res.end(Buffer.from('<h4>Not Supported</h4>', 'utf8'))
    next(req, res)
  }
}
```

上述代码首先根据请求 url 获取 proxyClient 对象 `let proxy = Client.get(req.url)`。然后判断请求协议版本号 `req.httpVersion` 进行不同的处理。

* `req.httpVersion === '2.0'`
  
  2.0 请求会先触发服务器的 request 事件，然后是 stream 事件，2.0 请求已在 stream.js 模块中处理过了，此处无需处理，忽略即可。

* `req.httpVersion === '1.1'`

  为了在后端程序中可获取到客户端真实 IP，在头部中增加 `x-real-ip` 及 `x-forwarded-for` 字段，然后判断当前请求如果是一个 websocket 握手请求（连接升级），则调用 _upgradeWebsocket 方法，否则会调用 _request_1_1 方法。

* 其他情况直接向浏览器返回不支持错误

本节内容会介绍 _upgradeWebsocket，而 _request_1_1 方法会在下一节介绍。

```javascript
function _upgradeWebsocket(req, res, next, proxy) {
  let headers = req.headers
  let isWebsocket = (
    (headers[HTTP2_HEADER_CONNECTION] || '').toLowerCase() === 'upgrade' &&
    (headers[HTTP2_HEADER_UPGRADE] || '').toLowerCase() === 'websocket'
  )
  if(isWebsocket) {
    console.debug(
      'websocket.onRequest-upgrade',
      headers,
      req._context.sid,
      req._context.reqid)

    let onupgrade = (response, socket, upgradeHead) => {
        res.writeHead(HTTP_STATUS_SWITCHING_PROTOCOLS, response.headers)
        res.end()
        res.socket.pipe(socket).pipe(res.socket)
        console.debug(
          'websocket.onUpgrade-response',
          headers,
          response.headers,
          req._context.sid,
          req._context.reqid)
        next(req, res)
    }

    proxy.get().request(headers, { endStream: true, onupgrade })
    .on('error', (err) => {
      console.error(
        'websocket.client.onError',
        req._context.sid,
        req._context.reqid,
        err)
      res.writeHead(HTTP_STATUS_INTERNAL_SERVER_ERROR, { 'content-type': 'text/html; charset=utf-8' })
      res.end()
      next(req, res, err)
    })
    .setTimeout(5000, () => {
      console.error(
        'websocket.client.onTimeout',
        req._context.sid,
        req._context.reqid)
      res.writeHead(HTTP_STATUS_GATEWAY_TIMEOUT, { 'content-type': 'text/html; charset=utf-8' })
      res.end()
      next(req, res)
    })
  }
  return isWebsocket
}
```

如上述代码，首先判断如果是 websocket 握手请求，则定义握手成功后的连接升级回调函数 `onupgrade`，`onupgrade` 回调函数中 socket 参数是升级后的套接字，代表 hawkey 与后端程序之间的 websoket 连接，而该函数中主要逻辑为

```javascript
res.writeHead(HTTP_STATUS_SWITCHING_PROTOCOLS, response.headers)
res.end()
res.socket.pipe(socket).pipe(res.socket)
```

即为，返回 `HTTP/1.1 101 Swiching Protocols`，及其他头部信息表示握手成功，连接已升级到 websocket。其中很关键的一句 `res.socket.pipe(socket).pipe(res.socket)` 把 _<浏览器-hawkey-后端程序>_ 之间建立了管道 pipe。

```javascript
// res.socket | socket
// res.socket 的输出（浏览器 -> hawkey），作为 socket 的输入，输入 socket 的数据会流向后端程序
res.socket.pipe(socket)

// socket | res.socket
// socket 的输出（后端程序 -> hawkey），作为 res.socket 的输入，输入 res.socket 的数据会流向浏览器
socket.pipe(res.socket)
```

websocket 帧数据会在这个管道里流起来，就好象 hawkey 不存在，浏览器与后端程序之间直接建立了 websocket 连接。这里的 pipe 操作和 Linux 系统中管道是一个意思，就像 `ls -lh |grep reverse-server | awk '{print $1,$4}'` 中的 `|` ，把上个程序的输出作为下个程序的输入。