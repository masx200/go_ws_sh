# go_ws_sh

go_ws_sh

### 介绍：全 Go 语言 WebSocket 远程 Shell 模拟器

这款远程 Shell 模拟器完全使用 Go 语言开发，客户端和服务端均基于 Go 构建。它利用
WebSocket
协议实现高效双向通信，并支持伪终端（PTY）窗口大小的实时同步，确保本地与远程环境的一致性。该模拟器提供了类似`gotty`的功能，如多用户会话共享、URL
链接分享和自定义命令启动，同时增强了安全性和跨平台兼容性。适用于运维管理、开发协作、教育培训及产品演示等多种场景，提供流畅、安全的远程命令行体验。

# 支持的身份认证方式

1.http basic 身份验证

```http
Authorization : Basic YWRtaW46cGFzcw==
```

```http
Authorization : Bearer b6e915c46
```

2.websocket protocol 身份验证

```http
sec-websocket-protocol:token=b6e915c46
```

```http
sec-websocket-protocol:username=admin;password=pass
```

# 用法

```
Usage of main.exe:
  -config string
        the configuration file
  -mode string
        server or client mode
  -serverip string
        server ip
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

# 代码解释

### 代码解释

#### 文件：HandleWebSocketProcess.go

##### 主要功能

该文件实现了 WebSocket 连接的处理逻辑，主要用于与客户端建立 WebSocket
连接，并通过该连接执行命令、发送和接收数据。以下是主要函数及其功能：

1. **`SendTextMessage` 函数**

   - **功能**：通过 WebSocket 连接发送文本消息。
   - **参数**：
     - `conn`: WebSocket 连接。
     - `typestring`: 消息类型。
     - `body`: 消息体。
     - `binaryandtextchannel`: 用于传递消息的通道。
   - **实现**：将消息编码为 JSON 格式并通过通道发送。

2. **`HandleWebSocketProcess` 函数**
   - **功能**：处理 WebSocket
     连接的整个生命周期，包括读取消息、执行命令并发送结果。
   - **参数**：
     - `session`: 包含要执行的命令和参数的会话信息。
     - `codec`: 用于编解码 Avro 消息的编解码器。
     - `conn`: 与客户端的 WebSocket 连接。
   - **实现**：
     - 建立 WebSocket 连接并初始化相关资源。
     - 读取客户端发送的消息，解析为 JSON 格式。
     - 处理不同类型的消息（如调整终端大小、执行命令等）。
     - 将命令输出通过 WebSocket 发送给客户端。
     - 处理命令结束后的清理工作。

##### 关键点

- 使用了 `github.com/hertz-contrib/websocket` 和 `github.com/linkedin/goavro/v2`
  库来处理 WebSocket 和 Avro 编解码。
- 通过 `console.New` 创建了一个虚拟终端，用于执行命令并捕获其输出。
- 使用了多个 goroutine 来并发处理读写操作，确保性能和响应性。
- 通过 `binaryandtextchannel` 通道来同步不同 goroutine 之间的消息传递。

#### 文件：pipe-std-ws-client.go

##### 主要功能

该文件实现了 WebSocket 客户端逻辑，负责连接到 WebSocket
服务器，并处理与服务器的通信，包括身份验证、消息编码和解码，以及与标准输入/输出的交互。

1. **`Client_start` 函数**

   - **功能**：从配置文件中读取配置并启动 WebSocket 客户端。
   - **参数**：
     - `config`: 配置文件路径。
   - **实现**：解析配置文件并调用 `pipe_std_ws_client` 函数。

2. **`pipe_std_ws_client` 函数**

   - **功能**：创建并管理 WebSocket 客户端连接。
   - **参数**：
     - `configdata`: 配置数据结构。
   - **实现**：
     - 配置 WebSocket 客户端连接（包括 TLS 设置）。
     - 连接到 WebSocket 服务器并发送初始消息（如终端大小）。
     - 处理来自服务器的消息（如命令输出）并将其写入标准输出。
     - 处理用户输入并将其发送到服务器。

3. **`sendMessageToWebsocketStdin` 函数**

   - **功能**：将标准输入的数据编码为 Avro 格式并通过 WebSocket 发送到服务器。
   - **参数**：
     - `data`: 输入数据。
     - `codec`: Avro 编解码器。
     - `binaryandtextchannel`: 用于传递消息的通道。
   - **实现**：将输入数据封装为 Avro 消息并通过通道发送。

4. **`configureWebSocketTLSCA` 函数**
   - **功能**：配置 WebSocket 客户端的 TLS 设置。
   - **参数**：
     - `x`: WebSocket Dialer 实例。
     - `configdata`: 配置数据结构。
   - **实现**：根据配置加载 CA 证书并设置 TLS 配置。

##### 关键点

- 使用了 `github.com/gorilla/websocket` 库来处理 WebSocket 连接。
- 支持 TLS 加密连接，确保通信安全。
- 使用了 `goavro/v2` 库进行 Avro 编解码，保证数据格式的一致性和高效传输。
- 通过 `TermboxPipe` 函数接管标准输入输出，实现与终端的交互。

### 总结

这两个文件共同实现了一个基于 WebSocket
的远程命令执行系统。`HandleWebSocketProcess.go` 负责服务器端的 WebSocket
处理逻辑，而 `pipe-std-ws-client.go` 则实现了客户端的 WebSocket
连接和交互逻辑。两者通过 Avro
编解码器进行数据交换，确保了数据格式的一致性和高效传输。

## 生成 protobuf 文件

```
protoc --go_out=./ --go_opt=Mwsmsg.proto=./go_ws_sh wsmsg.proto
```

## 服务端配置文件举例

### windows 系统

```json
{
  "token_file": "token_store.json",
  "credential_file": "credential_store.json",
  "sessions": [
    {
      "username": "admin",
      "path": "pwsh",
      "cmd": "pwsh",
      "args": ["-noProfile"],
      "dir": "C:\\Users\\Public"
    }
  ],
  "servers": [
    {
      "alpn": "h2",
      "port": "28443",
      "protocol": "https",
      "cert": "localhost.crt",
      "key": "localhost.key"
    },
    {
      "alpn": "h3",
      "port": "28443",
      "protocol": "https",
      "cert": "localhost.crt",
      "key": "localhost.key"
    }
  ]
}
```

### linux 系统

```json
{
  "token_file": "token_store.json",
  "credential_file": "credential_store.json",
  "sessions": [
    {
      "username": "admin",
      "path": "bash",
      "cmd": "bash",
      "args": ["-i"],
      "dir": "/root"
    }
  ],
  "servers": [
    {
      "alpn": "h2",
      "port": "28443",
      "protocol": "https",
      "cert": "localhost.crt",
      "key": "localhost.key"
    },
    {
      "alpn": "h3",
      "port": "28443",
      "protocol": "https",
      "cert": "localhost.crt",
      "key": "localhost.key"
    }
  ]
}
```

### 代码解释

这段代码展示了两个不同操作系统的服务端配置文件示例：Windows 和
Linux。配置文件使用 JSON 格式，定义了服务启动所需的各种参数。以下是具体解释：

#### 公共字段

- **token_folder**: 存储 token 的文件夹路径，值为 `"token_folder"`。

#### credentials（凭证）

- 定义了访问服务所需的用户名和密码。
- Windows 和 Linux 配置中均包含一个用户：
  - **username**: 用户名为 `"admin"`
  - **password**: 密码分别为 `"pass"` (Windows) 和 `"password"` (Linux)

#### sessions（会话）

- 定义了服务启动的命令行会话。
- **path**: 会话路径，分别指向 `pwsh` (Windows) 和 `bash` (Linux)。
- **cmd**: 启动命令，与 path 相同。
- **args**: 命令行参数，如 Windows 下的 `["-noProfile"]` 和 Linux 下的
  `["-i"]`。
- **dir**: 工作目录，分别是 `"C:\\Users\\Public"` (Windows) 和 `"/root"`
  (Linux)。

#### servers（服务器）

- 定义了服务监听的服务器配置。
- 每个系统下有两个服务器配置，分别支持 HTTP/2 (`h2`) 和 HTTP/3 (`h3`) 协议。
- **alpn**: 应用层协议协商，指定协议版本。
- **port**: 监听端口，均为 `"28443"`.
- **protocol**: 使用的传输协议，均为 `"https"`.
- **cert** 和 **key**: SSL/TLS 证书和私钥路径，在 Windows 中位于当前目录，而在
  Linux 中位于当前目录.

这些配置项确保了服务在不同操作系统上能够正确初始化并运行。
