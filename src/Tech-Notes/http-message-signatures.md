# 实现端到端 HTTP 消息真实、完整性验证

## 端到端 HTTP 消息是指什么?

```
+--------+       +-------+          +-------+       +--------+
|        |       |       |          |       |       |        |
| Client |-http->| Proxy |--https-->| Nginx |-http->| Server |
|        |       |       |          |       |       |        |
+---+----+       +---+---+          +---+---+       +----+---+
    |                |                  |                |
    |                |<-TLS end to end->|                |
    |                                                    |
    |<------------- http message end to end ------------>|
```

如上图，有些时候 Client 可能部署在 Proxy 后面，http 消息需要 Proxy 转发出去。更多的时候，业务服务器都会部署在像 Nginx 这样的反向代理后面。从 Client 发出 http 消息到 Server 接收到消息，即端到端 HTTP 消息。无论是从 Client 直接发出 http 消息到 Nginx 接收或从 Client 发出经 Proxy https 转发到 Nginx 接收，属于 TLS 通信端到端。

## 什么是真实、完整性验证？

是指 Client、Server 可以通过某种方法验证消息在传输过程中，没有被任何中间节点修改过。

## HTTPS 不能保证客户端或服务器接收到真实、完整的消息

Nginx 或 Proxy 等中间节点可以修改 http message，在不影响 http 语义的情况下，Client 或 Server 无法确认消息是否真实完整。

## 如何实现端到端 HTTP 消息真实、完整性验证

Client 在请求发出前对 request 进行签名，Server 端收到后可验证签名，然后返回签过名的 response，Client 同样验证签名。

## HTTP Message Signatures 协议

中间节点在转发 http message 时，调整 header 中的字段顺序、转换 Key 的大小写、添加键值对等很多行为都是被 http 协议允许的，所以对 http message 进行签名需要一个统一的 [HTTP Message Signatures](https://httpwg.org/http-extensions/draft-ietf-httpbis-message-signatures.html) 扩展协议（目前还是草案）。

## HTTP Message Signatures 实现

[Python 语言实现](https://github.com/pyauth/http-message-signatures)

## HTTP 请求签名实例

接下来我们基于上面的 Python 实现来编写一个 HTTP 请求签名实例程序，文中所有代码已放到 [Github http-request-signature](https://github.com/huoyijie/tech-notes-code) 目录下。

*前置条件*

已安装 Python 3.10+
已安装 IDE （如 vscode）

首先初始化项目（[参考 Flask 安装](https://flask.palletsprojects.com/en/2.3.x/installation/#virtual-environments)）

```bash
$ mkdir http-request-signature && cd http-request-signature
# Linux 确保已安装 python3.10-venv，Windows/MacOS 跳过
$ sudo apt install python3.10-venv
# 创建并激活虚拟环境
$ python3 -m venv .venv
$ . .venv/bin/activate
# 安装 Flask、http 消息签名扩展协议、requests
$ pip install Flask http-message-signatures requests
```

**密钥生成函数 & Magic number**

创建 app_key.py 文件，添加获取密钥相关方法

```python
import math
from http_message_signatures import HTTPSignatureKeyResolver

# 配置密钥生成函数
key_gen = lambda key_id: bytes(str(math.sqrt(2023)), 'utf-8')

class MyHTTPSignatureKeyResolver(HTTPSignatureKeyResolver):
  def resolve_public_key(self, key_id: str):
    return key_gen(key_id=key_id)

  def resolve_private_key(self, key_id: str):
    return key_gen(key_id=key_id)
```

注意 key_gen 密钥生成函数:

```python
# a = sqrt(2023)
# b = to_string(a)
# c = get_bytes(b)
key_gen = lambda key_id: bytes(str(math.sqrt(2023)), 'utf-8')
```

Client 与 Server 只要提前约定好一个 Magic number（如: 2023），就可以通过这个函数生成密钥。如果 Client 是运行在用户电脑上的浏览器或者手机上的 App，如何防止恶意用户通过反编译代码或者查看 JS 代码来找出这个密钥生成函数甚至这个 Magic number 呢？代码肯定要进行混淆压缩及其他反编译操作，最大程度隐匿好这个函数及 Magic number。

**Server**

```python
from flask import Flask, request, abort

from http_message_signatures import HTTPMessageVerifier, algorithms
from app_key import MyHTTPSignatureKeyResolver

verifier = HTTPMessageVerifier(signature_algorithm=algorithms.HMAC_SHA256, key_resolver=MyHTTPSignatureKeyResolver())

app = Flask(__name__)

@app.route('/hello/<greeting>', methods=['POST'])
def hello(greeting):
  try:
    print(request.headers)
    verifyReuslt = verifier.verify(request)
    print(verifyReuslt)
    # todo 验证 body 与 content-digest 一致
  except Exception as e:
    print(e)
    abort(400)

  return 'hello, %s | HELLO, %s!' %(greeting, request.json['HELLO'])

@app.errorhandler(400)
def bad_request(error):
  return 'Bad Request', 400
```

**Client**

```python
import requests, hashlib, http_sfv

from http_message_signatures import HTTPMessageSigner, algorithms
from app_key import MyHTTPSignatureKeyResolver

signer = HTTPMessageSigner(signature_algorithm=algorithms.HMAC_SHA256, key_resolver=MyHTTPSignatureKeyResolver())

request = requests.Request('POST', 'http://localhost:5000/hello/world', json={'HELLO': 'WORLD'})
request = request.prepare()
# 计算 request.body Digest
request.headers['Content-Digest'] = str(http_sfv.Dictionary({'sha-256': hashlib.sha256(request.body).digest()}))

# todo 签名过期时间
signer.sign(request, key_id='my_key_id', covered_component_ids=('@method', '@authority', '@target-uri', 'content-digest'))
print(request.headers)

response = requests.Session().send(request=request)
print(response.content)
```