# 桌面屏幕或窗口录屏并转为 GIF 图片工具推荐

## 录屏视频

桌面屏幕录屏工具有很多，一般都是保存为 mp4/webm 格式视频文件（如 Ubuntu 22.10 系统，按 `Print Screen` 键，既可以截屏也可以录屏）。把视频文件上传到视频平台也非常方便。但是有的场景直接使用视频可能不太方便，GIF 这种广为支持的动态图片格式可以一定程度上代替录屏等一些简单的视频。

## Markdown 嵌入视频

我的一个使用场景是，博客中偶尔会放一些视频，markdown 本身不支持视频，在 markdown 中插入视频需要一些额外的工作。

大多数人会把视频传到视频平台，比如 youtube，然后在 markdown 中插入一个视频封面图片链接，视频封面图上可以加一个播放按钮。点击图片链接后跳转至 youtube 平台播放视频。

还有一些人可能会通过自定义 markdown 插件，通过 js 嵌入其他视频平台的播放器。

或者直接把视频上传到 CDN，然后在 markdown 中引入 videojs 嵌入视频播放器。

如果对视频播放质量要求比较高，上面这些方法都可行。对于一些类似录屏的简单视频，还有另外一种简单的方法，就是把视频转换为 GIF 动态图片，在 markdown 中嵌入图片是非常方便的。

## 视频转 GIF

如果在搜索引擎中检索一下，可以看到非常多的视频转 GIF 在线工具。这些在线工具都会要求先把视频上传至服务器（存在隐私泄漏问题），然后由服务器完成转码工作。然后可以下载转换后的 GIF 图片。如果需要可以看下 [webm-to-gif](https://cloudconvert.com/webm-to-gif)。目测这些在线平台可能是通过 [ffmpeg](https://ffmpeg.org) 实现视频转码的。

还有一类需要下载安装的视频转码软件，感觉一般也都是封装的 ffmpeg，不需要上传视频，在自己的电脑上完成转码工作。

如果你电脑安装了 ffmpeg 工具，也可以直接通过它进行转换，非常便捷。转换命令如下:
```bash
$ ffmpeg -i input.webm -pix_fmt rgb8 output.gif
```

最后，还有一个基于浏览器的免安装的工具 [gifcap](https://github.com/joaomoreno/gifcap)，不需要上传视频，不需要服务器参与，仅在客户端运行。

## gifcap

![gifcap screenshot](https://cdn.huoyijie.cn/uploads/2023/09/screenshot.png)

Features:

* No installations, no bloatware, no updates: this works in any modern browser, including Google Chrome, Firefox, Edge and Safari;
* No server side, everything is 100% client-side only. All data stays in your machine, nothing gets uploaded to any server, the entire application is made of static files;
* PWA support makes it easy to add gifcap to your OS list of applications;
* Blazing fast GIF rendering powered by WASM, libimagequant and gifsicle;
* Highly optimized GIF file sizes, thanks to frame deduplication, boundary delta detection and lossy encoding;
* Entire screen recordings, or selection of single window;
Intuitive trimming UI
* Easy cropping via visual drag-and-drop

## [huoyijie.github.io/gifcap](https://huoyijie.github.io/gifcap)

我为了方便使用，把这个工具也部署到了 github pages 上，删除了一些统计等无关代码。工具使用起来也很简单，主要原理是基于 webrtc 的录屏 API 获取屏幕、窗口、浏览器 tab 页的录屏，然后通过 encoder.wasm 实现转码为 GIF 图片。

![infinity-record](https://cdn.huoyijie.cn/uploads/2023/09/infinity-record.gif)