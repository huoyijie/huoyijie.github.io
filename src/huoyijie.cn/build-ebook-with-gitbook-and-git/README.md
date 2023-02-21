# 使用 Git 与 Gitbook 创建管理电子书

这篇文章接上篇 [基于 Git 搭建代码托管服务器](https://huoyijie.cn/article/7d4e75e08ac011ebaa2ab110efea0133/)，Git 服务器也可以托管电子书。新建一个 git 项目，在根目录下为不同的电子书创建子目录，在子目录下按照 gitbook 模板编写 markdown，然后通过 git push 推送到 git server，触发服务器 hooks 脚本执行，脚本会调用 gitbook 程序生成 html 格式电子书。服务器上启动一个静态 HTTP Server，就可以把电子书发布出去了。如果要更新电子书内容，只需要修改相应的 markdown 内容，然后 push 到服务器，电子书就会自动更新。

有了 `Git + Gitbook + Markdown 编辑器`后就可以开始写博客了。这种方式和 `github pages` 是一样的，但部署速度更快，且可以定制 git hooks 脚本，还可以利用服务器渲染动态页面。`github pages` 虽是免费的，但只能部署静态的内容。

首先创建一个新的仓库，名为 `mybook`，登录服务器

```bash
$ cd /home/git
$ sudo mkdir mybook.git && cd mybook.git
$ sudo git init --bare
Initialized empty Git repository in /home/git/mybook.git/
```

修改仓库用户归属，确保后面可以 git 用户有权限推送代码

```bash
sudo chown -R git:git /home/git/mybook.git
```

在自己的电脑上克隆仓库

```bash
git clone git@huoyijie.cn:/home/git/mybook.git
正克隆到 'mybook'...
warning: 您似乎克隆了一个空仓库。
```

打开 vscode，打开 mybook 文件夹，目前是空项目。添加框架代码后来看一下项目结构。

```bash
$ tree -L 2
.
├── bin
│   └── www
├── mybook.json
├── mygitbook
│   ├── book.json
│   ├── README.md
│   └── SUMMARY.md
├── node_modules
│   ├── http-server
│   ├── shelljs
├── package.json
├── package-lock.json
├── public
│   ├── gitbook
│   └── index.html
└── scripts
    └── deploy.js
```

从上到下依次介绍下主要目录和文件

* `bin/www` 为 `http server` 启动脚本，看下内容

```bash
#!/bin/bash

DIR="public"
if [ ! -d "$DIR" ]; then
  mkdir $DIR
  echo "Home Page" > $DIR/index.html
fi

./node_modules/http-server/bin/http-server $DIR -a 127.0.0.1 -p 8001
```

静态文件存储目录为 `$DIR=public` ，监听 8001 端口。

* `mybook.json` 是自定义的关于当前仓库中 book 的配置文件，来看下内容

```javascript
{
    "deployBook": ["mygitbook"],
    "dropBook": []
}
```

每个 book 都有唯一的字符串 name，如当前这篇博客文章的名字为 mygitbook。这个配置文件目前还很简单，只定义了 2 个配置项。
1. deployBook: 需要发布到 http server 的 book name列表
2. dropBook: 需用从 http server 中删除的 book name列表

也就是说，新增的 book 要在 deployBook 配置一下才会真正发布到 http server，需要下线的 book 要在 dropBook 配置一下。修改完配置推送变更到服务器就会自动部署更新。

* mygitbook 目录里是当前文章真正的内容，可以通过 `gitbook init` 初始化。为了更方便的管理，我创建了 `gitbookgen.js` 脚本，可以快速生成电子书模板。
  
* node_modules 为当前项目依赖的 npm 包，主要是 http server 和 shelljs，后者可以方便调用系统 shell 命令。
  
* package.json 为 Node.js 项目配置文件

```javascript
{
  "name": "mybook",
  "version": "0.0.0",
  "private": true,
  "scripts": {
    "start": "bin/www"
  },
  "dependencies": {
    "http-server": "^0.12.3",
    "shelljs": "^0.8.4"
  }
}
```
* public 目录为 http server 静态内容存储的目录，后面会把电子书部署在这个目录下
  
* scripts/deploy.js 目录为部署 book 的脚本文件，会被 git server 上面的 hooks 脚本调用
  
```javascript
#!/usr/bin/env node

// run from project root dir with command './scripts/deploy.js'
const shell = require('shelljs')
const fs = require('fs')

console.log(`pwd: ${shell.pwd()}`)

if(shell.exec('git pull').code) {
  console.error('!> deploy error -> git pull failed')
  process.exit(1)
}

const publicDir = 'public'
const mybookJson = JSON.parse(fs.readFileSync('mybook.json'))
console.log(`> load mybook.json: `, mybookJson)

mybookJson.dropBook.forEach(book => {
  console.log(`> drop book: ${book}`)

  var tmp = book.split('|')
  book = tmp[0]
  var bookid = tmp[1] || ''

  if (!!bookid) {
    const bookDir = `${publicDir}/${bookid}`
    shell.rm('-rf', bookDir)
    shell.rm('-rf', book)
  } else {
    console.error(`!> delete book failed ->lack of param bookid.`)
    process.exit(1)
  }
})

mybookJson.deployBook.forEach(book => {
  console.log(`> deploy book: ${book}`)

  if (!fs.existsSync(book)) {
    console.error(`!> deploy error with ${book} -> not exist`)
    process.exit(1)
  } else {

    shell.cd(book)

    const bookJson = JSON.parse(fs.readFileSync(`book.json`))
    console.log(`> load book.json:`, bookJson)
    const bookDir = `${publicDir}/${bookJson.id}`
  
    if(!shell.exec('gitbook build').code) {
      shell.cd('-')
  
      if (!fs.existsSync(bookDir)) {
        shell.mkdir(bookDir)
      } else {
        shell.rm('-rf', `${bookDir}/*`)
      }
      
      shell.cp('-r', `${book}/_book/*`, bookDir)
      console.info(`> deploy succeed with ${book}`)
    } else {
      console.error(`!> deploy error with ${book}`)
      process.exit(1)
    }

  }

})

process.exit(0)
```

这个脚本主要做的事情就是等变更推送到服务器后，把配置在 deployBook 下的 book 发布出去，期间调用了 `gitbook build` 把 markdown 内容转为 html。把配置在 dropBook 下的 book 下线掉。这个脚本是被 git server 上的 `post-receive` hooks 调用的。`post-receive` 是在 git server 接受 push 后自动触发调用的。所以这个脚本是完成推送变更，自动部署更新 book 的关键。

下面来看一下存放在 git server 上的 `post-receive` hooks。文件存放路径为

```bash
$ ls -l /home/git/mybook.git/hooks/post-receive 
-rwxr-xr-- 1 git git 103 Mar 23 15:59 /home/git/mybook.git/hooks/post-receive
```

这个脚本的内容非常简单

```bash
#!/bin/sh

cd /home/git/vswork/mybook

unset GIT_DIR

./scripts/deploy.js

exec git-update-server-info
```

主要逻辑就是调用了 `./scripts/deploy.js` 部署脚本，其中 `/home/git/vswork/mybook` 是通过下面的命令克隆的，也是服务器上对应的博客工作目录。

```bash
$ cd /home/git/vswork
$ sudo git clone /home/git/mybook.git
$ cd mybook
$ npm install
```

安装好依赖后，修改下用户归属

```bash
sudo chown -R git:git mybook
```

到此为止，项目基本已经部署配置完成。接下来就是打开 vscode，编辑 mygitbook 电子书。编辑完成后 push 到服务器上完成文章自动部署更新。下面是 push 代码后的输出，可以看到 scripts/deploy.js 脚本在服务器上执行后的输出，有在自动部署 mygitbook。

```bash
$ git push
对象计数中: 4, 完成.
Delta compression using up to 8 threads.
压缩对象中: 100% (4/4), 完成.
写入对象中: 100% (4/4), 1.94 KiB | 1.94 MiB/s, 完成.
Total 4 (delta 2), reused 0 (delta 0)
remote: pwd: /home/git/vswork/mybook
remote: From /home/git/mybook
remote:    4fe0a5e..6e62e3d  master     -> origin/master
remote: Updating 4fe0a5e..6e62e3d
remote: Fast-forward
remote:  mygitbook/README.md | 108 ++++++++++++++++++++++++++++++++++++++++++++++++++--
remote:  1 file changed, 104 insertions(+), 4 deletions(-)
remote: > load mybook.json:  {
remote:   '// deployBook': "['firstbook','secondbook']",
remote:   deployBook: [ 'mygitbook' ],
remote:   '// dropBook': "['firstbook|f44b5ae08bb611ebbf5b4f5ee88ae545','secondbook|f44b5ae08bb611ebbf5b4f5ee88ae545']",
remote:   dropBook: []
remote: }
remote: > deploy book: mygitbook
remote: > load book.json: {
remote:   id: 'c0a4a7508bc011eb952e418fdf675f27',
remote:   name: 'mygitbook',
remote:   title: '使用 Git 与 Gitbook 创建管理电子书',
remote:   author: 'huoyijie',
remote:   description: '使用 Git 与 Gitbook 创建管理电子书',
remote:   introduce: '这里是内容介绍...',
remote:   coverImg: 'https://cdn.huoyijie.cn/ab/c0a4a7508bc011eb952e418fdf675f27/mygitbook-cover.jpg',
remote:   language: 'zh-hans',
remote:   plugins: [ '-lunr', '-search', '-sharing', '-fontsettings' ],
remote:   links: { sidebar: { '首页': 'https://huoyijie.cn/?from=article_mygitbook' } }
remote: }
remote: info: 7 plugins are installed 
remote: info: 2 explicitly listed 
remote: info: loading plugin "highlight"... OK 
remote: info: loading plugin "theme-default"... OK 
remote: info: found 1 pages 
remote: info: found 0 asset files 
remote: info: >> generation finished with success in 0.4s ! 
remote: > deploy succeed with mygitbook
To huoyijie.cn:/home/git/mybook.git
   4fe0a5e..6e62e3d  master -> master
```

OK，关于如何使用 Git 与 Gitbook 创建管理电子书的主要内容介绍完了。