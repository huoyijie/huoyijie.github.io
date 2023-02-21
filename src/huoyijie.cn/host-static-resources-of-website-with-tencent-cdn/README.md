# 网站图片视频接入CDN

我最近用之前买的云主机搭建了一个博客网站，所有环境都已经搭建好了，程序也部署上去了，功能都正常，不过网站上的图片加载速度非常慢，上传的视频则是一直卡顿，无法正常播放。服务器带宽最早选的是预付费 1Mbps，后来才改成按流量付费（带宽不限制），按流量付费对于流量小的网站更划算。测试下来，1Mbps 的带宽，当页面中图片比较多时，页面加载速度明显慢，手机拍摄的视频也无法正常播放。

为了保证图片显示和视频播放的正常体验，我决定把网站上所有图片和视频接入 CDN 服务，我使用的是腾讯云的服务，首先开通 CDN 服务，然后购买优惠的流量包。

![cdn package](https://cdn.huoyijie.cn/ab/ea33f3f07bc111ebb0c11fbcd78b590e/cdn-package.jpg)

购买完成后，大概需要5分钟时间流量包会生效。我买了2个100G的流量包，有效期是12个月。

![packages](https://cdn.huoyijie.cn/ab/ea33f3f07bc111ebb0c11fbcd78b590e/packages.png)

等流量包生效后，首先要去配置 cdn 接入域名，也就是通过哪个域名访问图片等资源，可以使用二级域名。比如我的域名是 huoyijie.cn，那么我可以使用 cdn.huoyijie.cn 作为接入域名。

进入CDN控制台，点击域名管理->添加域名

![接入域名配置](https://cdn.huoyijie.cn/ab/ea33f3f07bc111ebb0c11fbcd78b590e/config-cdn-domain.jpg)

接下来填入基本信息，绿色标出来的都填一下，最后点击确认提交。

![填写域名信息](https://cdn.huoyijie.cn/ab/ea33f3f07bc111ebb0c11fbcd78b590e/config-cdn-domain2.jpg)

下面是我的接入域名配置，可以参考一下。

![我的域名信息](https://cdn.huoyijie.cn/ab/ea33f3f07bc111ebb0c11fbcd78b590e/config-cdn-domain3.jpg)

有几个比较重要的点提一下，注意下CNAME cdn.huoyijie.cn.cdn.dnsv1.com，此域名是加速域名 (cdn.huoyijie.cn) CNAME 到 CDN 节点上的地址，直接访问此域名则无法获取正确资源信息。等所有配置都调整好后，可以通过域名 cdn.huoyijie.cn 访问 cdn 中的资源。

我选择的自有源站，回源协议暂时选的 HTTP，后面把回源服务的SSL证书配置好后，可以选择HTTPS。源站地址填写自己的回源服务 ip:port，也就是说当浏览器请求 CDN 的资源没有缓存时，会请求 ip:port 这个服务去请求源站的资源，然后再缓存到 CDN 里面，后续对该资源的访问只要没有过期，CDN 会直接返回给浏览器。

配置好域名后，接下来需要到域名注册服务商的网站增加一个配置，即加速域名 (cdn.huoyijie.cn) CNAME 到CDN 节点上的地址 cdn.huoyijie.cn.cdn.dnsv1.com。

![配置加速域名](https://cdn.huoyijie.cn/ab/ea33f3f07bc111ebb0c11fbcd78b590e/config-speed-domain.jpg)

找到域名后点击解析按钮进入解析记录列表，点击添加记录按钮，主要填3个信息。我的加速域名是 cdn.huoyijie.cn，所以主机记录填写二级域名 cdn 。记录类型选择 CNAME。记录值填写 cdn.huoyijie.cn.cdn.dnsv1.com。其他可默认，然后点击保存，大概需要几分钟生效时间。

![解析配置](https://cdn.huoyijie.cn/ab/ea33f3f07bc111ebb0c11fbcd78b590e/config-speed-domain2.jpg)

配置生效后可以 ping 域名验证下，可以看到加速域名已经解析到了 CDN 节点上。

```bash
$ ping cdn.huoyijie.cn
PING 48lbuckn.slt.cdntip.com (140.207.247.213) 56(84) bytes of data.
64 bytes from 140.207.247.213 (140.207.247.213): icmp_seq=1 ttl=58 time=5.10 ms
64 bytes from 140.207.247.213 (140.207.247.213): icmp_seq=2 ttl=58 time=4.80 ms
64 bytes from 140.207.247.213 (140.207.247.213): icmp_seq=3 ttl=58 time=10.2 ms
```

因为现在 CDN 节点上没有任何资源的缓存，现在通过加速域名 cdn.huoyijie.cn 请求资源都会被回源到上面配置过的回源服务器。如果你试着去访问一下加速域名下随便一个路径，会发现都是 404。因为还没有搭建自有的回源服务。

我是通过 nodejs 搭建了一个简单的回源服务，登录到回源服务器上面，配置好端口启动服务。再去访问加速域名下的资源，不再是返回 404 了，虽然首次加载资源还是比较慢(CDN 节点需要回源)，但是通过预热后，后面再去访问就非常快了。

我很快把网站上的资源调整了，通过cdn加速域名访问资源，发现 chrome 浏览器地址栏里出现了不安全提示。我的网站是支持 https 的，但是在网页里引用资源的 URL 无论是用 http 还是 https，浏览器都显示不安全。

![https not safe](https://cdn.huoyijie.cn/ab/ea33f3f07bc111ebb0c11fbcd78b590e/https-not-safe.png)

https 网站引用资源用 http 时显示不安全很好理解，通过 https 引用资源也显示不安全，分析了一下，原来是因为还没有给加速域名配置 SSL 证书。

再找到网站与域名下的 SSL 证书菜单，点击进入管理证书控制台

![ssl cert](https://cdn.huoyijie.cn/ab/ea33f3f07bc111ebb0c11fbcd78b590e/ssl-cert.jpg)

点击申请免费证书

![ssl cert](https://cdn.huoyijie.cn/ab/ea33f3f07bc111ebb0c11fbcd78b590e/ssl-cert2.jpg)

证书绑定域名填写加速域名，我的加速域名是 cdn.huoyijie.cn，填写好邮箱和备注点下一步，选择自动验证域名，然后提交申请。大概等 10 分钟左右，SSL 证书就会签发下来。

然后回到 CDN 控制台下的证书管理模块点击配置证书

![ssl cert](https://cdn.huoyijie.cn/ab/ea33f3f07bc111ebb0c11fbcd78b590e/ssl-cert3.jpg)

域名选择之前配置好的加速域名，证书来源选择腾讯托管证书，然后证书列表里会显示刚刚申请的 SSL 证书，选择上。回源协议选择 HTTP，点击提交，大概等几分钟配置生效。然后再访问网站，不安全提示已经没有了。

![https not safe](https://cdn.huoyijie.cn/ab/ea33f3f07bc111ebb0c11fbcd78b590e/https-safe.png)

上面介绍的是最基本的配置，CDN 控制台上还有很多高级配置项，包括安全等方面的，都可以研究一下。

![高级配置](https://cdn.huoyijie.cn/ab/ea33f3f07bc111ebb0c11fbcd78b590e/config-cdn-adv.jpg)

OK，主要内容都介绍完了，接下来要考虑的就是管理好回源服务器，通过工具能够快速把资源发布到回源服务器上面。