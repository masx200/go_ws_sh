package go_ws_sh

import (
	"context"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/golang-module/carbon/v2"
	password_hashed "github.com/masx200/go_ws_sh/password-hashed"
	"gorm.io/gorm"
)

func FormatTimeWithCarbon(t carbon.Carbon) string {
	return t.Format("Y年m月d日+H时i分s秒T时区")
}

// 新增删除用户处理函数声明
func DeleteCredentialHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
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
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
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

// 以下是示例处理函数，需要根据实际业务逻辑实现
func CreateTokenHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 实现创建令牌的具体逻辑
	authHandler := AuthorizationHandler(credentialdb, tokendb)
	authHandler(c, r)
}

func UpdateTokenHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 定义请求体结构体
	var req struct {
		Token struct {
			Identifier  string `json:"identifier"`
			Description string `json:"description"`
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
			return
		}
		log.Println("Username:", username)
	}

	// 查询要更新的令牌
	var token TokenStore
	if err := tokendb.Where(&TokenStore{Identifier: req.Token.Identifier}).First(&token).Error; err != nil {
		r.JSON(consts.StatusNotFound, map[string]any{
			"message":  "Error: Token not found",
			"username": username,
			"token": map[string]string{
				"identifier": req.Token.Identifier,
				"username":   username,
			},
		})
		return
	}

	// 更新令牌信息
	token.Description = req.Token.Description
	if err := tokendb.Save(&token).Error; err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	// 返回成功响应
	r.JSON(consts.StatusOK, map[string]any{
		"message":  "Token updated successfully",
		"username": username,
		"token": map[string]string{
			"identifier":  req.Token.Identifier,
			"description": req.Token.Description,
			"username":    username,
		},
	})
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

	ModifyPassword(r, credentialdb, tokendb)
}

func GetCredentialsHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 创建一个TokenInfo结构体，用于接收认证信息
	var body struct {
		Authorization CredentialsClient `json:"authorization"`
		Credential    struct {
			Username string `json:"username"`
		} `json:"credential"`
	}

	// 将请求参数绑定到TokenInfo结构体中
	err := r.BindJSON(&body)
	if err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	// 验证身份
	shouldReturn := Validatepasswordortoken(body.Authorization, credentialdb, tokendb, r)
	if shouldReturn {
		log.Println("用户登录失败:")
		return
	}

	// 查询所有用户的认证信息
	var credentials []CredentialStore

	if body.Credential.Username != "" {
		// 执行查询
		if err := credentialdb.Where("username =?", body.Credential.Username).Find(&credentials).Error; err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}
	} else {
		if err := credentialdb.Find(&credentials).Error; err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}
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

	username := body.Authorization.Username
	if username == "" {

		username, err = GetUsernameByTokenIdentifier(tokendb, body.Authorization.Identifier)
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
	var req struct {
		Authorization CredentialsClient
		Credential    struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"credential"`
	}

	if err := r.BindJSON(&req); err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}
	//检查NewPassword不为空
	if req.Credential.Password == "" || req.Credential.Username == "" {
		log.Println("Error: New password is empty or username is empty")
		r.AbortWithMsg("Error: New password is empty or username is empty", consts.StatusBadRequest)
		return
	}

	if IsUserExists(credentialdb, req.Credential.Username) {
		log.Println("Error: User already exists")
		r.AbortWithMsg("Error: User already exists", consts.StatusBadRequest)
		return
	}

	if credentialdb.Unscoped().Where("username = ?", req.Credential.Username).Delete(&CredentialStore{}).Error != nil {
		log.Println("Error: Failed to delete user")
		r.AbortWithMsg("Error: Failed to delete user", consts.StatusInternalServerError)
		return
	}
	var cred CredentialStore
	// if err := credentialdb.Where("username = ?", req.Authorization.Username).First(&cred).Error; err != nil {
	// 	r.AbortWithMsg("Error: Invalid credentials", consts.StatusUnauthorized)
	// 	return
	// }
	var reqcre CredentialsClient = req.Authorization
	// 验证旧密码
	// 假设已经有一个函数 ValidatePassword 用于验证密码
	shouldReturn := Validatepasswordortoken(reqcre, credentialdb, tokendb, r)
	if shouldReturn {
		return
	}
	// 更新密码
	newHashresult, err := password_hashed.HashPasswordWithSalt(req.Credential.Password, password_hashed.Options{Algorithm: "SHA-512"})

	if err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}
	cred.Hash = newHashresult.Hash
	cred.Salt = newHashresult.Salt
	cred.Algorithm = "SHA-512" // 假设使用 SHA-512 算法
	// credentialdb.Update()

	if err := credentialdb.Create(&CredentialStore{
		Username:  req.Credential.Username,
		Hash:      cred.Hash,
		Salt:      cred.Salt,
		Algorithm: cred.Algorithm,
	}).Error; err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

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

	r.JSON(consts.StatusOK, map[string]any{"message": "username and Password create successfully",

		"username": username,
		"credential": map[string]string{
			"username": req.Credential.Username,
		},
	})
}

func CreateSessionHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 定义请求体结构体
	var req struct {
		Session struct {
			Name string   `json:"name"`
			Cmd  string   `json:"cmd"`
			Args []string `json:"args"`
			Dir  string   `json:"dir"`
		} `json:"session"`
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

	// 检查 Name 是否为空
	if req.Session.Name == "" || req.Session.Cmd == "" || req.Session.Dir == "" {
		r.AbortWithMsg("Error: Name is empty or  Cmd or Dir is empty ", consts.StatusBadRequest)
		return
	}

	var err error
	username := req.Authorization.Username
	if username == "" {
		username, err = GetUsernameByTokenIdentifier(tokendb, req.Authorization.Identifier)
		if err != nil {
			log.Println("Error:", err)
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}
		log.Println("Username:", username)
	}
	// if IsSessionExists(sessiondb, req.Session.Name) {
	// 	log.Println("Error: Session already exists")
	// 	r.AbortWithMsg("Error: Session already exists", consts.StatusBadRequest)
	// 	return
	// }
	// 检查会话是否已存在
	var existingSession SessionStore
	if err := sessiondb.Where(&SessionStore{Name: req.Session.Name}).First(&existingSession).Error; err == nil {
		r.JSON(consts.StatusConflict, map[string]any{
			"message":  "Error: Session already exists",
			"username": username,
			"session": map[string]string{
				"name": req.Session.Name,
				// "username": username,
			},
		})
		return
	}
	if sessiondb.Unscoped().Where("name = ?", req.Session.Name).Delete(&SessionStore{}).Error != nil {
		log.Println("Error: Failed to delete session")
		r.AbortWithMsg("Error: Failed to delete session", consts.StatusInternalServerError)
		return
	}
	// 创建新的会话
	argsstringarray := StringSlice(req.Session.Args)
	var argsstring string
	argsbytes, err := argsstringarray.Value()
	if err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}
	argsstring = string(argsbytes)

	newSession := SessionStore{
		Name: req.Session.Name,
		Cmd:  req.Session.Cmd,
		Args: argsstring,
		Dir:  req.Session.Dir,
	}

	if err := sessiondb.Create(&newSession).Error; err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	// 返回成功响应
	r.JSON(consts.StatusOK, map[string]any{
		"message":  "Session created successfully",
		"username": username,
		"session": map[string]interface{}{
			"name":     req.Session.Name,
			"cmd":      req.Session.Cmd,
			"args":     req.Session.Args,
			"dir":      req.Session.Dir,
			"username": username,
		},
	})
}

func UpdateSessionHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 定义请求体结构体
	var req struct {
		Session struct {
			Name string   `json:"name"`
			Cmd  string   `json:"cmd"`
			Args []string `json:"args"`
			Dir  string   `json:"dir"`
		} `json:"session"`
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

	// 检查 Name 是否为空
	if req.Session.Name == "" {
		r.AbortWithMsg("Error: Name is empty", consts.StatusBadRequest)
		return
	}

	var err error
	username := req.Authorization.Username
	if username == "" {
		username, err = GetUsernameByTokenIdentifier(tokendb, req.Authorization.Identifier)
		if err != nil {
			log.Println("Error:", err)
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}
		log.Println("Username:", username)
	}

	// 查询要更新的会话
	var session SessionStore
	if err := sessiondb.Where(&SessionStore{Name: req.Session.Name}).First(&session).Error; err != nil {
		r.JSON(consts.StatusNotFound, map[string]any{
			"message":  "Error: Session not found",
			"username": username,
			"session": map[string]string{
				"name": req.Session.Name,
				// "username": username,
			},
		})
		return
	}

	// 更新会话信息
	session.Cmd = req.Session.Cmd
	argsstringarray := StringSlice(req.Session.Args)
	var argsstring string

	argsbytes, err := argsstringarray.Value()
	if err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}
	argsstring = string(argsbytes)
	session.Args = argsstring
	session.Dir = req.Session.Dir
	if err := sessiondb.Save(&session).Error; err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	// 返回成功响应
	r.JSON(consts.StatusOK, map[string]any{
		"message":  "Session updated successfully",
		"username": username,
		"session": map[string]interface{}{
			"name":     req.Session.Name,
			"cmd":      req.Session.Cmd,
			"args":     req.Session.Args,
			"dir":      req.Session.Dir,
			"username": username,
		},
	})
}

func DeleteSessionHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {
	// 定义请求体结构体
	var req struct {
		Session struct {
			Name string `json:"name"`
		} `json:"session"`
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

	// 检查 Name 是否为空
	if req.Session.Name == "" {
		r.AbortWithMsg("Error: Name is empty", consts.StatusBadRequest)
		return
	}

	var err error
	username := req.Authorization.Username
	if username == "" {
		username, err = GetUsernameByTokenIdentifier(tokendb, req.Authorization.Identifier)
		if err != nil {
			log.Println("Error:", err)
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}
		log.Println("Username:", username)
	}

	// 查询要删除的会话
	var session SessionStore
	if err := sessiondb.Where(&SessionStore{Name: req.Session.Name}).First(&session).Error; err != nil {
		r.JSON(consts.StatusOK, map[string]any{
			"message":  "Error: Session not found",
			"username": username,
			"session": map[string]string{
				"name": req.Session.Name,
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
			"name": req.Session.Name,
			// "username": username,
		},
	})
}

func GetSessionsHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, c context.Context, r *app.RequestContext) {

	// 实现显示会话的具体逻辑
	// 创建一个TokenInfo结构体
	var body struct {
		Authorization CredentialsClient `json:"authorization"`
		Session       struct {
			Name string `json:"name"`
		} `json:"session"`
	}

	// 将请求参数绑定到TokenInfo结构体中
	err := r.BindJSON(&body)
	if err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}
	log.Println(body)

	shouldReturn := Validatepasswordortoken(body.Authorization, credentialdb, tokendb, r)
	if shouldReturn {
		log.Println("用户登录失败:")
		return
	}
	log.Println("用户登录成功:")

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
	var sessions []Session
	if body.Session.Name != "" {
		sessions, err = ReadAllSessionsWithName(sessiondb, body.Session.Name)
		if err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}
	} else {
		sessions, err = ReadAllSessions(sessiondb)
		if err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}
	}

	r.JSON(
		consts.StatusOK,
		map[string]interface{}{
			"message":  "List of Sessions ok",
			"sessions": SessionsToMapSlice(sessions),
			"username": username,
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
