# agents.md

本文件为 Claude Code (claude.ai/code) 在此代码库中工作时提供指导。

## 项目概述

这是一个基于 Go 语言的 WebSocket 远程 Shell 终端系统，提供类似于 `gotty` 的多会话管理功能。该项目由一个管理 Shell 会话的服务器组件和一个通过 WebSocket 连接到这些会话的客户端组件组成。

## 构建和开发命令

### 基本构建和运行
```bash
# 构建主应用程序
go build -o main.exe main.go

# 使用 HTTP 运行服务器
go run main.go -config server-config.json -mode server

# 使用 HTTPS/TLS 运行服务器
go run main.go -config server-config-tls.json -mode server

# 使用 HTTP 运行客户端
go run main.go -config client-config.json -mode client

# 使用 HTTPS/TLS 运行客户端
go run main.go -config client-config-tls.json -mode client
```

### 测试
```bash
# 运行测试（如果有）
go test ./...

# 测试特定包
go test ./go_ws_sh
```

### Protocol Buffer 生成
```bash
# 从 protobuf 生成 Go 代码
protoc --go_out=./ --go_opt=Mwsmsg.proto=./go_ws_sh wsmsg.proto
```

## 架构概述

### 核心组件

1. **主应用程序** (`main.go`)
   - 支持服务器和客户端模式的入口点
   - 使用配置文件进行设置
   - 处理 panic 恢复

2. **核心包** (`go_ws_sh/`)
   - **服务器函数**: `Server_start()` 初始化并运行 WebSocket 服务器
   - **客户端函数**: `Client_start()` 连接到远程 WebSocket 会话
   - **身份验证**: 多种认证方法，包括密码、令牌和 WebSocket 协议头
   - **会话管理**: 数据库支持的会话、令牌和凭据存储

3. **数据库层**
   - 使用 SQLite 配合 GORM ORM
   - 三个主要存储：凭据、令牌和会话
   - 数据库文件：`credential_store.db`、`token_store.db`、`session_store.db`

4. **WebSocket 协议**
   - 使用 Avro 编码的自定义消息处理
   - PTY（伪终端）支持，实现真实终端仿真
   - 客户端和服务器之间的终端大小同步

### 主要特性

- **多用户支持**: 不同用户的多个并发会话
- **会话共享**: 通过 URL 链接共享终端会话
- **跨平台**: 支持 Windows (PowerShell) 和 Linux (bash)
- **安全性**: TLS 支持，多种身份验证方法
- **协议支持**: HTTP/1.1、HTTP/2、通过 QUIC 支持 HTTP/3

## 配置文件

### 服务器配置
- `server-config.json`: HTTP 服务器配置
- `server-config-tls.json`: HTTPS/TLS 服务器配置
- 定义初始凭据、会话和服务器端点

### 客户端配置
- `client-config.json`: HTTP 客户端配置
- `client-config-tls.json`: HTTPS/TLS 客户端配置
- 定义连接参数和身份验证

## 身份验证方法

1. **WebSocket 协议身份验证**
   ```
   sec-websocket-protocol:type%3Dtoken%26token%3Db6e915c46%26identifier%3D123456789%26username%3Dadmin
   sec-websocket-protocol:username%3Dadmin%26password%3Dpass%26type%3Dpassword
   ```

2. **基于令牌的身份验证**: 生成并使用令牌进行会话访问

3. **用户名/密码**: 传统的基于凭据的身份验证

## 开发说明

- 项目使用 Cloudwego Hertz 作为 HTTP 框架
- WebSocket 处理使用 `github.com/hertz-contrib/websocket` 和 `github.com/gorilla/websocket`
- 消息编码使用 Apache Avro，通过 `github.com/linkedin/goavro/v2`
- PTY 处理使用 `github.com/creack/pty` 和平台特定库
- 通过 `github.com/quic-go/quic-go` 支持 QUIC/HTTP3

## 文件结构

```
go_ws_sh/
├── main.go                 # 应用程序入口点
├── go_ws_sh/              # 核心包
│   ├── Server_start()     # 服务器初始化
│   ├── Client_start()     # 客户端连接
│   ├── GenerateRoutes.go  # HTTP 路由定义
│   ├── AuthorizationHandler.go
│   └── 各种处理器和存储
├── server-config*.json    # 服务器配置
├── client-config*.json    # 客户端配置
├── *.db                  # SQLite 数据库
├── localhost.crt/key     # TLS 证书
└── go-ws-sh-api/         # API 文档
```