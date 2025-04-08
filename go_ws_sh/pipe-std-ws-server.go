package go_ws_sh

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	// var listtokensHandler = ListTokensHandler(credentialdb, tokendb)
	// authHandler := AuthorizationHandler(credentialdb, tokendb)
	var routes []RouteConfig //{

	handlerGet := createhandlerauthorization(credentialdb, tokendb, func(w context.Context, r *app.RequestContext) {

		sessions, err := ReadAllSessions(sessiondb)
		if err != nil {
			r.AbortWithMsg("Internal Server Error"+"\n"+err.Error(), consts.StatusInternalServerError)
			return
		}
		var handlermap = map[string]func(w context.Context, r *app.RequestContext){}
		for _, session := range sessions {
			handlermap[session.Name] = createhandleWebSocket(session)
		}
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
	// handlerPost := createhandlerloginlogout(sessions, credentialdb, tokendb, func(w context.Context, r *app.RequestContext) {

	// 	r.AbortWithMsg("Not Found", consts.StatusNotFound)
	// 	// return

	// 	// handlermap[name](w, r)
	// })
	var initial_credentials = config.InitialCredentials
	var initial_sessions = config.InitialSessions
	routes = GenerateRoutesHttp(credentialdb, tokendb, sessiondb, initial_credentials, initial_sessions)
	// routes = append(gr, routes...)
	composedMiddleware := HertzCompose(
		MatchAndRouteMiddleware([]RouteConfig{
			{
				Path: "/tokens",

				MiddleWare: AuthorizationMiddleware(credentialdb, tokendb, sessiondb),
			},
			{
				Path: "/sessions",

				MiddleWare: AuthorizationMiddleware(credentialdb, tokendb, sessiondb),
			},
			{
				Path: "/credentials",

				MiddleWare: AuthorizationMiddleware(credentialdb, tokendb, sessiondb),
			},
		}),

		MatchAndRouteMiddleware(routes))
	handler := func(w context.Context, r *app.RequestContext) {

		Upgrade := strings.ToLower(r.Request.Header.Get("Upgrade"))
		Connection := strings.ToLower(r.Request.Header.Get("Connection"))

		if string(r.Method()) == consts.MethodGet && Connection == "upgrade" && Upgrade == "websocket" {
			handlerGet(w, r)
			return
		}

		composedMiddleware(w, r, func(c context.Context, r *app.RequestContext) {
			r.AbortWithMsg(
				"Method Not Allowed",

				consts.StatusMethodNotAllowed)
		})
		// shouldReturn := MatchAndHandleRoute(w, routes, r)
		// if shouldReturn {
		// 	return
		// }
		// if string(r.Method()) == consts.MethodPost {
		// 	handlerPost(w, r)
		// 	return
		// }

	}

	middlewares := []app.HandlerFunc{

		func(c context.Context, ctx *app.RequestContext) {
			routesmiddle := MatchAndRouteMiddleware([]RouteConfig{

				{
					Path:   "/sessions",
					Method: "COPY",
					MiddleWare: HertzCompose(AuthorizationMiddleware(credentialdb, tokendb, sessiondb), func(c context.Context, r *app.RequestContext, next HertzNext) {
						r.String(consts.StatusOK, "OK COPY")

					}),
				},
				{
					Path:   "/sessions",
					Method: "MOVE",
					MiddleWare: HertzCompose(AuthorizationMiddleware(credentialdb, tokendb, sessiondb), func(c context.Context, r *app.RequestContext, next HertzNext) {
						// 定义请求体结构体
						var body struct {
							Session struct {
								Name string `json:"name"`
							} `json:"session"`
							Authorization CredentialsClient `json:"authorization"`
							Destination   struct {
								Name string `json:"name"`
							} `json:"destination"`
						}

						// 绑定请求体
						if err := r.BindJSON(&body); err != nil {
							r.AbortWithMsg("Error: "+err.Error(), consts.StatusBadRequest)
							return
						}
						for _, session := range initial_sessions {
							if session.Name == body.Session.Name {
								r.AbortWithMsg("Error: Session is initial session,不允许删除", consts.StatusBadRequest)
								return
							}
						}
						var err error
						username := body.Authorization.Username
						if username == "" {
							username, err = GetUsernameByTokenIdentifier(tokendb, body.Authorization.Identifier)
							if err != nil {
								log.Println("Error:", err)
								r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
								return
							}
							log.Println("Username:", username)
						}
						var existingSession SessionStore
						if err := sessiondb.Where(&SessionStore{Name: body.Destination.Name}).First(&existingSession).Error; err == nil {
							r.JSON(consts.StatusConflict, map[string]any{
								"message":  "Error: Session already exists",
								"username": username,
								"session": map[string]string{
									"name": body.Session.Name,
									// "username": username,
								},
							})
							return
						}
						var newSession *SessionStore
						if newSession, err = MoveSession(sessiondb, body.Session.Name, body.Destination.Name); err != nil {
							log.Printf("Failed to move session: %v", err)
							r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
							return
						}
						var args []string
						// 将 Args 字段（字符串形式）反序列化为字符串切片
						if err := json.Unmarshal([]byte(newSession.Args), &args); err != nil {
							r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
							return
						}
						// 返回成功响应
						r.JSON(consts.StatusOK, map[string]any{
							"message":  "Session moved successfully",
							"username": username,
							"session": map[string]interface{}{
								"name":     newSession.Name,
								"cmd":      newSession.Cmd,
								"args":     args,
								"dir":      newSession.Dir,
								"username": username,
							},
						})

					}),
				},
			})
			routesmiddle(c, ctx, func(c context.Context, r *app.RequestContext) {
				r.Next(c)
			})
		},
	}
	for _, serverconfig := range config.Servers {

		tasks = append(tasks, createTaskServer(serverconfig,
			handler, middlewares...))
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

	// 创建自定义 Logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到控制台
		logger.Config{
			LogLevel:      logger.Info, // 设置日志级别为 Debug [[7]][[9]]
			SlowThreshold: time.Second, // 慢查询阈值（可选）
			Colorful:      true,
		},
	)
	credentialdb, err := gorm.Open(sqlite.Open(credentialFile), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	tokendb, err := gorm.Open(sqlite.Open(tokenFile), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	sessiondb, err := gorm.Open(sqlite.Open(sessionFile), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sessiondb.AutoMigrate(&SessionStore{})
	credentialdb.AutoMigrate(&CredentialStore{})
	tokendb.AutoMigrate(&TokenStore{})
	tokendb = tokendb.Debug()
	credentialdb = credentialdb.Debug()
	sessiondb = sessiondb.Debug()
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
