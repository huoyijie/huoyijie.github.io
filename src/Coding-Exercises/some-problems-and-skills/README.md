# Problems And Skills

本文会把开发过程中常常会遇到的一些小问题和对应的解决办法记录下来，尽量避免重复犯错。再就是记录整理一些我实践过的编程技巧，以及一些碎片化的问题和知识点方便后面查阅。

因为需要长期收集整理，可能会不定期更新...

首先附上一个刚刚遇到的关于 gitbook v3.2.3 的问题，切换章节时控制台不断报错

```javascript
Uncaught TypeError: Cannot read property 'split' of undefined
```

[fix split of undefined with gitbook theme](https://huoyijie.cn/revealjs/fix-split-of-undefined-with-gitbook-theme.html)

对于刚刚上线的网站，肯定会想直到有哪些外部来源跳转过来，可以通过记录请求 Referer 来分析一下
[Express网站分析Referer来源](https://huoyijie.cn/revealjs/analyze-referer-of-express-website.html)

[配置NODE.JS 代理](https://huoyijie.cn/revealjs/nodejs-http-proxy-middleware.html)

通过分析 access.log 可以了解到博客页面的 PV/UV 情况，这些数据最好的展现方式就是图表，可以更直观的看出数据的变化趋势。最近我用 [Chart.js](https://github.com/chartjs/Chart.js) 制作了博客的数据实时分析 Demo，[make statics with Chart.js](https://huoyijie.cn/revealjs/make-statics-with-chartsjs.html)