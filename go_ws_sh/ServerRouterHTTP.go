package go_ws_sh

import (
	"context"
	"log"

	"github.com/akrennmair/slice"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"
)

// GenerateRoutes 根据 openapi 文件生成路由配置
func GenerateRoutes(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB) []RouteConfig {
	routes := []RouteConfig{
		// /tokens POST
		{
			Headers: map[string]string{"x-HTTP-method-override": "POST"},
			Path:    "/tokens",
			Method:  "POST",
			MiddleWare: func(c context.Context, r *app.RequestContext, next HertzNext) {
				// 处理创建令牌的逻辑
				// 可以在这里添加从数据库查询、插入等操作
				// 示例：调用处理创建令牌的函数
				CreateTokenHandler(credentialdb, tokendb, sessiondb, c, r)
			},
		},
		// /tokens PUT
		{
			Path:   "/tokens",
			Method: "PUT",
			MiddleWare: func(c context.Context, r *app.RequestContext, next HertzNext) {
				// 处理修改令牌的逻辑
				UpdateTokenHandler(credentialdb, tokendb, sessiondb, c, r)
			},
		},
		// /tokens DELETE
		{
			Path:   "/tokens",
			Method: "DELETE",
			MiddleWare: func(c context.Context, r *app.RequestContext, next HertzNext) {
				// 处理删除令牌的逻辑
				DeleteTokenHandler(credentialdb, tokendb, sessiondb, c, r)
			},
		},
		// /tokens GET
		{
			Headers: map[string]string{"x-HTTP-method-override": "GET"},
			Path:    "/tokens",
			Method:  "POST",
			MiddleWare: func(c context.Context, r *app.RequestContext, next HertzNext) {
				// 处理显示令牌的逻辑
				GetTokensHandler(credentialdb, tokendb, sessiondb, c, r)
			},
		},
		// /credentials PUT
		{
			Path:   "/credentials",
			Method: "PUT",
			MiddleWare: func(c context.Context, r *app.RequestContext, next HertzNext) {
				// 处理修改密码的逻辑
				UpdateCredentialHandler(credentialdb, tokendb, sessiondb, c, r)
			},
		},
		// /credentials GET
		{
			Headers: map[string]string{"x-HTTP-method-override": "GET"},
			Path:    "/credentials",
			Method:  "POST",
			MiddleWare: func(c context.Context, r *app.RequestContext, next HertzNext) {
				// 处理显示用户的逻辑
				GetCredentialsHandler(credentialdb, tokendb, sessiondb, c, r)
			},
		},
		// /credentials POST
		{
			Headers: map[string]string{"x-HTTP-method-override": "POST"},
			Path:    "/credentials",
			Method:  "POST",
			MiddleWare: func(c context.Context, r *app.RequestContext, next HertzNext) {
				// 处理创建用户的逻辑
				CreateCredentialHandler(credentialdb, tokendb, sessiondb, c, r)
			},
		},
		// /credentials DELETE
		{
			Path:   "/credentials",
			Method: "DELETE",
			MiddleWare: func(c context.Context, r *app.RequestContext, next HertzNext) {
				// 处理删除用户的逻辑
				DeleteCredentialHandler(credentialdb, tokendb, sessiondb, c, r)
			},
		},
		// 可以根据 openapi 文件添加更多接口的路由配置
		// /sessions POST
		{
			Headers: map[string]string{"x-HTTP-method-override": "POST"},
			Path:    "/sessions",
			Method:  "POST",
			MiddleWare: func(c context.Context, r *app.RequestContext, next HertzNext) {
				// 处理创建会话的逻辑
				CreateSessionHandler(credentialdb, tokendb, sessiondb, c, r)
			},
		},
		// /sessions PUT
		{
			Path:   "/sessions",
			Method: "PUT",
			MiddleWare: func(c context.Context, r *app.RequestContext, next HertzNext) {
				// 处理修改会话的逻辑
				UpdateSessionHandler(credentialdb, tokendb, sessiondb, c, r)
			},
		},
		// /sessions DELETE
		{
			Path:   "/sessions",
			Method: "DELETE",
			MiddleWare: func(c context.Context, r *app.RequestContext, next HertzNext) {
				// 处理删除会话的逻辑
				DeleteSessionHandler(credentialdb, tokendb, sessiondb, c, r)
			},
		},
		// /sessions GET
		{
			Headers: map[string]string{"x-HTTP-method-override": "GET"},
			Path:    "/sessions",
			Method:  "POST",
			MiddleWare: func(c context.Context, r *app.RequestContext, next HertzNext) {
				// 处理显示会话的逻辑
				GetSessionsHandler(credentialdb, tokendb, sessiondb, c, r)
			},
		},
	}

	return routes
}

// 新增删除用户处理函数声明
func DeleteCredentialHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 实现删除用户的具体逻辑
	authHandler := AuthorizationHandler(credentialdb, tokendb)
	authHandler(c, r)
}

// 以下是示例处理函数，需要根据实际业务逻辑实现
func CreateTokenHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 实现创建令牌的具体逻辑
	authHandler := AuthorizationHandler(credentialdb, tokendb)
	authHandler(c, r)
}

func UpdateTokenHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 实现修改令牌的具体逻辑
}

func DeleteTokenHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 实现删除令牌的具体逻辑
}

func GetTokensHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 实现显示令牌的具体逻辑
	var listtokensHandler = ListTokensHandler(credentialdb, tokendb)
	listtokensHandler(c, r)
}

func UpdateCredentialHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 实现修改密码的具体逻辑
	authHandler := AuthorizationHandler(credentialdb, tokendb)
	authHandler(c, r)
}

func GetCredentialsHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 实现显示用户的具体逻辑
}

func CreateCredentialHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 实现创建用户的具体逻辑
}

func CreateSessionHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 实现创建会话的具体逻辑
}

func UpdateSessionHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 实现修改会话的具体逻辑
}

func DeleteSessionHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 实现删除会话的具体逻辑
}

func GetSessionsHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {

	sessions, err := ReadAllSessions(sessiondb)
	if err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}
	// 实现显示会话的具体逻辑
	// 创建一个TokenInfo结构体
	var credential TokenInfo = TokenInfo{}

	// 将请求参数绑定到TokenInfo结构体中
	err = r.BindJSON(&credential)
	if err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}
	log.Println(credential)
	shouldReturn := Validatepasswordortoken(credential, credentialdb, tokendb, r)
	if shouldReturn {
		log.Println("用户登录失败:")
		return
	}
	log.Println("用户登录成功:")
	r.JSON(
		consts.StatusOK,
		map[string]interface{}{
			"message": "List of Sessions ok",
			"list": slice.Map(sessions, func(session Session) string {
				return session.Name
			}),
			"username": credential.Username,
		},
	)
	// return


}
// RouteConfig 定义了路由的配置信息，包含路径、方法、中间件和头部信息。
// 
// 字段说明：
// - Path: 表示路由的路径，例如 "/api/v1/resource"。
// - Method: 表示 HTTP 请求方法，例如 "GET"、"POST" 等。
// - MiddleWare: 表示与该路由关联的中间件，类型为 HertzMiddleWare。
// - Headers: 表示与该路由关联的 HTTP 头部信息，以键值对的形式存储。
type RouteConfig struct {
	Path       string
	Method     string
	MiddleWare HertzMiddleWare
	Headers    map[string]string
}