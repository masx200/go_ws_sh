package go_ws_sh

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/akrennmair/slice"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/philippgille/gokv/file"
)

type TokenInfo struct {
	Token string `json:"token"`
}

func createhandlerloginlogout(Sessions []Session, TokenFolder string, credentials []Credentials, next func(w context.Context, r *app.RequestContext)) func(w context.Context, r *app.RequestContext) {

	var credentialsmap = map[string]bool{}

	for _, credential := range credentials {
		credentialsmap[credential.Username+":"+credential.Password] = true
	}
	var store, err = file.NewStore(file.Options{Directory: TokenFolder})
	return func(w context.Context, r *app.RequestContext) {
		var name = r.Param("name")
		if err != nil {
			log.Println("Error: " + err.Error())
			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
			return
		}
		if name == "list" {
			var credential TokenInfo = TokenInfo{}
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
			} else {
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
			}

			return
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

func ValidateToken(token string, store file.Store) (bool, map[string]string) {
	var result = map[string]string{}
	found, err := store.Get(token, &result)

	if err != nil || !found {
		log.Println("ValidateToken", token, result, found, err)
		return false, nil
	}
	return true, result
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
