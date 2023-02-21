# 三百行代码搭建一个简单的 SOCKS5 代理服务器

SOCKS是一种网络传输协议，主要用于客户端与外网服务器之间通讯的中间传递。当防火墙后的客户端要访问外部的服务器时，就跟SOCKS代理服务器连接。这个代理服务器控制客户端访问外网的资格，允许的话，就将客户端的请求发往外部的服务器。

这个协议最初由David Koblas开发，而后由NEC的Ying-Da Lee将其扩展到SOCKS4。最新协议是SOCKS5，与前一版本相比，增加支持UDP、验证，以及IPv6。根据OSI模型，SOCKS是会话层的协议，位于表示层与传输层之间。

SOCKS工作在比HTTP代理更低的层次：SOCKS使用握手协议来通知代理软件其客户端试图进行的SOCKS连接，然后尽可能透明地进行操作，而常规代理可能会解释和重写报头（例如，使用另一种底层协议，例如FTP；然而，HTTP代理只是将HTTP请求转发到所需的HTTP服务器）。虽然HTTP代理有不同的使用模式，HTTP CONNECT方法允许转发TCP连接；然而，SOCKS代理还可以转发UDP流量（仅SOCKS5），而HTTP代理不能。HTTP代理通常更了解HTTP协议，执行更高层次的过滤（虽然通常只用于GET和POST方法，而不用于CONNECT方法）。

_(以上摘自维基百科)_

# 协议部分

SOCKS5 协议是 1996 年发布的，迄今为止已经 25 年了，很多软件内部都支持 SOCKS5 代理。比如浏览器一般都会支持，方便有些用户通过 SOCK5 代理浏览网页。我在开发过程中也使用 CURL 进行过测试，也是支持的。也可以在系统设置里面找到代理设置。在软件开发测试领域，也有很多工具软件支持。需要测试的设备连上代理服务器后，所有流量都会经过相关工具软件，可以方便测试、调试软件。

如通过 Charles 可以查看 HTTP Request/Response 报文信息（非 HTTPS 网站）。可以统计 URL 访问次数、返回延时、数据包大小等非常有用的信息。代理服务器本身不会查看和修改所转发的数据，只是进行简单的转发。甚至不理解转发的内容，现在大部分网站都是基于 HTTPS 加密的，代理服务器没办法了解具体通信内容。当然如果是 HTTP 网站，数据都是明文传输，所有数据都可以看到。

