# 远程连接管理

本节内代码位于 `client.js` 文件

```javascript
const http = require('http'),
      config = require('./config')('server-config.json')

// 创建 client
function _client_1_1(authority) {
  // http/1.1
  var settings = {
    agent: new http.Agent({ keepAlive: true }),
    host: authority.split(':')[0],
    port: authority.split(':')[1] || 80,
    protocol: `${config.scheme}:`
  }

  var client = {
    authority: authority,
    settings: settings,

    request(headers, options) {
      headers['host'] = this.settings.host

      let opts = {
        headers: headers,
        method: headers['method'],
        path: headers['path']
      }
      
      for(let prop in this.settings) {
        opts[prop] = this.settings[prop]
      }

      let req = http.request(opts)
        .on('response', options.onres || console.debug)
        .on('upgrade', options.onupgrade || console.debug)
      if(options.endStream) {
        req.end()
      }
      return req
    },

    // 当不再使用时最好 destroy() Agent 实例，因为未使用的套接字会消耗操作系统资源
    close() {
      this.settings.agent.destroy()
    }
  }

  return client
}
```

如上述代码，调用 _client_1_1，传入 authority（ip:port），可针对指定 authority 创建一个远程连接管理对象，负责维持连接池、发送请求并处理响应。http.Agent 负责管理 HTTP 客户端的连接持久性和重用。它为给定的主机和端口维护一个待处理请求队列，为每个请求重用单独的套接字连接，直到队列为空，此时套接字被销毁或放入连接池，以便再次用于请求到同一个主机和端口。销毁还是放入连接池取决于 keepAlive 选项。

连接池中的连接已启用 TCP Keep-Alive，但服务器仍可能关闭空闲连接，在这种情况下，它们将从连接池中删除，并且当为该主机和端口发出新的 HTTP 请求时将建立新连接。 服务器也可以拒绝通过同一连接允许多个请求，在这种情况下，必须为每个请求重新建立连接，并且不能放入连接池。 Agent 仍将向该服务器发出请求，但每个请求都将通过新连接发生。当客户端或服务器关闭连接时，它将从连接池中删除。连接池中任何未使用的套接字都将被销毁，以便当没有未完成的请求时不用保持 Node.js 进程运行。

调用 client.request 方法，可把来自浏览器的请求转发至 client 对应的 authority（后端服务器进程）中，同时有响应时会调用 options.onres 或者 options.onupgrade。在后文中会看到，hawkey 在 serveRequest 时会通过调用 client.request 方法来转发请求。

有的同学可能还是有疑问，为什么 nginx 等不支持代理转发至 http/2.0 的后端服务呢？我想主要原因可能有如下几个

* http2 一般会要求走 https，而内网环境走 https 显然是不必要的。
* 就算 2.0 协议本身支持非 https，支持转发到 2.0 协议也会增加很多额外的开发维护成本，潜在的性能收益并不大，而因为 http2 协议改变了协议传输层，很多基于 http 协议的流量分析监听工具可能会失效，进而带来服务维护成本上升
* 内网环境下网络很稳定，http/1.1 连接支持 keep-alive，与连接池相比，http/2.0 并没有优势

所以大部分网站在部署 http2 时，只是在浏览器与网站接入层（如 nginx）之间进行了协议升级，端到端部署 http2 不太多。现在只支持代理转发到 http/1.1 的后端服务。

```javascript
// 公开 Client 对象
const Client = {
  proxyClients: config.proxy.map(proxy => {

    let proxyClient = {
      pathRegExp: proxy.pathRegExp,
      authority: proxy.authority,
      httpVersion: proxy.httpVersion,

      // 初始化 proxyClient 对象，创建 client 对象，后续匹配了 pathRegExp 的请求会通过此 client 进行代理转发
      init() {
        this._authorities_1_1 = this.authority.map(authority => _client_1_1(authority))
        return this
      },

      // 关闭 client 对应的连接池
      close() {
        this._authorities_1_1.forEach(authority => authority.close())
      },

      // 目前只配置 1 个后端服务，可升级为多台轮询，支持负载均衡
      get() {
        return this._authorities_1_1[0]
      },

      // 判断请求路径是否与当前 proxyClient 匹配，以决定是否转发该请求
      pathMatch(path) {
        return new RegExp(this.pathRegExp).test(path)
      }
    }

    return proxyClient.init()
  
  }),

  // 根据请求 path，找出对应的 proxyClient
  get(path) {
    for(let proxyClient of this.proxyClients) {
      if (proxyClient.pathMatch(path)) return proxyClient
    }
    return false
  },

  // 关闭所有 proxyClient
  close() {
    this.proxyClients.forEach(proxyClient => proxyClient.close())
  }
}

module.exports = Client
```

如上述代码，可以先查看代码注释。client.js 文件公开了 Client 对象有如下字段和方法

* proxyClients -> 可遍历该列表找到请求 path 对应的请求转发代理
* get(path) -> 遍历 proxyClients，找到请求 path 对应的请求转发代理
* close() -> 关闭所有远程连接池，释放资源

在后面的 serveRequest 和 serveStream 组件中，可以通过调用 Client.get(path) 获取 proxyClient 对象，再调用 proxyClient.get() 方法获取代理转发 client 对象。