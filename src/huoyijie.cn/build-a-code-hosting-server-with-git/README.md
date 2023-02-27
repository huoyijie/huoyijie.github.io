
# 基于 Git 搭建代码托管服务器

我个人的项目代码主要托管在 Github，最近发现无论是 clone 代码还是 push 代码，速度都很慢。如果在我的个人服务器上搭建一个 git server，不但可以托管我个人的私有项目，也可以为我个人的公开项目提供镜像，以提升 git 操作体验。git server 的部署非常简单，主要是在服务器上安装 git 及配置 SSH Keys。我尽可能还原了所有操作细节，可能会显得有点长。

首先登录服务器，安装 git

```bash
$ sudo apt-get install git
```

添加 git 组和 git 用户

```bash
$ sudo groupadd git
$ sudo useradd -m -d /home/git -g git git
```

新建 `authorized_keys` 文件

```bash
$ cd /home/git
$ sudo mkdir .ssh && cd .ssh
$ sudo touch authorized_keys
```

新建的几个目录或文件的权限非常关键，为避免后面遇到权限错误或安全问题，需确保只有 git 用户有读写权限

```bash
$ ls -ld /home/git
drwxr-xr-x 7 git git 4096 Mar 22 14:28 /home/git

$ sudo ls -ld /home/git/.ssh
drwx------ 2 git git 4096 Mar 22 15:32 /home/git/.ssh

$ sudo ls -l /home/git/.ssh
total 4
-rw------- 1 git git 747 Mar 22 15:32 authorized_keys
```

例如，我在尝试把 `/home/git` 权限从 `755` 调整为 `775`，克隆远程仓库时，git 用户会遇到下面的错误（其他配置正常）

```bash
$ git clone git@huoyijie.cn:newproject.git
正克隆到 'newproject'...
Permission denied (publickey).
fatal: 无法读取远程仓库。

请确认您有正确的访问权限并且仓库存在。
```

下面命令会在控制台输出我的电脑(不是服务器)公钥信息，可进行复制

```bash
$ cat ~/.ssh/id_rsa.pub
```

如果本地没有这个文件，可先创建新的公钥、私钥。

```bash
$ ssh-keygen -o
Generating public/private rsa key pair.
Enter file in which to save the key (/home/huoyijie/.ssh/id_rsa):
Created directory '/home/huoyijie/.ssh'.
Enter passphrase (empty for no passphrase):
Enter same passphrase again:
Your identification has been saved in /home/huoyijie/.ssh/id_rsa.
Your public key has been saved in /home/huoyijie/.ssh/id_rsa.pub.
The key fingerprint is:
d0:82:24:8e:d7:f1:bb:9b:33:53:96:93:49:da:9b:e3 huoyijie@mylaptop
```

首先 ssh-keygen 会确认密钥的存储位置（默认是 .ssh/id_rsa），然后它会要求你输入两次密钥口令。如果你不想在使用密钥时输入口令，将其留空即可。 然而，如果你使用了密码，那么请确保添加了 -o 选项，它会以比默认格式更能抗暴力破解的格式保存私钥。

现在来添加 SSH Keys，编辑服务器上 `authorized_keys` 文件

```bash
$ sudo vi authorized_keys
```

把自己电脑的公钥信息复制过来，写入该文件，每个公钥占一行，保存退出。现在从自己的电脑可以免密码登录服务器了。

需要注意的是，目前所有（获得授权的）开发者用户都能以系统用户 git 的身份登录服务器从而获得一个普通shell。如果你想对此加以限制，则需要修改 /etc/passwd 文件中（git 用户所对应）的 shell 值。借助一个名为 git-shell 的受限 shell 工具，你可以方便地将用户 git 的活动限制在与 Git 相关的范围内。

该工具随 Git 软件包一同提供。如果将 git-shell 设置为用户 git 的登录 shell（login shell），那么该用户便不能获得此服务器的普通 shell 访问权限。若要使用git-shell，需要用它替换掉 bash 或 csh，使其成为该用户的登录 shell。

为进行上述操作，首先你必须确保 git-shell 的完整路径名已存在于 /etc/shells 文件中:

```bash
$ cat /etc/shells # 查看 /etc/shells 中 git-shell 是否已存在，如果不存在运行下条命令
$ which git-shell # 区别系统已安装 git-shell 工具，复制该命令路径
$ sudo -e /etc/shells # 把 git-shell 命令路径加入 /etc/shells 文件最后面
```

现在你可以使用 `chsh <username> -s <shell>` 命令修改任一系统用户的 shell

```bash
$ sudo chsh git -s $(which git-shell)
```

这样，用户 git 就只能利用 SSH 连接对 Git 仓库进行推送和拉取操作，而不能登录机器并取得普通 shell。如果试图登录，会被拒绝：

```bash
$ ssh git@huoyijie.cn
fatal: Interactive git shell is not enabled.
hint: ~/git-shell-commands should exist and have read and execute access.
Connection to huoyijie.cn closed.
```

此时，用户仍可通过 SSH 端口转发来访问任何可达的 git 服务器。 如果你想要避免它，可编辑 authorized_keys 文件并在所有想要限制的公钥之前添加以下选项：

```conf
no-port-forwarding,no-X11-forwarding,no-agent-forwarding,no-pty
```

类似这样:

