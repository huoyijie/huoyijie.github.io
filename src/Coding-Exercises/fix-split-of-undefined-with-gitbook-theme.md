## Problem
![pic](https://cdn.huoyijie.cn/ab/fix-split-of-undefined-with-gitbook-theme/split-of-undefined-error.jpg)

今天在使用gitbook创建文档时遇到了一个错误，在undefined上面调用了split方法
---
![pic](https://cdn.huoyijie.cn/ab/fix-split-of-undefined-with-gitbook-theme/themejs.jpg)

错误来自theme.js文件，但是theme.js文件是经过压缩混淆过的代码，很难找到问题真正原因。
---
![pic](https://cdn.huoyijie.cn/ab/fix-split-of-undefined-with-gitbook-theme/format-themejs.jpg)

此时可以点下图片左下角括号图标格式化代码，然后输入关键词定位到出错的代码处，看起来像从url中提取出hash值
---
theme.js是复制gitbook主题的渲染，这个错误是在切换目录导航时出现的，现在必须找到源码修复后重现build出新的theme.js
---
```bash
$ git clone git@github.com:/GitbookIO/theme-default
```

gitbook的主题默认是<a href="https://github.com/GitbookIO/theme-default">theme-default</a>，代码clone下来
---
## 搜索代码
  
![search](https://cdn.huoyijie.cn/ab/fix-split-of-undefined-with-gitbook-theme/search-code.jpg)

在导航navigation.js代码中 $link.attr('href') 直接split，可能会出现这个问题
---
代码调整如下:

```js
function getChapterHash($chapter) {
  var $link = $chapter.children('a');
  var hash = $link.attr('href') && $link.attr('href').split('#')[1];

  if (hash) hash = '#'+hash;
  return (!!hash)? hash : '';
}
```
---
## Build
```
$ cd theme-default
$ npm install
$ ./src/build.sh
```

---
项目根目录下是新生成的_assets
```
$ tree -L 1
  .
  ├── _assets
  ├── _i18n
  ├── index.js
  ├── _layouts
  ├── LICENSE
  ├── node_modules
  ├── package.json
  ├── package-lock.json
  ├── preview.png
  ├── README.md
  └── src
```
---
已生成新的theme.js
```
$ tree -L 2
  .
  ├── ebook
  │   ├── ebook.css
  │   ├── epub.css
  │   ├── mobi.css
  │   └── pdf.css
  └── website
      ├── fonts
      ├── gitbook.js
      ├── images
      ├── style.css
      └── theme.js
```
---
用新生成的theme.js覆盖已安装的主题插件中的文件
```
cp _assets/website/theme.js ~/.gitbook/versions/3.2.3/node_modules/gitbook-plugin-theme-default/_assets/website/theme.js
```
---
清理gitbook缓存
```
rm -rf _book
```
重新运行
```
gitbook serve
```
---
## DONE
![done](https://cdn.huoyijie.cn/ab/fix-split-of-undefined-with-gitbook-theme/done.jpg)

再重试已经没有那个错误了。