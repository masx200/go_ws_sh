package go_ws_sh

import (
	"context"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"
)

func DeleteSessionHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext, initial_sessions []Session) {
	// 定义请求体结构体
	var body struct {
		Session struct {
			Name string `json:"name"`
		} `json:"session"`
		Authorization CredentialsClient `json:"authorization"`
	}

	// 绑定请求体
	if err := r.BindJSON(&body); err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusBadRequest)
		return
	}
	//检查是否为初始会话
	for _, session := range initial_sessions {
		if session.Name == body.Session.Name {
			r.AbortWithMsg("Error: Session is initial session,不允许删除", consts.StatusBadRequest)
			return
		}
	}
	// 验证身份
	validateFailure := Validatepasswordortoken(body.Authorization, credentialdb, tokendb, r)
	if validateFailure {
		return
	}

	// 检查 Name 是否为空
	if body.Session.Name == "" {
		r.AbortWithMsg("Error: Name is empty", consts.StatusBadRequest)
		return
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

	// 查询要删除的会话
	var session SessionStore
	if err := sessiondb.Where(&SessionStore{Name: body.Session.Name}).First(&session).Error; err != nil {
		r.JSON(consts.StatusOK, map[string]any{
			"message":  "Error: Session not found",
			"username": username,
			"session": map[string]string{
				"name": body.Session.Name,
				// "username": username,
			},
		})
		return
	}

	// 删除会话
	if err := sessiondb.Delete(&session).Error; err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	// 返回成功响应
	r.JSON(consts.StatusOK, map[string]any{
		"message":  "Session deleted successfully",
		"username": username,
		"session": map[string]string{
			"name": body.Session.Name,
			// "username": username,
		},
	})
}
