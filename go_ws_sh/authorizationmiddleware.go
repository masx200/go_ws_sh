package go_ws_sh

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"

	"github.com/masx200/go_ws_sh/types"
)

// AuthorizationMiddleware 定义身份验证中间件
func AuthorizationMiddleware(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB) types.HertzMiddleWare {
	return func(c context.Context, r *app.RequestContext, next types.HertzNext) {

		bearertoken := r.Request.Header.Get("authorization")

		if bearertoken != "" {
			bearertoken=strings.TrimPrefix(bearertoken, "Bearer ")
			decoded, err := base64.StdEncoding.DecodeString(bearertoken)
			if err != nil {
				r.AbortWithStatusJSON(consts.StatusUnauthorized, map[string]string{
					"message": "Error: Invalid token",
				})
				return
			}
			var cc types.CredentialsClient
			if err := json.Unmarshal(decoded, &cc); err != nil {
				r.AbortWithStatusJSON(consts.StatusUnauthorized, map[string]string{
					"message": "Error: Invalid token",
				})
				return
			}
			validateFailure := Validatepasswordortoken(cc, credentialdb, tokendb, r)
			if validateFailure {
				return
			}

			// 验证成功，调用下一个处理函数
			next(c, r)
			return
		}
		var req struct {
			Authorization types.CredentialsClient `json:"authorization"`
		}

		// 解析请求体中的 JSON 数据
		if err := r.BindJSON(&req); err != nil {
			r.AbortWithStatusJSON(consts.StatusUnauthorized, map[string]string{
				"message": "Error: Invalid request body",
			})
			return
		}

		// 调用 Validatepasswordortoken 函数进行身份验证
		validateFailure := Validatepasswordortoken(req.Authorization, credentialdb, tokendb, r)
		if validateFailure {
			return
		}

		// 验证成功，调用下一个处理函数
		next(c, r)
	}
}
