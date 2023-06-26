# 基于 OAuth2 Password Grant 实现用户认证

用户认证和授权可以保护服务器接口，保障服务器和用户数据安全，是非常重要的安全基础。但是开发一个安全的用户认证授权流程不是件容易的事情，如果考虑不周或者代码实现有漏洞，很容易被黑客针对攻击。很多大的互联网平台(如 Google、微信)等，在用户认证授权方面积累了非常多的最佳实践经验，他们纷纷对第三方提供了身份认证服务，一般都遵循 OAuth 协议。选择接入这些大平台的身份认证服务，可以简化安全方面的工作，把主要精力投入到业务开发工作中。

## OAuth2

[OAuth 2.0](https://oauth.net/2/) 是第三方认证授权的工业级标准协议，旨在简化 Web 应用、桌面应用、移动应用和其他智能设备的特定的用户认证授权开发流程。

```
+--------+                               +---------------+
|        |--(A)- Authorization Request ->|   Resource    |
|        |                               |     Owner     |
|        |<-(B)-- Authorization Grant ---|               |
|        |                               +---------------+
|        |
|        |                               +---------------+
|        |--(C)-- Authorization Grant -->| Authorization |
| Client |                               |     Server    |
|        |<-(D)----- Access Token -------|               |
|        |                               +---------------+
|        |
|        |                               +---------------+
|        |--(E)----- Access Token ------>|    Resource   |
|        |                               |     Server    |
|        |<-(F)--- Protected Resource ---|               |
+--------+                               +---------------+

                Figure 1: Abstract Protocol Flow
```

OAuth 协议定义了4种角色:

* Resource Owner

资源拥有者，终端用户

* Resource server

平台资源服务器，接受携带 Access Token 的请求并返回受保护的资源数据(如: 获取用户信息、订单信息等)

* Client

准备接入平台身份认证服务的第三方应用程序，获得终端用户授权后，可以通过 Access Token 向资源服务器发送资源请求

* Authorization server

平台身份认证服务提供的授权服务器，在认证用户和获取用户授权后，为 Client 发放 Access Token

OAuth 协议主要流程:

Client 应用中登录界面显示其他平台(如: Google)登录选项，用户点击跳转至 Google 登录页面，输入 Google 用户名密码进行认证，通过后会打开授权页面，显示 Client 应用申请访问用户的某些资源(如: 用户资料)，用户点击授权后会返回到 Client 应用，Client 在得到授权后可从 Authorization server 获取 Access Token，最后通过 Access Token 可以从 Resource server 获取用户资料。此时，Client 可为用户创建账户并登录。

## OAuth2 Password Grant

本文主旨不在于介绍上述 OAuth 的主要使用场景，而是想介绍一个 OAuth 协议不太推荐的遗存的认证授权方式，即通过用户名密码直接获取 Access Token。主要流程是 Client 应用提供登录页面，收集用户用户名、密码，然后提交至 authorization server，后者完成认证后直接返回 Access Token 给 Client 应用程序。**Client 应用可以接触到用户名、密码**，是这种方式不被推荐的原因，除非 Client 被高度信任(如: Client 为 iOS 操作系统)，否则不应该使用这种方式。

有一种情况比较适合 OAuth2 Password Grant，就是 Client、Resource server、Authorization server 都属于同一个应用本身，此时 Client 是可信任的。

```
+----------+
| Resource |
|  Owner   |
|          |
+----------+
    v
    |    Resource Owner
    (A) Password Credentials
    |
    v
+---------+                                  +---------------+
|         |>--(B)---- Resource Owner ------->|               |
|         |         Password Credentials     | Authorization |
| Client  |                                  |     Server    |
|         |<--(C)---- Access Token ---------<|               |
|         |    (w/ Optional Refresh Token)   |               |
+---------+                                  +---------------+

      Figure 2: Resource Owner Password Credentials Flow
```

Client 在获取到 Access Token 后，可携带 Access Token 访问 Resource server。下面我们来看看如何基于 [go-oauth2 server](https://github.com/go-oauth2/oauth2) 和 [oauth2 client](https://pkg.go.dev/golang.org/x/oauth2) 实现这种用户认证方式。

## 实现用户认证

文中所有代码已放到 [Github user-auth-with-oauth2](https://github.com/huoyijie/tech-notes-code) 目录下。

*前置条件*
* 已安装 Go 1.20+
* 已安装 IDE （如 vscode）

创建 user-auth-with-oauth2 项目

```bash
$ mkdir user-auth-with-oauth2 && cd user-auth-with-oauth2
$ go mod init user-auth-with-oauth2
```

**实现 OAuth2 Authorization Server**

创建 oauth2.go 文件，并添加下面代码

```go
package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// 用户模型
type User struct {
	Username, PasswordHash string
}

// 模拟数据库存储，真实应用需写入数据库表中
var users = map[string]User{
	"huoyijie": {
		Username: "huoyijie",
		// 原始密码: mypassword
		PasswordHash: "JDJhJDE0JElHWVpnTzdtd1pZbEVTQnAyY1VhTk9CVEJkcUcwV2xyMFZaWElKZ25EZlNjM0lqZHllc2E2",
	},
}

func decode(passwordHash string) (bytes []byte) {
	bytes, _ = base64.StdEncoding.DecodeString(passwordHash)
	return
}

// 获取密钥
func getSecretKey() []byte {
	key, err := hex.DecodeString("3e367a60ddc0699ea2f486717d5dcd174c4dee0bcf1855065ab74c348e550b78" /* Load key from somewhere, for example an environment variable */)
	if err != nil {
		log.Fatal(err)
	}
	return key
}

// 启动 oauth2 server
func runOAuth2(r *gin.Engine) {
	manager := manage.NewDefaultManager()
	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	// generate jwt access token
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("jwt", getSecretKey(), jwt.SigningMethodHS512))

	// client memory store
	clientStore := store.NewClientStore()
	clientStore.Set("100000", &models.Client{
		ID:     "100000",
		Secret: "575f508960a9415a97f05a070a86165b",
	})
	manager.MapClientStorage(clientStore)

	// 设置 oauth2 server
	srv := server.NewDefaultServer(manager)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})
	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	// 用户认证
	srv.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		if clientID == "100000" {
			// 验证用户存在且密码哈希比对成功
			if user, found := users[username]; found && bcrypt.CompareHashAndPassword(decode(user.PasswordHash), []byte(password)) == nil {
				userID = username
				return
			}
		}
		err = errors.ErrUnauthorizedClient
		return
	})

	// 返回 token
	r.POST("/oauth/token", func(c *gin.Context) {
		srv.HandleTokenRequest(c.Writer, c.Request)
	})
}
```

主要看负责启动 oauth2 server 的 runOAuth2 方法，代码首先创建 manager，设置基于内存的 Token Storage，设置生成 JWT 格式的 Access Token，然后设置基于内存的 Client 信息存储，预先配置了一个 ID 为 100000 的 Client，只有这里配置过的 Client 才可以通过 oauth2 server 获取 Access Token。

基于内存存储 Token 和 Client 信息容易设置，非常方便用来开发和测试，但是重启进程数据就没了，实践中可以采用插件替换为基于数据库或 Redis 的方案，具体参见[go-oauth2](https://github.com/go-oauth2/oauth2)。

设置完 manager 后，基于 manager 创建 srv，设置错误处理回调函数，设置用户名密码认证回调函数。用户认证回调函数首先判断 clientID 等于 100000，然后比对用户名和密码哈希，认证成功返回 userID，认证失败返回 error。

最后配置路由 `/oauth/token`，完全接管生成 Token 的请求。

运行 `go mod tidy` 安装依赖:

```bash
$ go mod tidy
```

**实现 OAuth2 Client**

创建 app.go，并添加如下代码:

```go
package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// 统一返回包装类型
type Result struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// 登录表单
type SigninForm struct {
	Username string `json:"username" binding:"required,alphanum,max=40"`
	Password string `json:"password" binding:"required,min=8,max=40"`
}

const (
	authServerURL = "http://127.0.0.1:8080"
)

// 通过 oauth2 server 预先配置的 Client 访问 Authorization server
var config = oauth2.Config{
	ClientID:     "100000",
	ClientSecret: "575f508960a9415a97f05a070a86165b",
	Endpoint: oauth2.Endpoint{
		TokenURL: authServerURL + "/oauth/token",
	},
}

// 启动 App
func runApp(r *gin.Engine) {
	// 登录获取 token
	r.POST("/signin", func(c *gin.Context) {
		form := &SigninForm{}
		if err := c.BindJSON(form); err != nil {
			return
		}

		// 通过提供用户名、密码获取 token
		token, err := config.PasswordCredentialsToken(context.Background(), form.Username, form.Password)
		if err != nil {
			if e, ok := err.(*oauth2.RetrieveError); ok {
				c.JSON(http.StatusOK, Result{
					Code:    e.ErrorCode,
					Message: e.ErrorDescription,
				})
				return
			}
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, Result{Data: token})
	})
}
```

这段代码定义了 config 变量，通过 oauth2 server 预先配置的 Client 访问 Authorization server，然后是最重要的 runApp 函数，定义了登录 Signin 接口，调用 [oauth2 client](https://pkg.go.dev/golang.org/x/oauth2) 库提供的 PasswordCredentialsToken 方法，使用用户名密码交换 Token。PasswordCredentialsToken 方法会向 oauth2 Authorization server 发送类似如下请求:

```
POST /oauth/token HTTP/1.1
 
grant_type=password
&username=huoyijie
&password=mypassword
&client_id=100000
&client_secret=575f508960a9415a97f05a070a86165b
```

oauth2 Authorization server 会验证 client_id & client_secret 是否合法，调用下面预先配置好的回调函数进行用户名密码认证。

```go
// oauth2.go
func runOAuth2() {
  // ...
	// 用户认证
	srv.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		if clientID == "100000" {
			// 验证用户存在且密码哈希比对成功
			if user, found := users[username]; found && bcrypt.CompareHashAndPassword(decode(user.PasswordHash), []byte(password)) == nil {
				userID = username
				return
			}
		}
		err = errors.ErrUnauthorizedClient
		return
	})
	// ...
}
```

认证成功后，会返回 Token 给 oauth client。

```bash
# 安装依赖
$ go mod tidy
```

```
e.ErrorCode undefined (type *"golang.org/x/oauth2".RetrieveError has no field or method ErrorCode
```

如果代码报上述错误，手动编辑 go.mod 调整 oauth2 client 库版本为 v0.9.0

```
require (
-	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
+	golang.org/x/oauth2 v0.9.0
)
```

然后再次执行 `go mod tidy`，报错消失。下面我们来测试一下: