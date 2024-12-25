<link rel="stylesheet" href="https://cdn.huoyijie.cn/npm/video.js@8.0.4/dist/video-js.min.css">
<script src="https://cdn.huoyijie.cn/npm/video.js@8.0.4/dist/video.min.js"></script>
<script>
    window.HELP_IMPROVE_VIDEOJS = false
</script>

# 从零开始制作树莓派小车

制作一个遥控小车是 2020 下半年时产生的一个想法，当时偶然接触到树莓派，感觉瞬间进入一个新世界。看着各种各样基于树莓派创建的有趣的项目层出不穷，手有点痒痒。我虽然学的是计算机专业，但对硬件、芯片等不太在行，决定还是从相对简单的遥控车项目上手。

以前买过遥控小车玩具，可能比较便宜的原因？可玩性都很差，也没有 DIY 的乐趣，孩子玩几下也就扔一边了。所以就想带着小孩从头开始制作一个遥控小车，让他全程参与这个过程，对他来说，过程远比最后的结果重要。

![car](https://cdn.huoyijie.cn/ab/3b8281b1e8aa6a1d8bc6718a4256b141/rpi-car.jpg)

虽然用树莓派4B来控制小车实在是有些大材小用，但是“用牛刀杀鸡”确实简单明了。树莓派 4B 本身性能强劲，基于 Linux 系统生态，再搭配 Python 语言丰富的社区软件库，使得编程任务非常简单，适合有创意的小伙伴儿基于它开发各种好玩的东西。

选定了遥控小车项目后，开始查阅相关资料，利用业余时间组装调试小车，前后差不多用了几个星期时间，小车终于可以通过电脑控制跑起来了，通过输入方向键控制小车前进、后退、左右转弯。键盘控制不太方便，想着等再有时间，要把蓝牙通信加上，这样可以用手机蓝牙遥控小车。时间一晃就来到了 2021 年，再次拿出吃灰已久的小车（还好还可以跑），把手机遥控加上了。

![car](https://cdn.huoyijie.cn/ab/3b8281b1e8aa6a1d8bc6718a4256b141/car.jpg)

上面是调试小车时拍的照片，线路没有焊接，组件是胶布粘的，看起来一碰就散，跑起来还挺稳的。下面是我们录制的视频。

<br><video id="playing_car_video" class="video-js" controls muted preload="auto" width="720" data-setup="{}">
  <source src="https://cdn.huoyijie.cn/ab/3b8281b1e8aa6a1d8bc6718a4256b141/playing_car@home.mp4" type="video/mp4">
</video><br>

<br><video
  id="playing_car_metro_video" class="video-js" controls muted preload="auto" width="400" data-setup="{}"
>
  <source src="https://cdn.huoyijie.cn/ab/3b8281b1e8aa6a1d8bc6718a4256b141/playing_car@metro.mp4" type="video/mp4">
</video><br>