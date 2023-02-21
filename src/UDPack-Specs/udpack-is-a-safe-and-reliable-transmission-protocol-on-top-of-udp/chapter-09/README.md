# 第九章 UDPack 应用之代理服务器

代理服务器一般在企业内部比较常见，可起到类似防火墙的作用。企业内部员工通过设置代理服务器上网，所有的网络请求都会经过代理服务器。可通过全局屏蔽一些与工作内容无关甚至有风险的网站，提高互联网访问的安全性。

代理协议主要有 SOCKS 和 HTTP 两种。SOCKS 是一种网络传输协议，主要用于客户端与外网服务器之间通讯的中间传递。当防火墙后的客户端要访问外部的服务器时，就跟 SOCKS 代理服务器连接。这个代理服务器控制客户端访问外网的资格，允许的话，就将客户端的请求发往外部的服务器。这个协议最初由 David Koblas 开发，而后由 NEC 的 Ying-Da Lee 将其扩展到 SOCKS4。最新协议是 SOCKS5，与前一版本相比，增加支持 UDP、验证，以及 IPv6。根据 OSI 模型，SOCKS 是会话层的协议，位于表示层与传输层之间。

SOCKS 工作在比 HTTP 代理更低的层次，SOCKS 使用握手协议来通知代理软件其客户端试图进行的 SOCKS 连接，然后尽可能透明地进行操作。SOCKS 代理不理解 HTTP 协议，不会试图解析流量内容，可确保用户隐私安全。而 HTTP 代理理解 HTTP 协议，通常会解释和重写报头，甚至可以分析流量，执行更高层次的过滤。

```
+--------------+          +--------------+          +-------------+
|              |  SOCKS   |              |          |             |
|   Browser    +----------> Proxy Server +----------> google.com  |
|              |   LAN    |              | Internet |             |
+--------------+          +--------------+          +-------------+
  PC/Mac/Phone              Proxy Server
```

如上图所示，一般 SOCKS 代理服务器是这样部署的: 电脑手机等设备配置好位于局域网内部的 Proxy Server，Browser 与 Proxy Server 之间建立 TCP 连接，并通过 SOCKS 协议中转上网。本章会基于 UDPack 协议实现一个 SOCKS 代理，完全采用另外一种部署结构，使得 Browser 与 SOCKS 代理之间建立 UDPack 会话，并通过 UDPack 协议中转流量，实现代理上网。

```
+----------------------------------------------+
|                                              |
|                                              |
|                                              |
| +----------------+         +---------------+ |       +---------------+          +-------------+
| |                |  SOCKS  |               | |UDPack |               |          |             |
| |    Browser     +---------> UDPack Client +-+-------> UDPack Server +----------> google.com  |
| |                |localhost|               | | LAN   |               | Internet |             |
| +----------------+         +---------------+ |       +---------------+          +-------------+
|                                              |         Proxy Server
|                                              |
|                                              |
+----------------------------------------------+
                PC/Mac/Phone
```

如上图所示，电脑手机等设备配置好位于本机(localhost)的 UDPack Client，Browser 与 UDPack Client 之间建立 TCP 连接，并通过 SOCKS 协议中转流量。UDPack Client 与 UDPack Server 之间建立 UDPack 会话，并通过 UDPack 协议中转流量，实现代理上网。

可以看到新的部署结构多了 UDPack Client 节点，是因为 Browser 本身不理解 UDPack 协议，必须通过 UDPack Client 将 SOCKS 协议转为 UDPack 协议。UDPack Client 与 UDPack Server 之间通过基于 UDP 实现的 UDPack 协议通信。