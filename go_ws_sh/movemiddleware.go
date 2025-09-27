package go_ws_sh

import (
	"context"
	"encoding/json"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"

	"github.com/masx200/go_ws_sh/types"
)


func MoveMiddleware(initial_sessions []types.Session, credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext, next types.HertzNext) {
	// 定义请求体结构体
	var body struct {
		Session struct {
			Name string `json:"name"`
		} `json:"session"`
		Authorization types.CredentialsClient `json:"authorization"`
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

}
