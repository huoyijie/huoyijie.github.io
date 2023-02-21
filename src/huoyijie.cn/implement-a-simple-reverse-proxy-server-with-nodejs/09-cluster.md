# 集群管理

Node.js 是一个基于Chrome V8 引擎的JavaScript 运行环境。Node.js 使用了一个事件驱动、非阻塞式I/O 的模型，使其轻量又高效。Node.js 没有使用操作系统的线程模型，而是选择了更轻量级的事件循环机制。换句话说，Node.js 进程是单线程的，而且不可能也不能阻塞，如果阻塞了唯一的线程程序就不响应了。单个线程的 Node 通常只能在一个 CPU 核心上执行，如果服务器有更多的核心也用不到。想想是不是觉得浪费资源？

其实 Node 通过另外一种方式解决了这个问题，就是可以启动多个进程。比如说有 4 个 CPU 核心，可以选择启动 4 个工作进程。下面介绍的代码位于 app.js 中。

```javascript
const config = require('./config')('server-config.json'),
      log = require('./log')(config.env, console),
      os = require('os'),
      numCPUs = os.cpus().length,
      fs = require('fs'),
      cluster = require('cluster'),
      worker = require('./worker'),
      pidFilePath = 'hawkey.pid'

if(cluster.isMaster) {
  // 主进程

  // 输出主进程 ID 到文件中，方便通过脚本读取该文件获取进程号，进行启动停止进程的操作
  fs.writeFileSync(pidFilePath, `${process.pid}`)
  console.info(`master.process [${process.pid}] running`)

  // fork 出 cpu 核心数量的 工作进程
  for(let i = 0; i < numCPUs; i++) {
    let wk = cluster.fork()
    console.info(`master.process [${process.pid}], lanch worker-${i} [${wk.process.pid}]`)
  }

  cluster.on('exit', (wk, code, signal) => {
    console.info(`worker.process [${wk.process.pid}] exit, code [${code}], signal [${signal}]`)
  })

  // 优雅关闭主进程，进而优雅的关闭所有工作进程
  // 主进程会给所有工作进程发送 shutdown 信号，工作进程收到该信息号后，会主动 shutdown
  // 等工作进程 shudown 后，主进程会 kill 掉工作进程
  let shutdown = () => {
    for(let id in cluster.workers) {
      let wk = cluster.workers[id]
      wk.send('shutdown', (err) => {
        if(!err) {
          if(!wk.isDead()) {
            wk.kill()
          }
        } else {
          console.error(`master send shutdown to worker [${wk.process.pid}] error`, err)
          wk.process.kill()
        }
      })
    }
    // delete .pid file
    fs.unlink(pidFilePath, (err) => {
      if (!err) {
        console.debug('delete .pid file')
      } else {
        console.error('delete .pid file error', err)
      }
    })
  }

  // 主进程注册信号处理程序
  // 收到 Ctrl+C 或者 kill 信号后，主进程会主动 shutdown，并退出
  process
    .on('unhandledRejection', (reason, promise) => {
      console.error(reason, 'Unhandled Rejection at Promise', promise, '@master')
    })
    .on('uncaughtException', (err, origin) => {
      console.error('Uncaught Exception thrown', err.stack, origin, 'stop master')
      shutdown()
      process.exitCode = 1
    })
    .once('SIGINT', function (code) {
      console.warn('SIGINT received...', code, 'stop master')
      shutdown()
    })
    .once('SIGTERM', function (code) {
      console.warn('SIGTERM received...', code, 'stop master')
      shutdown()
    })
} else {

  // 工作进程，worker 组件代码位于 worker.js，前面已经介绍过了。
  // worker 组件主要是启动服务器工作进程，并进入监听端口状态，完成具体的请求处理

  worker(config)
  console.info(`worker.process [${process.pid}] running`)
}
```

如上代码为主进程启动后，fork 出更多工作进程来完成实际的请求响应处理。代码上加了注释，比较好理解。