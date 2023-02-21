# 配置SSL安全传输

为邮箱服务器配置SSL，需要为邮箱域名申请证书，先去域名服务商控制台，找到域名 `mail.huoyijie.cn`，点击申请SSL证书。输入邮箱服务器域名 `mail.huoyijie.cn`，输入申请邮箱。点下一步，选择自动验证，然后等大概几秒中时间，证书会自动签发。

![申请ssl证书](https://cdn.huoyijie.cn/ab/84f0486086be11ebaf1339e97396ca47/apply-ssl-cert.jpg)

注：如果暂时不能申请证书，使用本地生成的自签名证书也可以先用来测试

证书签发下来后，可以下载到邮箱服务器上面，把公钥、私钥文件分别部署配置好。

证书文件: `/etc/ssl/certs/ssl-cert-mail.crt`
私钥文件: `/etc/ssl/private/ssl-cert-mail.key`

确保私钥文件权限 `640`，只有 root 用户才有权限读写。

```bash
$ sudo ls -l /etc/ssl/private/ssl-cert-mail.key
-rw-r----- 1 root root 1704 Mar 16 15:00 /etc/ssl/private/ssl-cert-mail.key
```

证书文件权限 `644`，证书是公开的，所以其他用户可读取一般也没有关系。

```bash
$ sudo ls -l /etc/ssl/certs/ssl-cert-mail.crt
-rw-r--r-- 1 root root 3735 Mar 16 14:59 /etc/ssl/certs/ssl-cert-mail.crt
```

先来配置 postfix，编辑配置文件 `/etc/postfix/main.cf`

```bash
sudo vi /etc/postfix/main.cf
```

修改下面几个配置

```conf
smtpd_tls_cert_file=/etc/ssl/certs/ssl-cert-mail.crt
smtpd_tls_key_file=/etc/ssl/private/ssl-cert-mail.key
smtpd_use_tls=yes
```

编辑配置文件 `/etc/postfix/master.cf`

```bash
sudo vi /etc/postfix/master.cf
```

找到 `smtps` 进程配置，打开注释

```conf
smtps     inet  n       -       y       -       -       smtpd
   -o syslog_name=postfix/smtps
   -o smtpd_tls_wrappermode=yes
```

修改完配置后，reload 配置

```bash
sudo service postfix reload
```

`smtps` 进程会打开 465 端口，现在检查一下

```bash
$ sudo netstat -nap |grep 465
tcp        0      0 0.0.0.0:465             0.0.0.0:*               LISTEN      17041/master        
tcp6       0      0 :::465                  :::*                    LISTEN      17041/master
```

接下来配置 dovecot

编辑配置文件 `/etc/dovecot/conf.d/10-ssl.conf`

```bash
sudo vi /etc/dovecot/conf.d/10-ssl.conf
```

找到下面3个配置项进行调整，注意下方配置文件路径时有 `<`

```conf
ssl = required
ssl_cert = </etc/ssl/certs/ssl-cert-mail.crt
ssl_key = </etc/ssl/private/ssl-cert-mail.key
```

修改完配置后，reload 配置

```bash
sudo service dovecot reload
```

`pop3s` 占用端口 995，`imaps` 占用端口 993，检查一下

```bash
$ sudo netstat -nap |grep 993
tcp        0      0 0.0.0.0:993             0.0.0.0:*               LISTEN      1/init
tcp6       0      0 :::993                  :::*                    LISTEN      1/init              
```

```bash
$ sudo netstat -nap |grep 995
tcp        0      0 0.0.0.0:995             0.0.0.0:*               LISTEN      17133/dovecot       
tcp6       0      0 :::995                  :::*                    LISTEN      17133/dovecot
```

最后需要调整一下邮件客户端 `Thunderbird` 的发件/收件服务器配置。

![配置SSL for Thunderbird](https://cdn.huoyijie.cn/ab/84f0486086be11ebaf1339e97396ca47/config-thunderbird-ssl.jpg)

打开邮箱帐号设置，接收 IMAP 配置了 993 端口。而发出SMTP 配置了 465 端口。点击完成，设置完成。

`mail.huoyijie.cn` 作为 `SMTP` 邮箱服务器服务端时，会接收来自其他域名邮箱服务器的邮件投递请求，也会接收来自邮件客户端 `Thunderbird` 的发送邮件请求。来自其他邮箱服务器的请求会走 `25` 端口，通过机会加密StartTLS机制来决定是否加密传输数据包括邮件。还记得在配置 SSL 那节时提到过，我配置了 `TrustAsia TLS RSA CA` 证书，包含公钥，可以和连接上来的客户端完成加密协商（客户端侧一般不验证证书）。来自邮件客户端 `Thunderbird` 的发件请求会走 `465` 端口，客户端侧是要验证证书，并完成加密协商，后面所有数据都是加密传输。

`mail.huoyijie.cn` 作为 `IMAP` 邮箱服务器服务端时，来自邮件客户端 `Thunderbird` 的取件请求会走 `993` 端口，客户端侧是要验证证书，并完成加密协商，后面所有数据都是加密传输。因为涉及到下载用户邮箱数据，所以 `IMAP` 服务端要求连接上来的客户端一定要进行用户认证，通过认证才能收取邮件。

`SMTP` 服务端默认没有要求用户认证，也就是说可以连接上服务器端就可以发送发件请求。后面会为 `SMTP` 服务端配置用户认证。