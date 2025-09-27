package go_ws_sh

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"
	"log"
	"github.com/masx200/go_ws_sh/types"
)


func Validatepasswordortoken(req types.CredentialsClient, credentialdb *gorm.DB, tokendb *gorm.DB, r *app.RequestContext) bool {
	if req.Type == "token" && req.Token != "" && req.Identifier != "" {
		log.Println("开始Token 认证")
		// Token 认证
		if ok, err := ValidateToken(req, tokendb); !ok {
			log.Println(err)
			r.AbortWithMsg("Error: Invalid credentials", consts.StatusUnauthorized)
			log.Println("Error: Invalid credentials")
			return true
		}
		log.Println("success: success credentials")
		return false
	}
	log.Println("开始password 认证")
	// 用户名密码认证
	var cred CredentialStore
	if err := credentialdb.Where("username = ?", req.Username).First(&cred).Error; err != nil {
		r.AbortWithMsg("Error: Invalid credentials", consts.StatusUnauthorized)
		return true
	}
	//用户名和密码都不为空
	if req.Username == "" || req.Password == "" {
		r.AbortWithMsg("Error: Username or password is empty", consts.StatusBadRequest)
		return true
	}
	// 验证密码
	// 这里需要实现具体的密码验证逻辑
	// 假设已经有一个函数 ValidatePassword 用于验证密码
	if ok, err := ValidatePassword(req.Password, cred.Hash, cred.Salt, cred.Algorithm); !ok {
		log.Println(err)
		r.AbortWithMsg("Error: Invalid credentials", consts.StatusUnauthorized)
		return true
	}
	return false
}
