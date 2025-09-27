package go_ws_sh

import (
	"context"
	"fmt"
	"log"
	randv2 "math/rand/v2"

	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"

	password_hashed "github.com/masx200/go_ws_sh/password-hashed"
)

// ListTokensHandler 列出所有令牌
func ListTokensHandler(credentialdb *gorm.DB, tokendb *gorm.DB) func(w context.Context, r *app.RequestContext) {
	return func(w context.Context, r *app.RequestContext) {
		tokendb = tokendb.Debug()
		var body struct {
			Authorization CredentialsClient `json:"authorization"`
			Token         struct {
				Identifier string `json:"identifier"`
			} `json:"token"`
		}

		if err := r.BindJSON(&body); err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}

		validateFailure := Validatepasswordortoken(body.Authorization, credentialdb, tokendb, r)
		if validateFailure {
			return
		}

		// 查询所有令牌
		var tokens []TokenStore
		if body.Token.Identifier != "" {
			log.Println("body.Token.Identifier:", body.Token.Identifier)
			if err := tokendb.Where("identifier =?", body.Token.Identifier).First(&tokens).Error; err != nil {
				r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
				return
			}
		} else {
			if err := tokendb.Find(&tokens).Error; err != nil {
				r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
				return
			}

		}

		// 构建响应数据
		var tokenList []map[string]string
		for _, token := range tokens {
			tokenList = append(tokenList, map[string]string{
				"identifier": token.Identifier,
				"username":   token.Username,
				// "algorithm":  token.Algorithm,
				"created_at":  FormatTimeWithCarbon(carbon.CreateFromStdTime((token.CreatedAt))),
				"updated_at":  FormatTimeWithCarbon(carbon.CreateFromStdTime(token.UpdatedAt)),
				"description": token.Description,
				// 注意：不建议返回哈希和盐值，这里仅为示例
				// "hash": token.Hash,
				// "salt": token.Salt,
			})
		}

		if body.Authorization.Username != "" {
			username := body.Authorization.Username
			r.JSON(consts.StatusOK, map[string]interface{}{
				"tokens":   tokenList,
				"username": username,
				"message":  "Tokens listed successfully",
			})
			return
		}
		username, err := GetUsernameByTokenIdentifier(tokendb, body.Authorization.Identifier)
		if err != nil {
			log.Println("Error:", err)
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		} else {
			log.Println("Username:", username)
		}
		r.JSON(consts.StatusOK, map[string]interface{}{
			"tokens":   tokenList,
			"username": username,
			"message":  "Tokens listed successfully",
		})
	}
}
func GetUsernameByTokenIdentifier(tokendb *gorm.DB, identifier string) (string, error) {
	var token TokenStore

	// 查询指定 identifier 的条目
	if err := tokendb.Where("identifier = ?", identifier).First(&token).Error; err != nil {
		// 如果没有找到对应的条目，返回错误
		return "", fmt.Errorf("token with identifier '%s' not found", identifier)
	}

	// 返回查询到的 username
	return token.Username, nil
}

// AuthorizationHandler 处理授权相关的请求
func AuthorizationHandler(credentialdb *gorm.DB, tokendb *gorm.DB) func(w context.Context, r *app.RequestContext) {
	return func(w context.Context, r *app.RequestContext) {
		// 获取请求方法
		method := r.Method()

		switch string(method) {
		case consts.MethodPost:
			CreateToken(r, credentialdb, tokendb)
		// case consts.MethodPut:
		// 	handlePut(r, credentialdb, tokendb)
		// case consts.MethodDelete:
		// 	handleDelete(r, credentialdb, tokendb)
		default:
			r.AbortWithMsg("Method Not Allowed", consts.StatusMethodNotAllowed)
		}
	}
}
func ValidateToken(reqcredential CredentialsClient, tokendb *gorm.DB) (bool, error) {
	var token TokenStore
	if err := tokendb.Where(&TokenStore{Identifier: reqcredential.Identifier}). // Username: reqcredential.Username,
											First(&token).Error; err != nil {

		return false, err
	}

	var hashresult, err = password_hashed.HashPasswordWithSalt(reqcredential.Token, password_hashed.Options{Algorithm: token.Algorithm,
		SaltHex: token.Salt,
	})
	if err != nil {
		return false, err
	}
	log.Println("hashresult", hashresult)
	log.Println("token", token)
	log.Println("credential", reqcredential)
	var hash = hashresult.Hash
	if hash != token.Hash {
		return false, fmt.Errorf("token is invalid")
	}
	return true, nil
}

