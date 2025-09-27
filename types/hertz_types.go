package types

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// HertzNext 定义了Hertz中间件中的下一个处理函数类型
type HertzNext func(c context.Context, r *app.RequestContext)

// HertzMiddleWare 定义了Hertz中间件的函数类型
type HertzMiddleWare func(c context.Context, r *app.RequestContext, next HertzNext)

// CredentialsClient 定义了客户端凭据的结构
type CredentialsClient struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Token      string `json:"token"`
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

// RouteConfig 定义了路由的配置信息，包含路径、方法、中间件和头部信息
type RouteConfig struct {
	Path       string
	Method     string
	MiddleWare HertzMiddleWare
	Headers    map[string]string
}

// InitialCredentials 定义了初始凭据的结构
type InitialCredentials []struct {
	Username string `json:"username"`
	Password string `json:"password"`
}