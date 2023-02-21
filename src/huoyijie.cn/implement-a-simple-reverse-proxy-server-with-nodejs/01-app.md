# 启动监听服务

hawkey 主要基于 Node.js 的 http2 模块实现，该模块实现了 http/2.0 协议。2.0 协议没有改变 http 协议的语义，主要修改了传输层实现。虽然说是没有改变语义，但是 2.0 协议并非向前兼容的，也就是说必须由 http 服务器和客户端软件自行兼容。

服务器能否只实现 2.0 协议呢？实际上要看情况，并非所有的请求都可以通过 2.0 协议发送，比如 websocket 请求一般会通过 http/1.1 协议发送，向服务器协商升级连接，当然也可能会有很少的人电脑安装的浏览器版本或者其他客户端比较老不支持 2.0 协议。如果不考虑上述场景可以只实现 2.0 协议，充分发挥 2.0 协议通信性能，又免去了兼容旧协议的痛苦。否则最好都支持下，服务器兼容性更好。

比如 nginx 就同时支持 http/0.9、http/1.0、http/1.1、http/2.0 请求。chrome 等主流浏览器也是如此。我在实现 hawkey 时也决定要兼容之前的协议。浏览器怎么知道用什么协议来向服务器发送请求呢？

http/2.0 协议虽然允许非 https 请求，但是所有浏览器在实现 http/2.0 时都要求服务器必须支持 https，这已经是事实上的标准。回答刚刚提出的问题，浏览器怎样和服务器协商协议版本呢？就是通过在 TLS 握手时，通过 ALPN 来协商协议版本。浏览器首先把自己支持的协议版本列表发给服务器，服务器从列表中选出自己支持的，并最终选定其中一个通知客户端完成协议协商。通常浏览器和服务器都会尽量采用高版本协议（http/2.0）通信。

我的电脑和服务器系统都是 ubuntu

```bash
$ uname -a
Linux xps 5.4.0-72-generic #80-Ubuntu SMP Mon Apr 12 17:35:00 UTC 2021 x86_64 x86_64 x86_64 GNU/Linux
```

开发基于 Node.js

```bash
$ node -v
v14.16.0
```

#### 创建服务器

```javascript
const http2 = require('http2')
const fs = require('fs')
const server = http2.createSecureServer({
  key: fs.readFileSync(config.pems.key),
  cert: fs.readFileSync(config.pems.cert),
  allowHTTP1: true
})
```
如上以兼容旧协议（allowHTTP1: true）模式创建 http2 服务器，同时通过读取文件提供了服务器的私钥 key 以及 证书 cert。服务器私钥 key 及 证书 cert 可以通过向 CA 机构申请免费证书获得，也可以先用自签名证书测试。我这边有生成自签名证书的脚本。

```bash
#!/bin/bash

name="cert-name"
localip="127.0.0.1"
dnsname=$localip

while(( $# > 1 ));
do
  case $1 in

  --name)
    if [ "$2" != "" ]; then
      name="$2"
    fi

    DIR=`dirname $0`
    if [ ! -d "keys" ]; then
      mkdir $DIR/keys
    fi

    if [ -f "$DIR/keys/$name-cert.pem" ]; then
      echo "$DIR/keys/$name-cert.pem exist, please check first"
      exit 1
    fi

    cd $DIR/keys

    openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes \
    -keyout $name-key.pem -out $name-cert.pem -extensions san -config \
    <(echo "[req]"; 
      echo distinguished_name=req; 
      echo "[san]"; 
      echo subjectAltName=DNS:$dnsname,IP:$localip
      ) \
    -subj "/CN=$dnsname"

    ls -l $name-*.pem

    exit 0
    ;;
  esac;
  shift 2
done

echo "Usage: ./certgen.sh --name cert-name"
```

上面脚本中有 2 个变量 $localip （127.0.0.1）和 $dnsname（127.0.0.1），适用于在本机测试。如果在测试服务器上测试可以分别填写测试服务器的 DNS 和 IP 地址。保存上述脚本到文件中，并且添加可执行权限，运行上述脚本生成证书，脚本会在当前目录创建 keys 目录存放新生成的私钥和证书文件。

```bash
$ ./certgen.sh --name mycert
```

#### 启动监听

