package go_ws_sh

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	// "net/http"
	"os"
	// "strings"
	// "unicode/utf8"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	// "github.com/hertz-contrib/websocket"
)

// handleWebSocket 处理WebSocket请求

// PromiseAll 接受一个函数切片，每个函数都会被并发调用。
// 每个函数应该没有参数并且返回一个接口和一个错误。
// 它返回一个通道，该通道将发送一个包含所有结果的切片或第一个错误。

func pipe_std_ws_server(config ConfigServer /* httpServeMux *http.ServeMux, handler func(w context.Context, r *app.RequestContext) */) {

	var handlermap = map[string]func(w context.Context, r *app.RequestContext){}
	for _, session := range config.Sessions {
		handlermap[session.Path] = createhandleWebSocket(session)
	}
	handlerGet := createhandlerauthorization(config.TokenFile, config.CredentialFile /* config */ /* httpServeMux */, func(w context.Context, r *app.RequestContext) {
		var name = r.Param("name")
		if handler2, ok := handlermap[name]; ok {

			handler2(w, r)
		} else {
			log.Println("Not Found shell path", name)
			r.AbortWithMsg("Not Found", consts.StatusNotFound)
			return
		}
		// handlermap[name](w, r)
	})
	// for _, session := range config.Sessions {
	// 	httpServeMux.HandleFunc("/"+session.Path, createhandleWebSocket(session))

	// }
	tasks := []func() (interface{}, error){}
	handlerPost := createhandlerloginlogout(config.Sessions, config.TokenFile, config.CredentialFile /* config */ /* httpServeMux */, func(w context.Context, r *app.RequestContext) {

		r.AbortWithMsg("Not Found", consts.StatusNotFound)
		// return

		// handlermap[name](w, r)
	})
	for _, serverconfig := range config.Servers {
		tasks = append(tasks, createTaskServer(serverconfig, func(w context.Context, r *app.RequestContext) {
			if string(r.Method()) == consts.MethodGet {
				handlerGet(w, r)
				return
			}
			if string(r.Method()) == consts.MethodPost {
				handlerPost(w, r)
				return
			}
			r.AbortWithMsg(
				"Method Not Allowed",

				consts.StatusMethodNotAllowed)
			// return
		}))
	}
	// 启动服务器
	result, ok := PromiseAll(tasks).Receive()
	if ok {
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

}

// 定义结构体以匹配JSON结构
type Credentials struct {
	Username  string `json:"username"`
	Hash      string `json:"hash"`
	Salt      string `json:"salt"`
	Algorithm string `json:"algorithm"`
}

type Session struct {
	Username string   `json:"username"`
	Path     string   `json:"path"`
	Cmd      string   `json:"cmd"`
	Args     []string `json:"args"`
	Dir      string   `json:"dir"`
}

type ServerConfig struct {
	Alpn     string `json:"alpn"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
	Cert     string `json:"cert"`
	Key      string `json:"key"`
}

type CredentialsStore []Credentials
type TokenStore []struct {
	Hash       string `json:"hash"`
	Salt       string `json:"salt"`
	Algorithm  string `json:"algorithm"`
	Identifier string `json:"identifier"`
	Username   string `json:"username"`
}
type ConfigServer struct {
	CredentialFile string         `json:"credential_file"`
	Sessions       []Session      `json:"sessions"`
	Servers        []ServerConfig `json:"servers"`

	TokenFile string `json:"token_file"`
}

func Server_start(config string) {
	configFile, err := os.Open(config)
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	jsonDecoder := json.NewDecoder(configFile)
	var configdata ConfigServer
	err = jsonDecoder.Decode(&configdata)
	if err != nil {
		log.Fatal(err)
	}
	// var httpServeMux = http.NewServeMux()

	pipe_std_ws_server(configdata /* httpServeMux, handler */)
}
