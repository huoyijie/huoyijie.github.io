# 密码学方法

## 对称加密

对称加密需要预先生成 `key`，`key` 的字节长度会确定具体的加密方法，对应关系如下:

| key 字节数 | 加密方法 |
|----|----|
| 16 | AES-128 |
| 24 | AES-192 |
| 32 | AES-256 |

一般来说，选择更长的 `key`，运算会慢一些，安全性会高一些。`NewEncryptionKey` 方法可随机生成 `32` 字节长度的 `key`。

```go
func NewEncryptionKey() *[32]byte {
	key := [32]byte{}
	if _, err := io.ReadFull(rand.Reader, key[:]); err != nil {
		panic(err)
	}
	return &key
}
```

选择 `AES-256-GCM` 算法进行加解密。

```go
func newGCM(key *[32]byte) (gcm cipher.AEAD, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return
	}
	gcm, err = cipher.NewGCM(block)
	return
}
```

下面是 `Encrypt` 加密方法。

```go
func Encrypt(plaintext []byte, gcm cipher.AEAD) []byte {
	// 随机生成字节 slice，使得每次的加密结果具有随机性
	nonce := randNonce(gcm.NonceSize())
	// Seal 方法第一个参数 nonce，会把 nonce 本身加入到加密结果
	return gcm.Seal(nonce, nonce, plaintext, nil)
}
```

最终加密结果 `ciphertext` 由 `nonce|ciphertext|tag` 三部分连接组成，不包括 `|`。下面是生成随机 `nonce` 的方法。

```go
func randNonce(nonceSize int) []byte {
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}
	return nonce
}
```

下面是 `Decrypt` 解密方法，`ciphertext` 是由 `Encrypt` 方法加密得到的加密结果，由 `nonce|ciphertext|tag` 三部分连接组成，不包括 `|`。

```go
func Decrypt(ciphertext []byte, gcm cipher.AEAD) ([]byte, error) {
	// 首先得到加密时使用的 nonce
	nonce := ciphertext[:gcm.NonceSize()]
	// 传入 nonce 并进行数据解密
	return gcm.Open(nil, nonce, ciphertext[gcm.NonceSize():], nil)
}
```

下面的代码使用上述加解密方法进行测试。连续两次对 `Hello, world!` 字符串进行加密，最后才进行解密。

```go
func main() {
	key := NewEncryptionKey()
	gcm, err := newGCM(key)
	if err != nil {
		fmt.Println(err)
		return
	}

	ciphertext := Encrypt([]byte("Hello, world!"), gcm)
	fmt.Println(hex.EncodeToString(ciphertext))

	ciphertext = Encrypt([]byte("Hello, world!"), gcm)
	fmt.Println(hex.EncodeToString(ciphertext))

	plaintext, err := Decrypt(ciphertext, gcm)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(plaintext))
}
```

程序输出:

```
3a1e10d33f73fcbce973a2765468286442142c97cb85c67d09c170988de0185235aad490b78c3a60dd
36fc70905be24708420a25b81ae8666b9da37129447591826ccb4d3b426d04a117138c281ed04f362b
Hello, world!
```

可以看到，连续两次都是对 `Hello, world!` 字符串进行加密，结果是不同的，这符合预期，确保了每次加密结果的随机性。而最后也成功解密了。

