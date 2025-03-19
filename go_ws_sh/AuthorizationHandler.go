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
	if err := tokendb.Where(&Tokens{Identifier: req.Identifier}).First(&token).Error; err != nil {

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

	if req.Type == "token" && req.Token != "" && req.Identifier != "" {
		// Token 认证
		if ok, err := ValidateToken(req, tokendb); !ok {
			log.Println(err)
			r.AbortWithMsg("Error: Invalid credentials", consts.StatusUnauthorized)
			return
		}
	}

	// 用户名密码认证
	var cred Credentials
	if err := tokendb.Where("username = ?", req.Username).First(&cred).Error; err != nil {
		r.AbortWithMsg("Error: Invalid credentials", consts.StatusUnauthorized)
		return
	}

	// 验证密码
	// 这里需要实现具体的密码验证逻辑
	// 假设已经有一个函数 ValidatePassword 用于验证密码
	if ok, err := ValidatePassword(req.Password, cred.Hash, cred.Salt, cred.Algorithm); !ok {
		log.Println(err)
		r.AbortWithMsg("Error: Invalid credentials", consts.StatusUnauthorized)
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

	node, err := snowflake.NewNode(randv2.Int64())
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
	r.JSON(consts.StatusOK, map[string]string{"token": hexString, "message": "Login successful"})
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
		Username    string `json:"username"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := r.BindJSON(&req); err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	var cred Credentials
	if err := tokendb.Where("username = ?", req.Username).First(&cred).Error; err != nil {
		r.AbortWithMsg("Error: Invalid credentials", consts.StatusUnauthorized)
		return
	}

	// 验证旧密码
	// 假设已经有一个函数 ValidatePassword 用于验证密码
	if !ValidatePassword(req.OldPassword, cred.Hash, cred.Salt, cred.Algorithm) {
		r.AbortWithMsg("Error: Invalid credentials", consts.StatusUnauthorized)
		return
	}

	// 更新密码
	newHash, newSalt := HashPasswordWithSalt(req.NewPassword, cred.Salt)
	cred.Hash = newHash
	cred.Salt = newSalt
	if err := tokendb.Save(&cred).Error; err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}
	r.JSON(consts.StatusOK, map[string]string{"message": "Password updated successfully"})
}

// handleDelete 处理 DELETE 请求，删除某个 Token
func handleDelete(r *app.RequestContext, credentialdb *gorm.DB, tokendb *gorm.DB) {
	var req struct {
		Token string `json:"token"`
	}

	if err := r.BindJSON(&req); err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}

	var token Tokens
	if err := tokendb.Where("hash = ?", req.Token).First(&token).Error; err != nil {
		r.AbortWithMsg("Error: Token not found", consts.StatusNotFound)
		return
	}

	if err := tokendb.Delete(&token).Error; err != nil {
		r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		return
	}
	r.JSON(consts.StatusOK, map[string]string{"message": "Token deleted successfully"})
}
