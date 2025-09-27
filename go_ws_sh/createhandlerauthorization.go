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
)

func createhandlerauthorization_websocket(credentialdb *gorm.DB, tokendb *gorm.DB, next func(w context.Context, r *app.RequestContext)) func(w context.Context, r *app.RequestContext) {

	return func(w context.Context, r *app.RequestContext) {

		Upgrade := strings.ToLower(r.Request.Header.Get("Upgrade"))
		Connection := strings.ToLower(r.Request.Header.Get("Connection"))

		if !strings.Contains(Connection, "upgrade") {
			log.Println("Not a upgrade request")
			r.NotFound()
			return
		}
		if !strings.Contains(Upgrade, "websocket") {
			log.Println("Not a websocket request")

			r.NotFound()
			return
		}

		if !r.IsGet() {
			log.Println("Not a get request")
			r.NotFound()

			return
		}

		proto := r.Request.Header.Get("Sec-Websocket-Protocol")
		if proto != "" {
			for _, str := range strings.Split(proto, ",") {
				decoded, err := url.QueryUnescape(str)
				if err != nil {
					log.Println("proto", str)
					fmt.Printf("Error parsing input: %v\n", err)
					r.SetStatusCode(consts.StatusUnauthorized)
					r.WriteString(err.Error())
					return
				}
				postData, err := url.ParseQuery(decoded)
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
				validateFailure := Validatepasswordortoken(req, credentialdb, tokendb, r)
				if validateFailure {
					log.Println("用户登录失败:")
					return
				}
				var ok = !validateFailure
				if ok {
					log.Println("用户登录成功:")
					next(w, r)
					return
				}

			}
		} else {

			r.Response.Header.Set("WWW-Authenticate", "Basic realm=\"go_ws_sh\"")
			r.SetStatusCode(consts.StatusUnauthorized)
			r.WriteString("Invalid credential Unauthorized")

			return

		}
	}
}

// func parseKeyValuePairs(input string) (map[string]string, error) {
// 	pairs := strings.Split(input, ";")
// 	result := make(map[string]string)

// 	for _, pair := range pairs {
// 		kv := strings.SplitN(pair, "=", 2)
// 		if len(kv) != 2 {
// 			return nil, fmt.Errorf("invalid pair: %s", pair)
// 		}

// 		key := strings.TrimSpace(kv[0])
// 		value := strings.TrimSpace(kv[1])

// 		if key == "" || value == "" {
// 			return nil, fmt.Errorf("empty key or value in pair: %s", pair)
// 		}

// 		result[key] = value
// 	}

// 	return result, nil
// }
