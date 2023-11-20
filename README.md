# 网络连接监控小程序

运行在光猫/路由器中的，监控网络通断情况的小程序。
主要代码使用ChatGPT编写，实现的比较简单。

# arm交叉编译编译命令

```powershell
# powershell
$env:GOOS = "linux"
$env:GOARCH = "arm"
$env:GOARM = "7"
go build -o net-conn-monitor-arm net-conn-monitor.go

```

# 消息推送
当前使用 [pushplus](https://www.pushplus.plus)推送。


# 运行
程序启动前，需要配置 `PUSH_TOKEN` 环境变量。