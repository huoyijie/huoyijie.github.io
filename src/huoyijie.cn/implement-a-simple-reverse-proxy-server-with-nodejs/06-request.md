# HTTP/1.1 请求处理

本节内容接上节内容，继续介绍 request.js 模块中的 _request_1_1 函数，主要是兼容之前的协议，使得浏览器或者其他客户端可以通过非 http/2.0 协议访问。

```javascript
function _request_1_1(req, res, next, proxy) {
  let headers = req.headers
  let chunks = []
  req
  .on('data', (chunk) => chunks.push(chunk))
  .on('end', () => {
    let onres = (response) => {
      let resChunks = []
      response
      .on('data', (chunk) => resChunks.push(chunk))
      .on('end', () => {
        res.writeHead(response.statusCode, response.headers)
        res.end(Buffer.concat(resChunks))
        next(req, res)
      })
    }
  
    let endStream = (chunks.length == 0)
    let reqStream = proxy.get().request(headers, { endStream: endStream, onres })
    if(!endStream) {
      reqStream.end(Buffer.concat(chunks))
    }
    reqStream
    .on('error', (err) => {
      res.writeHead(HTTP_STATUS_INTERNAL_SERVER_ERROR, { 'content-type': 'text/html; charset=utf-8' })
      res.end()
      next(req, res, err)
    })
    .setTimeout(5000, () => {
      res.writeHead(HTTP_STATUS_GATEWAY_TIMEOUT, { 'content-type': 'text/html; charset=utf-8' })
      res.end()
      next(req, res)
    })
  })
}
```

如上述代码，首先定义了代理转发后的回调函数 onres，然后代理转发请求 `proxy.get().request(headers, { endStream: endStream, onres })` 到后端程序，后端程序响应后回调 onres 函数，把后端程序的响应数据写回浏览器。