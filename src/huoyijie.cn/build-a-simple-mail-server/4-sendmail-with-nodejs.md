# 通过 Node.js 程序发送邮件

现在已经搭建好了邮箱服务器，如果有程序发送邮件的需求该如何实现呢？

`Node.js` 中有个 [nodemailer](https://nodemailer.com/) npm 包，处理程序发送邮件是非常容易的。首先安装包。

```bash
$ npm i nodemailer --save
```

现在写个 `Demo` 程序测试下

```javascript
'use strict';
const nodemailer = require('nodemailer');

async function main() {

  let testAccount = {
      user: 'huoyijie',
      pass: '******'
  };

  // create reusable transporter object using the default SMTP transport
  let transporter = nodemailer.createTransport({
    host: 'mail.huoyijie.cn',
    port: 465,
    secure: true, // true for 465, false for 25
    auth: {
      user: testAccount.user, // generated ethereal user
      pass: testAccount.pass, // generated ethereal password
    },
  });

  // send mail with defined transport object
  let info = await transporter.sendMail({
    from: `huoyijie <huoyijie@huoyijie.cn>`, // sender address
    to: 'author@huoyijie.cn, huorong@huoyijie.cn', // list of receivers
    subject: 'Hello, from Node.js', // Subject line
    text: 'Hello world?', // plain text body
    html: '<b>Hello world?</b>', // html body
  });

  console.log('Message sent: %s', info.messageId);
  // Message sent: <42a91923-be1e-0bab-cd2e-c404a59165fd@huoyijie.cn>
}

main().catch(console.error);
```

下面是邮箱服务器上已有的一个邮件帐号，密码是该帐号的系统登录密码。

```javascript
let testAccount = {
    user: 'huoyijie',
    pass: '******'
};
```

下面是邮件服务器域名、端口号，通过 SSL 建立连接

```json
host: 'mail.huoyijie.cn',
port: 465,
secure: true, // true for 465, false for 25
```

上面代码保存在 `sendmail.js` 中，下面来运行该程序

```bash
$ node sendmail.js 
Message sent: <42a91923-be1e-0bab-cd2e-c404a59165fd@huoyijie.cn>
Preview URL: false
```

打开 `Thunderbird` 客户端，已收到邮件！

![Node.js 程序中发送邮件](https://cdn.huoyijie.cn/ab/84f0486086be11ebaf1339e97396ca47/sendmail-with-nodejs.jpg)