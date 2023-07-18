# 基于 Websocket 和 React 实现一个 IM 原型

本文主要基于 Websocket、React、msgpack、Tailwind CSS 等技术实现一个即时通信 IM 原型，支持发送/接收实时消息、离线消息，支持显示用户在线状态，支持显示未读消息数等功能。主要是探索如何通过 websocket+msgpack 实现客户端与服务器的双向通信。前端复用了[基于 Server Sent Events 和 React 实现一个 IM 原型](https://huoyijie.cn/docsifys/Tech-Notes/chat-with-sse-react)里的代码，采用 React 简化了应用状态与 UI 之间的同步，并引入时下非常流行的 Tailwind CSS 库。

![chat-with-sse-react](https://cdn.huoyijie.cn/uploads/2023/07/chat-with-sse-react.png)

## Github 项目地址

本文代码在 [tech-notes-code/chat-with-websocket](https://github.com/huoyijie/tech-notes-code) 目录下，代码注释很详细，可以边看文章边看代码。

## Websocket

![websocket logo](https://cdn.huoyijie.cn/uploads/2023/07/websocket-logo.jpg)

WebSocket 是一种与 HTTP 不同的协议。两者都位于 OSI 模型的应用层，并且都依赖于传输层的 TCP 协议。虽然它们不同，但是 WebSocket 可通过 HTTP 端口 80 和 443 进行工作，并支持 HTTP 代理，从而使其与 HTTP 协议兼容。为了实现兼容性，WebSocket握手使用 HTTP Upgrade 头从 HTTP 协议升级为 WebSocket 协议。

WebSocket 协议支持 Web 浏览器（或其他客户端应用程序）与 Web 服务器之间的交互，具有较低的开销，便于实现客户端与服务器的实时数据传输。服务器可以通过标准化的方式来实现，而无需客户端首先请求内容，并允许消息在保持连接打开的同时来回传递。通过这种方式，可以在客户端和服务器之间进行双向持续对话。*更多介绍可参考[WebSocket](https://zh.wikipedia.org/zh-cn/WebSocket)*

## 客户端与服务器通信

```
                         (huoyijie)
                          +------+
                          |      |
+------+                  |Chrome|
|      |<---websocket---->|      |
| HTTP |        |         +------+
|      |    双向通信连接
|Server|        |         +------+
|      |<---websocket---->|      |
+------+                  |Chrome|
                          |      |
                          +------+
                           (jack)
```

每个用户打开客户端后，会自动连接服务器并发送升级协议请求，服务器确认升级协议后，双方建立 websocket 连接，后续客户端与服务器可通过 websocket 连接实现双向通信。从 HTTP 到 Websocket 的协议升级请求也称作 Websocket 握手。主要过程如下:

**客户端发送如下 GET 请求**

```bash
GET /chat HTTP/1.1
Host: example.com:8000
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
Sec-WebSocket-Version: 13
```

**服务器返回如下响应完成 Websocket 握手**

```bash
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
```

## msgpack

![msgpack introduce](https://cdn.huoyijie.cn/uploads/2023/07/msgpack-intro.png)

[msgpack](https://github.com/msgpack/msgpack) 是一种高效的二进制序列化格式，很像 json，但是运行速度更快、序列化后数据更小，支持各种语言。

## 服务器

**Web 服务器**

服务器使用 [Gin](https://github.com/gin-gonic/gin) Web 框架搭建 Web 服务器。

**Websocket**

服务器通过 [gorilla/websocket](https://github.com/gorilla/websocket) 库支持 Websocket 协议。

**msgpack**

服务器使用 [vmihailenco/msgpack](https://github.com/vmihailenco/msgpack) 库实现消息序列化/反序列化。

## 客户端

**msgpack**

前端使用 [msgpack/msgpack-javascript](https://github.com/msgpack/msgpack-javascript) 库实现消息序列化/反序列化。

**React**

前端使用的 [React](https://zh-hans.react.dev/) 框架，通过定义函数式组件把 UI 拆分成不同的嵌套组件，在每个组件内部控制自己的状态和样式。应用组件思想构建 UI，有点像搭乐高很有趣。一方面 React 可以帮助你根据应用状态自动更新 UI 展示，无需手动操作 DOM 元素。另一方面没有了长长的 html 代码片段，而是拆分成了很多的 jsx 组件文件，每个组件单独一个文件，代码结构更清晰。

```bash
$ tree -l
├── public
│   ├── images
│   │   ├── huoyijie.svg
│   │   ├── jack.svg
│   │   └── rose.svg
│   └── js
│       ├── App.jsx
│       ├── ChatBoxHeader.jsx
│       ├── ChatBoxInput.jsx
│       ├── ChatBox.jsx
│       ├── ChatBoxMessageList.jsx
│       ├── Chat.jsx
│       ├── Header.jsx
│       ├── Message.jsx
│       ├── PaperPlane.jsx
│       ├── UserList.jsx
│       ├── Users.jsx
│       └── ws.js
├── client.go
├── data.go
├── hub.go
├── main.go
└── templates
    └── index.htm
```

上述是代码目录结构，`public/js` 目录下是所有的函数组件，下图标记了主要的组件。

![react-components](https://cdn.huoyijie.cn/uploads/2023/07/chat-with-sse-react-components.png)

**Tailwind CSS**

CSS 样式是用时下非常流行的 [Tailwind (A utility-first CSS framework)](https://tailwindcss.com/) 写的，可以通过非常丰富的内置 class 精细的控制页面样式。

## 运行

![chat-with-rose](https://cdn.huoyijie.cn/uploads/2023/07/chat-with-rose.png)

![chat-with-sse-react](https://cdn.huoyijie.cn/uploads/2023/07/chat-with-sse-react.png)