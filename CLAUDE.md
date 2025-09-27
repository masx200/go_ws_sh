# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based WebSocket remote shell terminal system that provides multi-session management similar to `gotty`. The project consists of a server component that manages shell sessions and a client component that connects to these sessions via WebSocket connections.

## Build and Development Commands

### Basic Build and Run
```bash
# Build the main application
go build -o main.exe main.go

# Run server with HTTP
go run main.go -config server-config.json -mode server

# Run server with HTTPS/TLS
go run main.go -config server-config-tls.json -mode server

# Run client with HTTP
go run main.go -config client-config.json -mode client

# Run client with HTTPS/TLS
go run main.go -config client-config-tls.json -mode client
```

### Testing
```bash
# Run tests (if any)
go test ./...

# Test specific package
go test ./go_ws_sh
```

### Protocol Buffer Generation
```bash
# Generate Go code from protobuf
protoc --go_out=./ --go_opt=Mwsmsg.proto=./go_ws_sh wsmsg.proto
```

## Architecture Overview

### Core Components

1. **Main Application** (`main.go`)
   - Entry point that supports server and client modes
   - Uses configuration files for setup
   - Handles panic recovery

2. **Core Package** (`go_ws_sh/`)
   - **Server Functions**: `Server_start()` initializes and runs the WebSocket server
   - **Client Functions**: `Client_start()` connects to remote WebSocket sessions
   - **Authentication**: Multiple auth methods including password, token, and WebSocket protocol headers
   - **Session Management**: Database-backed session, token, and credential storage

3. **Database Layer**
   - Uses SQLite with GORM ORM
   - Three main stores: credentials, tokens, and sessions
   - Database files: `credential_store.db`, `token_store.db`, `session_store.db`

4. **WebSocket Protocol**
   - Custom message handling with Avro encoding
   - PTY (pseudo-terminal) support for real terminal emulation
   - Terminal size synchronization between client and server

### Key Features

- **Multi-user Support**: Multiple concurrent sessions with different users
- **Session Sharing**: Share terminal sessions via URL links
- **Cross-platform**: Works on Windows (PowerShell) and Linux (bash)
- **Security**: TLS support, multiple authentication methods
- **Protocol Support**: HTTP/1.1, HTTP/2, HTTP/3 via QUIC

## Configuration Files

### Server Configuration
- `server-config.json`: HTTP server configuration
- `server-config-tls.json`: HTTPS/TLS server configuration
- Defines initial credentials, sessions, and server endpoints

### Client Configuration
- `client-config.json`: HTTP client configuration
- `client-config-tls.json`: HTTPS/TLS client configuration
- Defines connection parameters and authentication

## Authentication Methods

1. **WebSocket Protocol Authentication**
   ```
   sec-websocket-protocol:type%3Dtoken%26token%3Db6e915c46%26identifier%3D123456789%26username%3Dadmin
   sec-websocket-protocol:username%3Dadmin%26password%3Dpass%26type%3Dpassword
   ```

2. **Token-based Authentication**: Generate and use tokens for session access

3. **Username/Password**: Traditional credential-based authentication

## Development Notes

- The project uses Cloudwego Hertz as the HTTP framework
- WebSocket handling uses both `github.com/hertz-contrib/websocket` and `github.com/gorilla/websocket`
- Message encoding uses Apache Avro via `github.com/linkedin/goavro/v2`
- PTY handling uses `github.com/creack/pty` and platform-specific libraries
- QUIC/HTTP3 support via `github.com/quic-go/quic-go`

## File Structure

```
go_ws_sh/
├── main.go                 # Application entry point
├── go_ws_sh/              # Core package
│   ├── Server_start()     # Server initialization
│   ├── Client_start()     # Client connection
│   ├── GenerateRoutes.go  # HTTP route definitions
│   ├── AuthorizationHandler.go
│   └── Various handlers and stores
├── server-config*.json    # Server configurations
├── client-config*.json    # Client configurations
├── *.db                  # SQLite databases
├── localhost.crt/key     # TLS certificates
└── go-ws-sh-api/         # API documentation
```