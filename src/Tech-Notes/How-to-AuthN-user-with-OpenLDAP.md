# 基于 OpenLDAP、React、Next.js、MUI 实现用户认证

## Next.js 介绍

[Next.js](https://nextjs.org/) 是一个用于构建现代 React 应用程序的流行开源框架。它提供了一组强大的工具和约定，使得 React 应用的开发变得更加简单和高效。

## MUI 介绍

[Material-UI](https://mui.com/) 是一个基于 React 的流行开源组件库，用于构建符合 [Google Material Design](https://m3.material.io/) 规范的用户界面。MUI 提供了一套丰富而灵活的 React 组件，使开发人员能够轻松创建具有现代外观和交互的 Web 应用程序。这个库提供了各种组件，包括按钮、表单、导航、对话框等，这些组件遵循 Material Design 的设计原则。

## LDAP/OpenLDAP 介绍

LDAP（轻量目录访问协议，Lightweight Directory Access Protocol）是一种用于访问和维护分布式目录信息服务的协议。LDAP 最初设计用于提供一种轻量级的、低开销的访问目录信息的方法。目录服务是一种特殊的数据库，用于存储和检索组织中的信息，例如用户、计算机、组织单元等。LDAP 在企业网络中广泛应用，可用来实现统一身份认证，让用户在多个服务中使用相同的身份和密码，也常常配合企业内部实现单点登录(SSO)。

[OpenLDAP](https://www.openldap.org/) 是一个开源的 LDAP 实现，而 Microsoft Active Directory（AD）是一个广泛使用的 LDAP 目录服务的实例，用于 Windows 环境中的身份认证和目录服务。

录入用户信息时，还是会写入数据库中，同时导出一份数据写入 openLDAP 中，并保持数据随时同步。后面系统进行用户认证时，直接与 openLDAP 交互即可。

## 安装配置 OpenLDAP

```bash
# Ubuntu 22.04
$ sudo apt install slapd
$ sudp apt install ldap-utils
```

进行初始设置，注意会提示设置管理员密码，后面管理和程序连接时会用到。

```bash
$ sudo dpkg-reconfigure slapd
```

设定域

![ldap-1](https://cdn.huoyijie.cn/uploads/2023/12/ldap-1.png)

设定组织名称

![ldap-2](https://cdn.huoyijie.cn/uploads/2023/12/ldap-2.png)

设定管理员密码

![ldap-3](https://cdn.huoyijie.cn/uploads/2023/12/ldap-3.png)

检查 openldap 已正常运行
```bash
$ ldapsearch -x -LLL -D cn=admin,dc=huoyijie,dc=cn -W -b dc=huoyijie,dc=cn
Enter LDAP Password: ******
dn: dc=huoyijie,dc=cn
objectClass: top
objectClass: dcObject
objectClass: organization
o: huoyijie
dc: huoyijie
```

添加测试用户，创建一个 user.ldif 文件，编辑内容:
```bash
# 添加 users 组
dn: ou=users,dc=huoyijie,dc=cn
objectClass: organizationalUnit
objectClass: top
ou: users

# 添加用户 huoyijie
dn: uid=huoyijie,ou=users,dc=huoyijie,dc=cn
objectClass: top
objectClass: person
objectClass: organizationalPerson
objectClass: inetOrgPerson
uid: huoyijie
cn: Yijie Huo
sn: Huo
userPassword: $2a$10$Mx.Uz6ybst6d13jzaQX6XOMOOq4BsUpHSxB5PzNfaRj65jxJNy4dO
mail: huoyijie@huoyijie.cn
```

运行 ldapadd 命令
```bash
$ ldapadd -x -D "cn=admin,dc=huoyijie,dc=cn" -W -f user.ldif
Enter LDAP Password: 
adding new entry "ou=users,dc=huoyijie,dc=cn"

adding new entry "uid=huoyijie,ou=users,dc=huoyijie,dc=cn"
```

* -x: 使用简单身份验证。
* -D: 指定用于绑定到 LDAP 服务器的管理员 DN。
* -W: 提示输入管理员密码
* -f: 指定包含 LDIF 数据的文件

在执行该命令后，系统将提示输入管理员密码，然后将 LDIF 文件中的用户数据添加到 LDAP 中。

查询验证新用户
```bash
$ ldapsearch -x -D "cn=admin,dc=huoyijie,dc=cn" -W -b "ou=users,dc=huoyijie,dc=cn" "(uid=huoyijie)"
Enter LDAP Password: 
# extended LDIF
#
# LDAPv3
# base <ou=users,dc=huoyijie,dc=cn> with scope subtree
# filter: (uid=huoyijie)
# requesting: ALL
#

# huoyijie, users, huoyijie.cn
dn: uid=huoyijie,ou=users,dc=huoyijie,dc=cn
objectClass: top
objectClass: person
objectClass: organizationalPerson
objectClass: inetOrgPerson
uid: huoyijie
cn: Yijie Huo
sn: Huo
userPassword:: JDJhJDEwJE14LlV6Nnlic3Q2ZDEzanphUVg2WE9NT09xNEJzVXBIU3hCNVB6TmZ
 hUmo2NWp4Sk55NGRP
mail: huoyijie@huoyijie.cn

# search result
search: 2
result: 0 Success

# numResponses: 2
# numEntries: 1
```

* -x: 使用简单身份验证。
* -D: 指定用于绑定到 LDAP 服务器的管理员 DN。
* -W: 提示输入管理员密码。
* -b: 指定搜索的起始点（Base DN）
* "(uid=huoyijie)": 指定 LDAP 搜索过滤器，这里是根据 uid 进行搜索