```javascript
const uuid = require('uuid')
server.listen(443, '127.0.0.1', onListening)
  // 有新连接时触发，为连接生成唯一 ID
  .on('connection', (socket /** net.Socket */) => {
    socket.id = uuid.v1().replace(/-/g, '')
    console.debug('server.onConnection', socket.id)
  })
  // 为连接创建 session 后触发
  .on('session', onSession)
  .on('sessionError', (err) => {
    console.error('server.onSessionError', err)
  })
  // 所有 http 请求（包括 http2）都会触发
  .on('request', onRequest)
  // 所有 http2 请求触发
  .on('stream', onStream)
  .on('error', (err) => {
    console.error('server.onError', err)
  })
  .on('unknownProtocol', (socket) => {
    console.debug('server.onUnknownProtocol', socket.id)
    socket.end()
  })
```

上述代码启动了服务器监听，代码上有注释，解释了几个重要的事件回调函数

* onListening

```javascript
function onListening() {
  console.info(`server listening on ${config.host}:${config.port}`)
}
```

* onSession

```javascript
function onSession(session) {
  session
  // 创建新 session，把之前创建连接时的连接 ID 作为此会话 ID
  .on('connect', (_, socket /** tls.TLSSocket */) => {
    sessions.add(session)
    session.id = session.socket.id = socket._parent.id
    console.debug('server.onSession.connect', session.id, `${sessions.size} sessions`)
  })
  .on('error', (err) => {
    console.error('server.onSession.error', session.id, err)
  })
  .on('frameError', (type, code, id) => {
    console.error('server.onSession.frameError', type, code, id, session.id)
  })
  .on('goaway', () => {
    console.debug('server.onSession.goaway', session.id)
  })
  .on('close', () => {
    sessions.delete(session)
    console.debug('server.onSession.closed', session.id, `${sessions.size} sessions`)
  })
  .on('localSettings', (settings) => {
    console.debug('server.onSession.localSettings', session.id, settings)
  })
  .on('remoteSettings', (settings) => {
    console.debug('server.onSession.remoteSettings', session.id, settings)
  })
}
```

* `onRequest` http/1.1 请求处理，后面分析
* `onStream` http/2.0 请求处理，后面分析

#### 自动重定向

服务器开启了 https，必须通过 https 访问。如果浏览器中输入 http 访问会显示错误。此时可以开启自动把 http 请求重定向到相应的 https 地址。301 表示资源地址永久转移，同时加了 connection: close 头信息，响应后连接马上断开。

```javascript
/** HTTP Server 301 redirect */
const http = require('http')
  httpServer = http.createServer((req, res) => {
    res.writeHead(301, {
      'Location': `https://127.0.0.1:443${req.url}`,
      'Connection': 'close'
    })
    res.end()
  }).listen(80, '127.0.0.1')
```

#### 请求处理

```javascript
const Client = require('./client'),
      serveRequest = require('./serve/request')(Client),
      serveStatic = require('./serve/static')(config),
      serveStream = require('./serve/stream')(Client, config)
```

如上引入远程连接管理 Client 及 3 个相应的请求处理组件 serveRequest、serveStatic、serveStream。接下来定义 2 个关键的 server 回调函数，就是前面内容提到过的 onRequest 及 onStream。

```javascript
/** HTTP/1.1 */
function onRequest(req, res, next) {
  console.debug(
    'server.onRequest',
    req.url,
    `[http/${req.httpVersion}]`,
    req._context.sid,
    req._context.reqid)
  serveRequest(req, res, next)
}
```

如上述代码，所有进来的 http/1.1 请求都委托给 serveRequest 组件进行处理。其中 next 回调可先忽略。serveRequest 会在后文中展开介绍。

```javascript
/** HTTP/2.0 */
function onStream(stream, headers, next) {
  streams.add(stream)

  const reqPath = headers[HTTP2_HEADER_PATH]
  console.debug('server.onStream', reqPath, stream.id, stream._context.sid, stream._context.reqid, `${streams.size} streams`)

  stream
  .on('close', () => {
    streams.delete(stream)
    console.debug('stream.onClose', reqPath, stream.id, stream._context.sid, stream._context.reqid, `${streams.size} streams`)
  })
  .on('frameError', (type, code, id) => {
    console.error('stream.onFrameError', type, code, id, stream.id, stream._context.sid, stream._context.reqid)
    next(stream, headers, code)
  })
  .on('error', (err) => {
    console.error('stream.onError', stream.id, stream._context.sid, stream._context.reqid, err)
    next(stream, headers, err)
  })
  .on('aborted', () => {
    console.error('stream.onAborted', stream.id, stream._context.sid, stream._context.reqid)
    next(stream, headers)
  })
  // 设置超时处理
  .setTimeout(5000, () => {
    console.error('stream.onTimeout', stream.id, stream._context.sid, stream._context.reqid)
    stream.close()
    next(stream, headers)
  })

  // hack Http2Session.respond, ignore if already closed or desdroyed
  for(let fname of ['respond', 'respondWithFile']) {
    stream[`safe_${fname}`] = ((...args) => {
      try {
        stream[fname](...args)
      } catch(err) {
        console.error(
          `server.onStream.safe_${fname}.error`,
          stream.id,
          stream._context.sid,
          stream._context.reqid)
      }
    }).bind(stream)
  }

  if(serveStatic(stream, headers, next)) {
    console.debug('server.onStream.serveStatic', stream.id, stream._context.sid, stream._context.reqid)
  } else {
    serveStream(stream, headers, next)
  }
}
```

http/2.0 协议要求一个 authority（ip+port） 只需建立一个长连接，并在此连接上维持一个会话 session，后续所有的请求都经由此 session 发送。如上述代码，2.0 请求会触发 stream 事件，并传入 stream 和 headers 参数。stream 是双向流，可读写。headers 中是本次请求的头部信息。上述代码的重点在最后一句，就是

```javascript
if(serveStatic(stream, headers, next)) {
  console.debug('server.onStream.serveStatic', stream.id, stream._context.sid, stream._context.reqid)
} else {
  serveStream(stream, headers, next)
}
```

根据请求路径判断，如果是静态资源请求则委托给 serveStatic 组件进行处理，其他请求都委托给 serveStream 组件进行处理。serveStatic 和 serveStream 会在后文中展开介绍。

#### 关闭服务器

```javascript
/**
  * server shutdown gracefully
  */
