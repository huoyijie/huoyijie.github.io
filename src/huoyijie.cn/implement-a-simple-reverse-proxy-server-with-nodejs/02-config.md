# 配置文件

配置文件部分比较简单，配置项很少，位于 server-config.json 中

```javascript
{
  "port": 8443,
  "host": "0.0.0.0",
  "pems": {
    "key": "keys/hawkey-key.pem",
    "cert": "keys/hawkey-cert.pem"
  },
  "accessLog": "logs/access.log",
  "serverRoot": "public",
  "staticDir": "static",

  "proxy": [
    {
      "pathRegExp": "/",
      "authority": ["127.0.0.1:7999"],
      "httpVersion": "1.1"
    }
  ],

  "env": "development"
}
```

* port -> 启动监听端口，如 443
* host -> 启动地址，如 0.0.0.0
* pems -> 配置服务器证书及私钥
  * key -> 私钥文件
  * cert -> 证书文件
* accessLog -> 流水日志
* serverRoot -> 服务器根目录
* staticDir -> 静态文件根目录，位于 serverRoot 下
* proxy -> 反向代理 upstream 列表
  * pathRegExp -> 根据路径规则进行正则匹配。
  * authority -> 可以配置后端服务器 ip:port 列表，可以配置多个。
  * httpVersion -> 为后端 web 服务支持的协议版本。一般后端服务都不支持 http/2.0，主要原因是主流服务器如 nginx 等一般都不支持转发到 http/2.0 的后端服务，http/2.0 没有办法做到端到端整条链路的部署，通常只是在浏览器与 nginx 之间采用 http/2.0，nginx 与后端服务之间走 http1 协议。我的博客程序基于 Express 4，也是不支持 http/2.0 的。所以 hawkey 与博客程序之间采用 1.1 协议通信。
* env -> 启动环境，不同环境具有不同的环境变量和行为。如生产环境（production）输出更少的日志，而开发环境（development）可输出 debug 日志，方便分析问题。

该配置文件通过 config.js 进行加载

```javascript
const fs = require('fs')

module.exports = (conf) => {
  let config = JSON.parse(fs.readFileSync(conf))
  config.scheme = 'http'
  config.hostname = 'huoyijie.cn'
  config.httpPort = 80
  return config
}
```

其中 `conf = 'server-config.json'`，加载配置很简单，加载完成后可以通过 config 对象访问配置信息。也可以在 config.js 中配置一些不需要用户参与调整的配置项，如

* scheme -> 后端服务只需提供 http 服务，其他诸如加解密、加解压缩等操作可在 hawkey 中完成
* hostname -> 服务器域名，本地测试可填写 127.0.0.1
* httpPort -> 80，本地测试可随便填写，该端口上会启动一个 http server，把所有 http 请求自动转发到 https 相应路径上