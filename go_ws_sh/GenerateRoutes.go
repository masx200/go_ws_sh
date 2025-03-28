package go_ws_sh

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

// GenerateRoutesHttp 根据 openapi 文件生成路由配置
func GenerateRoutesHttp(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB) []RouteConfig {
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
