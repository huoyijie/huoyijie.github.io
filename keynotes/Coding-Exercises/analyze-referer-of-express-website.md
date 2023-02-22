## Analyze Referer Of Express Website

Express网站分析Referer来源
---
## 先建表

| 字段 | 说明 |
|-----|:-----:|
| path | 用户访问的网页地址 |
| referer | 来源 |
---
## 定义models.Referer
```js
sequelize.define('referer', {
  path: {
      type: Sequelize.STRING,
      allowNull: false
  },
  referer: {
      type: Sequelize.STRING,
      allowNull: false
  }
})
```
---
## 加一个Express Filter
```js [8,9]
const myFilter = (req, res, next) => {
  if (req.path === '/' 
    || (req.path.startsWith('/blog/') && req.path.includes('/play'))
    || (req.path.startsWith('/article/') && !req.path.includes('/gitbook'))
    || (req.path.startsWith('/revealjs/') && req.path.includes('.html'))) {

    models.Referer.create({
      path: req.path.trim(),
      referer: (req.headers.referer || '').trim()
    });

  }
  next();
};

此Filter把满足条件的请求，对应的path和referer存入表中
```
---
## Express配置Filter

```js [3]
// 自定义全局过滤器
const myFilter = require('./filters/myFilter');
app.use(myFilter);

于app.js文件中配置全局filter(注意filter的顺序)
```
---
通过记录的数据，可以分析网站有哪些外部来源及其占比，也可以具体分析某个网址的来源。