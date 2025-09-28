package routes

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	randv2 "math/rand/v2"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"

	password_hashed "github.com/masx200/go_ws_sh/password-hashed"
)

type TokenStore struct {
	Identifier string         `json:"identifier" gorm:"primarykey;unique;index;not null"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	Hash      string `json:"hash" gorm:"index;not null"`
	Salt      string `json:"salt" gorm:"index;not null"`
	Algorithm string `json:"algorithm" gorm:"index;not null"`

	Username    string `json:"username" gorm:"index;not null"`
	Description string `json:"description" gorm:"index;not null"`
}

func (TokenStore) TableName() string {
	return strings.ToLower("TokenStore")
}

func CreateTokenHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB) func(c context.Context, r *app.RequestContext, next HertzNext) {
	return func(c context.Context, r *app.RequestContext, next HertzNext) {
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

		numBytes := 120
		hexString, err := generateHexKey(numBytes)
		if err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}

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

		id := node.Generate()
		Identifier = id.String()
		newToken := TokenStore{
			Description: req.Token.Description,
			Hash:        hashresult.Hash,
			Salt:        hashresult.Salt,
			Algorithm:   "SHA-512",
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
			"message":  "Login successful",
			"username": username,
		})
	}
}

func UpdateTokenHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB) func(c context.Context, r *app.RequestContext, next HertzNext) {
	return func(c context.Context, r *app.RequestContext, next HertzNext) {
		var body struct {
			Token struct {
				Identifier  string `json:"identifier"`
				Description string `json:"description"`
				Username    string `json:"username"`
			} `json:"token"`
			Authorization CredentialsClient `json:"authorization"`
		}

		if err := r.BindJSON(&body); err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}

		if body.Token.Identifier == "" {
			r.AbortWithMsg("Error: Identifier is empty", consts.StatusBadRequest)
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

		var token TokenStore
		if err := tokendb.Where(&TokenStore{Identifier: body.Token.Identifier}).First(&token).Error; err != nil {
			r.JSON(consts.StatusNotFound, map[string]any{
				"message":  "Error: Token not found",
				"username": username,
				"token": map[string]string{
					"identifier": body.Token.Identifier,
					"username":   username,
				},
			})
			return
		}

		token.Description = body.Token.Description
		token.Username = body.Token.Username
		if err := tokendb.Save(&token).Error; err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}

		r.JSON(consts.StatusOK, map[string]any{
			"message":  "Token updated successfully",
			"username": username,
			"token": map[string]string{
				"identifier":  body.Token.Identifier,
				"description": body.Token.Description,
				"username":    username,
			},
		})
	}
}

func DeleteTokenHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB) func(c context.Context, r *app.RequestContext, next HertzNext) {
	return func(c context.Context, r *app.RequestContext, next HertzNext) {
		var body struct {
			Token struct {
				Identifier string `json:"identifier"`
			} `json:"token"`
			Authorization CredentialsClient `json:"authorization"`
		}

		if err := r.BindJSON(&body); err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}

		if body.Token.Identifier == "" {
			r.AbortWithMsg("Error: Identifier is empty", consts.StatusBadRequest)
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
		}

		if err := tokendb.Where("identifier = ?", body.Token.Identifier).Delete(&TokenStore{}).Error; err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}

		r.JSON(consts.StatusOK, map[string]any{
			"message":  "Token deleted successfully",
			"username": username,
		})
	}
}

func GetTokensHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB) func(c context.Context, r *app.RequestContext, next HertzNext) {
	return func(c context.Context, r *app.RequestContext, next HertzNext) {
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

		var tokenList []map[string]string
		for _, token := range tokens {
			tokenList = append(tokenList, map[string]string{
				"identifier":  token.Identifier,
				"username":    token.Username,
				"created_at":  FormatTimeWithCarbon(carbon.CreateFromStdTime(token.CreatedAt)),
				"updated_at":  FormatTimeWithCarbon(carbon.CreateFromStdTime(token.UpdatedAt)),
				"description": token.Description,
			})
		}

		var username string
		var err error
		if body.Authorization.Username != "" {
			username = body.Authorization.Username
		} else {
			username, err = GetUsernameByTokenIdentifier(tokendb, body.Authorization.Identifier)
			if err != nil {
				log.Println("Error:", err)
				r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
				return
			}
			log.Println("Username:", username)
		}

		r.JSON(consts.StatusOK, map[string]interface{}{
			"tokens":   tokenList,
			"username": username,
			"message":  "Tokens listed successfully",
		})
	}
}

func GenerateTokenRoutes(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB) []RouteConfig {
	return []RouteConfig{
		{
			Headers:    map[string]string{"x-HTTP-method-override": "POST"},
			Path:       "/tokens",
			Method:     "POST",
			MiddleWare: CreateTokenHandler(credentialdb, tokendb, sessiondb),
		},
		{
			Path:       "/tokens",
			Method:     "PUT",
			MiddleWare: UpdateTokenHandler(credentialdb, tokendb, sessiondb),
		},
		{
			Path:       "/tokens",
			Method:     "DELETE",
			MiddleWare: DeleteTokenHandler(credentialdb, tokendb, sessiondb),
		},
		{
			Headers:    map[string]string{"x-HTTP-method-override": "GET"},
			Path:       "/tokens",
			Method:     "POST",
			MiddleWare: GetTokensHandler(credentialdb, tokendb, sessiondb),
		},
	}
}

func GetUsernameByTokenIdentifier(tokendb *gorm.DB, identifier string) (string, error) {
	var token TokenStore

	if err := tokendb.Where("identifier = ?", identifier).First(&token).Error; err != nil {
		return "", fmt.Errorf("token with identifier '%s' not found", identifier)
	}

	return token.Username, nil
}

func FormatTimeWithCarbon(t carbon.Carbon) string {
	return t.Format("Y年m月d日+H时i分s秒T时区")
}

func generateHexKey(numBytes int) (string, error) {
	// 创建一个字节数组来保存随机字节
	randomBytes := make([]byte, numBytes)

	// 使用crypto/rand读取随机字节
	n, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	if n != numBytes {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	// 将字节切片转换为16进制字符串
	hexString := hex.EncodeToString(randomBytes)

	return hexString, nil
}
