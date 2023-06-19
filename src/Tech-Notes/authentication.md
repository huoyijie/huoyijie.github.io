# Authentication

## 用户身份识别

除非网站上所有的数据和功能都是公开的，否则你一定会需要用户、角色和权限系统。不同的用户会有不同的权限，你需要识别出请求来自哪个用户。

想对用户请求进行识别，就需要前端在向服务器发送请求时携带的身份标识信息，一般均采用 Token 机制来实现，通过 Cookie/Header 发送。Token 可以是通过对称密钥加密的，也可以是采用私钥签名的。主要目的是阻止第三方伪造 Token 以其他用户身份访问网站。

* 对称加密

```
Token = 对称加密([用户ID,时间戳...], 密钥)
```

密钥是保密的，只存在于服务器上，Token 的生成与解析由服务器负责，客户端也无法解密 Token。Token 中通常只有 userId 和 timestamp。

* 非对称加密 (JWT)

```
Identify = base64.encode([用户ID,时间戳...])
Token = Identify + 签名(哈希(Identify), 私钥)
```

私钥是保密的，只存在于服务器上，Token 的生成由服务器负责，通常由原始用户身份信息 Identify 和签名两部分组成，用户身份 JSON 字符串，通过 base64 算法进行编码作为 Identify，由于 Identify 长度不固定，可用与客户端约定好的哈希算法进行运算，再用私钥进行签名。

这种方法生成的 Token 并不是加密的，客户端甚至第三方拿到 Token 后都可以解析出其中的用户身份信息。

客户端持有相对应的公钥，可对 Token 进行签名验证，验证此 Token 是否是具有私钥的服务器生成的，并可以从 Token 中解析出原始用户信息。Token 可携带很多用户信息，甚至可包含用户角色权限等信息。

```
[Identify, signature] = split(Token)
签名验证(哈希(Identify), signature, 公钥)
[用户ID,时间戳...] = base64.decode(Identify)
```

对称加解密、非对称加解密、签名、哈希等算法都属于密码学范畴，之前有介绍过 Golang 语言的各种密码学方法。[密码学方法](https://huoyijie.cn/gitbooks/writing-a-cloud-encrypting-file-system-with-golang-and-FUSE/latest/crypto.html)

* 随机生成 UUID 作为 Token

还有一种方法是随机生成 UUID 作为 Token，但是为了防止第三方伪造，需要服务器端保存所有有效的 Token 集合。服务器再接收到请求中携带的 Token 后需要进行有效性验证。

## 如何进行用户认证并下发 Token

用户认证的原理在于要求你输入只有你自己知道的秘密信息，并由服务器进行校验。

* username/password
* email/password

添加用户时，密码上传到服务端后，确保所有中间服务器不记录、不打印日志，并经过 bcrypt 哈希后存储于数据库中。服务器没有任何理由保存用户原始密码。