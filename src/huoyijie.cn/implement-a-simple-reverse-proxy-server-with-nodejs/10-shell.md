# 启停脚本管理

* 启动服务器脚本 `start.sh`

```bash
sudo node app.js > logs/console.log 2> logs/error.log &
```

* 关闭服务器脚本 `stop.sh`

```bash
sudo kill `cat hawkey.pid`
```

其中 hawkey.pid 记录的是主进程 ID，主进程启动后会把 process.pid 写入该文件，主进程退出前会删除该文件。

* 重启服务器脚本 `restart.sh`

```bash
./stop.sh && ./start.sh
```