# go_ws_sh

go_ws_sh

### 简短介绍：全Go语言WebSocket远程Shell模拟器

这款远程Shell模拟器完全使用Go语言开发，客户端和服务端均基于Go构建。它利用WebSocket协议实现高效双向通信，并支持伪终端（PTY）窗口大小的实时同步，确保本地与远程环境的一致性。该模拟器提供了类似`gotty`的功能，如多用户会话共享、URL链接分享和自定义命令启动，同时增强了安全性和跨平台兼容性。适用于运维管理、开发协作、教育培训及产品演示等多种场景，提供流畅、安全的远程命令行体验。

```
Usage of main.exe:
  -config string
        the configuration file
  -mode string
        server or client mode
```

```
go run main.go -config server-config-tls.json -mode server
```

```
go run main.go -config client-config-tls.json -mode client
```

```
go run main.go -config server-config.json -mode server
```

```
go run main.go -config client-config.json -mode client
```

## 配置文件查看文件

[`server-config-tls.json`](server-config-tls.json)和[`client-config-tls.json`](client-config-tls.json)

[`server-config.json`](server-config.json)和[`client-config.json`](client-config.json)

## 参考资料

https://learn.microsoft.com/zh-cn/windows/console/console-virtual-terminal-sequences

https://pkg.go.dev/github.com/nsf/termbox-go

https://pkg.go.dev/github.com/runletapp/go-console
