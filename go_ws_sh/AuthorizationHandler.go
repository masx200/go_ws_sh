package go_ws_sh

import (
	"context"
	"fmt"
	"log"
	randv2 "math/rand/v2"

	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"

	password_hashed "github.com/masx200/go_ws_sh/password-hashed"
)

// ListTokensHandler 列出所有令牌
func ListTokensHandler(credentialdb *gorm.DB, tokendb *gorm.DB) func(w context.Context, r *app.RequestContext) {
	return func(w context.Context, r *app.RequestContext) {
		var req CredentialsClient

		if err := r.BindJSON(&req); err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}

		shouldReturn := Validatepasswordortoken(req, credentialdb, tokendb, r)
		if shouldReturn {
			return
		}

		// 查询所有令牌
		var tokens []Tokens
		if err := tokendb.Find(&tokens).Error; err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}

		// 构建响应数据
		var tokenList []map[string]string
		for _, token := range tokens {
			tokenList = append(tokenList, map[string]string{
				"identifier": token.Identifier,
				"username":   token.Username,
				"algorithm":  token.Algorithm,
				// 注意：不建议返回哈希和盐值，这里仅为示例
				// "hash": token.Hash,
				// "salt": token.Salt,
			})
		}

		r.JSON(consts.StatusOK, map[string]interface{}{
			"tokens":  tokenList,
			"message": "Tokens listed successfully",
		})
	}
}

// AuthorizationHandler 处理授权相关的请求
func AuthorizationHandler(credentialdb *gorm.DB, tokendb *gorm.DB) func(w context.Context, r *app.RequestContext) {
	return func(w context.Context, r *app.RequestContext) {
		// 获取请求方法
		method := r.Method()

		switch string(method) {
		case consts.MethodPost:
			handlePost(r, credentialdb, tokendb)
		case consts.MethodPut:
			handlePut(r, credentialdb, tokendb)
		case consts.MethodDelete:
			handleDelete(r, credentialdb, tokendb)
		default:
			r.AbortWithMsg("Method Not Allowed", consts.StatusMethodNotAllowed)
		}
	}
}
func ValidateToken(req CredentialsClient, tokendb *gorm.DB) (bool, error) {
	var token Tokens
	if err := tokendb.Where(&Tokens{Identifier: req.Identifier,
		Username: req.Username,
	}).First(&token).Error; err != nil {

		return false, err
	}

	var hashresult, err = password_hashed.HashPasswordWithSalt(req.Token, password_hashed.Options{Algorithm: token.Algorithm,
		SaltHex: token.Salt,
	})
	if err != nil {
		return false, err
	}
	var hash = hashresult.Hash
	if hash != token.Hash {
		return false, fmt.Errorf("token is invalid")
	}
	return true, nil
}

// handlePost 处理 POST 请求，支持用户名密码认证和创建新的 Token
func handlePost(r *app.RequestContext, credentialdb *gorm.DB, tokendb *gorm.DB) {
	var req CredentialsClient

	if err := r.BindJSON(&req); err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	shouldReturn := Validatepasswordortoken(req, credentialdb, tokendb, r)
	if shouldReturn {
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
		fmt.Println(err)
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	// Generate a snowflake ID.
	id := node.Generate()
	Identifier = id.String()
	newToken := Tokens{
		Hash:       hashresult.Hash,
		Salt:       hashresult.Salt,
		Algorithm:  "SHA-512", // 假设使用 SHA-512 算法
		Identifier: Identifier,
		Username:   req.Username,
	}
	if err := tokendb.Create(&newToken).Error; err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}
	r.JSON(consts.StatusOK, map[string]string{"token": hexString, "message": "Login successful",

		"identifier": Identifier, "username": req.Username, "type": "token",
	})
}

func Validatepasswordortoken(req CredentialsClient, credentialdb *gorm.DB, tokendb *gorm.DB, r *app.RequestContext) bool {
	if req.Type == "token" && req.Token != "" && req.Identifier != "" {
		// Token 认证
		if ok, err := ValidateToken(req, tokendb); !ok {
			log.Println(err)
			r.AbortWithMsg("Error: Invalid credentials", consts.StatusUnauthorized)
			return true
		}
	}

	// 用户名密码认证
	var cred Credentials
	if err := credentialdb.Where("username = ?", req.Username).First(&cred).Error; err != nil {
		r.AbortWithMsg("Error: Invalid credentials", consts.StatusUnauthorized)
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
		return false, fmt.Errorf("token is invalid")
	}
	return true, nil
}

// handlePut 处理 PUT 请求，修改用户名密码
func handlePut(r *app.RequestContext, credentialdb *gorm.DB, tokendb *gorm.DB) {
	var req struct {
		CredentialsClient
		Username    string `json:"username"`
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	}

	if err := r.BindJSON(&req); err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}
	//检查NewPassword不为空
	if req.NewPassword == "" {

		r.AbortWithMsg("Error: New password is empty", consts.StatusBadRequest)
	}
	var cred Credentials
	if err := credentialdb.Where("username = ?", req.Username).First(&cred).Error; err != nil {
		r.AbortWithMsg("Error: Invalid credentials", consts.StatusUnauthorized)
		return
	}
	var reqcre CredentialsClient = CredentialsClient{
		Username:   req.Username,
		Password:   req.Password,
		Type:       req.Type,
		Token:      req.Token,
		Identifier: req.Identifier,
	}
	// 验证旧密码
	// 假设已经有一个函数 ValidatePassword 用于验证密码
	shouldReturn := Validatepasswordortoken(reqcre, credentialdb, tokendb, r)
	if shouldReturn {
		return
	}
	// 更新密码
	newHashresult, err := password_hashed.HashPasswordWithSalt(req.NewPassword, password_hashed.Options{Algorithm: "SHA-512"})

	if err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}
	cred.Hash = newHashresult.Hash
	cred.Salt = newHashresult.Salt
	cred.Algorithm = "SHA-512" // 假设使用 SHA-512 算法
	if err := credentialdb.Save(&cred).Error; err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}
	r.JSON(consts.StatusOK, map[string]string{"message": "Password updated successfully", "username": req.Username})
}

// handleDelete 处理 DELETE 请求，删除某个 Token
func handleDelete(r *app.RequestContext, credentialdb *gorm.DB, tokendb *gorm.DB) {
	var req CredentialsClient

	if err := r.BindJSON(&req); err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	shouldReturn := Validatepasswordortoken(req, credentialdb, tokendb, r)
	if shouldReturn {
		return
	}
	var data map[string]interface{}
	if err := r.BindJSON(&data); err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}
	//检查Identifier不为空

	if data["delete_identifier"] == nil {
		r.AbortWithMsg("Error: Identifier is empty", consts.StatusBadRequest)
		return
	}
	if data["delete_identifier"].(string) == "" {
		r.AbortWithMsg("Error: Identifier is empty", consts.StatusBadRequest)
		return
	}
	var token Tokens
	if err := tokendb.Where(&Tokens{Identifier: data["delete_identifier"].(string), Username: req.Username}).First(&token).Error; err != nil {
		r.JSON(consts.StatusOK, map[string]string{
			"message":           "Error: Token not found",
			"username":          req.Username,
			"delete_identifier": data["delete_identifier"].(string),
		})

		//r.AbortWithMsg("Error: Token not found", consts.StatusNotFound)
		return
	}

	if err := tokendb.Delete(&token).Error; err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}
	r.JSON(consts.StatusOK, map[string]string{"message": "Token deleted successfully",
		"username":          req.Username,
		"delete_identifier": data["delete_identifier"].(string),
	})
}
