package go_ws_sh

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	// "github.com/hertz-contrib/websocket"
)

type RouteConfig struct {
	Path    string
	Method  string
	Handler func(context.Context, *app.RequestContext)
}

// handleWebSocket 处理WebSocket请求

// PromiseAll 接受一个函数切片，每个函数都会被并发调用。
// 每个函数应该没有参数并且返回一个接口和一个错误。
// 它返回一个通道，该通道将发送一个包含所有结果的切片或第一个错误。

func pipe_std_ws_server(config ConfigServer, credentialdb *gorm.DB, tokendb *gorm.DB) {
	var listtokensHandler = ListTokensHandler(credentialdb, tokendb)
	authHandler := AuthorizationHandler(credentialdb, tokendb)
	var routes = []RouteConfig{

		{
			Path:    "/authorization",
			Method:  "POST",
			Handler: authHandler,
		},

		{
			Path:    "/tokens",
			Method:  "POST",
			Handler: listtokensHandler,
		},
		{
			Path:    "/authorization",
			Method:  "PUT",
			Handler: authHandler,
		},
		{
			Path:    "/authorization",
			Method:  "DELETE",
			Handler: authHandler,
		},
	}

	err := EnsureCredentials(config, credentialdb)
	if err != nil {
		log.Fatal(err)
		return
	}
	var handlermap = map[string]func(w context.Context, r *app.RequestContext){}
	for _, session := range config.Sessions {
		handlermap[session.Path] = createhandleWebSocket(session)
	}
	handlerGet := createhandlerauthorization(credentialdb, tokendb, func(w context.Context, r *app.RequestContext) {
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
	handlerPost := createhandlerloginlogout(config.Sessions, credentialdb, tokendb, func(w context.Context, r *app.RequestContext) {

		r.AbortWithMsg("Not Found", consts.StatusNotFound)
		// return

		// handlermap[name](w, r)
	})
	for _, serverconfig := range config.Servers {
		tasks = append(tasks, createTaskServer(serverconfig, func(w context.Context, r *app.RequestContext) {

			fmt.Println("Request FullURI:", string(r.URI().FullURI()))
			fmt.Println("Request Method:", string(r.Method()))
			fmt.Println("Request Headers:")
			fmt.Println("{")
			r.Request.Header.VisitAll(func(key, value []byte) {
				fmt.Println(string(key), ":", string(value))
			})
			fmt.Println("}")
			//routes
			for _, route := range routes {
				if route.Path == string(r.Path()) && route.Method == string(r.Method()) {
					route.Handler(w, r)
					return
				}
			}

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

type ConfigServer struct {
	CredentialFile string         `json:"credential_file"`
	Sessions       []Session      `json:"sessions"`
	Servers        []ServerConfig `json:"servers"`

	TokenFile string `json:"token_file"`

	// 添加初始化用户名和初始密码字段
	InitialUsername string `json:"initial_username"`
	InitialPassword string `json:"initial_password"`
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
		return
	}

	credentialFile := configdata.CredentialFile
	if credentialFile == "" {
		credentialFile = "credential_store.db"
	}

	tokenFile := configdata.TokenFile
	if tokenFile == "" {
		tokenFile = "token_store.db"
	}
	credentialdb, err := gorm.Open(sqlite.Open(credentialFile), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	tokendb, err := gorm.Open(sqlite.Open(tokenFile), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	credentialdb.AutoMigrate(&CredentialStore{})
	tokendb.AutoMigrate(&TokenStore{})
	pipe_std_ws_server(configdata, credentialdb, tokendb)
}
