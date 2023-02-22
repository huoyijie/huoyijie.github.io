## Node.js http proxy middleware

先安装 http-proxy-middleware

```bash
$ npm install http-proxy-middleware
```
---
## 代理配置

```js
// 代理配置
const { createProxyMiddleware } = require('http-proxy-middleware');
// proxy middleware options
const options = {
  target: 'http://127.0.0.1:8000', // target host
  changeOrigin: true,
  pathRewrite: function (path, req) { return path.replace('/revealjs', '/') }
};
// create the proxy (without context)
const revealjsProxy = createProxyMiddleware(options);
app.use('/revealjs', revealjsProxy);
```
---
前端的项目遇到 /revealjs 路径会转发到服务器的 8000端口进程，并由后者负责处理该http请求。
---
options.pathRewrite可以配置路径重写规则，比如在转发到后面的进程中时去掉 /revealjs 前缀。
---
## 配置完成

后续浏览器访问带 /revealjs 前缀的路径会自动代理转发到 8000 端口