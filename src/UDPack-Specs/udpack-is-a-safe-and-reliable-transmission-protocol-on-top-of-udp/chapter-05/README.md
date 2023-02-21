# 第五章 包分片(Fragment)

本章主要描述长度超过 MAX_PAYLOAD_LEN 的 Packet 的分片发送与接收。数据包状态参见 [4.1 数据包(Packet)状态](chapter-04/4.1-packet-state.md)。

发送端首先创建 Packet 对象，状态为 FRAGMENTING，更新当前时间为 Packet 发送时间。把待发送加密数据进行分片(Fragment)处理，每个 Fragment 长度不超过 MAX_PAYLOAD_LEN。每个分片会分配递增序号(fragId)，是一个 4 字节长无符号整数，从 1 开始。分片数据会被封装到 Fragment 中。