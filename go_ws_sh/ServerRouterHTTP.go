package go_ws_sh

import (
	"context"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

func FormatTimeWithCarbon(t carbon.Carbon) string {
	return t.Format("Y年m月d日+H时i分s秒T时区")
}

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
	// authHandler := AuthorizationHandler(credentialdb, tokendb)
	// authHandler(c, r)
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
	// 定义请求体结构体
	var req struct {
		Token struct {
			Identifier string `json:"identifier"`
			Username   string `json:"username"`
		} `json:"token"`
		Authorization CredentialsClient `json:"authorization"`
	}

	// 绑定请求体
	if err := r.BindJSON(&req); err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	// 验证身份
	shouldReturn := Validatepasswordortoken(req.Authorization, credentialdb, tokendb, r)
	if shouldReturn {
		return
	}
	// log.Println(req)
	// 检查 Identifier 是否为空
	if req.Token.Identifier == "" {
		r.AbortWithMsg("Error: Identifier is empty", consts.StatusBadRequest)
		return
	}
	var err error
	username := req.Authorization.Username
	if username == "" {

		username, err = GetUsernameByTokenIdentifier(tokendb, req.Authorization.Identifier)
		if err != nil {
			log.Println("Error:", err)
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		} else {
			log.Println("Username:", username)
		}

	}
	// 查询要删除的令牌
	var token TokenStore
	if err := tokendb.Where(&TokenStore{Identifier: req.Token.Identifier}).
		// Username: req.Token.Username

		First(&token).Error; err != nil {

		r.JSON(consts.StatusOK, map[string]any{
			"message":  "Error: Token not found",
			"username": username,
			"token": map[string]string{
				"identifier": req.Token.Identifier,
				"username":   username,
			},
		})
		return
	}

	// 删除令牌
	if err := tokendb.Delete(&token).Error; err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	// 返回成功响应
	r.JSON(consts.StatusOK, map[string]any{
		"message":  "Token deleted successfully",
		"username": username,
		"token": map[string]string{
			"identifier": req.Token.Identifier,
			"username":   username,
		},
	})
}

func GetTokensHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 实现显示令牌的具体逻辑
	var listtokensHandler = ListTokensHandler(credentialdb, tokendb)
	listtokensHandler(c, r)
}

func UpdateCredentialHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 实现修改密码的具体逻辑
	// authHandler := AuthorizationHandler(credentialdb, tokendb)
	// authHandler(c, r)

	handlePut(r, credentialdb, tokendb)
}

func GetCredentialsHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 创建一个TokenInfo结构体，用于接收认证信息
	var credential struct {
		Authorization CredentialsClient `json:"authorization"`
	}

	// 将请求参数绑定到TokenInfo结构体中
	err := r.BindJSON(&credential)
	if err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	// 验证身份
	shouldReturn := Validatepasswordortoken(credential.Authorization, credentialdb, tokendb, r)
	if shouldReturn {
		log.Println("用户登录失败:")
		return
	}

	// 查询所有用户的认证信息
	var credentials []CredentialStore
	if err := credentialdb.Find(&credentials).Error; err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	// 构建响应数据
	var credentialList []map[string]string
	for _, cred := range credentials {
		credentialList = append(credentialList, map[string]string{
			"username":   cred.Username,
			"created_at": FormatTimeWithCarbon(carbon.CreateFromStdTime(cred.CreatedAt)),
			"updated_at": FormatTimeWithCarbon(carbon.CreateFromStdTime(cred.UpdatedAt)),
			// 注意：不建议返回哈希和盐值，这里仅为示例
			// "hash": cred.Hash,
			// "salt": cred.Salt,
		})
	}

	username := credential.Authorization.Username
	if username == "" {

		username, err = GetUsernameByTokenIdentifier(tokendb, credential.Authorization.Identifier)
		if err != nil {
			log.Println("Error:", err)
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		} else {
			log.Println("Username:", username)
		}

	}
	// 返回成功响应
	r.JSON(consts.StatusOK, map[string]interface{}{
		"credentials": credentialList,
		"username":    username,
		"message":     "Credentials listed successfully",
	})
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
	var credential struct {
		Authorization CredentialsClient `json:"authorization"`
	}

	// 将请求参数绑定到TokenInfo结构体中
	err = r.BindJSON(&credential)
	if err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}
	log.Println(credential)
	shouldReturn := Validatepasswordortoken(credential.Authorization, credentialdb, tokendb, r)
	if shouldReturn {
		log.Println("用户登录失败:")
		return
	}
	log.Println("用户登录成功:")
	r.JSON(
		consts.StatusOK,
		map[string]interface{}{
			"message":  "List of Sessions ok",
			"sessions": SessionsToMapSlice(sessions),
			"username": credential.Authorization.Username,
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
