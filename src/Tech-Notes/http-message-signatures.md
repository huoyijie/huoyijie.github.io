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

是指在 Server 收到 http 消息后，通过签名验证的方式，确认该消息真实且完整。

## HTTPS 不能保证客户端或服务器接收到真实、完整的消息

## 如何实现端到端 HTTP 消息真实、完整性验证

## HTTP Message Signatures 协议

## HTTP Message Signatures 实现

## HTTP 请求签名实例