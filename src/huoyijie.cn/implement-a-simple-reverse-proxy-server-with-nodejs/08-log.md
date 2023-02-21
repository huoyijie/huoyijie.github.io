# 流水日志

```log
4/21/2021, 1:31:59 PM   200     7       960b76c0a24c11eb851d71230b772f89        e48fb4d1a26211eb851d71230b772f89        h2      /article/a2e12610a0f811ebafc9c93cfc0a76ae/06-request.html
4/21/2021, 1:50:35 PM   200     9       7dbce4a0a26511eb851d71230b772f89        7de665a1a26511eb851d71230b772f89        h2      /article/a2e12610a0f811ebafc9c93cfc0a76ae/02-config.html
4/21/2021, 2:00:58 PM   200     13      f13972d0a26611eb851d71230b772f89        f1832600a26611eb851d71230b772f89        1.1     /?from=blog
```

流水日志输出如上，输出字段依次为时间戳、响应码、耗时（毫秒）、会话 ID、请求 ID、协议版本、请求路径。

http2 的优势之一就是，每个域只需建立一个连接，并基于此连接建立会话，后续所有请求都通过此会话来进行。为了跟踪方便，每个会话分配了唯一 ID，在所有日志（包括 debug 日志）输出中都会携带会话 ID，方便跟踪一个 session 上的所有活动。因为每个 session 会发送不受限制的请求数量，因此还需要一个请求 ID 来具体跟踪每一个请求。请求 ID 会在具体某一个请求过程中的所有日志中输出，方便跟踪具体请求的所有活动。

```javascript
server.listen(config.port, config.host, onListening)
  .on('connection', (socket /** net.Socket */) => {
    // 分配连接 ID，也是会话 ID
    socket.id = uuid.v1().replace(/-/g, '')
    console.debug('server.onConnection', socket.id)
  })
  // onSession 回调函数中会读取连接 ID，并赋值给会话 ID
  .on('session', onSession)
  .on('request', (req, res) => {
    const now = new Date()
    req._context = {
      timestamp: now.toLocaleString(),
      startTime: now.getTime(),
      path: req.url,
      version: req.httpVersion,
      sid: req.socket._parent.id,
      // 分配请求 ID
      reqid: uuid.v1().replace(/-/g, ''),
      format() {
        return `${this.timestamp}\t${this.code}\t${this.cost}\t${this.sid}\t${this.reqid}\t${this.version}\t${this.path}\n`
      }
    }
    const next = (req, res, err) => {
      req._context.endTime = Date.now()
      req._context.code = res.statusCode
      req._context.cost = req._context.endTime - req._context.startTime
      console.debug('req.access.log', req._context)
      fs.appendFile(config.accessLog, req._context.format(), (err) => {
        if (err) {
          console.error('req.access.log error', req._context.sid, err)
        }
      })
    }
    onRequest(req, res, next)
  })
  .on('stream', (stream, headers) => {
    const now = new Date()
    stream._context = {
      timestamp: now.toLocaleString(),
      startTime: now.getTime(),
      path: headers[HTTP2_HEADER_PATH],
      version: stream.session.alpnProtocol,
      sid: stream.session.id,
      // 分配请求 ID
      reqid: uuid.v1().replace(/-/g, ''),
      format() {
        return `${this.timestamp}\t${this.code}\t${this.cost}\t${this.sid}\t${this.reqid}\t${this.version}\t${this.path}\n`
      }
    }
    const next = (stream, headers, err) => {
      stream._context.endTime = Date.now()
      stream._context.code = stream.sentHeaders ? stream.sentHeaders[HTTP2_HEADER_STATUS] : HTTP_STATUS_INTERNAL_SERVER_ERROR
      stream._context.cost = stream._context.endTime - stream._context.startTime
      console.debug('stream.access.log', stream._context)
      fs.appendFile(config.accessLog, stream._context.format(), (err) => {
        if (err) {
          console.error('stream.access.log error', stream._context.sid, err)
        }
      })
    }
    onStream(stream, headers, next)
  })
```

上述代码关键部分已加上了注释，整体逻辑很简单。就是在启动服务器监听时，注册 request 及 stream 事件处理函数，并在这 2 个函数中定义 next 回调函数，由 next 回调函数最后完成流水日志的输出。next 函数会在向浏览器写回数据时调用。