package go_ws_sh

import (
	"context"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"
)

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
	validateFailure := Validatepasswordortoken(req.Authorization, credentialdb, tokendb, r)
	if validateFailure {
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
