# 第四章 数据包(Packet)

UDPack 是数据传输协议，发送与接收 Packet 是协议的核心部分。前面提到过，既可以通过 Session，也可以通过 Stream 发送与接收数据。Session 和 Stream 都有自己的发送与接收数据缓存，除了这个差别，其他都是一样的。并且每个 Stream 都有自己独立的发送与接收缓存，互不影响。后面内容会描述如何发送与接收数据包，如果没有特别说明，不特指 Session 或 Stream。但无论是通过 Session 或者 Stream，它们都必须处于 READY 状态。

Session 握手时协商了对称加密参数(algorithm, key, iv)，待发送的数据通过已协商的 algorithm, key, iv 进行对称加密，然后发送出去，另一端通过同样的参数解密还原数据。由于 MTU 的原因，每次发送的数据包长度不能超过 MTU，否则可能会被路由中的中间节点默默丢弃。协议定义了 Payload 数据最大长度 MAX_PAYLOAD_LEN，如果发送的数据长度超过 MAX_PAYLOAD_LEN 需要进行分片处理。本章主要介绍数据长度小于 MAX_PAYLOAD_LEN 的数据包的发送与接收。