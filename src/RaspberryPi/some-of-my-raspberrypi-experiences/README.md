# 树莓派的一些使用经验

现在越来越多的开发者喜欢上了基于 Linux 生态系统的树莓派，开发者可以使用它 DIY 各种有趣的项目，编程语言的选择很多，c、c++、nodejs、python 等非常自由。我在[弈杰围棋智能终端](https://huoyijie.cn/article/98e27e7051ba11ecb154451bde618cf8/)项目中就使用了 raspberry pi zero，期间遇到了很多问题，在这篇文章里，我想把这些经验总结下来，方便后面查阅。

## 1.自制树莓派开/关机按钮

![Raspberry Pi Pinout](https://cdn.huoyijie.cn/ab/3b8281b1e8aa6a1d8bc6718a4256b141/GPIO-Pinout.png)

如上图所示，针脚 GPIO 3(也就是 pin 5)具有重启系统的功能，如果把 GPIO 3 和 任意一个 Ground 针脚通过自复位按钮开关连接起来，就可以通过这个按钮控制树莓派关机(进入 halt 状态)，或者重启(从 halt 状态恢复)。树莓派的这个设计和 PC 电脑是一样的，当电脑关机后，其实也是进入了 halt 状态，尽管耗电量很小，但是主板仍然是供电的。此时去按开机键，就会使得电脑从 halt 状态中恢复(开机)。

自复位开关可以是按下去使得电路闭合，松手后开关自动复位电路断开，电脑的开机按钮就是自复位开关。在淘宝上一搜有很多卖的。

除非特殊情况，不要通过直接断开电源或者插电源来进行重启系统，可能会对存储文件系统造成破坏，强烈建议给树莓派加上开机键。

GPIO 3 上的重启功能是固定的，不能更换其他针脚。除了重启功能 GPIO 3 针脚默认还有 shutdown 功能。处于运行状态的树莓派，按下开机键，系统会进入 halt 状态，安全优雅地 shutdown。

通过针脚图可以发现，GPIO 3 属于 I2C，当系统开启了 I2C 后，shutdown 功能会和 I2C 冲突(重启功能不会冲突，因为 halt 状态 I2C 是不工作的)。此时可以把 GPIO 3 针脚的 shutdown 功能通过配置切换到其他 GPIO 针脚(我使用了 GPIO 27)。

简易电路图
```
            ----------
GPIO 3 -----|        |
            | switch |------ Ground
GPIO 27 ----|        |
            ----------
```

按照上面的电路图焊接按钮，按钮下方一般有2个引脚，其中一个引脚焊接两根杜邦线，然后两根杜邦线分别插入 GPIO 3 和 GPIO 27，另外一个引脚焊接一根杜邦线，插入树莓派的 Ground 针脚（有很多，可以就近选择一个）。

编辑 /boot/config.txt 文件，增加如下一行:

```conf
dtoverlay=gpio-shutdown,gpio_pin=27
# 默认 gpio_pin=3 可以 safe shutdown，但是如果开启了 I2C会失效，此时调整配置为 gpio 27
```

`sudo shutdown -r now` 重启系统测试，此时应该可以通过按钮进行安全关机和开机了。

自制开机按钮的想法参考了 [www.raspberrypi.org](https://www.raspberrypi.org/forums/viewtopic.php?f=91&t=217442&sid=0db2744db6f854cdac4e5af7df2a0e7d&start=25)

树莓派针脚图可以参考 [pinout.xyz](https://pinout.xyz/)

## 2. 连接 LED 显示系统状态

我有个需求，就是系统开机后，LED 指示灯亮，关机后，LED 指示灯灭。树莓派本身就支持这个功能，只需要通过配置打开，然后在连接好电路。

简易电路图
```
                                      ---------
Ground - PIN 6  ----电阻(330欧)---- (-)|  LED  |
TXD    - PIN 8  ------------------ (+)|       |
                                      ---------
```
LED 等可以选择额定电压 2V 左右的，淘宝上有很多卖的。LED 灯的短的针脚是负极，负极上焊接一个 330 欧姆的电阻，然后焊接到杜邦线上插入到树莓派的 PIN 6 Ground 针脚。其中电阻是用来保护树莓派避免被过大的电流烧坏。LED 灯的长的针脚是正极，可以焊接到杜邦线上然后插入到树莓派的 TXD PIN 8 针脚。焊接好后，可以用绝缘胶带缠好。

然后编辑 /boot/config.txt，增加如下一行
```conf
enable_uart=1
```

然后重启系统进行测试。

连接 LED 灯实现参考了 [howchoo.com](https://howchoo.com/g/ytzjyzy4m2e/build-a-simple-raspberry-pi-led-power-status-indicator)

## 3. raspberry pi zero w 配置 wifi

我的项目里使用了 pi zero w，这个开发板没有 USB，所以不方便连接鼠标键盘。一般都是通过 ssh 登录远程控制，但是新的开发板显然没有 wifi 配置，也没有开启 ssh 服务，无法通过 ssh 控制，也不能通过鼠标键盘输入。树莓派系统本身已经考虑到这个问题了，就是可以在制作系统盘时，通过一点儿小技巧来实现系统自动配置 wifi 和开启 ssh。

制作系统启动盘时，一般选择 TF 卡，然后去官网下载系统镜像文件，然后通过 rpi imager 软件把镜像写入 TF card。此时可以在 TF 卡的 /boot 根目录下放置一个 wpa_supplicant.conf 文件，其中配置 wifi 用户名密码信息，文件名最好复制，如果键盘输入一定仔细检查避免拼写错误。

wpa_supplicant.conf 文件内容如下:

```conf
country=CN
ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
update_config=1

network={
  ssid="YOUR WIFI SSID"
  psk="YOUR WIFI PASSWORD"
  scan_ssid=1
  key_mgmt=WPA-PSK
}
```

再在 /boot 目下放置一个内容为空的文件，命名为 ssh。

然后把这张启动盘插入树莓派，系统启动后会读取这两个文件，自动配置 wifi 信息和开启 ssh 服务。此时可以通过电脑远程 ssh 登录树莓派，进行远程操作。

配置 wifi 参考了 [www.raspberrypi-spy.co.uk](https://www.raspberrypi-spy.co.uk/2017/04/manually-setting-up-pi-wifi-using-wpa_supplicant-conf/)

## 4. raspberry pi zero 安装 node

raspberry pi zero cpu 架构为 arm6，Node 官方预编译的版本早已不在支持，必须自己编译安装 node，好在他们还维护了一个[非官方支持网站](https://unofficial-builds.nodejs.org)，有很多非官方支持的预编译版本。我采用了 node-v8.16.0-linux-armv6l.tar.xz

```bash
$ wget https://unofficial-builds.nodejs.org/download/release/v8.16.0/node-v8.16.0-linux-armv6l.tar.xz
$ tar -xvJf node-v8.16.0-linux-armv6l.tar.xz
$ sudo ln -s /home/pi/node-v8.16.0-linux-armv6l/bin/node /usr/local/bin/node
$ sudo ln -s /home/pi/node-v8.16.0-linux-armv6l/bin/npm /usr/local/bin/npm
```

更换源
```bash
$ npm config set registry https://registry.npm.taobao.org
```

pkg 打包，这个有点复杂，pkg 对 arm6 的支持要追溯到 pkg@4.3.8 版本，之后的版本不支持 armv6，之前的版本存在 bugs 无法打包成功。我在尝试过很多次之后，唯有 pkg@4.3.8 版本可以打包成功。pkg@4.3.8 支持到 node-v8.16.0，这也是我上面选择安装 node-v8.16.0 的原因

```bash
$ npm i pkg@4.3.8 -g -u
$ sudo ln -s /home/pi/node-v8.17.0-linux-armv6l/bin/pkg /usr/local/bin/pkg
```

编辑 /home/pi/node-v8.16.0-linux-armv6l/lib/node_modules/pkg/node_modules/pkg-fetch/lib-es5/index.js 文件第 96 行，nodeVersion 值修改为 'v8.16.0'

```javascript
fetched = (0, _places.localPlace)({ from: 'fetched', arch: arch, nodeVersion: 'v8.16.0', platform: platform, version: _package.version });
```

下载 [预编译 node8](https://github.com/yao-pkg/pkg-binaries/releases/tag/node8)
选择 https://github.com/yao-pkg/pkg-binaries/releases/download/node8/fetched-v8.16.0-linux-armv6

```bash
$ cp fetched-v8.16.0-linux-armv6 ～/.pkg-cache/v2.5/fetched-v8.16.0-linux-armv6
```

然后可以进入到项目目录运行打包命令。

## 5. apt 源指向国内镜像

```bash
$ echo "update apt source"
$ sudo echo "\
deb http://mirrors.tuna.tsinghua.edu.cn/raspbian/raspbian/ buster main non-free contrib rpi
deb-src http://mirrors.tuna.tsinghua.edu.cn/raspbian/raspbian/ buster main non-free contrib rpi
" > /etc/apt/sources.list

$ sudo echo "\
deb http://mirrors.tuna.tsinghua.edu.cn/raspberrypi/ buster main ui
" > /etc/apt/sources.list.d/raspi.list

$ echo "apt update ..."
$ sudo apt update
$ echo "apt update done"
```

## 6. UPS 不间断电源

购买树莓派开发板时，一般只有电源线，在插电源线供电的情况下，如果意外断点或者电源线被拔出来了，系统会突然断掉可能造成文件系统损坏。如果加装一个 UPS 模块就可以很好的解决这个问题。UPS 模块一般由锂电池和电路板组成，支持一边充电一边放电，电路板上也提供了额外的功能，比如可以实时获得电池供电电压、电流，据此推算电池剩余使用时间等等。

每个树莓派型号，市面上都有对应设计好的 UPS 模块，比如 raspberry pi zero 相对应的 UPS 提供这样的功能:

+ 采用弹簧顶针接口设计，适用于 Raspberry Pi Zero 系列主板
+ 板载锂电池充电芯片，支持动态路径管理，供电更稳定
+ 板载升压芯片，可提供稳定 5V 电压输出
+ 可通过 I2C 接口通信，测量电池电压、电流、功率和剩余电量等参数，实时检测模块工作状态
+ 板载电池保护电路，防止过充、过放、过流和短路，工作稳定更安全
+ 板载充电指示灯，方便查看电池工作状态

## 7. 系统设置

系统联网后需要调整时区，默认时区会导致进程中在格式化时间字符串时遇到时区问题
```bash
$ sudo timedatectl set-timezone Asia/Shanghai
```

如果遇到 vi 打开文件后中文乱码，可调整 vim 文件编码
```bash
$ set fileencodings=utf-8
$ set termencoding=utf-8
$ set encoding=utf-8
```

如何开启 SPI 接口
+ Run sudo raspi-config.
+ Use the down arrow to select 5 Interfacing Options
+ Arrow down to P4 SPI.
+ Select yes when it asks you to enable SPI,
+ Also select yes if it asks about automatically loading the kernel module.
+ Use the right arrow to select the <Finish> button.
+ Select yes when it asks to reboot.

如何开启 I2C 接口
+ Run sudo raspi-config.
+ Use the down arrow to select 5 Interfacing Options
+ Arrow down to P5 I2C.
+ Select yes when it asks you to enable I2C
+ Also select yes if it asks about automatically loading the kernel module.
+ Use the right arrow to select the <Finish> button.
+ Select yes when it asks to reboot.

开启 SPI/I2C 接口参考了 [raspberry-pi-spi-and-i2c-tutorial](https://learn.sparkfun.com/tutorials/raspberry-pi-spi-and-i2c-tutorial/all)


## 8. 通过蓝牙与树莓派通信

如何与基于树莓派搭建的设备通信？

通常都会需要从另一个设备访问树莓派的信息或者控制树莓派，一般有两种方式。推荐的方式是通过蓝牙与树莓派直接通信，这种方式简单直接，限制少。另一种方式是让设备都处于同一个 wifi 下，通过 http 协议或者 tcp 协议通信。

有一个场景，如果拿到基于树莓派的设备后，如果不能够通过直接读写 TF 卡的方式预先配置好 wifi，那就必须通过蓝牙去设置 wifi，否则树莓派没有办法联网。

我是基于 Node 构建蓝牙 BLE 通信程序，树莓派作为 peripheral，需要先安装 [bleno](https://github.com/noble/bleno) 和 bluetooth-hci-socket npm 包
```bash
$ npm i @abandonware/bleno
$ npm i @abandonware/bluetooth-hci-socket
```

如果使用电脑作为 central 设备，可以安装 [noble](https://github.com/noble/noble) npm 包
```bash
$ npm i @abandonware/noble
```

上方是安装了 bleno 和 noble 的一个开发分支 @abandonware/bleno 和 @abandonware/noble，因为直接安装它们在运行时会报错，这个报错在开发分支上修复了。因为项目已不再积极维护了，所以在做技术选型时也需要慎重考虑一下，可以采用更新的技术框架。

如果用手机作为 central 设备，推荐通过微信小程序来连接树莓派，小程序有专门支持蓝牙的 [BLE API](https://developers.weixin.qq.com/miniprogram/dev/api/device/bluetooth-ble/wx.writeBLECharacteristicValue.html)。