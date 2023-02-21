# UDPack 2.0 协议

UDPack 2.0 是一种基于 UDP 的安全可靠数据包传输协议。UDPack 为上层应用程序提供了面向数据包(Packet)的安全、可靠的双向数据传输能力。

相比 [UDPack 1.0](https://huoyijie.cn/article/1b11c0d06ef311eda78619d87c80c179/) 协议，2.0 协议进行了极大的简化，某些地方借鉴了 `TCP` 协议中的概念。这个协议的目标是在 `TCP` 协议无法使用或者更适宜 `UDP` 协议的时候，提供一个安全可靠的传输协议，实现要比 `QUIC` 协议更简单、更轻量。

为了更清晰了解协议变化，下面给出帧格式定义:

## 二进制帧格式

### 1.0 协议帧格式

```
0                   1                   2                   3
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+---------------------------------------------------------------+
|                        Session Id (32)                        |
+---------------------------------------------------------------+
|                        timestamp ...                          |
+ - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
|                    ... timestamp (64)                         |
+-+-+-+-----+-+-+-------+-------+-------------------------------+
|S|P|F| RSV |S|F|  RSV  |OPCODE |          Payload len          |
|T|A|R| (3) |L|I|  (4)  | (4)   |             (16)              |
|R|C|A|     |O|N|       |       |                               |
|E|K|G|     |W| |       |       |                               |
+-+-+-+-----+-+-+-------+-------+-------------------------------+
|                        Stream Id (32)                         |
|                       (if STRE set to 1)                      |
+---------------------------------------------------------------+
|                        Packet Id (32)                         |
|                       (if PACK set to 1)                      |
+---------------------------------------------------------------+
|                        Fragment Id (32)                       |
|                       (if FRAG set to 1)                      |
+---------------------------------------------------------------+
:                        Payload Data ...                       :
+---------------------------------------------------------------+
```

### 2.0 协议帧格式

```
+--------------------+-----------------+
|      required      |    optional     |
+--------------------+-------+---------+
|           header           |  body   |
+------+------+------+-------+---------+
| uuid | kind | flag |pack id|pack data|
+------+------+------+-------+---------+
|  8   |  1   |  1   |   4   |   len   |
+------+------+------+-------+---------+
```

下面是二进制帧字段说明:

| 字段 | 长度(单位: 字节) | 说明 |
|----|----|----|
| uuid | 8 | 唯一标识一个连接会话递增 ID，由 SonyFlake 算法生成，类似 SnowFlake 算法 |
| kind | 1 | 帧类型 |
| flag | 1 | 帧位标识，不同帧类型可能会利用位 0~7 进行特殊设定 |
| pack id | 4 | 在一个连接会话内递增 ID，唯一标识一个数据包 |
| pack data| len | 载荷数据 |
点击了解 [SonyFlake](https://github.com/sony/sonyflake) 算法

下面是帧类型定义和说明:

```rust
#[derive(Debug)]
pub enum FrameKind {
  Open(bool),
  Ready,
  Ping,
  Pong,
  Data { pack_id: u32, bytes: Bytes },
  Ack { pack_id: u32 },
  Sync { pack_id: u32 },
  Slow { pack_id: u32 },
  Lost { pack_id: u32 },
  Shutdown,
  Close,
  Error,
}
```

| 帧类型 | 说明 |
|----|----|
| Open | 请求建立会话连接 |
| Ready | 已建立连接 |
| Ping | 心跳检查请求 |
| Pong | 心跳检查应答 |
| Data | 数据帧 |
| Ack | Ack 帧 |
| Sync | 与发送方同步接收进度帧 |
| Slow | 通知发送方降速帧 |
| Lost | 通知发送方数据丢失帧 |
| Shutdown | 数据发送完毕，准备关闭连接 |
| Close | 关闭连接 |
| Error | 发生错误 |

### 主要变化

1. 唯一标识连接的 4 字节 `Session Id` 改为 8 字节 `uuid`。
2. 删除 `timestamp` 时间戳，互联网主机时间存在较大误差，差距甚至达到百毫秒，该字段意义不大。
3. `OPCODE` 字段及各种 `Flag` 位，由 1 字节 `kind` 和 1 字节 `flag` 取代。.
4. 删除 `Payload len` 字段
5. 删除 `Stream Id` 字段
6. 删除 `Fragment Id` 字段

2.0 协议废除了 `Stream` 与 `Fragment` 概念，使得协议得到了极大简化。但是，需要配合上层 `SocketIO` 协议来实现发送任意长度的数据包。`SocketIO` 协议在后面会介绍。

1.0 协议内会自动对数据包进行分片与组装，传输层传输的分片 `Fragment` 最大长度为 65535，实际传输的 `Packet` 理论上没有大小限制，协议具体实现会设定最大分片数量，进而限定了 `Packet` 长度上限。

2.0 协议因为没有分片设定，不可能发送任意长度数据包，因此会把数据包划分成小的数据包发送出去，并由数据接收方进行缓冲与组装，但是 2.0 协议没有规定数据包大小，因此在 `UDPack` 协议层无法实现数据组装。这是进行协议简化带来的弊端。

数据接收端会定期发送 `Sync` 帧，与数据发送端进行同步，快速同步接收进度，防止因为 `Ack` 帧丢失造成不必要的超时数据重发。`Sync` 与 `Ack` 的区别是，后者只针对单个数据包进行确认，前者会通知发送端小于某个序号的 `Packet` 已全部收到。

如果发送端发送数据过快造成接收端缓存已满，会发送 `Slow` 帧给发送端，通知其降低数据包发送速度。如果发送端自己检测到数据包重发比例超过一定阈值后，也会主动降低数据包发送速度。发送端的发送速率主要是靠滑动窗口实现的，每个发送周期内只发送窗口大小的数据包数量，如果网络情况良好，滑动窗口会快速变大，以尽快提升数据包发送速率。如果收到 `Slow` 帧或者发送端重试率超过阈值，都会减小滑动窗口。需要控制速率的主要是 `Data` 帧，具体实现时，可尽量优先发送控制帧。

为了上层协议可以按照顺序接收数据包，需要按照序号顺序返回数据给上层协议，如果某个序号的数据包丢失，即时收到了在这个序号之后的数据包，也只能缓存在接收缓存中。此时，可通过发送 `Lost` 帧给发送端，通知其尽快重新发送数据包，尽量避免发送端等到发送超时后重试，因为在这段时间内，接收端接收到的数据都只能先缓存，可能会造成缓存满，然后通知发送端降速。接收端一段时间内不能够向上交付数据，会造成"堵塞"。

具体在实现协议时，可以对待发送的数据帧进行优先级排队，主要的目标是降低整体数据包重试率，即一段时间内，总的数据包重试次数除以总的正在发送数据包数量。重试率越低，说明网络情况越好，可以更快的交付更多的数据。

连接期间，遇到不可恢复错误，则可向对方发送 `Error` 帧，连接双方会尽快断开连接，并可重新进行连接。

## 2.0 协议实现

我曾采用 `nodejs` 实现了 1.0 协议，并且经过了大量的测试，并基于该协议实现了一个科学上网工具。但是，实现 1.0 协议过程中，我发现自己对 `js` 语言的掌握有限，总是不能用最好的方式去写，加上 `js` 语言本身的弊端，代码量上来后可读性较差，而且 `nodejs` 本身也没有啥性能优势。

实现 2.0 协议时，我最早尝试通过 C 语言来实现，初步写好协议代码后，还是选择通过实现科学上网工具来测试，刚启动程序时协议正常运行，过一段时间就会发生内存段错误闪退。我发现以我的水平，手动管理内存是很难不出现内存错误的，而且还很多，调试过程真的很难，很难定位到问题，就决定放弃使用 C 语言了。

后来，我了解到了 `Rust` 这个能够最大程度上保证`内存安全`，而性能上又非常理想的语言。`Rust` 提供了大量的新颖的现代编程语言特性，提供内存安全又不引入垃圾收集器，但也出现了变量所有权、借用等等很多复杂的概念，上手难度不小。

最后，我还是决定采用 `Rust` 实现 2.0 协议，并重新实现了那个科学上网工具。边学习 `Rust` 边写代码的过程还是挺痛苦的，很多新概念理解起来挺困难，尤其是变量所有权、引用和协程等。下面是基于 `UDPack` 实现的客户/服务端通信程序代码。

### Server

```rust
use rust_udpack::{Transport, Udpack};
use std::io;
use tokio::{signal, task::JoinHandle};

#[tokio::main]
async fn main() -> io::Result<()> {
  let mut udpack: Udpack = Udpack::new("0.0.0.0:8080").await?;

  loop {
    tokio::select! {
      res = udpack.accept() => {
        let _handle: JoinHandle<io::Result<()>> = tokio::spawn(async move {
          let mut transport: Transport = res.unwrap();

          while let Some(bytes) = transport.read().await {
            println!("{:?}", bytes);
            transport.write(bytes).await?;
          }
          Ok(())
        });
      }
      _ = signal::ctrl_c() => {
        println!("ctrl-c received!");
        udpack.shutdown().await?;
        break;
      }
    }
  }
  Ok(())
}
```

### Client

```rust
use bytes::Bytes;
use rust_udpack::{Transport, Udpack};
use std::io;
use tokio::{
  signal,
  time::{interval, Duration},
};

#[tokio::main]
async fn main() -> io::Result<()> {
  let udpack: Udpack = Udpack::new("0.0.0.0:0").await?;
  let mut transport: Transport = udpack.connect("127.0.0.1:8080").await?;
  let mut interval = interval(Duration::from_secs(3));

  loop {
    tokio::select! {
      res = transport.read() => {
        if let Some(bytes) = res {
          println!("{:?}", bytes);
        }
      }
      _ = interval.tick() => {
        transport.write(Bytes::from("Hello, world!")).await?;
      }
      _ = signal::ctrl_c() => {
        println!("ctrl-c received!");
        udpack.shutdown().await?;
        return Ok(());
      }
    };
  }
}
```

可以看到，主要有两个概念，`Udpack` 和 `Transport`。通过 `Udpack` 建立 endpoint，然后客户端连接服务端，返回 `Transport`，通过 `Transport`，可以双向读写数据。

## SocketIO 协议

前面提到过，UDPack 2.0 协议不能够发送任意长度的数据包，原因在于 Internet 存在客观的数据长度限制，这是信号载体物理性质决定的，一次调制的信号长度超过一定的大小，各种干扰因素可能会使信号失真。虽然 `IP` 协议支持数据分片，但这个是兜底的，会造成传输性能下降。所以 `UDPack` 1.0 协议支持 `Fragment` 概念，协议会自动对数据包进行分片和组装。为了简化协议，2.0 协议中没有 `Fragment` 概念了，通过 2.0 协议发送的数据包同样会被自动分割成小的数据包按照顺序依次发送，此时，必须借助上层额外的 `SocketIO` 协议来辅助进行数据包组装。

SocketIO 主要基于 `LengthDelimitedCodec` 实现，简单来说，就是在数据包前加一个 `length` 字段。数据接收端按照序号缓存收到的 `packet`，然后通过读取 `length` 字段得到当前数据包长度，再读取相应的数据，然后返回给上层应用。

另一方面，`UDPack` 2.0 协议本身没有对通信数据进行加密，因此，`SocketIO` 层增加了数据自动加解密处理。

`SocketIO` 协议格式如下:

```
+--------+----------------------+
| length | ecrypted packet data |
+--------+----------------------+
```

下面是基于 `SocketIO` 实现的客户/服务端通信程序代码。

### Server

```rust
use rust_socketio::SocketIO;
use std::io;

const SECRET_KEY: &[u8; 32] = b"ade2ff15798d44959d2846974bbf0bb3";
const SECRET_IV: &[u8; 16] = b"bd3c01bfb8c2edca";

#[tokio::main]
async fn main() -> io::Result<()> {
  let mut io = SocketIO::new("0.0.0.0:8080", *SECRET_KEY, *SECRET_IV).await?;
  loop {
    tokio::select! {
      res = io.accept() => {
        let mut socket = res.unwrap();
        tokio::spawn(async move {
          loop {
            match socket.read().await {
              Ok(Some(bytes)) => {
                if let Err(e) = socket.write(bytes).await {
                  eprintln!("socket.write failed; err = {:?}", e);
                }
              }
              Ok(None) => return,
              Err(e) => {
                eprintln!("socket.read failed; err = {:?}", e);
                return;
              }
            };
          }
        });
      }
      _ = tokio::signal::ctrl_c() => {
        println!("ctrl-c received!");
        io.shutdown().await?;
        break;
      }
    }
  }
  Ok(())
}
```

### Client

```rust
use bytes::Bytes;
use rust_socketio::SocketIO;
use std::io;
use tokio::{time, time::Duration};

const SECRET_KEY: &[u8; 32] = b"ade2ff15798d44959d2846974bbf0bb3";
const SECRET_IV: &[u8; 16] = b"bd3c01bfb8c2edca";

#[tokio::main]
async fn main() -> io::Result<()> {
  let io: SocketIO = SocketIO::new("0.0.0.0:0", *SECRET_KEY, *SECRET_IV).await?;
  let mut socket = io.connect("127.0.0.1:8080").await?;
  let mut interval = time::interval(Duration::from_secs(3));
  loop {
    tokio::select! {
      res = socket.read() => {
        if let Some(bytes) = res? {
          eprintln!("received {:?}", bytes);
        } else {
          println!("EOF");
          io.shutdown().await?;
          break;
        }
      }
      _ = interval.tick() => {
        socket.write(Bytes::from("Hello, world!")).await?;
      }
      _ = tokio::signal::ctrl_c() => {
        println!("ctrl-c received!");
        io.shutdown().await?;
        break;
      }
    }
  }
  Ok(())
}
```

## 关于 UDPack 2.0

因为文章篇幅所限，没有办法展开详细看代码。我已经在 Github 上面开源了相关项目，如果你对源代码感兴趣，欢迎下载查看，有问题可以交流探讨！UDPack 2.0 项目是基于业余兴趣，为了学习的目的而建立的，代码实现可能不够严谨，测试可能也不够充分，请勿用于商业场景，非常感谢大家的支持！

### UDPack 2.0 实现

[rust-udpack](https://github.com/huoyijie/rust-udpack)

### SocketIO 实现

[rust-socketio](https://github.com/huoyijie/rust-socketio)

### Socks5 实现

[rust-socks5](https://github.com/huoyijie/rust-socks5)

### VPN Proxy 实现

[rust-proxy](https://github.com/huoyijie/rust-proxy)

### fork of SonyFlake

[rust-sonyflake](https://github.com/huoyijie/rust-sonyflake)