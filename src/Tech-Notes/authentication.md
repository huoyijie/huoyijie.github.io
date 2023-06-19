# Authentication

## 用户身份识别

除非网站上所有的数据和功能都是公开的，否则你一定会需要用户、角色和权限系统。不同的用户会有不同的权限，你需要识别出请求来自哪个用户。

想对用户请求进行识别，就需要前端在向服务器发送请求时携带的身份标识信息，一般均采用 Token 机制来实现，通过 Cookie/Header 发送。Token 可以是通过对称密钥加密的，也可以是采用私钥签名的。主要目的是阻止第三方伪造 Token 以其他用户身份访问网站。

* 对称加密

```
Token = 对称加密([用户ID,时间戳...], 密钥)
```

密钥是保密的，只存在于服务器上，Token 的生成与解析由服务器负责，客户端也无法解密 Token。Token 中通常只有 userId 和 timestamp。

* 非对称加密 (JWT)

```
Identify = base64.encode([用户ID,时间戳...])
Token = Identify + 签名(哈希(Identify), 私钥)
```

私钥是保密的，只存在于服务器上，Token 的生成由服务器负责，通常由原始用户身份信息 Identify 和签名两部分组成，用户身份 JSON 字符串，通过 base64 算法进行编码作为 Identify，由于 Identify 长度不固定，可用与客户端约定好的哈希算法进行运算，再用私钥进行签名。

这种方法生成的 Token 并不是加密的，客户端甚至第三方拿到 Token 后都可以解析出其中的用户身份信息。

客户端持有相对应的公钥，可对 Token 进行签名验证，验证此 Token 是否是具有私钥的服务器生成的，并可以从 Token 中解析出原始用户信息。Token 可携带很多用户信息，甚至可包含用户角色权限等信息。

```
[Identify, signature] = split(Token)
签名验证(哈希(Identify), signature, 公钥)
[用户ID,时间戳...] = base64.decode(Identify)
```