function shutdown() {
  console.debug(`close worker [${process.id}]`)
  // close server
  server.close()
  // close all sessions
  for(let session of sessions) {
    if(!session.closed) {
      session.close()
    }
  }
  // close http server
  httpServer.close()
  // close client
  Client.close()
}
```

如上先定义 shutdown 函数。http2 会保持和服务器的长连接，并在此连接上创建 session，后续所有请求都是经由同一个连接和 session 进行发送和响应。所以如果想正常停止服务器进程，首先要调用 server.close()，会停止服务器创建新的连接和 session。但是已存在的 session 还是可以继续发送新请求进来，还需要在所有 session 上调用 session.close() 关闭会话，这不会马上关闭会话，而是需要等该会话上所有请求都处理完成后才关闭。

此时已存在的连接和 session 还会继续处理之前已经进来的请求，直到所有请求处理完毕然后关闭所有会话和连接，退出服务器进程。这样在关闭服务器进程过程中，所有的请求都会被处理，浏览器端也不会报错。但是到这台服务器的新的连接请求会被拒绝。如果有2台机器的情况，可以先停掉一台发布新代码启动进程，然后再发布另一台。发布服务器期间基本不会报错。

然后调用了 httpServer.close()，停止用来自动重定向的 httpServer。然后关闭了其他资源 Client.close()。

为什么要如此关闭服务器呢？当然可以暴力强关服务器，但不利于服务器操作系统即时释放资源，而且如此的话客户端浏览器不知道连接已断开，还会继续等待比较长的一段时间，或者页面直接报错体验很差。

```javascript
process
  .on('unhandledRejection', (reason, promise) => {
    console.error(reason, 'Unhandled Rejection at Promise', promise, `@worker [${process.id}]`)
  })
  .on('uncaughtException', (err, origin) => {
    console.error('Uncaught Exception thrown', err.stack, origin, `close worker [${process.id}]`)
    shutdown()
    process.exitCode = 1
  })
  .on('message', (msg) => {
    if(msg === 'shutdown') {
      console.warn('message from master received...', msg, `close worker [${process.id}]`)
      shutdown()
    }
  })
  .once('SIGINT', function (code) {
    console.warn('SIGINT received...', code, `close worker [${process.id}]`)
    shutdown()
  })
  .once('SIGTERM', function (code) {
    console.warn('SIGTERM received...', code, `close worker [${process.id}]`)
    shutdown()
  })
```

上述代码在进程 process 上注册回调函数

1. unhandledRejection
async 函数内报错且未在 promise 上调用 .catch() 时触发，此处进行兜底处理
2. uncaughtException
程序中有未捕获的 Error 时触发，此时输出错误信息方便修复，同时关闭服务器
3. SIGINT
Ctrl+C 信号触发，关闭服务器
4. SIGTERM
kill 信号触发，关闭服务器(非强关进程 `kill -9`)，此时可通过 `kill pid` 正常关闭服务器进程