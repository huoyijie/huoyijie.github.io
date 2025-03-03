<link rel="stylesheet" href="https://cdn.huoyijie.cn/npm/video.js@8.0.4/dist/video-js.min.css">
<script src="https://cdn.huoyijie.cn/npm/video.js@8.0.4/dist/video.min.js"></script>
<script>
    window.HELP_IMPROVE_VIDEOJS = false
</script>

# 弈杰围棋智能终端

弈杰围棋智能终端 (`YijieGo`) 是我和儿子一起设计开发的基于人工智能的围棋对弈终端。最早接触围棋也是因为儿子报了围棋兴趣课，他经常和我提到想和电脑下棋，而市面上的人机对弈软件要么是非 `AI` 很弱的，要么是 `AI` 付费的。我也不想让他总看着电脑或者平板甚至是手机下棋，再加上各种基于 `AI` 围棋程序大行其道，我们就想一起制作一个基于人工智能的围棋终端设备。我们的分工很简单，他作为用户提需求、测试和体验，我来写程序。

这个项目在写这篇文章的一年前就已经完成了，期间他也经常拿出来对弈，虽然总是输给 `AI`，也总是需要鼓励他一下，但是他总算没有灰心，一直在下。下面是我们录制的一段 `YijieGo` 的 `AI` 分析模式的视频，可以基于当前对弈情况推荐最佳选点。

<br><video id="video-1" class="video-js" controls muted preload="auto" width="720" data-setup="{}">
  <source src="https://cdn.huoyijie.cn/ab/98e27e7051ba11ecb154451bde618cf8/yijiego-1.mp4" type="video/mp4">
</video><br>

`YijieGo` 支持 9x9 13x13 19x19 路围棋，出现的方块为推荐的选点，粗线方块为最佳选点，里面的数字代表下在此处带来的目数上的收益或者损失。

<br><video id="video-2" class="video-js" controls muted preload="auto" width="720" data-setup="{}">
  <source src="https://cdn.huoyijie.cn/ab/98e27e7051ba11ecb154451bde618cf8/yijiego-2.mp4" type="video/mp4">
</video><br>

没有选择 LCD 液晶屏而是墨水屏，主要是考虑墨水屏对眼睛更友好，但是需要自己通过 `C` 语言编写显示驱动程序。因为这款墨水屏本身不支持触摸功能，因此需要通过手机来辅助落子，棋盘上每个交叉点都有坐标，如 `D3`。在手机上输入该坐标，可以在此处落子。

棋盘左上角会显示设备运行状态，总共有 4 个图标，从左到右依次代表:

1. 是否正在计算
2. 手机是否已通过蓝牙连接到设备
3. 是否已连接 WIFI
4. 设备剩余电量

## 硬件部分

我曾经整理过制作 `YijieGo` 所需的所有材料表格，和大家分享一下:

| 材料 | 价格 | 备注 |
| ---- | ---- | ---- | ---- |
| raspberry pi zero wh | ¥371.70 | Raspberry Pi Zero WH x1 |
| SanDisk 32GB TF Card | ¥30.90  | 最高读取速度 120MB/秒，写入速度 10MB/秒 |
| UPS 不间断电源  | ¥123.90 | UPS HAT (C) × 1, 1000mAh 锂电池 × 1, 螺丝包 × 1 |
| 6inch HD e-Paper HAT  | ¥598.50 | 6inch HD e-Paper × 1, e-Paper IT8951 Driver HAT (B) × 1, 6inch e-Paper Adapter (B) × 1, 40PIN FFC(同向) × 1, USB线 type A公口 转micro公口 × 1, RPi铜柱包(2PCS套) × 1, PH2.0 8PIN连接线 20cm × 1 |
| PBS-110 自复位开关带线 | ¥3 | 制作设备开机键 |
| 3mm LED 红黄绿各20个 | ¥2.18 | 60 个总价 |
| 330欧姆 1/4W金属膜电阻包 | ¥9.04 | 总价，调节经过 LED 灯电流等 |
| M2.3*6*4.6*0.9 200个 | ¥5.75 | 总价，固定芯片外壳 |
| M1.4*3*2.3*0.8 100个 | ¥2.40 | 总价，固定芯片外壳 |
| 服务器 | - | 普通服务器，KataGo 模型推理通过 CPU 计算 |
| 3D打印机 | - | 普通 3D 打印机 |
| HC-PLA直径1.75mm 1kg 黑色 | ¥68 | 单个外壳 33g + 64g + 96g 200g以内，不损耗的话，可以制作5套外壳 |

