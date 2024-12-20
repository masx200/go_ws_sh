package go_ws_sh

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
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
func createhandler(credentials []Credentials /* config Config, */, next func(w context.Context, r *app.RequestContext) /* httpServeMux *http.ServeMux */) func(w context.Context, r *app.RequestContext) {

	var credentialsmap = map[string]bool{}

	for _, credential := range credentials {
		credentialsmap[credential.Username+":"+credential.Password] = true
	}
	return func(w context.Context, r *app.RequestContext) {
		fmt.Println("Request Method:", string(r.Method()))
		fmt.Println("Request Headers:")
		fmt.Println("{")
		r.Request.Header.VisitAll(func(key, value []byte) {
			fmt.Println(string(key), ":", string(value))
		})
		fmt.Println("}")
		//check crediential
		auth := r.Request.Header.Get("Authorization")
		if auth == "" {
			log.Println("No Authorization header")
			r.Response.Header.Set("WWW-Authenticate", "Basic realm=\"go_ws_sh\"")
			r.SetStatusCode(consts.StatusUnauthorized)
			// r.AbortWithMsg("No Authorization header", consts.StatusUnauthorized)
			r.WriteString("No Authorization header")
			return
		}
		if !strings.HasPrefix(auth, "Basic ") {
			log.Println("No Basic auth")
			r.SetStatusCode(consts.StatusUnauthorized)
			r.Response.Header.Set("WWW-Authenticate", "Basic realm=\"go_ws_sh\"")
			r.WriteString("No Basic auth")
			return
		}

		credential := strings.TrimPrefix(auth, "Basic ")
		var rawcredential []byte
		if rawcredential2, err := base64.StdEncoding.DecodeString(credential); err != nil {
			r.Response.Header.Set("WWW-Authenticate", "Basic realm=\"go_ws_sh\"")
			r.SetStatusCode(consts.StatusUnauthorized)
			r.WriteString(err.Error())
			return
		} else {
			rawcredential = rawcredential2
		}
		// fmt.Printf("credential: %v\n", string(rawcredential))
		if _, ok := credentialsmap[string(rawcredential)]; !ok {
			log.Println("Invalid credential", credential)
			r.Response.Header.Set("WWW-Authenticate", "Basic realm=\"go_ws_sh\"")
			r.SetStatusCode(consts.StatusUnauthorized)
			r.WriteString("Invalid credential")
			// r.AbortWithMsg("Invalid credential", consts.StatusUnauthorized)
			return
		}

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

		if !r.IsGet() /* != http.MethodGet */ {
			log.Println("Not a get request")
			r.NotFound()
			//http.NotFound(w, r)
			return
		}
		//httpServeMux.ServeHTTP(w, r)
		next(w, r)
	}

}
