# 基于 HTTP SSE(Server Sent Event) 实现话题订阅机制

...

## Server Sent Event

...

## Topic & Subscribe

...

## SSE 话题订阅实例

文中所有代码已放到 [Github subscribe-topic-with-sse](https://github.com/huoyijie/tech-notes-code) 目录下。

*前置条件*

* 已安装 Go 1.20+
* 已安装 IDE （如 vscode）

创建 subscribe-topic-with-sse 项目

```bash
$ mkdir subscribe-topic-with-sse && cd subscribe-topic-with-sse
$ go mod init subscribe-topic-with-sse
```

新增 mydb.go 文件，添加测试用户和话题数据:

```go
// mydb.go
package main

// 所有注册用户
var myUsers = map[string]bool{
	"huoyijie": true,
	"jack":     true,
}

// 所有注册话题
var myTopics = map[string]bool{
	"chatgpt": true,
	"robot":   true,
}
```

新增 types.go 文件，添加如下几个类型:

```go
// types.go
package main

// 在某个话题上发布消息时会创建消息对象
type Message struct {
	Topic string `json:"topic" binding:"required"`
	Data  string `json:"data" binding:"required"`
}

// 向客户端发布消息的 Channel
type ClientChan chan Message

// 客户端，封装用户、订阅话题列表及客户端 Channel
type Client struct {
	User   string
	Topics []string
	C      ClientChan
}
```

新增 event.go，添加最核心的 Event 类型:

```go
// event.go
package main

import "log"

type Event struct {
	Message chan Message

	// New client connections
	NewClients chan Client

	// Closed client connections
	ClosedClients chan Client

	// All client connections
	Clients map[string]ClientChan

	// 所有用户的订阅话题
	// 格式: map["huoyijie"]map["chatgpt"]true
	Subscribs map[string]map[string]bool
}

// 创建 Event 对象
func NewEvent() (event *Event) {
	event = &Event{
		Message:       make(chan Message),
		NewClients:    make(chan Client),
		ClosedClients: make(chan Client),
		Clients:       make(map[string]ClientChan),
		Subscribs:     make(map[string]map[string]bool),
	}
  // 启动单独一个协程内更新 Clients/Subscribs，避免并发读写 map
	go event.listen()
	return
}

// 启动单独一个协程内更新 Clients/Subscribs，避免并发读写 map
func (event *Event) listen() {
	for {
		select {
		// Add new available client
		case client := <-event.NewClients:
			event.Clients[client.User] = client.C
			// 订阅话题列表
			topics := map[string]bool{}
			for _, topic := range client.Topics {
				if _, ok := myTopics[topic]; ok {
					topics[topic] = true
				}
			}
			event.Subscribs[client.User] = topics
			log.Printf("Client added. %d registered clients", len(event.Clients))

		// Remove closed client
		case client := <-event.ClosedClients:
			delete(event.Clients, client.User)
			close(client.C)
			log.Printf("Removed client. %d registered clients", len(event.Clients))

		// Forward message to client
		case eventMsg := <-event.Message:
			for user, topics := range event.Subscribs {
				if _, ok := topics[eventMsg.Topic]; ok {
					event.Clients[user] <- eventMsg
				}
			}
		}
	}
}
```

新增 subscribe.go 文件，添加客户端订阅话题 Handler:

```go
// subscribe.go
package main

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func subscribe(event *Event) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.Query("user")
		topics := c.Query("topics")

		// 判断用户是否存在，订阅话题是否为空
		if _, ok := myUsers[user]; !ok || len(topics) == 0 {
			c.AbortWithError(http.StatusBadRequest, errors.New("bad request"))
			return
		}

		// 创建新客户端
		client := Client{
			User:   user,
			Topics: strings.Split(topics, "|"),
			C:      make(ClientChan),
		}
		event.NewClients <- client

		// 如果连接端开，删除该客户端
		defer func() {
			event.ClosedClients <- client
		}()

		// Stream message to client
		c.Stream(func(w io.Writer) bool {
      // 发送该 client Channel 的消息通过 SSE Stream 发给浏览器
			if message, ok := <-client.C; ok {
				c.SSEvent("message", message)
				return true
			}
			return false
		})
	}
}
```

新增 publish.go 文件，添加发布话题消息 Handler:

```go
// publish.go
package main

import "github.com/gin-gonic/gin"

func publish(event *Event) gin.HandlerFunc {
	return func(c *gin.Context) {
		msg := Message{}
		if err := c.BindJSON(&msg); err != nil {
			return
		}

		// 发布话题
		event.Message <- msg
	}
}
```

新增 home.go，添加 web 页面:

```go
// home.go
package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func home() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.Query("user")
		topics := c.Query("topics")

		// 判断用户是否存在，订阅话题是否为空
		if _, ok := myUsers[user]; !ok || len(topics) == 0 {
			c.AbortWithError(http.StatusBadRequest, errors.New("bad request"))
			return
		}

		c.Writer.WriteString(fmt.Sprintf(`
			<!doctype html>
			<html lang="en">

			<head>
					<meta charset="UTF-8">
					<title>Server Sent Event</title>
			</head>

			<body>
			<div>Topic Messages:</div>
			<div class="event-data"></div>
			</body>

			<script src="https://code.jquery.com/jquery-1.11.1.js"></script>
			<script>
					// EventSource object of javascript listens the streaming events from our go server and prints the message.
					var stream = new EventSource("/subscribe?user=%s&topics=%s");
					stream.addEventListener("message", function(e){
							$('.event-data').append(e.data + "</br>")
					});
			</script>

			</html>`, user, topics))
	}
}
```

请求 `http://host:port/?user=huoyijie&topics=chatgpt|robot` 时会返回上述页面，注意上面 js 代码部分:

```javascript
var stream = new EventSource("/subscribe?user=%s&topics=%s");
stream.addEventListener("message", function(e){
    $('.event-data').append(e.data + "</br>")
});
```

页面打开后，会马上请求 `/subscribe?user=huoyijie&topics=chatgpt|robot` 实现用户 `huoyijie` 对话题列表 `chatgpt|robot` 的订阅。然后当某个话题有新消息时，就会自动通过 SSE Stream 发送消息到 web 页面。

最后添加 main.go 文件，添加 Handler 注册代码:

```go
// main.go
package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	event := NewEvent()

	r := gin.Default()
	r.GET("/", home())

	// 订阅话题消息
	r.GET("subscribe", subscribe(event))

	// 发布话题消息
	r.POST("publish", publish(event))
	r.Run(":8000")
}
```

接下来安装依赖和启动服务器:

```bash
# 下载安装依赖
$ go mod tidy

# 运行应用
$ go run .
```

打开一个新浏览器 Tab 页面，访问 URL: `http://localhost:8000/?user=huoyijie&topics=chatgpt`

![Client 1](https://cdn.huoyijie.cn/uploads/2023/07/sse-client1.png)

打开另一个新 Tab 页面，访问 URL: `http://localhost:8000/?user=jack&topics=chatgpt|robot`

![Client 2](https://cdn.huoyijie.cn/uploads/2023/07/sse-client2.png)

现在有 2 个用户分别订阅了话题，现在我们来测试下发布话题消息:

```bash
# 向 robot 话题发布消息
$ curl -d '{"topic": "robot", "data": "hi robot"}' http://localhost:8000/publish
```

只有用户 `jack` 收到了 `hi robot` 消息，因为只有 `jack` 订阅了 `robot`。

![received robot](https://cdn.huoyijie.cn/uploads/2023/07/sse-client-1-robot.png)

![not received robot](https://cdn.huoyijie.cn/uploads/2023/07/sse-client1.png)

```bash
# 向 chatgpt 话题发布消息
curl -d '{"topic": "chatgpt", "data": "hi chatgpt"}' http://localhost:8000/publish
```

2个用户都收到了 `hi chatgpt` 消息，因为他们都订阅了 `chatgpt`。

![received chatgpt](https://cdn.huoyijie.cn/uploads/2023/07/sse-client-chatgpt-1.png)

![received chatgpt](https://cdn.huoyijie.cn/uploads/2023/07/sse-client-chatgpt-2.png)