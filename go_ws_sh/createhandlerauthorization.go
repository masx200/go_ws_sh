package go_ws_sh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"
	// "github.com/philippgille/gokv/file"
)

// createHandler takes a slice of Credentials and a next function to generate a new handler function.
// This handler function will authenticate the request based on the provided credentials and, if authenticated, call the next function.
// Parameters:
//
//	credentials - A slice of Credentials used for authentication.
//	next - A function to execute if the authentication is successful.
//
// Returns:
//
//	A function that takes a context and a RequestContext, performs authentication, and calls the next function if successful.
func createhandlerauthorization(credentialdb *gorm.DB, tokendb *gorm.DB, next func(w context.Context, r *app.RequestContext)) func(w context.Context, r *app.RequestContext) {

	return func(w context.Context, r *app.RequestContext) {
		//		if TokenFile == "" {
		//			log.Println("Error: " + "TokenFile is empty")
		//			r.AbortWithMsg("Error:  "+"TokenFile is empty", consts.StatusInternalServerError)
		//			return
		//		}
		//		if err != nil {
		//			log.Println("Error: " + err.Error())
		//			r.AbortWithMsg("Error: "+err.Error(), consts.StatusInternalServerError)
		//			return
		//		}
		fmt.Println("Request Method:", string(r.Method()))
		fmt.Println("Request Headers:")
		fmt.Println("{")
		r.Request.Header.VisitAll(func(key, value []byte) {
			fmt.Println(string(key), ":", string(value))
		})
		fmt.Println("}")

		Upgrade := strings.ToLower(r.Request.Header.Get("Upgrade"))
		Connection := strings.ToLower(r.Request.Header.Get("Connection"))
		//if !tokenListContainsValue(r.Request.Header, "Connection", "upgrade") {
		if !strings.Contains(Connection, "upgrade") {
			log.Println("Not a upgrade request")
			r.NotFound() //http.NotFound(w, r)
			return
		}
		if !strings.Contains(Upgrade, "websocket") {
			log.Println("Not a websocket request")
			// if !tokenListContainsValue(r.Header, "Upgrade", "websocket") {
			r.NotFound() //http.NotFound(w, r)
			return
		}
		//
		if !r.IsGet() /* != http.MethodGet */ {
			log.Println("Not a get request")
			r.NotFound()
			//http.NotFound(w, r)
			return
		}
		//check crediential
		//Sec-Websocket-Protocol
		proto := r.Request.Header.Get("Sec-Websocket-Protocol")
		if proto != "" {
			for _, str := range strings.Split(proto, ",") {
				postData, err := url.ParseQuery(str)
				if err != nil {
					log.Println(err)
					r.SetStatusCode(consts.StatusBadRequest)
					r.WriteString("Bad Request")
					return
				}
				log.Println(postData)
				var req CredentialsClient
				req.Token = postData.Get("token")
				req.Type = postData.Get("type")
				req.Username = postData.Get("username")
				req.Identifier = postData.Get("identifier")
				req.Password = postData.Get("password")
				shouldReturn := Validatepasswordortoken(req, credentialdb, tokendb, r)
				if shouldReturn {
					log.Println("用户登录失败:")
					return
				}
				var ok = !shouldReturn
				if ok {
					log.Println("用户登录成功:")
					next(w, r)
					return
				}

			}
		} else {

			//			auth := r.Request.Header.Get("Authorization")
			//			if auth == "" {
			//				log.Println("No Authorization header")
			//				r.Response.Header.Set("WWW-Authenticate", "Basic realm=\"go_ws_sh\"")
			//				r.SetStatusCode(consts.StatusUnauthorized)
			//				// r.AbortWithMsg("No Authorization header", consts.StatusUnauthorized)
			//				r.WriteString("No Authorization header")
			//				return
			//			}
			//			if strings.HasPrefix(auth, "Bearer ") {
			//				token := strings.TrimPrefix(auth, "Bearer ")
			//				if token == "" {
			//					r.AbortWithMsg("Error: Unauthorized token is empty", consts.StatusUnauthorized)
			//					return
			//				}
			//				if ok, result := ValidateToken(token, store); !ok {
			//					r.AbortWithMsg("Error: Unauthorized token is invalid", consts.StatusUnauthorized)
			//					return
			//				} else if slices.Contains(slice.Map(credentials, func(credential Credentials) string { return credential.Username }), result["username"]) {
			//					log.Println("用户登录成功:" + result["username"] + ":" + token)
			//					next(w, r)
			//					return
			//				} else {
			//					r.AbortWithMsg("Error: Unauthorized token is invalid", consts.StatusUnauthorized)
			//					return
			//				}
			//			}
			//			if !strings.HasPrefix(auth, "Basic ") {
			//				log.Println("No Basic auth")
			//				r.SetStatusCode(consts.StatusUnauthorized)
			//				r.Response.Header.Set("WWW-Authenticate", "Basic realm=\"go_ws_sh\"")
			//				r.WriteString("No Basic auth Unauthorized")
			//				return
			//			}
			//
			//			credential := strings.TrimPrefix(auth, "Basic ")
			//			var rawcredential []byte
			//			if rawcredential2, err := base64.StdEncoding.DecodeString(credential); err != nil {
			//				r.Response.Header.Set("WWW-Authenticate", "Basic realm=\"go_ws_sh\"")
			//				r.SetStatusCode(consts.StatusUnauthorized)
			//				r.WriteString(err.Error())
			//				return
			//			} else {
			//				rawcredential = rawcredential2
			//			}
			//			// fmt.Printf("credential: %v\n", string(rawcredential))
			//			if _, ok := credentialsmap[string(rawcredential)]; !ok {
			//				log.Println("Invalid credential", credential)
			r.Response.Header.Set("WWW-Authenticate", "Basic realm=\"go_ws_sh\"")
			r.SetStatusCode(consts.StatusUnauthorized)
			r.WriteString("Invalid credential Unauthorized")
			// r.AbortWithMsg("Invalid credential", consts.StatusUnauthorized)
			return
			//			}
			//

			//			//httpServeMux.ServeHTTP(w, r)
			//			log.Println("用户登录成功:" + string(rawcredential))
			//			next(w, r)
		}
	}
}

// parseKeyValuePairs 解析形如 "key1=value1; key2=value2; key3=value3" 的字符串
// 并返回一个映射。
func parseKeyValuePairs(input string) (map[string]string, error) {
	pairs := strings.Split(input, ";")
	result := make(map[string]string)

	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid pair: %s", pair)
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		if key == "" || value == "" {
			return nil, fmt.Errorf("empty key or value in pair: %s", pair)
		}

		result[key] = value
	}

	return result, nil
}
