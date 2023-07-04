# 基于 TOTP 实现多因素(Multi-factor)认证



## 多因素(Multi-factor)认证

[多重要素验证](https://zh.wikipedia.org/zh-cn/%E5%A4%9A%E9%87%8D%E8%A6%81%E7%B4%A0%E9%A9%97%E8%AD%89)，或多因子认证、多因素验证、多因素认证，是一种资源访问控制的方法，用户要通过两种以上的认证机制之后，才能得到授权，使用资源。这种认证方式可以提高安全性。例如，用户要输入PIN码，插入银行卡，最后再经指纹比对，通过这三种认证方式，才能获得授权。

## OTP/HOTP/TOTP

- OTP

[OTP](https://zh.wikipedia.org/zh-cn/%E4%B8%80%E6%AC%A1%E6%80%A7%E5%AF%86%E7%A2%BC)（One-time Password），又称动态密码或单次有效密码，是指计算机系统或其他数字设备上只能使用一次的密码，有效期为只有一次登录会话或交易。

- HOTP

[HOTP](https://zh.wikipedia.org/zh-cn/%E5%9F%BA%E4%BA%8E%E6%95%A3%E5%88%97%E6%B6%88%E6%81%AF%E9%AA%8C%E8%AF%81%E7%A0%81%E7%9A%84%E4%B8%80%E6%AC%A1%E6%80%A7%E5%AF%86%E7%A0%81%E7%AE%97%E6%B3%95)（HMAC-based One-time Password）是一种基于散列消息验证码（HMAC）的一次性密码（OTP）算法。

- TOTP

[TOTP](https://zh.wikipedia.org/zh-cn/%E5%9F%BA%E4%BA%8E%E6%97%B6%E9%97%B4%E7%9A%84%E4%B8%80%E6%AC%A1%E6%80%A7%E5%AF%86%E7%A0%81%E7%AE%97%E6%B3%95)（Time-based One-Time Password）是一种根据预共享的密钥与当前时间计算一次性密码的算法。它已被互联网工程任务组接纳为RFC 6238标准，成为主动开放认证（OATH）的基石，并被用于众多多重要素验证系统当中。

TOTP基于HOTP实现，它结合一个私钥与当前时间戳，使用一个密码散列函数来生成一次性密码。由于网络延迟与时钟不同步可能导致密码接收者不得不尝试多次遇到正确的时间来进行身份验证，时间戳通常以30秒为间隔，从而避免反复尝试。

## OTP Golang 实现

## 多重要素(Multi-factor)认证实例