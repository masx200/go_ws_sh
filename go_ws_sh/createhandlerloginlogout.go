package go_ws_sh

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"slices"

	"github.com/akrennmair/slice"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"
	// "github.com/philippgille/gokv/file"
)

type TokenInfo struct {
	Token string `json:"token"`
}

// 创建一个登录登出处理函数
func createhandlerloginlogout(Sessions []Session, credentialdb *gorm.DB, tokendb *gorm.DB, next func(w context.Context, r *app.RequestContext)) func(w context.Context, r *app.RequestContext) {

	// 返回一个处理函数
	return func(w context.Context, r *app.RequestContext) {
		// 获取请求参数
		var name = r.Param("name")
		// 如果TokenFile为空，则返回错误
		if TokenFile == "" {
			log.Println("Error: " + "TokenFile is empty")
			r.AbortWithMsg("Error:  "+"TokenFile is empty", consts.StatusInternalServerError)
			return
		}
		// 如果创建文件存储失败，则返回错误
		if err != nil {
			log.Println("Error: " + err.Error())
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}
		// 如果请求参数为list，则处理列表请求
		if name == "list" {
			// 创建一个TokenInfo结构体
			var credential TokenInfo = TokenInfo{}
			// 将请求参数绑定到TokenInfo结构体中
			var err = r.BindJSON(&credential)
			if err != nil {
				r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
				return
			}
			var token = credential.Token
			if token == "" {
				r.AbortWithMsg("Error: Unauthorized token is empty", consts.StatusUnauthorized)
				return
			}
			if ok, result := ValidateToken(token, store); !ok {
				r.AbortWithMsg("Error: Unauthorized token is invalid", consts.StatusUnauthorized)
				return
			} else if slices.Contains(slice.Map(credentials, func(credential Credentials) string { return credential.Username }), result["username"]) {
				r.JSON(
					consts.StatusOK,
					map[string]interface{}{
						"message": "List of Sessions ok",
						"list": slice.Map(Sessions, func(session Session) string {
							return session.Path
						}),
						"username": result["username"],
					},
				)
				return
			} else {
				r.AbortWithMsg("Error: Unauthorized token is invalid", consts.StatusUnauthorized)
				return
			}

			// return
		}
		if name == "login" {
			if err != nil {
				r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
				return
			}
			var credential Credentials = Credentials{}
			var err = r.BindJSON(&credential)
			if err != nil {
				r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
				return
			}
			var rawcredential = credential.Username + ":" + credential.Password
			if _, ok := credentialsmap[(rawcredential)]; !ok {
				log.Println("Invalid credential", credential)
				r.Response.Header.Set("WWW-Authenticate", "Basic realm=\"go_ws_sh\"")
				r.SetStatusCode(consts.StatusUnauthorized)
				r.WriteString("Invalid credential Unauthorized")
				// r.AbortWithMsg("Invalid credential", consts.StatusUnauthorized)
				return
			}
			numBytes := 120
			hexString, err := generateHexKey(numBytes)
			if err != nil {
				r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
				return
			}
			if err := store.Set(hexString, map[string]string{"username": credential.Username}); err != nil {
				log.Println("Error: " + err.Error())
				r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
				return
			}
			r.JSON(consts.StatusOK, map[string]string{"token": hexString,
				"message": "Login successful", "username": credential.Username,
			})
			return

		} else if name == "logout" {
			if err != nil {
				r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
				return
			}
			var credential TokenInfo = TokenInfo{}
			var err = r.BindJSON(&credential)
			if err != nil {
				r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
				return
			}
			var token = credential.Token
			if err := store.Delete(token); err != nil {
				r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
				return
			}
			r.JSON(consts.StatusOK, map[string]string{"message": "Logout successful", "token": token})
			return
		}

		next(w, r)
		// return

	}
}

func generateHexKey(length int) (string, error) {
	// 创建一个字节数组来保存随机字节
	randomBytes := make([]byte, length)

	// 使用crypto/rand读取随机字节
	n, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	if n != length {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	// 将字节切片转换为16进制字符串
	hexString := hex.EncodeToString(randomBytes)

	return hexString, nil
}
