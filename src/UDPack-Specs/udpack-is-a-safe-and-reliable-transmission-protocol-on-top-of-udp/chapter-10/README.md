# 第十章 UDPack 应用之实时双向通信框架

应用 UDPack 协议实现一个实时双向通信框架是很容易的，首先要考虑的是引入一个通用的对象序列化协议，可在对象与字节数组之间相互转换。而 [msgpack](http://msgpack.org/) 正是这样的一个对象序列化协议，类似 JSON，但是速度更快，序列化后数据更小。msgpack 有各种语言的实现版本，也包含 javascript 版本。

下面给出该通信框架(udpack.io)的全部代码:

```javascript
const EventEmitter = require('events')
const UDPack = require('udpack')()
const log4js = require('log4js')
const logger = log4js.getLogger('udpack.io')
const { Encoder, Decoder } = require('@msgpack/msgpack')
const encoder = new Encoder()
const decoder = new Decoder()

const STATE = {
  READY: 0,
  CLOSE: -1
}

class Socket extends EventEmitter {
  constructor(session, opts) {
    opts = opts || {}
    super(opts)
    this._$session = session
    this.opts = opts
    session._$udpackio_socket = this
    this.state = STATE.READY
  }

  close() {
    if (this.state !== STATE.CLOSE) {
      this._$session.close()
      delete this._$session._$udpackio_socket
      delete this._$session
    }
    return this
  }
}

class Channel extends EventEmitter {
  constructor(stream, opts) {
    opts = opts || {}
    super(opts)
    this._$stream = stream
    this.opts = opts
    stream._$udpackio_channel = this
    this.state = STATE.READY
  }

  send(event, data, cb) {
    let encoded = encoder.encode([event, data])
    let buffer = Buffer.from(encoded)
    logger.debug(this._$stream.session.netSessionIdStr(), this._$stream.id, 'length', buffer.length)
    this._$stream.write(buffer, cb)
  }

  close() {
    if (this.state !== STATE.CLOSE) {
      this._$stream.close()
      delete this._$stream._$udpackio_channel
      delete this._$stream
    }
    return this
  }
}

class Server extends EventEmitter {
  constructor(opts) {
    opts = opts || {}
    super(opts)
    this.opts = opts
    const { port, host, pubkey, prikey, authToken } = opts
    this.udpack = new UDPack({ port, host, pubkey, prikey, authToken })

    this.udpack.on('listening', () => {
      logger.debug('listenling on', this.udpack.endpoint)
      this.emit('listening')
    }).on('session', (session) => {
      let socket = new Socket(session)
      session.on('goaway', (code) => {
        logger.debug(`[session-${session.netSessionIdStr()}] goaway [${code}]`)
        socket.close().emit('close')
      })
      this.emit('connection', socket)
    }).on('stream', (stream) => {
      let channel = new Channel(stream)
      stream.on('data', (buffer) => {
        logger.debug(stream.session.netSessionIdStr(), stream.id, 'length', buffer.length)
        try {
          let [event, data] = decoder.decode(buffer)
          channel.emit(event, data)
        } catch (e) {
          logger.error(stream.session.netSessionIdStr(), stream.id, buffer, e)
        }
      }).on('close', () => {
        channel.close().emit('close')
      })
      this.emit('channel', channel)
    }).on('speed', (speed) => {
      this.emit('speed', speed)
    })

    this.state = STATE.READY
  }

  close() {
    if (this.state !== STATE.CLOSE) {
      this.udpack.close()
    }
    return this
  }
}

class Client extends EventEmitter {
  constructor(opts) {
    opts = opts || {}
    super(opts)
    this.opts = opts
    const { port, host, pubkey, token } = opts
    this.udpack = new UDPack({ port, host, pubkey, token })

    this.udpack.on('listening', () => {
      logger.debug('listenling on', this.udpack.endpoint)
      this.emit('listening')
    })
  }

  connect(serverPort, serverAddr, cb) {
    this.udpack.on('session', (session) => {
      let socket = new Socket(session)
      session.on('goaway', (code) => {
        logger.debug(`[session-${session.netSessionIdStr()}] goaway [${code}]`)
        socket.close().emit('close')
        this.udpack.connect(serverPort, serverAddr, () => {
          logger.debug('reconnect to server...')
        })
      })
      this.emit('connection', socket)
    })
    .on('connect_error', (session) => {
      this.emit('connect_error')
    })
    .on('speed', (speed) => {
      this.emit('speed', speed)  
    })
    .connect(serverPort, serverAddr, cb)
  }

  openChannel(cb) {
    this.udpack.openStream((err, stream) => {
      if (err) {
        cb(err)
      } else {
        let channel = new Channel(stream)
        stream.on('data', (buffer) => {
          logger.debug(stream.session.netSessionIdStr(), stream.id, 'length', buffer.length)
          try {
            let [event, data] = decoder.decode(buffer)
            channel.emit(event, data)
          } catch (e) {
            logger.error(stream.session.netSessionIdStr(), stream.id, buffer, e)
          }
        }).on('close', () => {
          channel.close().emit('close')
        })
        cb(null, channel)
      }
    })
  }

  close() {
    if (this.state !== STATE.CLOSE) {
      this.udpack.close()
    }
    return this
  }
}

module.exports = {
  Client, Server
}
```

可以看到该框架把 UDPack 协议中的 Session 和 Stream 对象分别封装为了 Socket 和 Channel 对象。并用 UDPack 对象分别实现了 Server 以及 Client 对象。

下面给出一段利用上述框架实现文件上传的示例代码:

## 文件上传服务端

```javascript
const fs = require('fs')
const { Server } = require('udpack.io')

const env = 'development'
const level = env === 'development' ? 'debug' : 'info'
const log4js = require('log4js')
log4js.configure({
  appenders: {
    udpack: { type: 'file', filename: 'logs/udpack.log' },
    'udpack.io': { type: 'file', filename: 'logs/udpack.io.log' },
    server: { type: 'file', filename: 'logs/server.log' },
  },
  categories: {
    udpack: { appenders: ['udpack'], level },
    'udpack.io': { appenders: ['udpack.io'], level },
    default: { appenders: ['server'], level },
  }
})
const logger = log4js.getLogger('server')

let port = 45456,
  host = '0.0.0.0',
  pubkey = fs.readFileSync(`${__dirname}/keys/pubkey.pem`),
  prikey = fs.readFileSync(`${__dirname}/keys/prikey.pem`),
  authToken = (token) => (token === 'ea970f20b09311ebbe1d11acebc10dbf')

const opts = { port, host, pubkey, prikey, authToken }
const server = new Server(opts)

let total = 0
let received = 0
console.time('used')

let cb = (err) => {
  if (err) {
    logger.error(err)
    server.close()
  }
}

server.on('channel', (channel) => {
  channel.on('file', (data) => {
    let writeStream = fs.createWriteStream(data.filename)
    writeStream.on('error', cb)
    total = data.filesize

    channel.on('chunk', (data) => {
      logger.debug('received data ', data.chunk.length, ' bytes')
      received += data.chunk.length
      writeStream.write(data.chunk, cb)
    }).on('finish', () => {
      logger.debug('recv finished')
      console.timeEnd('used')
      writeStream.end()
      channel.send('done', {}, cb)
    }).on('done', () => {
      server.close()
    })
  })
}).on('speed', (speed) => {
  console.clear()
  console.debug(`[${((received / (total || 0.00001)) * 100).toFixed(2)}%]`, `[recv.speed ${speed.recv} KB/s]`)
})
```

通过 `const opts = { port, host, pubkey, prikey, authToken }` 创建 Server 对象实例后，监听 channel 事件，并通过 channel 实例读取上传文件数据。这里用了自定义的文件上传协议。当客户端向服务端上传文件时，会触发 channel 对象的 file 事件，此时可打开一个本地文件。然后监听 chunk 和 finish 事件，通过 chunk 事件可获取文件数据流。finish 事件代表文件传输已完成。

## 文件上传客户端

```javascript
const fs = require('fs')
const { Client } = require('udpack.io')
const path = require('path')

const env = 'development'
const level = env === 'development' ? 'debug' : 'info'
const log4js = require('log4js')
log4js.configure({
  appenders: {
    udpack: { type: 'file', filename: 'logs/udpack-client.log' },
    'udpack.io': { type: 'file', filename: 'logs/udpack.io-client.log' },
    client: { type: 'file', filename: 'logs/client.log' },
  },
  categories: {
    udpack: { appenders: ['udpack'], level },
    'udpack.io': { appenders: ['udpack.io'], level },
    default: { appenders: ['client'], level },
  }
})
const logger = log4js.getLogger('client')

let serverAddr = '127.0.0.1',
  serverPort = 45456,
  port = 0, // os will select a random port
  host = '0.0.0.0',
  pubkey = fs.readFileSync(`${__dirname}/keys/pubkey.pem`),
  token = 'ea970f20b09311ebbe1d11acebc10dbf'

const opts = { port, host, pubkey, token }
const client = new Client(opts)

client.connect(serverPort, serverAddr, (err) => {
  if (err) {
    logger.error(err)
    client.close()
  }
})

let total = 0
let sent = 0
console.time('used')

let cb = (err) => {
  if (err) {
    logger.error(err)
    client.close()
  }
}

client.on('connection', (socket) => {
  client.openChannel((err, channel) => {
    if (err) {
      cb(err)
    } else {
      let args = process.argv.slice(2)
      let pathname = args[0]
      let filename = path.basename(pathname)
      let filesize = fs.statSync(pathname).size
      total = filesize
      channel.send('file', { filename, filesize }, (err) => {
        if (err) {
          cb(err)
        } else {
          const readStream = fs.createReadStream(pathname, { highWaterMark: 1024 * 1024 })
          readStream.on('data', (chunk) => {
            readStream.pause()
            channel.send('chunk', { chunk }, (err) => {
              if (err) {
                cb(err)
              } else {
                logger.debug('send data ', chunk.length, ' bytes')
                sent += chunk.length
                readStream.resume()
              }
            })
          }).on('end', () => {
            logger.debug('send finished')
            console.timeEnd('used')
            channel.send('finish')
          }).on('error', cb)
    
          channel.on('done', () => {
            channel.send('done', {}, (err) => {
              if (err) {
                logger.error(err)
              }
              client.close()
            })
          })
        }
      })
    }
  })
}).on('connect_error', () => {
  logger.debug('connect error')
  client.close()
}).on('speed', (speed) => {
  console.clear()
  console.debug(`[${((sent / (total || 0.00001)) * 100).toFixed(2)}%]`, `[send.speed ${speed.send} KB/s]`)
})
```

通过 `const opts = { port, host, pubkey, token }` 创建 Client 实例后，可调用 `client.connect` 连接服务器端，当连接成功后会触发 connection 事件，并获得 socket 对象实例。触发 connection 事件代表连接已完成，此时可通过调用 client.openChannel 方法打开新的 channel，最后通过此 channel 对象与服务端进行通信。

客户端打开 channel 后，首先读取本地文件，并向服务器发送 file 事件，传递文件名和文件大小。然后开始向服务器发送 chunk 事件，传输文件数据。完成后向服务器发送 finish 事件。