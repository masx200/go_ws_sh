package go_ws_sh

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	// "github.com/philippgille/gokv/file"
	"github.com/masx200/go_ws_sh/types"
)

type TokenInfo = types.CredentialsClient

// 创建一个登录登出处理函数
// func createhandlerloginlogout(Sessions []Session, credentialdb *gorm.DB, tokendb *gorm.DB, next func(w context.Context, r *app.RequestContext)) func(w context.Context, r *app.RequestContext) {

// 	// 返回一个处理函数
// 	return func(w context.Context, r *app.RequestContext) {
// 		// 获取请求参数
// 		var name = r.Param("name")
// 		// 如果TokenFile为空，则返回错误

// 		// 如果请求参数为list，则处理列表请求
// 		if name == "sessions" {
// 			// 创建一个TokenInfo结构体
// 			var credential TokenInfo = TokenInfo{}

// 			// 将请求参数绑定到TokenInfo结构体中
// 			var err = r.BindJSON(&credential)
// 			if err != nil {
// 				r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
// 				return
// 			}
// 			log.Println(credential)
// 			validateFailure := Validatepasswordortoken(credential, credentialdb, tokendb, r)
// 			if validateFailure {
// 				log.Println("用户登录失败:")
// 				return
// 			}
// 			log.Println("用户登录成功:")
// 			r.JSON(
// 				consts.StatusOK,
// 				map[string]interface{}{
// 					"message": "List of Sessions ok",
// 					"list": slice.Map(Sessions, func(session Session) string {
// 						return session.Name
// 					}),
// 					"username": credential.Username,
// 				},
// 			)
// 			return

// 			// return
// 		}
// 		if name == "login" {
// 			handlePost(r, credentialdb, tokendb)
// 			return

// 		} else if name == "logout" {
// 			handleDelete(r, credentialdb, tokendb)
// 			return
// 		}

// 		next(w, r)
// 		// return

// 	}
// }

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

// // filterSessionsByUsername 根据输入的 username 过滤 Sessions 数组
// func filterSessionsByUsername(Sessions []Session, username string) []Session {
// 	var filteredSessions []Session
// 	for _, session := range Sessions {
// 		if session.Username == username {
// 			filteredSessions = append(filteredSessions, session)
// 		}
// 	}
// 	return filteredSessions
// }
