# 安装 dovecot

Dovecot 是一个开源的 IMAP 和 POP3 邮件服务器，支持 Linux/Unix 系统。POP / IMAP 是 MUA（邮件用户代理） 从邮件服务器中读取邮件时使用的协议。其中，POP3 是从邮件服务器中下载邮件存起来，IMAP4 则是将邮件留在服务器端直接对邮件进行管理、操作。

下面来安装 dovecot

```bash
$ sudo apt-get install dovecot-common dovecot-imapd dovecot-pop3d -y
```

安装完成后，配置下存储邮箱用户邮件的目录，编辑配置文件

```bash
sudo vi /etc/dovecot/conf.d/10-mail.conf
```

找到 `mail_location` 配置项

```conf
#   mail_location = maildir:~/Maildir
#   mail_location = mbox:~/mail:INBOX=/var/mail/%u
#   mail_location = mbox:/var/mail/%d/%1n/%n:INDEX=/var/indexes/%d/%1n/%n
#
# <doc/wiki/MailLocation.txt>
#
mail_location = mbox:~/mail:INBOX=/var/mail/%u
```

修改为

```conf
mail_location = maildir:~/Maildir
```

在上一节中，配置 postfix 时，我没设置为邮件存储在 `~/Maildir` 目录中，并 reload 配置

```bash
sudo service dovecot reload
```

查看下端口

```bash
$ sudo netstat -nap |grep dovecot
tcp        0      0 0.0.0.0:110             0.0.0.0:*               LISTEN      15150/dovecot       
tcp        0      0 0.0.0.0:143             0.0.0.0:*               LISTEN      15150/dovecot       
tcp6       0      0 :::110                  :::*                    LISTEN      15150/dovecot       
tcp6       0      0 :::143                  :::*                    LISTEN      15150/dovecot
```

110 是 POP3 的端口，143 是 IMAP 的端口。都正常启动了。

接下来还一个步骤比较重要，就是配置邮件域名解析，否则发送邮件服务器是没有办法路由和转发邮件的。因为没有办法把邮件地址 `author@huoyijie.cn` 中的域名 `huoyijie.cn` 对应到哪个邮件服务器上面。而且邮件客户端软件是需要设置 SMTP 服务器域名以及 IMAP/POP3 服务器域名的。

首先来到域名注册商的控制台，找到域名 `huoyijie.cn`，点击添加 `MX`记录，如下图

![添加 MX 记录](https://cdn.huoyijie.cn/ab/84f0486086be11ebaf1339e97396ca47/dns-mx-record.jpg)

这是一条邮件交换(MX)记录，其中主机记录 @ 代表直接解析 `huoyijie.cn`，记录值为 `mail.huoyijie.cn`，意思是填写 `xxx@huoyijie.cn` 地址时，表示该邮件的目标邮箱服务器实际是 `mail.huoyijie.cn`。

此时还需要添加一条 A 记录，把邮件服务器域名 `mail.huoyijie.cn` 解析到 其IP地址，如下图

![添加 A 记录](https://cdn.huoyijie.cn/ab/84f0486086be11ebaf1339e97396ca47/dns-a-record.jpg)

主机记录填 `mail`，记录类型为 A，记录值为邮箱服务器 IP 地址。加好后等几分钟时间生效。此时 ping 下域名可以正常解析为 IP地址了。

```bash
$ ping mail.huoyijie.cn
PING mail.huoyijie.cn (62.234.116.248) 56(84) bytes of data.
64 bytes from 62.234.116.248 (62.234.116.248): icmp_seq=1 ttl=52 time=29.6 ms
64 bytes from 62.234.116.248 (62.234.116.248): icmp_seq=2 ttl=52 time=28.9 ms
64 bytes from 62.234.116.248 (62.234.116.248): icmp_seq=3 ttl=52 time=29.4 ms
```

接下来选一个邮件客户端软件，我使用了电脑自带的 Thunderbird 邮件客户端程序，配置一下邮箱账户。

![添加邮件帐号](https://cdn.huoyijie.cn/ab/84f0486086be11ebaf1339e97396ca47/config-thunderbird-1.jpg)

点击“跳过并使用已有的电子邮箱”

![添加邮件帐号](https://cdn.huoyijie.cn/ab/84f0486086be11ebaf1339e97396ca47/config-thunderbird-2.jpg)

输入显示名字、电子邮件地址以及密码。这个密码就是邮箱服务器系统帐号 author 的登录密码，在配置 postfix 那节中，有添加 author 这个用户，也设置了密码。点击“继续”按钮

![添加邮件帐号](https://cdn.huoyijie.cn/ab/84f0486086be11ebaf1339e97396ca47/config-thunderbird-3.jpg)

软件会自动检查邮箱服务器的各项配置，大概需要几秒钟时间。上图是自动检测的配置，可以点击“手动配置”按钮查看一下

![添加邮件帐号](https://cdn.huoyijie.cn/ab/84f0486086be11ebaf1339e97396ca47/config-thunderbird-4.jpg)

此时点击“完成”按钮，邮箱帐号设置完成，可以用这个帐号收发邮件了。

![发邮件](https://cdn.huoyijie.cn/ab/84f0486086be11ebaf1339e97396ca47/mailserver-cover.jpg)

也可以从别的域的邮箱发过来，比如用qq邮箱发过来也可以。但是现在一般情况下是不能够直接往别的域发邮件的，因为会直接被其他邮箱服务提供商拒绝掉，邮件会被退原来。原因主要是那些大的邮箱服务提供商为了避免垃圾邮件，选择直接拒绝一些“信用”度比较低的邮箱服务器转发过来的邮件。

所以为了能够向其他域名发送邮件，还需要更多的配置和调优，以提升自己的邮箱服务器的“信用度”。