最开始我租了价格不菲的 GPU 服务器，配置非常高，运行 KataGo 推理非常快，相比普通服务器，能在更短的时间内能够给出更佳的选点，但是确实太贵了，作为个人项目没有必要上 GPU 服务器。后面我还是使用了自己的博客服务器，虽然推理速端慢了些，也足够用了。

树莓派芯片和墨水屏确实有点贵，而且墨水屏尺寸有点小，也不支持触摸功能，有点遗憾。最早的时候想的是至少 12 寸以上的大屏，甚至能够像电脑显示器一样可以挂在墙上，通过手指触摸下棋，体验会好太多了。UPS 不间断电源也比较贵，我买的那款电源经过测试，电池单独供电的时候，不间断运行可以用一个半小时左右吧，确实容量上也不太理想。

这个灰色的盒子，是我根据墨水屏、芯片、电池、连线、开关等等花了差不多一个星期设计出来的，使用的是开源的 Blender 3D 模型设计软件。确实没有学过模型设计，设计的比较粗糙，不过刚好能够把所有东西封装在里面，然后用螺丝把芯片、屏幕和外壳进行固定。因为 3D 打印机打印速度很慢，所以模型调整起来非常花时间，哪里不合适还要再调整模型，然后再打印出来再试，这个过程要反复进行。因为外壳设计的不合理，安装的时候碎了一个屏，还有一个屏线弄断了，第三个屏才安装好，太费钱了。

### 连接树莓派和墨水屏控制器(IT8951)

![连接树莓派和墨水屏控制器](https://cdn.huoyijie.cn/ab/98e27e7051ba11ecb154451bde618cf8/it8951-connect-to-rpi.jpg)

图中的不是 raspberry pi zero，但 BCM 引脚是一样的。下面的表格是两块板子相对应的引脚。

| IT8951引脚 | Pi BCM引脚 | 描述 |
| ---- | ---- | ---- |
| 5V | 5V | 电源正(5V电源输入) |
| GND | GND | 电源地 |
| MISO | P9 | SPI通信MISO引脚 |
| MOSI | P10 | SPI通信MOSI引脚 |
| SCK | P11 | SPI通信SCK引脚 |
| CS | P8 | SPI片选引脚（低电平有效） |
| RST | P17 | 外部复位引脚（低电平复位） |
| HRDY | P24 | 忙状态输出引脚（低电平表示忙） |

## 软件部分

![arch](https://cdn.huoyijie.cn/ab/98e27e7051ba11ecb154451bde618cf8/arch.svg)

上面的图是 YijieGo 总体架构图，主要分为 Server、智能围棋终端、手机遥控器三个部分。

### 1.Server

主要运行 AI Engine 和 Website，后者主要是记录对弈数据等。AI Engine 是一个基于 socket.io 和 KataGo 搭建的人工智能围棋引擎服务。[KataGo](https://github.com/lightvector/KataGo) 是开源的，它的网站上有训练好的模型，可以选择最好的模型下载到服务器上面，经过简单的配置就好了。

### 2.智能围棋终端

首先需要用 TF 卡制作系统镜像，到树莓派官网上下载 raspios-lite 最新的镜像文件，并通过 rpi imager 把镜像烧录至 TF 系统卡。制作好的系统卡插入 rpi zero，即可启动系统。(ps: 关于树莓派，我之前写过几篇介绍文章，大家感兴趣可以往前翻一翻。)

围棋终端上主要运行两个服务，`BLE Controller` 和 `Chess`，前者是设备控制中心，后者是控制墨水屏显示围棋盘、落子、提子、显示推荐选点、电量等等的驱动程序。

### 3.手机遥控器

通过微信小程序实现，通过 Bluetooth BLE 与围棋终端进行连接和操控。

## 最后

这个项目最后得以完成，离不开我的儿子，他提出了很多很好的建议，而且也负责测试和体验。(ps: 他也是受益人)

如果你对该项目感兴趣，或者有什么疑问，欢迎联系我 yijie.huo@foxmail.com 交流探讨。