package go_ws_sh

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	// "net/http"
	"os"
	"strings"

	// "unicode/utf8"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/websocket"
)

// handleWebSocket 处理WebSocket请求

// PromiseAll 接受一个函数切片，每个函数都会被并发调用。
// 每个函数应该没有参数并且返回一个接口和一个错误。
// 它返回一个通道，该通道将发送一个包含所有结果的切片或第一个错误。

func pipe_std_ws_server(config Config /* httpServeMux *http.ServeMux, handler func(w context.Context, r *app.RequestContext) */) {

	var handlermap = map[string]func(w context.Context, r *app.RequestContext){}
	for _, session := range config.Sessions {
		handlermap[session.Path] = createhandleWebSocket(session)
	}
	handler := createhandler(config /* httpServeMux */, func(w context.Context, r *app.RequestContext) {
		var name = r.Param(":name")
		if handler, ok := handlermap[name]; ok {
			handler(w, r)
		} else {
			r.AbortWithMsg("Not Found", consts.StatusNotFound)
			return
		}
		// handlermap[name](w, r)
	})
	// for _, session := range config.Sessions {
	// 	httpServeMux.HandleFunc("/"+session.Path, createhandleWebSocket(session))

	// }
	tasks := []func() (interface{}, error){}

	for _, serverconfig := range config.Servers {
		tasks = append(tasks, createTaskServer(serverconfig, handler))
	}
	// 启动服务器
	result := <-PromiseAll(tasks)

	switch v := result.(type) {
	case error:
		fmt.Printf("Error: %v\n", v)
	case []interface{}:
		fmt.Println("All tasks completed successfully. Results:")
		for _, res := range v {
			fmt.Printf("%v\n", res)
		}
	default:
		fmt.Println("Unexpected result type")
	}
}

func createhandleWebSocket(session Session) func(w context.Context, r *app.RequestContext) {
	return func(w context.Context, r *app.RequestContext) {
		codec, err := create_msg_codec()
		if err != nil {
			log.Println(err)
			// r.SetStatusCode(http.StatusInternalServerError)
			r.AbortWithMsg("Internal Server Error"+"\n"+err.Error(), consts.StatusInternalServerError)
			// w.WriteHeader(http.StatusInternalServerError)
			// w.Write([]byte("Internal Server Error" + "\n" + err.Error()))
			return
		}

		upgrader := websocket.HertzUpgrader{
			CheckOrigin: func(ctx *app.RequestContext) bool {
				return true // 允许跨域
			},
			EnableCompression: true,
		}
		err = upgrader.Upgrade(r, func(conn *websocket.Conn) {
			HandleWebSocketProcess(session, codec, conn)
		})
		// 升级HTTP连接到WebSocket协议
		// conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			r.AbortWithMsg("Internal Server Error"+"\n"+err.Error(), consts.StatusInternalServerError)
			// w.WriteHeader(http.StatusInternalServerError)
			// w.Write([]byte("Internal Server Error" + "\n" + err.Error()))
			return
		}
		// 设置标准输入、输出和错误流
		// 启动命令
		// 处理标准输出和错误流
		// 等待命令执行完毕
		//进程结束
		// 处理标准输入流
		//读取out_queue,并转换
		//将数据转换为二进制
		// 发送消息到WebSocket连接
		//读取out_queue,并转换
		//将数据转换为二进制
		// 发送消息到WebSocket连接
		// 循环读取WebSocket消息
		//连接结束
		// 将消息发送回客户端
		// err = conn.WriteMessage(mt, message)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }

	}
}

// 定义结构体以匹配JSON结构
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session struct {
	Path string   `json:"path"`
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

type Server struct {
	Port string `json:"port"`
	Tls  bool   `json:"tls"`
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

type Config struct {
	Credentials []Credentials `json:"credentials"`
	Sessions    []Session     `json:"sessions"`
	Servers     []Server      `json:"servers"`
}

func Server_start(config string) {
	configFile, err := os.Open(config)
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	jsonDecoder := json.NewDecoder(configFile)
	var configdata Config
	err = jsonDecoder.Decode(&configdata)
	if err != nil {
		log.Fatal(err)
	}
	// var httpServeMux = http.NewServeMux()

	pipe_std_ws_server(configdata /* httpServeMux, handler */)
}
func createhandler(config Config, next func(w context.Context, r *app.RequestContext) /* httpServeMux *http.ServeMux */) func(w context.Context, r *app.RequestContext) {
	return func(w context.Context, r *app.RequestContext) {

		Upgrade := strings.ToLower(r.Request.Header.Get("Upgrade"))
		Connection := strings.ToLower(r.Request.Header.Get("Connection"))
		//if !tokenListContainsValue(r.Request.Header, "Connection", "upgrade") {
		if !strings.Contains(Connection, "upgrade") {
			r.NotFound() //http.NotFound(w, r)
			return
		}
		if !strings.Contains(Upgrade, "websocket") {
			// if !tokenListContainsValue(r.Header, "Upgrade", "websocket") {
			r.NotFound() //http.NotFound(w, r)
			return
		}

		if !r.IsGet() /* != http.MethodGet */ {
			r.NotFound()
			//http.NotFound(w, r)
			return
		}
		//httpServeMux.ServeHTTP(w, r)
		next(w, r)
	}

}
