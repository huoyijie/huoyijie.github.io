# My daily keynotes

网站通过 [keynote](https://github.com/huoyijie/keynote) 生成。更多内容请关注 [我的网站](https://huoyijie.cn)。

## 添加 gitbook

```bash
npm init -y

npm i gitbook-plugin-prism gitbook-plugin-code --save

gitbook build ./ latest
```
## 本地测试

```bash
keynote
```

## 生成静态网站

```bash
keynote -gen
```