点击 [The Go Playground](https://go.dev/play/p/I8QC2cQ2TZj) 去测试代码。

下面是完整程序清单:

```go
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func randNonce(nonceSize int) []byte {
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}
	return nonce
}

func newGCM(key *[32]byte) (gcm cipher.AEAD, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return
	}
	gcm, err = cipher.NewGCM(block)
	return
}

func NewEncryptionKey() *[32]byte {
	key := [32]byte{}
	if _, err := io.ReadFull(rand.Reader, key[:]); err != nil {
		panic(err)
	}
	return &key
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

func main() {
	key := NewEncryptionKey()
	gcm, err := newGCM(key)
	if err != nil {
		fmt.Println(err)
		return
	}

	ciphertext := Encrypt([]byte("Hello, world!"), gcm)
	fmt.Println(hex.EncodeToString(ciphertext))

	ciphertext = Encrypt([]byte("Hello, world!"), gcm)
	fmt.Println(hex.EncodeToString(ciphertext))

	plaintext, err := Decrypt(ciphertext, gcm)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(plaintext))
}
```

## 非对称加密

`NewRsaKey` 方法可随机生成 `2048` 位长度 `RSA` 公/私密钥

```go
func NewRsaKey() (priv *rsa.PrivateKey, err error) {
	// 2048-bit
	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	return
}
```

`EncryptOAEP` 方法采用 `RSA-OAEP` 算法进行加密

```go
func EncryptOAEP(pub *rsa.PublicKey, plaintext []byte, label []byte) (ciphertext []byte, err error) {
	// label 可以确保一个场景的加密数据，不能够在另一个场景中解密，在一个场景下用相同的 label 进行加密和解密
	ciphertext, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, plaintext, label)
	return
}
```

`DecryptOAEP` 方法对采用 `RSA-OAEP` 算法加密的数据进行解密

```go
func DecryptOAEP(priv *rsa.PrivateKey, ciphertext []byte, label []byte) (plaintext []byte, err error) {
	// label 可以确保一个场景的加密数据，不能够在另一个场景中解密，在一个场景下用相同的 label 进行加密和解密
	plaintext, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, ciphertext, label)
	return
}
```

下面的代码使用上述加解密方法进行测试。连续两次对 `Hello, world!` 字符串进行加密，最后才进行解密。

```go
func main() {
	priv, err := NewRsaKey()
	if err != nil {
		fmt.Println(err)
		return
	}

	ciphertext, err := EncryptOAEP(&priv.PublicKey, []byte("Hello, world!"), []byte("test data 1"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hex.EncodeToString(ciphertext))

	ciphertext, err = EncryptOAEP(&priv.PublicKey, []byte("Hello, world!"), []byte("test data 2"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hex.EncodeToString(ciphertext))

	plaintext, err := DecryptOAEP(priv, ciphertext, []byte("test data 2"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(plaintext))
}
```

程序输出:

```
3ed39b1f8877acbb62ddd4374f1ae8a5e94f38e95f8e554a2718dd016780cebe4a7b70542b5a159864d2f3a7ddefe5f4b4f6f8c22e8feeb241c10d324529992f99f5d0829fe888dfcdbe0601bf55777afb20fa20d6666cffb29b3328df6d72d34d9a87d37e5d17c3b6f568e1c1f49f1e91679bda1c54caf25a51f562aa08097f6aab4d7be41e7881a4123d88af9e1df4d8751759105665c2b61fe1ad360a37be3a0f3b8ba269386bb6cf607f49bc4363468ae713c892cf3303fd64a3cf691b104b72f9015b7bc73a2a1f986aaeb2aa7d5bcba1262fc52310890142d388b48842a3ac5dcac120f415650da6a4f069dbdee48dcb4a03c41b75e066cbbd69c32627
304d56f58d2ea878dfc8d01844cad6cc7df406233903343763e80e3a0974ca7a70f1ae43d30be0175ff09323720b22286c6be77d880cda7bc99096b5f2dffa539ef8a5f42f37b1c03e4c54cbce04570f0ab6f98214b2b7a687496bb48302ba07f8cb013642cc41f2428fa01dfd692cfbac25e48529426be58748d715e82ad3129060a5c85adbadefa59d290e97b386e57382aba8add37a248f39a4394d0ca3244d33e08beae92b1dcb3f48108aecc38974087244f0e9c5223bbcd3c449863715cab2ebde417e6740282b9998b82ff0ac097e926b257485b9e533c579cc9dfa964bef881b59749ea478e634606a7fba549de4bde270df35c8f37c60bf04478b92
Hello, world!
```

可以看到，连续两次都是对 `Hello, world!` 字符串进行加密，结果是不同的，这符合预期，确保了每次加密结果的随机性。而最后也成功解密了。

点击 [The Go Playground](https://go.dev/play/p/GVbNAEeM5vj) 去测试代码。

下面是完整程序清单:

```go
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func NewRsaKey() (priv *rsa.PrivateKey, err error) {
	// 2048-bit
	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	return
}

func EncryptOAEP(pub *rsa.PublicKey, plaintext []byte, label []byte) (ciphertext []byte, err error) {
	// label 可以确保一个场景的加密数据，不能够在另一个场景中解密，在一个场景下用相同的 label 进行加密和解密
	ciphertext, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, plaintext, label)
	return
}

func DecryptOAEP(priv *rsa.PrivateKey, ciphertext []byte, label []byte) (plaintext []byte, err error) {
	// label 可以确保一个场景的加密数据，不能够在另一个场景中解密，在一个场景下用相同的 label 进行加密和解密
	plaintext, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, ciphertext, label)
	return
}

func main() {
	priv, err := NewRsaKey()
	if err != nil {
		fmt.Println(err)
		return
	}

	ciphertext, err := EncryptOAEP(&priv.PublicKey, []byte("Hello, world!"), []byte("test data 1"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hex.EncodeToString(ciphertext))

	ciphertext, err = EncryptOAEP(&priv.PublicKey, []byte("Hello, world!"), []byte("test data 2"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hex.EncodeToString(ciphertext))

	plaintext, err := DecryptOAEP(priv, ciphertext, []byte("test data 2"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(plaintext))
}
```

## 信息摘要

很多场景都需要信息摘要算法，比如数据签名，首先通过信息摘要算法 `sha256` 得到 `32` 字节长度字节 slice，然后在对该字节 slice 进行签名。

点击 [The Go Playground](https://go.dev/play/p/QcUaSSrr7o3) 去测试代码。

程序清单:

```go
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {
	data := []byte("Hello, world!")
	digest := sha256.Sum256(data)
	fmt.Println(hex.EncodeToString(digest[:]))
}
```

## 签名验证

`NewSigningKey` 是生成 `Ecdsa` 签名验证公/私钥方法

```go
func NewSigningKey() (*ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return key, err
}
```

`Sign` 是数据签名方法，首先通过 `sha256` 算法得到数据信息摘要，然后通过 `Ecdsa` 私钥完成了数据签名。

```go
func Sign(data []byte, privkey *ecdsa.PrivateKey) ([]byte, error) {
	// hash message
	digest := sha256.Sum256(data)

	// sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privkey, digest[:])
	if err != nil {
		return nil, err
	}

	// encode the signature {R, S}
	// big.Int.Bytes() will need padding in the case of leading zero bytes
	params := privkey.Curve.Params()
	curveOrderByteSize := params.P.BitLen() / 8
	rBytes, sBytes := r.Bytes(), s.Bytes()
	signature := make([]byte, curveOrderByteSize*2)
	copy(signature[curveOrderByteSize-len(rBytes):], rBytes)
	copy(signature[curveOrderByteSize*2-len(sBytes):], sBytes)

	return signature, nil
}
```

`Verify` 是签名验证方法，首先通过 `sha256` 算法得到数据信息摘要，然后对数据摘要和签名进行了校验。

```go
func Verify(data, signature []byte, pubkey *ecdsa.PublicKey) bool {
	// hash message
	digest := sha256.Sum256(data)

	curveOrderByteSize := pubkey.Curve.Params().P.BitLen() / 8

	r, s := new(big.Int), new(big.Int)
	r.SetBytes(signature[:curveOrderByteSize])
	s.SetBytes(signature[curveOrderByteSize:])

	return ecdsa.Verify(pubkey, digest[:], r, s)
}
```

下面是签名验证测试代码

```go
func main() {
	priv, err := NewSigningKey()
	if err != nil {
		fmt.Println(err)
	}

	signature, err := Sign([]byte("Hello, world!"), priv)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hex.EncodeToString(signature))

	if ok := Verify([]byte("Hello, world!"), signature, &priv.PublicKey); ok {
		fmt.Println("verified ok")
	}
}
```

程序输出:

```
eb5a28450f221e888237249e03b07acff0cbb524bc8f2082d5ef996bed4eb957cfcdb9ba1aa103954ac427468565cb8e904111a5fd1c589df94d9309b24e9e3f
verified ok
```

可以看到，首先对 `Hello, world!` 字符串进行了签名，然后又成功验证了签名。

点击 [The Go Playground](https://go.dev/play/p/3-KXceMpHNg) 去测试代码。

程序清单:

```go
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

func NewSigningKey() (*ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return key, err
}

func Sign(data []byte, privkey *ecdsa.PrivateKey) ([]byte, error) {
	// hash message
	digest := sha256.Sum256(data)

	// sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privkey, digest[:])
	if err != nil {
		return nil, err
	}

	// encode the signature {R, S}
	// big.Int.Bytes() will need padding in the case of leading zero bytes
	params := privkey.Curve.Params()
	curveOrderByteSize := params.P.BitLen() / 8
	rBytes, sBytes := r.Bytes(), s.Bytes()
	signature := make([]byte, curveOrderByteSize*2)
	copy(signature[curveOrderByteSize-len(rBytes):], rBytes)
	copy(signature[curveOrderByteSize*2-len(sBytes):], sBytes)

	return signature, nil
}

func Verify(data, signature []byte, pubkey *ecdsa.PublicKey) bool {
	// hash message
	digest := sha256.Sum256(data)

	curveOrderByteSize := pubkey.Curve.Params().P.BitLen() / 8

	r, s := new(big.Int), new(big.Int)
	r.SetBytes(signature[:curveOrderByteSize])
	s.SetBytes(signature[curveOrderByteSize:])

	return ecdsa.Verify(pubkey, digest[:], r, s)
}

func main() {
	priv, err := NewSigningKey()
	if err != nil {
		fmt.Println(err)
	}

	signature, err := Sign([]byte("Hello, world!"), priv)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hex.EncodeToString(signature))

	if ok := Verify([]byte("Hello, world!"), signature, &priv.PublicKey); ok {
		fmt.Println("verified ok")
	}
}
```

## 密码存储认证

点击 [The Go Playground](https://go.dev/play/p/UcNg7KyB8Xh) 去测试代码。

程序清单:

```go
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "mypassword"
	passHash := sha256.Sum256([]byte(password))
	passHashAndBcrypt, err := bcrypt.GenerateFromPassword(passHash[:], 14)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hex.EncodeToString(passHashAndBcrypt))

	if err := bcrypt.CompareHashAndPassword(passHashAndBcrypt, passHash[:]); err == nil {
		fmt.Println("verified ok!")
	}
}
```

上面的程序首先对原始密码 `mypassword` 进行了 `sha256` 哈希，得到的 `passHash` 可以用来对私钥进行加密。然后再对 `passHash` 调用 `bcrypt.GenerateFromPassword` 哈希，得到 `passHashAndBcrypt`。两次的哈希都是不可逆的，因此在密钥文件中存储 `passHashAndBcrypt` 是安全的，然后程序对 `passHashAndBcrypt` 和 `passHash` 调用 `bcrypt.CompareHashAndPassword` 完成了密码认证，如果该方法返回的 err 为 nil，代表密码认证通过。