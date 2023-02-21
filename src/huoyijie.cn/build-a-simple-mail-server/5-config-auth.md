# 配置邮箱用户认证

前面提到过 `SMTP` 发件服务端默认没有要求用户认证，也就是说可以连接上服务器端就可以发送发件请求。`IMAP` 收件服务端是要求必须进行用户认证的，接下来调整下配置，让 `SMTP` 也进行用户认证，并且让认证方式和 `IMAP` 收件服务一致。

`IMAP` 收件服务基于 dovecot 搭建的，可以配置 postfix 基于 dovecot 来进行用户认证。首先打开 dovecot 配置文件。配置文件位于目录 `/etc/dovecot/` 下。

```conf
conf.d/10-master.conf:
    service auth {
      ...
      unix_listener /var/spool/postfix/private/auth {
        mode = 0660
        # Assuming the default Postfix user and group
        user = postfix
        group = postfix        
      }
      ...
    }

conf.d/10-auth.conf
    auth_mechanisms = plain login
```

```bash
sudo service dovecot reload
```

编辑 postfix 加入配置文件 `/etc/postfix/main.cf`

```conf
smtpd_sasl_type = dovecot
smtpd_sasl_path = private/auth
smtpd_sasl_auth_enable = yes
```

```bash
sudo service postfix reload
```

打开 `Thunderbird` 删除之前建立的邮箱帐号，重新新建帐号测试。此时没有邮箱帐号密码，无法接收邮件，也不能再发送邮件。

如果想看邮件服务日志，可以查看文件，里面可以看到邮件客户端连接邮件服务器，请求登录认证、登录成功或失败以及发送邮件等日志信息。

```bash
tail -f /var/log/mail.log
```