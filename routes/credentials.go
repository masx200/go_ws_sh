package routes

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"

	"github.com/masx200/go_ws_sh/go_ws_sh"
	password_hashed "github.com/masx200/go_ws_sh/password-hashed"
)

func (CredentialStore) TableName() string {
	return strings.ToLower("CredentialStore")
}

type CredentialStore struct {
	gorm.Model
	Username  string `json:"username" gorm:"index;unique;not null"`
	Hash      string `json:"hash" gorm:"index;not null"`
	Salt      string `json:"salt" gorm:"index;not null"`
	Algorithm string `json:"algorithm" gorm:"index;not null"`
}

func UpdateCredentialHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB) func(c context.Context, r *app.RequestContext, next HertzNext) {
	return func(c context.Context, r *app.RequestContext, next HertzNext) {
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

		if req.Credential.Password == "" || req.Credential.Username == "" {
			r.AbortWithMsg("Error: New password is empty or username is empty", consts.StatusBadRequest)
			return
		}

		var cred CredentialStore
		if err := credentialdb.Where("username = ?", req.Credential.Username).First(&cred).Error; err != nil {
			r.AbortWithMsg("Error: Invalid credentials", consts.StatusUnauthorized)
			return
		}

		newHashresult, err := password_hashed.HashPasswordWithSalt(req.Credential.Password, password_hashed.Options{Algorithm: "SHA-512"})
		if err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}
		cred.Hash = newHashresult.Hash
		cred.Salt = newHashresult.Salt
		cred.Algorithm = "SHA-512"

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

		r.JSON(consts.StatusOK, map[string]any{"message": "Password updated successfully",
			"username": username,
			"credential": map[string]string{
				"username": req.Credential.Username,
			},
		})
	}
}

func GetCredentialsHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB) func(c context.Context, r *app.RequestContext, next HertzNext) {
	return func(c context.Context, r *app.RequestContext, next HertzNext) {
		credentialdb = credentialdb.Debug()
		var body struct {
			Authorization CredentialsClient `json:"authorization"`
			Credential    struct {
				Username string `json:"username"`
			} `json:"credential"`
		}

		err := r.BindJSON(&body)
		if err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}

		var credentials []CredentialStore

		if body.Credential.Username != "" {
			log.Println("查询用户:", body.Credential.Username)
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

		var credentialList []map[string]string
		for _, cred := range credentials {
			credentialList = append(credentialList, map[string]string{
				"username":   cred.Username,
				"created_at": FormatTimeWithCarbon(carbon.CreateFromStdTime(cred.CreatedAt)),
				"updated_at": FormatTimeWithCarbon(carbon.CreateFromStdTime(cred.UpdatedAt)),
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

		r.JSON(consts.StatusOK, map[string]interface{}{
			"credentials": credentialList,
			"username":    username,
			"message":     "Credentials listed successfully",
		})
	}
}

func CreateCredentialHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB) func(c context.Context, r *app.RequestContext, next HertzNext) {
	return func(c context.Context, r *app.RequestContext, next HertzNext) {
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
		newHashresult, err := password_hashed.HashPasswordWithSalt(req.Credential.Password, password_hashed.Options{Algorithm: "SHA-512"})
		if err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}
		cred.Hash = newHashresult.Hash
		cred.Salt = newHashresult.Salt
		cred.Algorithm = "SHA-512"

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
}

func DeleteCredentialHandler(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, initial_credentials InitialCredentials) func(c context.Context, r *app.RequestContext, next HertzNext) {
	return func(c context.Context, r *app.RequestContext, next HertzNext) {
		var body struct {
			Authorization CredentialsClient `json:"authorization"`
			Credential    struct {
				Username string `json:"username"`
			} `json:"credential"`
		}

		if err := r.BindJSON(&body); err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}

		if body.Credential.Username == "" {
			r.AbortWithMsg("Error: Username is empty", consts.StatusBadRequest)
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

		if err := credentialdb.Where("username = ?", body.Credential.Username).Delete(&CredentialStore{}).Error; err != nil {
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}

		r.JSON(consts.StatusOK, map[string]any{
			"message":  "Credential deleted successfully",
			"username": username,
		})
	}
}

func GenerateCredentialRoutes(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, initial_credentials InitialCredentials) []RouteConfig {
	return []RouteConfig{
		{
			Path:   "/credentials",
			Method: "PUT",
			MiddleWare: UpdateCredentialHandler(credentialdb, tokendb, sessiondb),
		},
		{
			Headers: map[string]string{"x-HTTP-method-override": "GET"},
			Path:    "/credentials",
			Method:  "POST",
			MiddleWare: GetCredentialsHandler(credentialdb, tokendb, sessiondb),
		},
		{
			Headers: map[string]string{"x-HTTP-method-override": "POST"},
			Path:    "/credentials",
			Method:  "POST",
			MiddleWare: CreateCredentialHandler(credentialdb, tokendb, sessiondb),
		},
		{
			Path:   "/credentials",
			Method: "DELETE",
			MiddleWare: DeleteCredentialHandler(credentialdb, tokendb, sessiondb, initial_credentials),
		},
	}
}

func IsUserExists(credentialdb *gorm.DB, username string) bool {
	var user CredentialStore
	if err := credentialdb.Where("username = ?", username).First(&user).Error; err != nil {
		return false
	}
	return true
}

type InitialCredentials  =go_ws_sh.InitialCredentials