// CreateToken 处理 POST 请求，支持用户名密码认证和创建新的 Token
func CreateToken(r *app.RequestContext, credentialdb *gorm.DB, tokendb *gorm.DB) {
	var req struct {
		Authorization CredentialsClient `json:"authorization"`
		Token         struct {
			Username    string `json:"username"`
			Description string `json:"description"`
		} `json:"token"`
	}

	if err := r.BindJSON(&req); err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	validateFailure := Validatepasswordortoken(req.Authorization, credentialdb, tokendb, r)
	if validateFailure {
		return
	}
	numBytes := 120
	hexString, err := generateHexKey(numBytes)
	if err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}
	// 创建新的 Token
	hashresult, err := password_hashed.HashPasswordWithSalt(hexString, password_hashed.Options{Algorithm: "SHA-512"})
	if err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}
	var Identifier string

	node, err := snowflake.NewNode(randv2.Int64() % 1024)
	if err != nil {
		log.Println(err)
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	// Generate a snowflake ID.
	id := node.Generate()
	Identifier = id.String()
	newToken := TokenStore{
		Description: req.Token.Description,
		Hash:        hashresult.Hash,
		Salt:        hashresult.Salt,
		Algorithm:   "SHA-512", // 假设使用 SHA-512 算法
		Identifier:  Identifier,
		Username:    req.Token.Username,
	}
	if err := tokendb.Create(&newToken).Error; err != nil {
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
	r.JSON(consts.StatusOK, map[string]any{
		"token": map[string]string{
			"identifier":  Identifier,
			"username":    username,
			"description": req.Token.Description,
			"token":       hexString,
		},

		"message": "Login successful",

		"username": username,
	})
}

func Validatepasswordortoken(req CredentialsClient, credentialdb *gorm.DB, tokendb *gorm.DB, r *app.RequestContext) bool {
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

func ValidatePassword(Password, Hash, Salt, Algorithm string) (bool, error) {
	var hashresult, err = password_hashed.HashPasswordWithSalt(Password, password_hashed.Options{Algorithm: Algorithm,
		SaltHex: Salt,
	})
	if err != nil {
		return false, err
	}
	var hash = hashresult.Hash
	if hash != Hash {
		return false, fmt.Errorf("password is invalid")
	}
	return true, nil
}

// ModifyPassword 处理 PUT 请求，修改用户名密码
func ModifyPassword(r *app.RequestContext, credentialdb *gorm.DB, tokendb *gorm.DB) {
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

		r.AbortWithMsg("Error: New password is empty", consts.StatusBadRequest)
	}
	var cred CredentialStore
	if err := credentialdb.Where("username = ?", req.Credential.Username).First(&cred).Error; err != nil {
		r.AbortWithMsg("Error: Invalid credentials", consts.StatusUnauthorized)
		return
	}
	var reqcre CredentialsClient = req.Authorization
	// 验证旧密码
	// 假设已经有一个函数 ValidatePassword 用于验证密码
	validateFailure := Validatepasswordortoken(reqcre, credentialdb, tokendb, r)
	if validateFailure {
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

	if err := credentialdb.Model(&CredentialStore{}).Where("username = ?", req.Credential.Username).Updates(CredentialStore{

		Hash:      cred.Hash,
		Salt:      cred.Salt,
		Algorithm: cred.Algorithm,
	}).Error; err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusNotFound)
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

// // handleDelete 处理 DELETE 请求，删除某个 Token
// func handleDelete(r *app.RequestContext, credentialdb *gorm.DB, tokendb *gorm.DB) {
// 	var req CredentialsClient

// 	if err := r.BindJSON(&req); err != nil {
// 		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
// 		return
// 	}

// 	validateFailure := Validatepasswordortoken(req, credentialdb, tokendb, r)
// 	if validateFailure {
// 		return
// 	}
// 	var data map[string]interface{}
// 	if err := r.BindJSON(&data); err != nil {
// 		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
// 		return
// 	}
// 	//检查Identifier不为空

// 	if data["delete_identifier"] == nil {
// 		r.AbortWithMsg("Error: Identifier is empty", consts.StatusBadRequest)
// 		return
// 	}
// 	if data["delete_identifier"].(string) == "" {
// 		r.AbortWithMsg("Error: Identifier is empty", consts.StatusBadRequest)
// 		return
// 	}
// 	var token TokenStore
// 	if err := tokendb.Where(&TokenStore{Identifier: data["delete_identifier"].(string), Username: req.Username}).First(&token).Error; err != nil {
// 		r.JSON(consts.StatusOK, map[string]string{
// 			"message":           "Error: Token not found",
// 			"username":          req.Username,
// 			"delete_identifier": data["delete_identifier"].(string),
// 		})

// 		//r.AbortWithMsg("Error: Token not found", consts.StatusNotFound)
// 		return
// 	}

// 	if err := tokendb.Delete(&token).Error; err != nil {
// 		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
// 		return
// 	}
// 	r.JSON(consts.StatusOK, map[string]string{"message": "Token deleted successfully",
// 		"username":          req.Username,
// 		"delete_identifier": data["delete_identifier"].(string),
// 	})
// }
