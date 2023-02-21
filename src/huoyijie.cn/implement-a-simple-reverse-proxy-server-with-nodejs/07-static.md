# 实现静态服务器

本节内容实现一个简单的静态文件服务组件 serveStatic，实现对静态资源的请求处理。代码位于文件 `serve/static.js` 中。

```javascript
module.exports = (config) => {
  let serveStatic = (stream, headers, next) => _serveStatic(stream, headers, next, config)
  return serveStatic
}
```

该模块参数为 config 配置组件对象，返回 serveStatic 组件，该组件实现基于 _serveStatic 方法。

```javascript
function _serveStatic(stream, headers, next, config) {
  // 请求静态文件路径
  const reqPath = headers[HTTP2_HEADER_PATH]

  // 请求以 /static 开头（可配置）代表静态资源请求
  if(reqPath.split(path.sep)[1] === config.staticDir) {
    console.debug(
      'server.onStream-static',
      stream.id, stream._context.sid,
      stream._context.reqid)

    // 获取静态文件相对 serverRoot 的绝对路径
    const fullPath = path.join(config.serverRoot, reqPath)
    // 获取请求文件类型，得出 content-type
    const responseMimeType = mime.lookup(fullPath)

    // 读取文件进行响应
    stream.safe_respondWithFile(fullPath, { 'content-type': responseMimeType }, {
      waitForTrailers: true,
      statCheck: (stat, resHeaders) => {
        resHeaders['last-modified'] = stat.mtime.toUTCString()
        resHeaders['cache-control'] = 'public, max-age=0'
        if(Date.parse(headers['if-modified-since']) === Date.parse(stat.mtime)) {
          stream.safe_respond({ ':status': HTTP_STATUS_NOT_MODIFIED })
          stream.end()
          next(stream, headers)
          return false; // Cancel the send operation
        }
      },
      onError: (err) => {
        console.debug(
          'server.onStream-static.respondToStreamError',
          stream.id,
          stream._context.sid,
          stream._context.reqid,
          err)
        if (err.code === 'ENOENT') {
          stream.safe_respond({ ':status': HTTP_STATUS_NOT_FOUND })
        } else {
          stream.safe_respond({ ':status': HTTP_STATUS_NOT_FOUND })
        }
        stream.end()
        next(stream, headers, err)
      }
    })
    stream.on('wantTrailers', () => {
      console.debug(
        'server.onStream-static.wantTrailers',
        stream.id,
        stream._context.sid,
        stream._context.reqid)
      next(stream, headers)
      stream.close()
    })

    return true
  }

  return false
}
```

代码上加了注释，整体上比较简单，值得一提的地方是 statCheck 回调。会在读取文件发送回浏览器之前做检查，如果从浏览器上次访问之后没有修改，则返回 304 给浏览器，并且不需返回文件内容。

```javascript
statCheck: (stat, resHeaders) => {
  resHeaders['last-modified'] = stat.mtime.toUTCString()
  resHeaders['cache-control'] = 'public, max-age=0'
  if(Date.parse(headers['if-modified-since']) === Date.parse(stat.mtime)) {
    stream.safe_respond({ ':status': HTTP_STATUS_NOT_MODIFIED })
    stream.end()
    next(stream, headers)
    return false; // Cancel the send operation
  }
}
```

上述代码会返回该文件 last-modified 时间给浏览器，并告诉浏览器缓存策略 cache-control 为，在确定使用缓存前来服务器请求确认下。浏览器在请求资源时，请求头中携带 if-modified-since，时间是服务器返回的 last-modified，如果服务器返回 304，则可以直接使用缓存文件。否则服务器会返回新的文件内容。