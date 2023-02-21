# EFS 加密文件系统

## 进入 EFS 网站

[进入 EFS 网站](https://huoyijie.cn/efs/)

## 系统介绍

[EFS (Encrypting File System)](https://huoyijie.cn/efs/) 是基于网络通信、加密、[FUSE](https://john-millikin.com/the-fuse-protocol) 等技术采用 `golang` 语言实现的网络加密文件系统，文件数据加密存储于远程服务器上面，一方面可以确保数据的隐私安全，只有数据所有者拥有密钥和密码，另一方面也可以用来对重要数据做远程备份。

EFS 由客户端、服务端软件组成，客户端程序支持 Windows/Mac/Linux 系统，运行在用户电脑上，服务端软件运行在远程服务器上。用户只需在电脑上安装好客户端，注册帐号并生成密钥，就可以正常使用 EFS 客户端了，运行 EFS 客户端后，电脑上相当于挂载了一个“远程”移动硬盘，使用起来非常简单。

EFS 客户端基于 FUSE (Filesystem in Userspace) 实现，在不同的平台上依赖不同的系统组件。文档的后面部分也会介绍如何安装系统依赖和客户端软件。

## 如何安装

用户需要在电脑上安装客户端软件 efscli 及其依赖组件。其中 efscli 是预编译的二进制可执行文件，下载后修改权限可直接运行。

### Windows

* [WinFsp](https://winfsp.dev/)

    点击 [WinFsp 进入下载页面](https://winfsp.dev/rel/)， 可以看到 Download WinFsp Installer 下载链接，点击下载，或者直接点击 [winfsp-1.12.22301.msi](https://github.com/winfsp/winfsp/releases/download/v1.12/winfsp-1.12.22301.msi) 下载。下载完成后，运行 WinFsp 安装器，按照提示进行安装。

* 下载 EFS 客户端软件 

    点击下载 [efscli-windows-amd64.exe](https://cdn.huoyijie.cn/efs/downloads/efscli-0.5.0-windows-amd64.exe)

### Mac

* [macFUSE](https://osxfuse.github.io/)

    点击 [macFUSE 进入下载页面](https://osxfuse.github.io/)，可以在左侧看到 macFUSE 4.4.1 下载链接，点击下载，或者直接点击 [macfuse-4.4.1.dmg](https://github.com/osxfuse/osxfuse/releases/download/macfuse-4.4.1/macfuse-4.4.1.dmg) 下载。下载完成后，运行 macfuse-4.4.1.dmg 安装器，按照提示进行安装。

* 下载 EFS 客户端软件

    点击下载 [efscli-darwin-amd64](https://cdn.huoyijie.cn/efs/downloads/efscli-0.5.0-darwin-amd64)

### Linux

常见的 Linux 发行版在内核中都内置 FUSE 支持，所以不需要安装额外的依赖组件。

* 下载 EFS 客户端软件

    点击下载 [efscli-linux-amd64](https://cdn.huoyijie.cn/efs/downloads/efscli-0.5.0-linux-amd64)

## 如何使用

* 点击 [注册帐号](https://huoyijie.cn/efs/signup)

    注册暂时不需设置登录密码，通过邮件接收验证码后自动登录，后续如果登录状态失效后，可再通过邮件验证码进行登录。

* 点击 [申请代码](https://huoyijie.cn/efs/auth/acc/keysettings)

    此处申请的代码，在后面生成客户端密钥时会用到，通过这个代码可以把客户端密钥关联到登录帐号，所以该代码需要妥善保管，确保不要泄漏给他人。如果代码已经使用过，则会立刻失效，不能够再生成新的密钥。

* 运行 efscli 客户端软件

    在上一节安装环节，已经下载了相应平台的客户端程序，Mac/Linux 平台需要设置可执行权限

    ```bash
    # Mac
    # 设置可执行权限
    $ chmod a+x ./efscli-darwin-amd64

    # 运行 efscli 客户端
    $ ./efscli-darwin-amd64
    ```

    或者

    ```bash
    # Linux
    # 设置可执行权限
    $ chmod a+x ./efscli-linux-amd64

    # 运行 efscli 客户端
    $ ./efscli-linux-amd64
    ```

    Windows 可以在命令行中运行，也可以双击 `efscli-0.5.0-windows-amd64.exe` 文件直接运行。

* 首次运行，生成密钥文件

    首次运行需要生成本地密钥文件，此时需要输入网站上申请的`代码`，并设置用来保护密钥文件的密码，该密码会用来对密钥文件本身内容进行加密。生成密钥文件后，会自动登录 EFS 系统。

* 生成密钥后，再次运行

    此时只需要输入刚刚设置的密码，对密钥文件进行解密，然后会自动登录 EFS 系统。

* 登录 EFS 系统

    类似插入移动硬盘到电脑，系统中会自动挂载，Mac/Linux 会把虚拟文件系统空间挂载在 `~/efscloud` 目录下，Windows 会出现一个新的盘符，如 `Z:\`。此时，就可以想使用移动硬盘一样，正常读写 EFS 文件系统。

* 修改昵称

    点击右上角帐号昵称，显示下拉菜单，进入帐号信息页面，可以修改帐号昵称。

* 设置密钥标签

    存在多个密钥的情况下，通过给密钥设置标签，客户端登录时，可以明确读取的是哪个密钥。点击右上角帐号昵称，显示下拉菜单，进入管理密钥页面，可以设置密钥标签。

## 安全隐私

* 通信安全

    efscli 客户端与服务端软件会建立 TLS 安全连接，并进行通信，因此通信过程中所有数据会被加密保护。

* 数据加密

    文件系统中主要的敏感数据为目录名称、文件名称、文件内容，在 efscli 客户端中，这些敏感数据会被加密，然后发送到服务器端，数据以加密的形式存储在服务器上。
    
    efscli 客户端从服务端读取数据时，读取到的数据是密文，也会对数据进行自动解密，然后再把解密后的数据返回给操作系统。

    从用户的角度看，存放在 EFS 中的目录和文件与正常使用的磁盘或者移动硬盘没有任何区别，但是其数据却是加密存储的。

* 密钥安全

    首次生成密钥，需要输入`代码`，为了正常生成密钥，请确保不要以任何方式泄漏给任何人。当使用该`代码`生成密钥文件后，此`代码`就立刻失效了。

    生成的密钥文件，密钥内容本身是会通过一个密码进行加密的，因此虽然密钥文件存储在用户的电脑上，即使其他人通过一定手段拿到密钥文件，没有密码，也是无法读取 EFS 文件系统数据的。

    但是，请保护好密钥文件，不要以任何方式泄漏给任何人。为了确保密钥文件不丢失，最好对该密钥文件做好备份工作，比如保存到自己的其他工作电脑或者其他移动硬盘中，但是请确保这些设备都是安全的。备份密钥文件的意义在于，如果用户电脑硬盘损坏，或者意外删除密钥文件，还可以通过备份文件恢复。

    密钥文件不会以任何形式，任何方式发送到服务器端，遵循`零信任`原则，efscli 客户端不相信服务端程序，不会把密钥文件上传到服务器，即使密钥本身也是加密的。

    只有用户设置的密码，才可以正常读取密钥文件。

* 密码安全

    生成密钥文件时，用户设置了密码，这个密码及这个密码的任何变体(如 Hash)不会被发往服务器端。

    efscli 客户端软件也不会存储该密码明文，仅仅把密码的 `bcrypt Hash` 变体存储在了密钥文件中，而这是安全的，因为 `Hash` 算法是不可逆的，不可能通过 `Hash` 值得到原始密码。efscli 客户端存储密码的 `bcrypt Hash` 变体的原因是用来进行密码认证，确保用户输入了正确的密码。

    运行 efscli 客户端软件时，当且仅当用户输入了正确的密码，才能够正常读取密钥文件内容，进而登录 EFS 系统，读写 EFS 文件系统数据。

    用户需要确保不以任何方式把密码泄漏给任何人，应该仅把密码记在大脑中。这样做的好处是，只有用户本身知道密码，进而使用密钥文件，进而登录 EFS 系统读写数据。

    [EFS 网站](https://huoyijie.cn/efs/) 登录帐号不设置密码的原因在于，避免用户把登录帐号密码和保护密钥文件的密码设置为相同的密码，避免在用户不知情的情况下，意外泄漏密码。

* 隐私安全

    用户在不以任何形式泄漏密码和密钥文件的情况下，也就拥有了确保隐私安全的存储在远程服务器上的“移动硬盘”，可以把相对隐私的数据存储在 EFS 文件系统里面。

## 技术实现

* efscli 客户端架构图

    ![efscli 客户端](https://cdn.huoyijie.cn/efs/umls/arch-1.0.0.svg)

* privkey 密钥文件

    ![privkey 密钥文件](https://cdn.huoyijie.cn/efs/umls/privkey.svg)

| 字段 | 说明 |
|-----|-----|
| KeyId | 密钥ID |
| AccId | 帐号ID |
| Email | 登录帐号 |
| PassHashAndBcrypt | bcrypt(sha256(密码明文), 14) |
| EcdsaPublicKey | 签名验证公钥 |
| EcdsaPrivateKey | aes.Encrypt(签名验证私钥, sha256(密码明文)) |
| RsaPublicKey | 非对称加密公钥 |
| RsaPrivateKey | aes.Encrypt(非对称加密私钥, sha256(密码明文)) |
| AesPrivateKey | rsa.Encrypt(对称加密密钥, RsaPublicKey) |

`密码明文`: 是生成密钥过程中，用户的输入密码。`PassHashAndBcrypt` 是密码明文先经过 `sha256` 哈希，再经过 `bcrypt` 哈希得到。是用来进行密码验证的，在 efscli 客户端软件内校验用户输入密码是否正确。也就是说密码明文及密码的任何变体都不会上传至服务器。

生成密钥文件的过程中，efscli 客户端软件会随机生成 3 个密钥，分别是 `ecdsa` 公私钥、`rsa` 公私钥和 `aes` 密钥。

`ecdsa` 密钥是用来进行进行客户端认证的，首次连接服务器生成密钥时，`ecdsa` 公钥会上传到服务器端，后续连接服务器的登录请求参数必须经过签名，并由服务端进行签名验证。

`rsa` 密钥是用来加密客户端敏感数据的，目前主要用来加密 `aes` 密钥。未来可能会出现其他类型敏感数据，也会通过 `rsa` 密钥加密。

`aes` 密钥是用来对文件系统数据进行加密的，包括目录名称、文件名称以及文件数据本身。

其中敏感数据 `EcdsaPrivateKey`、`RsaPrivateKey`、`AesPrivateKey` 都是加密存储的。

`KeyId` 是生成的密钥文件的全局唯一标识，由服务端进行分配。

`AccId` 及 `Email` 是密钥文件关联的帐号信息，生成密钥时会要求用户输入`代码`，代码是必须登录网站后，通过`申请代码`得到。

* 首次登录EFS系统

    ![首次登录EFS系统](https://cdn.huoyijie.cn/efs/umls/signup.svg)

* 非首次登录EFS系统

    ![非首次登录EFS系统](https://cdn.huoyijie.cn/efs/umls/signin.svg)

* File Transfer Protocol

    `ftp`协议运行于 `TLS` 安全连接之上，`ftp`客户端与服务端通过`rpc`进行请求与应答，请求/应答通过`protobuf`进行序列化与反序列化。在EFS系统里，文件是被分成`4KB`大小的块的集合来读或者写的。

    ![文件块](https://cdn.huoyijie.cn/efs/umls/fileblock.svg)

    如上图，是一个`410KB`文件的块划分，如果是写文件，`ftp`协议会在客户端内部将数据块进行加密，然后再传输至服务端进行存储。如果是读文件，客户端从服务器端直接读取的是加密块，`ftp`协议会自动进行解密，然后返回原始数据。

    可以看到，`ftp`协议会自动完成数据的加解密，这个过程是在 efscli 客户端软件内部完成的，服务端全程没有参与，服务端读写的文件块数据都是加密的。

    明确的说，`ftp`协议是对目录名称、文件名称以及文件块数据自动进行加密。采用 `aes` 密钥对数据进行对称加密。而前面已经论述过，`aes` 密钥加密存储于密钥文件中，必须有密码才能够解密。

* FUSE (Filesystem in Userspace)

    ![FUSE](https://cdn.huoyijie.cn/efs/umls/fuse-1.0.1.svg)

## 特别说明

目前 EFS 文件系统处于 `Beta` 测试阶段，为了确保文件数据不丢失，请确保存储在 EFS 中的文件数据在其他设备中有备份。目前每个帐号只能申请一个代码，每个代码可以用来生成一个密钥文件，每个密钥文件对应 `4GB` 云存储空间。

如果你对该项目感兴趣，或者有任何疑问，可以随时联系作者 `yijie.huo@foxmail.com`。

## Github

[huoyijie](https://github.com/huoyijie)