本文主要介绍 [SOCKS Protocol Version 5](https://www.ietf.org/rfc/rfc1928.txt)协议，原文比较简单，规定的是应用程序如何与代理服务器进行握手，握手成功后就行流量的正向及反向转发。接下来会先介绍下 SOCKS5 协议，然后通过 Node.js 实现一个简单的代理服务器。

## 1.协议从协商认证方法开始
```
/**
 * The client connects to the server, and sends a version identifier/method selection message:
 * +----+----------+----------+
 * |VER | NMETHODS | METHODS  |
 * +----+----------+----------+
 * | 1  |    1     | 1 to 255 |
 * +----+----------+----------+
 */
```
客户端应用程序连接上服务器，发送协议版本号、客户端支持的认证方法列表给服务器。上面注释中数字代表字节长度，比如 VER 为 1 字节长，后文中出现类似的数字都代表字节长度，不再赘述。

* VER 字段是 SOCKS 协议版本号，传 X'05'，1 字节长
* NMETHODS 字段是客户端支持的认证方法数量，每种认证方法用 1 字节进行编码，所以也决定了 METHODS 字段的长度
* METHODS 依次写入支持的认证方法编码

## 2.服务器根据客户端上报的方法列表选择，回复方法编码
```
/**
 * The server selects from one of the methods given in METHODS, and sends a METHOD selection message
 * +----+--------+
 * |VER | METHOD |
 * +----+--------+
 * | 1  |   1    |
 * +----+--------+
 */
```

* VER 字段是 SOCKS 协议版本号，传 X'05'，1 字节长
* METHOD 字段是服务器选择的认证方法编码, 1 字节长

下面定义了认证方法的编码 `X'00'` 代表 1 字节长的十六进制数字

* X'00' NO AUTHENTICATION REQUIRED // 不需要认证
* X'01' GSSAPI // GSSAPI 认证，在 [RFC 1961](https://tools.ietf.org/html/rfc1961) 里规定
* X'02' USERNAME/PASSWORD // 用户名密码认证，在 [RFC 1929](https://www.ietf.org/rfc/rfc1929.txt) 里规定
* X'03' to X'7F' IANA ASSIGNED // 未分配
* X'80' to X'FE' RESERVED FOR PRIVATE METHODS // 保留
* X'FF' NO ACCEPTABLE METHODS // 服务器发现客户端上报的方法列表都不合适时回复 X'FF'

一般服务器必须实现 GSSAPI 方法，建议实现 USERNAME/PASSWORD。比较安全的环境也可以不用认证。

## 3.客户端认证成功后发送请求信息
```
/**
 * The SOCKS request is formed as follows:
 * +----+-----+-------+------+----------+----------+
 * |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
 * +----+-----+-------+------+----------+----------+
 * | 1  |  1  | X'00' |  1   | Variable |    2     |
 * +----+-----+-------+------+----------+----------+
 */
```

* VER 字段是 SOCKS 协议版本号，传 X'05'，1 字节长
* CMD 字段是请求命令，1字节长
  * CONNECT X'01' 代表连接目标服务器命令，本文会实现 CONNECT 请求命令
  * BIND X'02' 暂时不介绍
  * UDP ASSOCIATE X'03' 暂时不介绍
* RSV 保留字节，传 X'00'，1字节长
* ATYP 目标地址类型，1字节长
  * IP V4 address: X'01' IPV4
  * DOMAINNAME: X'03' 域名
  * IP V6 address: X'04' IPV6
* DST.ADDR 客户端想要请求的目标地址
* DST.PORT 客户端想连接的目标端口号，2 字节长，网络字节序

DST.ADDR 字段根据地址类型不同，长度不同。
* 如果IPV4 是固定 4 字节长。
* 如果是域名，则首个字节代表域名的长度（字节数），接下来的可变长度是域名字符串

## 4.服务端评估该请求并回复

```
/**
 * The server evaluates the request, and returns a reply formed as follows:
 * +----+-----+-------+------+----------+----------+
 * |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
 * +----+-----+-------+------+----------+----------+
 * | 1  |  1  | X'00' |  1   | Variable |    2     |
 * +----+-----+-------+------+----------+----------+
 */
```

* VER 字段是 SOCKS 协议版本号，传 X'05'，1 字节长
* REP 回复字段:
    * X'00' 成功
    * X'01' general SOCKS server failure
    * X'02' connection not allowed by ruleset
    * X'03' Network unreachable
    * X'04' Host unreachable
    * X'05' Connection refused
    * X'06' TTL expired
    * X'07' 请求命令不支持
    * X'08' 地址类型不支持
    * X'09' to X'FF' 未分配
* RSV 保留字段，传 X'00'
* ATYP 地址类型
    * IP V4 address: X'01'
    * DOMAINNAME: X'03'
    * IP V6 address: X'04'
* BND.ADDR 代理服务器连接目标服务器的地址
* BND.PORT 代理服务器连接目标服务器的端口，2 字节长，网络字节序

服务器可以评估是否允许该请求，如果允许则代表客户端连接目标服务器( DST.ADDR:DST.PORT )，连接成功后回复客户端 REP 为成功，握手过程完成，后面就可以开始代理转发数据了。如果评估不允许该请求或者因为其他原因不能连接上目标服务器，则需返回对应回复信息。

_(以上为协议部分)_

# 代码部分

## consts.js 常量定义文件

```javascript
/** socks 5 */
const SOCKS_VERSION = 0x05,

      STATE = {
        METHOD_NEGOTIATION: 0x00,
        AUTHENTICATION: 0x01,
        REQUEST_CONNECT: 0x02,
        PROXY_FORWARD: 0x03
      },

      /**
       * o  X'00' NO AUTHENTICATION REQUIRED
       * o  X'01' GSSAPI
       * o  X'02' USERNAME/PASSWORD
       * o  X'03' to X'7F' IANA ASSIGNED
       * o  X'80' to X'FE' RESERVED FOR PRIVATE METHODS
       * o  X'FF' NO ACCEPTABLE METHODS
       */
      METHODS = {
        NO_AUTH: [0x00, 'no_auth'],
        GSSAPI: [0x01, 'gssapi'],
        USERNAME_PASSWD: [0x02, 'username_password'],
        NO_ACCEPTABLE: [0xFF, 'no_acceptable_methods'],

        get(method) {
          switch(method) {
            case this.NO_AUTH[0]:
              return this.NO_AUTH
            case this.GSSAPI[0]:
              return this.GSSAPI
            case this.USERNAME_PASSWD[0]:
              return this.USERNAME_PASSWD
          }
          console.error(`method [${method}] is not supported`)
          return false
        }
      },

      /**
       * o  CONNECT X'01'
       * o  BIND X'02'
       * o  UDP ASSOCIATE X'03'
       */
      REQUEST_CMD = {
        CONNECT: [0x01, 'connect'],
        BIND: [0x02, 'bind'],
        UDP_ASSOCIATE: [0x03, 'udp_associate'],

        get(cmd) {
          switch(cmd) {
            case this.CONNECT[0]:
              return this.CONNECT
            case this.BIND[0]:
              return this.BIND
            case this.UDP_ASSOCIATE[0]:
              return this.UDP_ASSOCIATE
          }
          console.error(`cmd [${cmd}] is not supported`)
          return false
        }
      },

      /** reserved byte value */
      RSV = 0x00,

      /**
       * o  IP V4 address: X'01'
       * o  DOMAINNAME: X'03'
       * o  IP V6 address: X'04'
       */
      ATYP = {
        IPV4: [0x01, 'ipv4'],
        FQDN: [0x03, 'domain name'],
        IPV6: [0x04, 'ipv6'],

        get(atyp) {
          switch(atyp) {
            case this.IPV4[0]:
              return this.IPV4
            case this.FQDN[0]:
              return this.FQDN
            case this.IPV6[0]:
              return this.IPV6
          }
          console.error(`atpy [${atyp}] is not supported`)
          return false
        }
      },

      /**
       * o  X'00' succeeded
       * o  X'01' general SOCKS server failure
       * o  X'02' connection not allowed by ruleset
       * o  X'03' Network unreachable
       * o  X'04' Host unreachable
       * o  X'05' Connection refused
       * o  X'06' TTL expired
       * o  X'07' Command not supported
       * o  X'08' Address type not supported
       * o  X'09' to X'FF' unassigned
       */
      REP = {
        SUCCEEDED: [0x00, 'succeeded'],
        GENERAL_FAILURE: [0x01, 'general SOCKS server failure'],
        NOT_ALLOWED: [0x02, 'connection not allowed by ruleset'],
        NETWORK_UNREACHABLE: [0x03, 'Network unreachable'],
        HOST_UNREACHABLE: [0x04, 'Host unreachable'],
        CONNECTION_REFUSED: [0x05, 'Connection refused'],
        TTL_EXPIRED: [0x06, 'TTL expired'],
        COMMAND_NOT_SUPPORTED: [0x07, 'Command not supported'],
        ADDRESS_TYPE_NOT_SUPPORTED: [0x08, 'Address type not supported']
      },

      /**
       * The VER field contains the current version of the subnegotiation, which is X'01'
       * username/password auth version
       */
      USERNAME_PASSWD_AUTH_VERSION = 0x01,

      /**
       * auth status
       */
      AUTH_STATUS = {
        SUCCESS: 0x00,
        FAILURE: 0X01
      }

module.exports = () => {

  return {
    SOCKS_VERSION: SOCKS_VERSION,
    STATE: STATE,
    METHODS: METHODS,
    REQUEST_CMD: REQUEST_CMD,
    RSV: RSV,
    ATYP: ATYP,
    REP: REP,
    USERNAME_PASSWD_AUTH_VERSION, USERNAME_PASSWD_AUTH_VERSION,
    AUTH_STATUS: AUTH_STATUS
  }
  
}
```

## config.js 程序配置

```javascript
const consts = require('../consts')(),
      app = {
        port: 3000,
        host: '0.0.0.0',
        auth_method: consts.METHODS.NO_AUTH
      }

module.exports = () => {
  return app
}
```

如上，代理服务器监听在 3000 端口，采用无认证方式

## app.js 程序启动

```javascript
const net = require('net'),
      server = new net.createServer(),
      config = require('./config')(server),
      Proxy = require('./proxy')(server)

server.listen(config.port, config.host)
      .on('listening', () => {
        console.log(`simple-proxy server listening on ${config.port}`)
      })
      .on('close', () => {
        console.log('simple-proxy server closed')
      })
      .on('error', err => {
        console.error('simple-proxy server throw error', err)
      })
      .on('connection', socket => {
        var proxy = Proxy(socket)
        // data package come in
        socket.on('data', buf => {
          proxy.handle(buf)
        })
        .on('end', () => {
          console.log(`socket ${proxy._session.id} end`)
        })
        .on('close', hadError => {
          console.log(`socket ${proxy._session.id} closed with error ${hadError}`)
        })
        .on('error', err => {
          console.error(`socket ${proxy._session.id} throw error`, err)
        })
        .on('timeout', () => {
          console.log(`socket ${proxy._session.id} timeout`)
        })
      
      })
```
可以看到其中依赖了核心 Proxy 类完成代理功能

## proxy.js 核心代理实现逻辑

下面代码是核心代理逻辑，也实现了简单的用户名密码认证方式。如果去掉用户名密码认证和注释，核心代码就 300 行左右。Proxy 类主要是实现了 SOCKS5 握手协议，与目标服务器建立连接进行数据转发。每个方法上面都有注释，应该比较容易看懂。

```javascript
const net = require('net'),
      dns = require('dns'),
      { assert } = require('console'),
      uuid = require('uuid'),

      utils = require('../utils'),
      consts = require('../consts')(),
      config = require('./config')()

function Proxy(socket) {
  return {
    /**
     * proxy socket
     */
    _socket: socket,

    /**
     * session
     */
    _session: {
      id: uuid.v1(),
      buffer: Buffer.alloc(0),
      offset: 0,
      state: consts.STATE.METHOD_NEGOTIATION
    },

    /**
     * The client connects to the server, and sends a version identifier/method selection message:
     * +----+----------+----------+
     * |VER | NMETHODS | METHODS  |
     * +----+----------+----------+
     * | 1  |    1     | 1 to 255 |
     * +----+----------+----------+
     */
    parseMethods() {
      let buf = this._session.buffer
      let offset = this._session.offset
  
      var checkNull = offset => {
        return typeof buf[offset] === undefined
      }
  
      if(checkNull(offset)) {
        return false
      }
      let socksVersion = buf[offset++]
      assert(socksVersion == consts.SOCKS_VERSION, `socket ${this._session.id} only support socks version 5, got [${socksVersion}]`)
      if(socksVersion != consts.SOCKS_VERSION) {
        this._socket.end()
        return false
      }
      
      if(checkNull(offset)) {
        return false
      }
      let methodLen = buf[offset++]
      assert(methodLen >= 1 && methodLen <= 255, `socket ${this._session.id} methodLen's value [${methodLen}] is invalid`)
  
      if(checkNull(offset + methodLen - 1)) {
        return false
      }
      let methods = []
      for(let i = 0; i < methodLen; i++) {
        let method = consts.METHODS.get(buf[offset++])
        if (!!method) {
          methods.push(method)
        }
      }
  
      console.log(`socket ${this._session.id} SOCKS_VERSION: ${socksVersion}`)
      console.log(`socket ${this._session.id} METHODS: `, methods)
  
      this._session.offset = offset
  
      return methods
    },

    /** socks server select auth method */
    selectMethod(methods) {
      let method = consts.METHODS.NO_ACCEPTABLE
      for(let i = 0; i < methods.length; i++) {
        if (methods[i] == config.auth_method) {
          method = config.auth_method
        }
      }
    
      console.log(`SELECT METHOD [${method}]`)

      this._session.method = method
    
      return method
    },

    /**
     * The server selects from one of the methods given in METHODS, and sends a METHOD selection message
     * +----+--------+
     * |VER | METHOD |
     * +----+--------+
     * | 1  |   1    |
     * +----+--------+
     * @param {*} method auth method selected
     */
    replyMethod(method) {
      this._socket.write(Buffer.from([consts.SOCKS_VERSION, method[0]]))
    },

    /**
     * This begins with the client producing a Username/Password request:
     * +----+------+----------+------+----------+
     * |VER | ULEN |  UNAME   | PLEN |  PASSWD  |
     * +----+------+----------+------+----------+
     * | 1  |  1   | 1 to 255 |  1   | 1 to 255 |
     * +----+------+----------+------+----------+
     */
    parseUsernamePasswd() {
      let buf = this._session.buffer
      let offset = this._session.offset
  
      var req = {}

      var checkNull = offset => {
        return typeof buf[offset] === undefined
      }
  
      if(checkNull(offset)) {
        return false
      }
      let authVersion = buf[offset++]
      assert(authVersion == consts.USERNAME_PASSWD_AUTH_VERSION,
        `socket ${this._session.id} only support auth version ${consts.USERNAME_PASSWD_AUTH_VERSION}, got [${authVersion}]`)
      if(authVersion != consts.USERNAME_PASSWD_AUTH_VERSION) {
        this._socket.end()
        return false
      }

      if(checkNull(offset)) {
        return false
      }
      let uLen = buf[offset++]
      assert(uLen >= 1 && uLen <= 255, `socket ${this._session.id} got wrong ULEN [${uLen}]`)
      if(uLen >= 1 && uLen <= 255) {
        if(checkNull(offset + uLen - 1)) {
          return false
        }
        req.username = buf.slice(offset, offset + uLen).toString('utf8')
        offset += uLen
      } else {
        this._socket.end()
        return false
      }

      if(checkNull(offset)) {
        return false
      }
      let pLen = buf[offset++]
      assert(pLen >= 1 && pLen <= 255, `socket ${this._session.id} got wrong PLEN [${pLen}]`)
      if(pLen >= 1 && pLen <= 255) {
        if(checkNull(offset + pLen - 1)) {
          return false
        }
        req.passwd = buf.slice(offset, offset + pLen).toString('utf8')
        offset += pLen
      } else {
        this._socket.end()
        return false
      }

      this._session.offset = offset

      return req
    },

    /**
     * The server verifies the supplied UNAME and PASSWD, and sends the following response:
     *  +----+--------+
     *  |VER | STATUS |
     *  +----+--------+
     *  | 1  |   1    |
     *  +----+--------+
     */
    replyAuth(succeeded) {
      let reply = [
        consts.USERNAME_PASSWD_AUTH_VERSION,
        succeeded ? consts.AUTH_STATUS.SUCCESS : consts.AUTH_STATUS.FAILURE
      ]
      if (succeeded) {
        this._socket.write(Buffer.from(reply))
      } else {
        this._socket.end(Buffer.from(reply))
      }
    },

    /**
     * The SOCKS request is formed as follows:
     * +----+-----+-------+------+----------+----------+
     * |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
     * +----+-----+-------+------+----------+----------+
     * | 1  |  1  | X'00' |  1   | Variable |    2     |
     * +----+-----+-------+------+----------+----------+
     */
    parseRequests() {
      let buf = this._session.buffer
      let offset = this._session.offset
    
      let req = {}
    
      var checkNull = offset => {
        return typeof buf[offset] === undefined
      }
    
      if(checkNull(offset)) {
        return false
      }
      let socksVersion = buf[offset++]
      assert(socksVersion == consts.SOCKS_VERSION, `socket ${this._session.id} only support socks version 5, got [${socksVersion}]`)
      if(socksVersion != consts.SOCKS_VERSION) {
        this._socket.end()
        return false
      }
    
      if(checkNull(offset)) {
        return false
      }
      req.cmd = consts.REQUEST_CMD.get(buf[offset++])
      if(!req.cmd || req.cmd != consts.REQUEST_CMD.CONNECT) {
        // 不支持的 cmd || 暂时只支持 connect
        this._socket.end()
        return false
      }
    
      if(checkNull(offset)) {
        return false
      }
      req.rsv = buf[offset++]
      assert(req.rsv == consts.RSV, `socket ${this._session.id} rsv should be ${consts.RSV}`)
    
      if(checkNull(offset)) {
        return false
      }
      req.atyp = consts.ATYP.get(buf[offset++])
      if(!req.atyp) {
        // 不支持的 atyp
        this._socket.end()
        return false
      } else if(req.atyp == consts.ATYP.IPV4) {
        let ipLen = 4
        if(checkNull(offset + ipLen - 1)) {
          return false
        }
        req.ip = `${buf[offset++]}.${buf[offset++]}.${buf[offset++]}.${buf[offset++]}`
      } else if(req.atyp == consts.ATYP.FQDN) {
        if(checkNull(offset)) {
          return false
        }
        let domainLen = buf[offset++]
        if(checkNull(offset + domainLen - 1)) {
          return false
        }
        req.domain = buf.slice(offset, offset + domainLen).toString('utf8')
        offset += domainLen
      } else {
        // 其他暂时不支持
        this._socket.end()
        return false
      }
    
      let portLen = 2
      if(checkNull(offset + portLen - 1)) {
        return false
      }
      req.port = buf.readUInt16BE(offset)
      offset += portLen
    
      console.log(`socket ${this._session.id} parse requests succeeded`, req)

      this._session.offset = offset
    
      return req
    },

    /**
     * The server evaluates the request, and returns a reply formed as follows:
     * +----+-----+-------+------+----------+----------+
     * |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
     * +----+-----+-------+------+----------+----------+
     * | 1  |  1  | X'00' |  1   | Variable |    2     |
     * +----+-----+-------+------+----------+----------+
     * @param {*} req client requests
     */
    dstConnect(req) {
      let dstHost = req.domain || req.ip
      dns.lookup(dstHost, { family: 4 }, (err, ip) => {
        if (err || !ip) {
          // failure reply
          let reply = [
            consts.SOCKS_VERSION,
            consts.REP.HOST_UNREACHABLE[0],
            consts.RSV,
            consts.ATYP.IPV4[0]
          ]
          .concat(utils.ipbytes('127.0.0.1')) // ip: 127.0.0.1
          .concat([0x00, 0x00]) // port: 0x0000
          // close connection
          this._socket.end(Buffer.from(reply))
        } else {
          // connect target host
          const dstSocket = net.createConnection({
            port: req.port, // port from client's requests
            host: ip // ip from dns lookup of socks proxy server
          })
    
          dstSocket.on('connect', () => {
            // success reply
            let bytes = [
              consts.SOCKS_VERSION,
              consts.REP.SUCCEEDED[0],
              consts.RSV,
              consts.ATYP.IPV4[0]
            ]
            // dstSocket.localAddress or default 127.0.0.1
            .concat(utils.ipbytes(dstSocket.localAddress || '127.0.0.1'))
            // default port 0x00
            .concat([0x00, 0x00])
            
            let reply = Buffer.from(bytes)
    
            // use dstSocket.localPort override default port 0x0000
            reply.writeUInt16BE(dstSocket.localPort, reply.length - 2)
    
            this._socket.write(reply)
    
            // pipe for proxy forward
            this._socket.pipe(dstSocket).pipe(this._socket)
          })
          .on('error', err => {
            console.error(`socket ${this._session.id} -> dstSocket`, err)
          })
          .on('end', () => {
            console.log(`socket ${this._session.id} -> dstSocket end`)
          })
          .on('close', () => {
            console.log(`socket ${this._session.id} -> dstSocket close`)
          })
    
          // save dstSocket to session
          this._session.dstSocket = dstSocket
        }
      })
    },

    /**
     * called by socket's 'data' event listener
     * @param {Buffer} buf data buffer
     */
    handle(buf) {
      // before proxy forward phase, otherwise do nothing
      if(this._session.state < consts.STATE.PROXY_FORWARD) {
        // append data to session.buffer
        this._session.buffer = Buffer.concat([this._session.buffer, buf])
        
        // discard processed bytes and move on to the next phase
        const discardProcessedBytes = (nextState) => {
          this._session.buffer = this._session.buffer.slice(this._session.offset)
          this._session.offset = 0
          this._session.state = nextState
        }
      
        switch(this._session.state) {
          case consts.STATE.METHOD_NEGOTIATION:
            let methods = this.parseMethods()
            if(!!methods) { // read complete data
              let method = this.selectMethod(methods)
              this.replyMethod(method)
              switch(method) {
                case consts.METHODS.USERNAME_PASSWD:
                  discardProcessedBytes(consts.STATE.AUTHENTICATION)
                  break
                case consts.METHODS.NO_AUTH:
                  discardProcessedBytes(consts.STATE.REQUEST_CONNECT)
                  break
                case consts.METHODS.NO_ACCEPTABLE:
                  this._socket.end()
                  break
                default:
                  this._socket.end()
              } 
            }
            break
          // curl www.baidu.com --socks5 127.0.0.1:3000 --socks5-basic --proxy-user  oiuytre:yhntgbrfvedc
          case consts.STATE.AUTHENTICATION:
            // add gssapi support
            // need check this._session.method for parse data
            let userinfo = this.parseUsernamePasswd()
            if(!!userinfo) { // read complete data
              let succeeded = (
                userinfo.username === config.username &&
                userinfo.passwd === config.passwd
              )
              discardProcessedBytes(
                succeeded ? consts.STATE.REQUEST_CONNECT : consts.STATE.AUTHENTICATION
              )
              this.replyAuth(succeeded)
            }
            break
          case consts.STATE.REQUEST_CONNECT:
            let req = this.parseRequests()
            if(!!req) { // read complete data
              this.dstConnect(req)
              discardProcessedBytes(consts.STATE.PROXY_FORWARD)
            }
            break
          case consts.STATE.PROXY_FORWARD:
          default:
            console.log(`handle state [${this._session.state}]`, this._session)
        }
      }
    }
  }
}

module.exports = server => {
  return Proxy
}
```

启动服务

```bash
$ node app.js
```

此时打开 Chrome 浏览器，安装 `SwitchyOmega` 插件，配置代理

![配置代理](https://cdn.huoyijie.cn/ab/43d8648090fe11ebafa22393c04133cd/SwitchyOmega.png)

配置完成后，后面所有的请求都是通过 SOCKS5 服务器 127.0.0.1:3000 转发完成的。本文主要内容介绍完了，如果想要一份完整的项目代码，可以与我联系。