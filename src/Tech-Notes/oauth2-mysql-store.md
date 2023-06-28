# 基于 mysql 实现 OAuth2 Token&Client 存储

```bash
$ mkdir auth2-mysql-store && cd auth2-mysql-store
$ go mod init auth2-mysql-store
```

```bash
$ go run .
panic: Error 1049 (42000): Unknown database 'oauth2-db'
```