package go_ws_sh

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"
	"log"
)

// 新增删除用户处理函数声明
func DeleteCredentialHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext, initial_credentials InitialCredentials) {
	// 定义请求体结构体
	var req struct {
		Authorization CredentialsClient `json:"authorization"`
		Credential    struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"credential"`
	}

	// 绑定请求体
	if err := r.BindJSON(&req); err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusBadRequest)
		return
	}
	//检查是否是初始用户

	for _, ic := range initial_credentials {

		if req.Authorization.Username == ic.Username {
			r.AbortWithMsg("Error: 初始用户不允许删除", consts.StatusBadRequest)
			return
		}
	}
	// 验证身份
	shouldReturn := Validatepasswordortoken(req.Authorization, credentialdb, tokendb, r)
	if shouldReturn {
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
	// 检查要删除的用户是否存在
	var cred CredentialStore
	if err := credentialdb.Where("username = ?", req.Credential.Username).First(&cred).Error; err != nil {
		r.JSON(consts.StatusOK, map[string]any{
			"message":  "Error: User not found",
			"username": username,
			"credential": map[string]string{
				"username": req.Credential.Username,
			},
		})
		return
	}
	// 计算用户数量,如果
	// 计算用户数量,用户数量小于等于1,则不允许删除
	var count int64
	if err := credentialdb.Model(&CredentialStore{}).Count(&count).Error; err != nil || count <= 1 {
		log.Println("Error:", err)
		log.Println("count:", count)
		r.JSON(consts.StatusBadRequest, map[string]any{
			"message":  "用户数量小于等于1,则不允许删除",
			"username": username,
			"credential": map[string]string{
				"username": req.Credential.Username,
			},
		})
		return
	}

	// 删除用户
	if err := credentialdb.Delete(&cred).Error; err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	// 返回成功响应
	r.JSON(consts.StatusOK, map[string]any{
		"message":  "User deleted successfully",
		"username": username,
		"credential": map[string]string{
			"username": req.Credential.Username,
		},
	})
}
