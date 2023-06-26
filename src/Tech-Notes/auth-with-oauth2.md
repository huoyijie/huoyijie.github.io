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

Client 在获取到 Access Token 后，可携带 Access Token 访问 Resource server。下面我们来看看如何实现这种用户认证方式。

## 实现用户认证

