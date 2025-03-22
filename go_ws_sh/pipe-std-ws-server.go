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

// RequestLoggerMiddleware 是一个 Hertz 中间件，用于记录请求的完整 URI、方法和头信息
func RequestLoggerMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		log.Println("Request FullURI:", string(ctx.URI().FullURI()))
		log.Println("Request Method:", string(ctx.Method()))
		log.Println("Request Headers:")
		log.Println("{")
		ctx.Request.Header.VisitAll(func(key, value []byte) {
			log.Println(string(key), ":", string(value))
		})
		log.Println("}")
		// 继续处理请求
		ctx.Next(c)
	}
}

// handleWebSocket 处理WebSocket请求

// PromiseAll 接受一个函数切片，每个函数都会被并发调用。
// 每个函数应该没有参数并且返回一个接口和一个错误。
// 它返回一个通道，该通道将发送一个包含所有结果的切片或第一个错误。

func pipe_std_ws_server(config ConfigServer, credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB) {
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

	sessions, err := ReadAllSessions(sessiondb)
	if err != nil {
		return
	}

	var handlermap = map[string]func(w context.Context, r *app.RequestContext){}
	for _, session := range sessions {
		handlermap[session.Name] = createhandleWebSocket(session)
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
	handlerPost := createhandlerloginlogout(sessions, credentialdb, tokendb, func(w context.Context, r *app.RequestContext) {

		r.AbortWithMsg("Not Found", consts.StatusNotFound)
		// return

		// handlermap[name](w, r)
	})
	handler := func(w context.Context, r *app.RequestContext) {

		shouldReturn := MatchAndHandleRoute(w, routes, r)
		if shouldReturn {
			return
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

	}
	for _, serverconfig := range config.Servers {

		tasks = append(tasks, createTaskServer(serverconfig,
			handler))
	}
	// 启动服务器
	result, ok := PromiseAll(tasks).Receive()
	if ok {
		switch v := result.(type) {
		case error:
			fmt.Printf("Error: %v\n", v)
		case []interface{}:
			log.Println("All tasks completed successfully. Results:")
			for _, res := range v {
				fmt.Printf("%v\n", res)
			}
		default:
			log.Println("Unexpected result type")
		}
	}

}

func MatchAndHandleRoute(w context.Context, routes []RouteConfig, r *app.RequestContext) bool {
	for _, route := range routes {
		// 检查 Path 是否为空，若不为空则进行匹配
		pathMatch := route.Path == "" || route.Path == string(r.Path())
		// 检查 Method 是否为空，若不为空则进行匹配
		methodMatch := route.Method == "" || route.Method == string(r.Method())

		headersMatch := true
		if len(route.Headers) > 0 {
			for key, value := range route.Headers {
				headerValue := string(r.Request.Header.Peek(key))
				if headerValue != value {
					headersMatch = false
					break
				}
			}
		}

		// 如果 Path、Method 和 Headers 都匹配，则执行处理函数
		if pathMatch && methodMatch && headersMatch {
			route.Handler(w, r)
			return true
		}
	}
	return false
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

	sessionFile := configdata.SessionFile
	if sessionFile == "" {
		sessionFile = "session_store.db"
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
	sessiondb, err := gorm.Open(sqlite.Open(sessionFile), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sessiondb.AutoMigrate(&SessionStore{})
	credentialdb.AutoMigrate(&CredentialStore{})
	tokendb.AutoMigrate(&TokenStore{})
	err = EnsureSessions(configdata, sessiondb)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = EnsureCredentials(configdata, credentialdb)
	if err != nil {
		log.Fatal(err)
		return
	}
	pipe_std_ws_server(configdata, credentialdb, tokendb, sessiondb)
}