```bash
$ cat ~/.ssh/authorized_keys
no-port-forwarding,no-X11-forwarding,no-agent-forwarding,no-pty ssh-rsa
AAAAB3NzaC1yc2EAAAADAQABAAABAQCB007n/ww+ouN4gSLKssMxXnBOvf9LGt4LojG6rs6h
PB09j9R/T17/x4lhJA0F3FR1rP6kYBRsWj2aThGw6HXLm9/5zytK6Ztg3RPKK+4kYjh6541N
YsnEAZuXz0jTTyAUfrtU3Z5E003C4oxOj6H0rfIF1kKI9MAQLMdpGW1GYEIgS9EzSdfd8AcC
IicTDWbqLAcU4UpkaX8KyGlLwsNuuGztobF8m72ALC/nLF6JLtPofwFBlgc+myivO7TCUSBd
LQlgMVOFq1I2uPWQOkOWQAHukEOmfjy2jctxSDBQ220ymjaNsHT4kgtZg2AYYgPqdAv8JggJ
ICUvax2T9va5 gsg-keypair
```

现在，网络相关的 Git 命令依然能够正常工作，但是开发者用户已经无法得到一个普通 shell 了。

那么应该如何将已有的仓库(如存储于 Github)导入自己的 git 服务器呢？我电脑上有一个放在 Github 的项目 revealjs

```bash
$ cd ~/vswork/revealjs
$ git clone --bare revealjs revealjs.git
克隆到纯仓库 'revealjs.git'...
完成。
```

首先进入项目目录，并执行 clone `--bare` 操作，生成了 revealjs 仓库的一个裸仓，也就是目录 revealjs.git，按照惯例，裸仓库的目录以 .git 结尾。`--bare` 是限定克隆出裸仓。裸仓 `只包含 .git 目录`，克隆裸仓的命令效果类似于:

```bash
$ cp -rf revealjs/.git revealjs.git
```

然后把裸仓放入到服务器上，因为 git 用户本身已限制了 ssh 登录后的权限，只能进行 git 相关操作。所以需要另一个有权限登录 shell 的用户 abc 把 revealjs.git 上传到 `/home/git` 目录下

```bash
$ scp -r revealjs.git abc@huoyijie.cn:/home/abc
```

用 abc 用户登入远程服务器，把仓库移动到 git 用户可以读写的目录(`/home/git`)下，并更改刚刚上传的裸仓目录归属及权限

```bash
sudo mv /home/abc/revealjs.git /home/git/revealjs.git
sudo chown -R git:git /home/git/revealjs.git
sudo chmod 755 /home/git/revealjs.git
```

然后可以从远程服务器 clone 这个仓库了，就像从 Github 上面克隆一样的。

```bash
$ git clone git@huoyijie.cn:revealjs.git
正克隆到 'revealjs'...
remote: Counting objects: 165, done.
remote: Compressing objects: 100% (159/159), done.
remote: Total 165 (delta 61), reused 0 (delta 0)
接收对象中: 100% (165/165), 1.68 MiB | 127.00 KiB/s, 完成.
处理 delta 中: 100% (61/61), 完成.
```

如果一个用户（如 git 用户），可 SSH 连接到一个服务器，并且有 /home/git/revealjs.git 目录的写权限，那么该用户将自动拥有推送权限。

那么应该如何创建一个完全新的仓库呢？

首先通过 abc 用户登录服务器，取得 shell，并创建新的裸仓

```bash
$ cd /home/git
$ sudo mkdir newproject.git && cd newproject.git
$ sudo git init --bare
Initialized empty Git repository in /home/git/newproject.git/
```

修改裸仓用户归属

```bash
sudo chown -R git:git /home/git/newproject.git
```

然后客户端就可以远程操作这个仓库了，先配置一下客户端 git 用户身份 (非服务器用户身份)，填好 `user.name & user.email`

```bash
git config --global -e
```

此时可以克隆仓库了，不过现在还是一个空仓库

```bash
$ git clone git@huoyijie.cn:newproject.git
正克隆到 'newproject'...
warning: 您似乎克隆了一个空仓库。
```

注: 克隆代码库时，可以使用相对 `/home/git` 的相对路径 `git clone git@huoyijie.cn:newproject.git`，也可以使用绝对路径 `git clone git@huoyijie.cn:/home/git/newproject.git`。

如果本地已有项目，也可以直接让本地目录关联到远程仓库

```bash
$ tree -L 2 project
project
└── readme.md
```

如上，我本地有个目录 project，现在把这个目录绑定到远程仓库 newproject 上面。

```bash
$ cd project
$ git init
$ git add .
$ git commit -m 'initial commit'
$ git remote add origin git@huoyijie.cn:/home/git/newproject.git
$ git push origin master
对象计数中: 3, 完成.
写入对象中: 100% (3/3), 214 bytes | 214.00 KiB/s, 完成.
Total 3 (delta 0), reused 0 (delta 0)
To huoyijie.cn:/home/git/newproject.git
 * [new branch]      master -> master
```

现在新加的 readme.md 文件已经 push 到远程了。

OK，到现在为止，基本可以满足我的代码托管需求了，比较适合几个人以内的小团队管理代码，当然如果有很多权限等方面的定制需求，可以考虑更复杂一些的开源方案，如 Gitlab。

本文参考了 [Pro Git book](https://git-scm.com/book/zh/v2) 这本书，可以在 `https://git-scm.com` 网站上查看或下载。[下载这本书](https://cdn.huoyijie.cn/ab/7d4e75e08ac011ebaa2ab110efea0133/progit_v2.1.55.pdf)