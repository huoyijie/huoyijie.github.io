# 安装 postfix

postfix 是 Wietse Venema 在 IBM 的 GPL 协议之下开发的 MTA（邮件传输代理）软件。实现了 SMTP 协议，用来发送或者中转电子邮件。我选择了基于 postfix 搭建邮箱服务器。

首先确保已有域名、公网 IP 及服务器，对于小的网站，服务器可同时运行其他 Web 服务。我会基于 Ubuntu 系统来介绍如何安装和配置邮件服务器，其他 Linux 发行版是类似的。

```bash
$ cat /etc/issue
Ubuntu 18.04.1 LTS
```

有一件事必须先做，就是修改服务器 hostname，hostname 标识了邮件服务器的身份，一定要用全限定域名。我的域名是 huoyijie.cn，惯例是用 mail.huoyijie.cn 来标识本邮件服务器。

```bash
$ sudo hostnamectl set-hostname mail.huoyijie.cn
```

设置完成再检查确认下

```bash
$ hostname -f
mail.huoyijie.cn
```

然后先安装 postfix

```bash
$ sudo apt-get update
$ sudo apt-get install postfix -y
```

安装过程中主要会询问 2 个问题
1. 选择 Internet Site，通过 STMP 收发邮件

![install postfix 1](https://cdn.huoyijie.cn/ab/84f0486086be11ebaf1339e97396ca47/install-postfix-1.jpg)

安装过程中会停在上图这里，这里主要是让你了解下几个选项的含义，按Tab键选中<确定>按钮，敲回车键进入下一步

![install postfix 2](https://cdn.huoyijie.cn/ab/84f0486086be11ebaf1339e97396ca47/install-postfix-2.jpg)

默认会选中 Internet Site 选项，如果不是，按上下箭头调整。调整完成后按 Tab 键选中 <确定> 按钮并敲回车键。

2. 填写用来限定邮箱地址的域名

![install postfix 3](https://cdn.huoyijie.cn/ab/84f0486086be11ebaf1339e97396ca47/install-postfix-3.jpg)

我的邮件地址是 author@huoyijie.cn，这里填写 @ 后面的域名 huoyijie.cn

然后就是等待安装完成。安装完成后，服务已经启动，监听 25 端口。

```bash
$ sudo netstat -nap |grep :25
tcp        0      0 0.0.0.0:25              0.0.0.0:*               LISTEN      13332/master        
tcp6       0      0 :::25                   :::*                    LISTEN      13332/master
```

配置文件主要是 2 个

1. /etc/postfix/main.cf
2. /etc/postfix/master.cf

正常配置不需要调整，就可以发送邮件了。服务器的系统用户都分配了邮件地址，可以收发邮件。我是用 huoyijie 用户登录服务器。接下来发邮件测试下，可以使用随 postfix 一起安装的 sendmail 工具。

```bash
$ echo "test email" | sendmail huoyijie@huoyijie.cn
```

然后安装 mailutils 工具使用命令来查看邮件

```bash
$ sudo apt install mailutils
```

```bash
$ mail
"/var/mail/huoyijie": 1 message 1 new
>N   1 huoyijie           三 3月 17 11:3  10/348   
? 
Return-Path: <huoyijie@huoyijie.cn>
X-Original-To: huoyijie@huoyijie.cn
Delivered-To: huoyijie@huoyijie.cn
Received: by xps (Postfix, from userid 1000)
	id 15F4063007D2; Wed, 17 Mar 2021 11:33:17 +0800 (CST)
Message-Id: <20210317033317.15F4063007D2@xps>
Date: Wed, 17 Mar 2021 11:33:17 +0800 (CST)
From: huoyijie@huoyijie.cn (huoyijie)

test email
```

可以看到已经可以正常发送邮件了。连续按 2 次 Ctrl+D 退出 mail 命令后可以看到下面的消息

```log
Saved 1 message in /home/huoyijie/mbox
```

postfix 默认是以 mbox 格式存储邮件，每个邮箱用户的所有的邮件都放到一个文件里面，如果恰好同时在接收邮件和读取邮件时会对文件加锁，比较影响性能。所以可以选择另外一种 Maildir 格式来存储邮件，一个用户的所有邮件都存储在一个特定目录下面，每个邮件单独一个文件。

修改配置文件 `/etc/postfix/main.cf`，在文件最后加上一行 `home_mailbox = Maildir/`

```bash
$ sudo service postfix reload
```

修改完 reload 一下使修改生效。会在用户 home 目录下创建 `MailDir` 文件夹存放该用户的所有邮件。如果继续发邮件测试，就可以在 `~/Maildir/new` 目录下发现新发送的邮件

```bash
$ ls ~/Maildir/
cur  new  tmp
$ $ ls ~/Maildir/new
1615882507.Vfc01I449d3M362690.mail.huoyijie.cn
```

邮件服务器系统用户即为邮件用户，所以现在可以创建一个新的邮箱用户 `author@huoyijie.cn`

可以把所有邮箱用户的数据指定在同一个目录，加入同一个用户组，方便管理邮件数据。

```bash
$ sudo mkdir /home/mailusers
$ sudo groupadd mailusers
$ sudo useradd -m -d /home/mailusers/author -g mailusers  author
$ sudo passwd author
```

所有用户集中在 `/home/mailusers` 目录下，添加 `author` 用户，并添加密码，邮箱帐号配置好了。