对称加解密、非对称加解密、签名、哈希等算法都属于密码学范畴，之前有介绍过 Golang 语言的各种密码学方法。[密码学方法](https://huoyijie.cn/gitbooks/writing-a-cloud-encrypting-file-system-with-golang-and-FUSE/latest/crypto.html)

* 随机生成 UUID 作为 Token

还有一种方法是随机生成 UUID 作为 Token，但是为了防止第三方伪造，需要服务器端保存所有有效的 Token 集合。服务器再接收到请求中携带的 Token 后需要进行有效性验证。

## 用户认证并下发 Token

添加用户时，密码上传到服务端后，确保所有中间服务器不记录、不打印日志，并经过 bcrypt 哈希后存储于数据库中。服务器没有任何理由保存用户原始密码。

用户认证的原理在于要求你输入只有你自己知道的秘密信息（如用户名、密码），并由服务器进行校验。登录表单提交时，不需要提交原始密码，而只提交密码哈希，避免用户密码泄漏。

服务器收到登录请求后，可与数据库中存储的用户名、密码哈希进行比对，以确定认证是否成功。认证成功后，可下发 Token。

## 根据 Token 自动对 API 请求进行认证

客户端收到服务器下发的 Token 后，可写入存储中，如写入 Cookie 或者 localStorage。后续的 API 请求可携带该 Token。服务器可以通过请求拦截器实现自动认证，拦截器在具体的请求前或后自动执行，可在请求被处理前，解析 Token 并校验有校性，然后获取登录用户信息并写入请求上下文中。

## 用户认证实例

*前置条件*
* 已安装 Go 1.20+
* IDE （如 vscode）

```bash
$ go version
go version go1.20.1 linux/amd64

$ mkdir user-auth && cd user-auth
$ go mod init user-auth
$ touch main.go
```

编辑 main.go 文件

```go
package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.Writer.WriteString("Hello, world!")
	})
	r.Run("0.0.0.0:8080")
}
```

本实例采用 web 框架 Gin，运行 `go mod tity` 自动下载安装依赖
```bash
go mod tidy
```

执行 main.go，并通过浏览器访问 `http://localhost:8080`，浏览器会输出 `Hello, world!`。

为了更容易用户认证部分代码，本实例进行了简化处理，不依赖数据库，注册用户数据放入内存 map 中。注意：本实例不支持并发注册和登录，并发的情况下会出错。

编辑 main.go 添加如下代码:

```go
// 统一返回包装类型
type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// 注册表单
type SignupForm struct {
	Username string `json:"username" binding:"required,alphanum,max=40"`
	Password string `json:"password" binding:"required,min=8,max=40"`
}

// 用户模型
type User struct {
	Username, PasswordHash string
}

// 模拟数据库存储，读写 map 未加锁，不支持并发注册登录
var users = map[string]User{}
```

然后修改 main 方法，添加如下代码:

```go
func main() {
  // ...
  r.POST("signup", func(c *gin.Context) {
    form := &SignupForm{}
    if err := c.BindJSON(form); err != nil {
      return
    }

    // check username unique
		if _, found := users[form.Username]; found {
			c.JSON(http.StatusOK, Result{
				Code:    -10000,
				Message: "用户已存在",
			})
			return
		}

    // calc password bcrypt hash bytes
    passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(form.Password), 14)
    if err != nil {
      c.AbortWithError(http.StatusInternalServerError, err)
      return
    }

    // base64 encode
    passwordHash := base64.StdEncoding.EncodeToString(passwordHashBytes)

    user := User{
      form.Username,
      passwordHash,
    }
    // 日志仅供调试
    fmt.Println("user:", user)

    // write new user to storage
    users[user.Username] = user
    // 日志仅供调试
    fmt.Println("users:", users)

    c.JSON(http.StatusOK, Result{
      Data: form.Username,
    })
  })
  // ...
}
```

运行 `go mod tity` 自动下载安装依赖，然后通过命令 `curl -d '{"username":"huoyijie","password":"mypassword"}' http://localhost:8080/signup` 发送注册请求。

注册成功后会返回用户名。可以看到代码中先对密码进行了 bcrypt 哈希，然后进行 base64 编码才写入存储。换句话说，服务器不会保存用户原始密码，也没有任何办法可以通过密码哈希值逆向得到。

新增 token.go 文件，存放生成与解析 Token 相关方法。

```go
// token.go
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
	"time"
)

func GetSecretKey() *[32]byte {
	key, err := hex.DecodeString("3e367a60ddc0699ea2f486717d5dcd174c4dee0bcf1855065ab74c348e550b78")
	if err != nil {
		log.Fatal(err)
	}
	return (*[32]byte)(key)
}

func NewGCM(key *[32]byte) (gcm cipher.AEAD, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return
	}
	gcm, err = cipher.NewGCM(block)
	return
}

func randNonce(nonceSize int) []byte {
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}
	return nonce
}

func Encrypt(plaintext []byte, gcm cipher.AEAD) []byte {
	// 随机生成字节 slice，使得每次的加密结果具有随机性
	nonce := randNonce(gcm.NonceSize())
	// Seal 方法第一个参数 nonce，会把 nonce 本身加入到加密结果
	return gcm.Seal(nonce, nonce, plaintext, nil)
}

func Decrypt(ciphertext []byte, gcm cipher.AEAD) ([]byte, error) {
	// 首先得到加密时使用的 nonce
	nonce := ciphertext[:gcm.NonceSize()]
	// 传入 nonce 并进行数据解密
	return gcm.Open(nil, nonce, ciphertext[gcm.NonceSize():], nil)
}

// 应该用 User ID 生成 Token
func GenerateToken(username string) (token string, err error) {
	gcm, err := NewGCM(GetSecretKey())
	if err != nil {
		return
	}

	bytes := make([]byte, len(username)+8)
	binary.BigEndian.PutUint64(bytes, uint64(time.Now().Unix()))
	copy(bytes[8:], []byte(username))
	token = base64.StdEncoding.EncodeToString(Encrypt(bytes, gcm))
	return
}

// 解析 Token
func ParseToken(token string) (username string, expired bool, err error) {
	gcm, err := NewGCM(GetSecretKey())
	if err != nil {
		return
	}

	tokenBytes, _ := base64.StdEncoding.DecodeString(token)
	bytes, err := Decrypt(tokenBytes, gcm)
	if err != nil {
		return
	}

	genTime := binary.BigEndian.Uint64(bytes)
	expired = time.Since(time.Unix(int64(genTime), 0)) > 30*24*time.Hour
	username = string(bytes[8:])
	return
}
```

编辑 main.go 文件，增加登录方法

```go
// main.go
func main() {
  // ...
  r.POST("signin", func(c *gin.Context) {
		form := &SigninForm{}
		if err := c.BindJSON(form); err != nil {
			return
		}

		decode := func(passwordHash string) (bytes []byte) {
			bytes, _ = base64.StdEncoding.DecodeString(passwordHash)
			return
		}

		if user, found := users[form.Username]; !found || bcrypt.CompareHashAndPassword(decode(user.PasswordHash), []byte(form.Password)) != nil {
			c.JSON(http.StatusOK, Result{
				Code:    -10001,
				Message: "用户或密码错误",
			})
			return
		} else if token, err := GenerateToken(user.Username); err != nil {
			c.JSON(http.StatusOK, Result{
				Code:    -10002,
				Message: "生成 Token 失败",
			})
			return
		} else {
			c.JSON(http.StatusOK, Result{
				Data: token,
			})
		}
	})
  // ...
}
```

运行服务器

```bash
$ cd user-auth
$ go run .
```

发送注册请求

```bash
$ curl -d '{"username":"huoyijie","password":"mypassword"}' http://localhost:8080/signup
{"code":0,"data":"huoyijie"}
```

再次发送同样的注册请求

```bash
$ curl -d '{"username":"huoyijie","password":"mypassword"}' http://localhost:8080/signup
{"code":-10000,"message":"用户已存在"}
```

发送错误的登录请求

```bash
$ curl -d '{"username":"notexist","password":"mypassword"}' http://localhost:8080/signin
{"code":-10001,"message":"用户或密码错误"}
```

发送正确的登录请求

```bash
$ curl -d '{"username":"huoyijie","password":"mypassword"}' http://localhost:8080/signin
{"code":0,"data":"IUqArBlhegws+ojRMZS/SD+ZKWnm6dNcWgHlFfzyFunkly2/jJLq90WCb/M="}
```

返回的 data 字段就是新生成的 Token，客户端可以写入存储中（如 Cookie 或 localStorage）。

以上所有代码已放到 [Github](https://github.com/huoyijie/tech-notes-code) user-auth 目录下。