
# 语音远程遥控树莓派小车

本文接 [微信小程序(蓝牙BLE)远程遥控树莓派小车](https://huoyijie.cn/article/769dba20801c11eb983f31fd884051bb/) 那篇文章，在实现手机端微信小程序与树莓派小车之间通过蓝牙 BLE 通信的基础上，增加了小程序端的语音识别能力，让小车可以接收语音命令。本文实践了微信公众平台开放的“微信同声传译插件”，这个插件可以让小程序轻松获得语音识别和不同语言之间翻译的能力。

首先 Clone 代码([代码地址](https://github.com/huoyijie/raspberrypi-car))

```bash
git clone git@github.com:huoyijie/raspberrypi-car.git
```

目录说明

```bash
$ tree -L 3
.
├── ble
│   ├── car.py
│   ├── miniclientctl
│   │   ├── app.js
│   │   ├── app.json
│   │   ├── app.wxss
│   │   ├── index
│   │   ├── project.config.json
│   │   └── sitemap.json
│   ├── miniclientctl_2.0_voice
│   │   ├── app.js
│   │   ├── app.json
│   │   ├── app.wxss
│   │   ├── components
│   │   ├── image
│   │   ├── index
│   │   ├── project.config.json
│   │   ├── sitemap.json
│   │   └── utils
│   └── rpicarctl
│       ├── characteristic.js
│       ├── main.js
│       ├── node_modules
│       ├── package.json
│       └── package-lock.json
├── car.py
├── LICENSE
└── README.md
```

本文主要是针对小程序进行调整，其他部分程序都不用调整。可以看到增加了 miniclientctl_2.0_voice 目录，是增加了语音遥控的小程序代码。

首先在微信公众平台上面注册一个小程序，把 `project.config.json` 中的 `appid` 替换成刚刚在公众平台获得的 `appid`。

![配置appid](https://cdn.huoyijie.cn/ab/ad271380823c11ebb2910fd63136fd48/config-appid.jpg)


然后在 `微信公众平台 → 设置 → 第三方设置 → 插件管理` 中 添加微信同声传译插件 (`wx069ba97219f66d99`)，这样小程序就可以使用该插件了。

![配置插件](https://cdn.huoyijie.cn/ab/ad271380823c11ebb2910fd63136fd48/config-plugin.jpg)

然后在项目文件app.json中配置该插件

![配置插件](https://cdn.huoyijie.cn/ab/ad271380823c11ebb2910fd63136fd48/config-plugin2.jpg)

编辑app.js文件，增加获取授权代码`getRecordAuth`

```javascript
App({
  onLaunch: function () {

  },

  // 权限询问
  getRecordAuth: function() {
    wx.getSetting({
      success(res) {
        console.log("succ")
        console.log(res)
        if (!res.authSetting['scope.record']) {
          wx.authorize({
            scope: 'scope.record',
            success() {
                // 用户已经同意小程序使用录音功能，后续调用 wx.startRecord 接口不会弹窗询问
                console.log("succ auth")
            }, fail() {
                console.log("fail auth")
            }
          })
        } else {
          console.log("record has been authed")
        }
      }, fail(res) {
          console.log("fail")
          console.log(res)
      }
    })
  }
})
```

小程序启动后，初始化完成后会先加载首页代码，即 index/index.js，可以在加载 index.js 时触发刚刚的授权请求。编辑 index.js 文件，增加 onLoad 方法，里面会增加请求授权和初始化插件代码。

```javascript
onLoad() {
  this.initRecord()
  app.getRecordAuth()
}
```

在看插件初始化代码前，先看下本次主要新增的代码有哪些

![语音遥控代码](https://cdn.huoyijie.cn/ab/ad271380823c11ebb2910fd63136fd48/voice-code.jpg)

主要有 `streamRecord`、`streamRecordEnd`、`initRecord` 3 个方法，分别对应开始录音、录音结束以及初始化插件的代码。语音识别完成后的回调函数其实是在 initRecord 里面设置的，因此当语音识别为文字时，可以在此处转换为对小车发出的指令。具体看下代码

```javascript
initRecord() {
  // 识别结束事件
  manager.onStop = (res) => {

    let text = res.result.trim()
    console.log(`>> 树莓派小车接手指令: ${text}`)

    if(text == '') {
      this.stop() // 指令异常先停止
      this.showRecordEmptyTip()
      return
    }

    if (text === '前进。' || text === 'GO.') {
      this.forward()
    } else if (text === '后退。' || text === 'Back.') {
      this.backward()
    } else if (text === '停止。' || text === 'Stop.') {
      this.stop()
    } else if (text === '左转。' || text === 'Left.') {
      this.left()
    } else if (text === '右转。' || text === 'Right.') {
      this.right()
    } else {
      this.stop()
    }

    this.setData({
      recording: false,
      bottomButtonDisabled: false,
    })
  }

  // 识别错误事件
  manager.onError = (res) => {

    this.setData({
      recording: false,
      bottomButtonDisabled: false,
    })

  }
}
```

`manage.onStop` 会在语音识别为文字后调用，`res.result` 即为识别后的文字，此时可以判断出如果是指定的小车指令，则给小车发出相应指令。

最后是小程序界面上增加了录音按钮，长按按钮录音，放开按钮后触发语音识别。

![语音遥控小程序](https://cdn.huoyijie.cn/ab/ad271380823c11ebb2910fd63136fd48/voice-remotectl.jpg)

OK，主要内容介绍完了，如果还有疑问可以看一下项目的代码，语音识别部分代码不多比较容易懂。文中语音识别部分参考了微信开源项目 [Face2FaceTranslator](https://github.com/Tencent/Face2FaceTranslator/) 的代码。