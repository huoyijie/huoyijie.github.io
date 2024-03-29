
# 目录

* [前言](README.md)
* [第一章 UDPack 协议](chapter-01/README.md)
    * [1.1 基本概念](chapter-01/1.1-basic-concept.md)
    * [1.2 二进制帧格式](chapter-01/1.2-binary-protocal.md)
    * [1.3 UDPack Interface](chapter-01/1.3-udpack-interface.md)
* [第二章 会话(Session)管理](chapter-02/README.md)
    * [2.1 会话 ID](chapter-02/2.1-session-id.md)
    * [2.2 会话状态](chapter-02/2.2-session-state.md)
    * [2.3 连接握手](chapter-02/2.3-handshake.md)
    * [2.4 关闭会话](chapter-02/2.4-goaway.md)
    * [2.5 会话保持与过期清理](chapter-02/2.5-keepalive.md)
    * [2.6 断线重连](chapter-02/2.6-session-reconnect.md)
    * [2.7 Session Interface](chapter-02/2.7-session-interface.md)
* [第三章 双向流(Stream)](chapter-03/README.md)
    * [3.1 流状态](chapter-03/3.1-stream-state.md)
    * [3.2 打开流](chapter-03/3.2-openstream.md)
    * [3.3 关闭流](chapter-03/3.3-closestream.md)
    * [3.4 Stream Interface](chapter-03/3.4-stream-interface.md)
* [第四章 数据包(Packet)](chapter-04/README.md)
    * [4.1 数据包(Packet)状态](chapter-04/4.1-packet-state.md)
    * [4.2 发送数据包(Packet)](chapter-04/4.2-send-packet.md)
    * [4.3 接收数据包(Packet)](chapter-04/4.3-recv-packet.md)
    * [4.4 周期检查发送缓存](chapter-04/4.4-check-send-packs.md)
    * [4.5 周期检查接收缓存](chapter-04/4.5-check-recv-packs.md)
* [第五章 包分片(Fragment)](chapter-05/README.md)
    * [5.1 分片(Fragment)状态](chapter-05/5.1-fragment-state.md)
    * [5.2 发送分片(Fragment)](chapter-05/5.2-send-fragment.md)
    * [5.3 接收分片(Fragment)](chapter-05/5.3-recv-fragment.md)
    * [5.4 周期检查分片(Fragment)发送缓存](chapter-05/5.4-check-send-frags.md)
    * [5.5 周期检查分片(Fragment)接收缓存](chapter-05/5.5-check-recv-frags.md)
* [第六章 心跳检测](chapter-06/README.md)
* [第七章 拥塞控制与流控制](chapter-07/README.md)
* [第八章 协议 Node.js 实现](chapter-08/README.md)
* [第九章 UDPack 应用之代理服务器](chapter-09/README.md)
    * [9.1 SOCKS 5 代理协议](chapter-09/9.1-socks5.md)
    * [9.2 创建 UDPack 对象](chapter-09/9.2-udpack.md)
    * [9.3 自定义代理协议](chapter-09/9.3-custom-protocol.md)
    * [9.4 代理客户端](chapter-09/9.4-client.md)
    * [9.5 代理服务端](chapter-09/9.5-server.md)
    * [9.6 应用程序打包](chapter-09/9.6-pkg.md)
* [第十章 UDPack 应用之实时双向通信框架](chapter-10/README.md)
* [附录](chapter-11/README.md)