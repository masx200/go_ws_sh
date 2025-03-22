package go_ws_sh

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func MatchAndRouteMiddleware(routes []RouteConfig) HertzMiddleWare {
	return func(c context.Context, r *app.RequestContext, next HertzNext) {
		if MatchAndHandleRoute(c, routes, r,next) {
			return
		}
		next(c, r)
	}

}
func MatchAndHandleRoute(w context.Context, routes []RouteConfig, r *app.RequestContext,next HertzNext) bool {
	for _, route := range routes {
		// 检查 Path 是否为空，若不为空则进行匹配
		pathMatch := route.Path == "" || route.Path == string(r.Path())
		// 检查 Method 是否为空，若不为空则进行匹配
		methodMatch := route.Method == "" || route.Method == string(r.Method())

		headersMatch := true
		if len(route.Headers) > 0 {
			for key, value := range route.Headers {
				headerValue := string(r.Request.Header.Peek(key))
				if headerValue != value {
					headersMatch = false
					break
				}
			}
		}

		// 如果 Path、Method 和 Headers 都匹配，则执行处理函数
		if pathMatch && methodMatch && headersMatch {
			route.MiddleWare(w, r,next)
			return true
		}
	}
	return false
}
