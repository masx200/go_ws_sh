package go_ws_sh

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)




type RouteConfig struct {
	Path    string
	Method  string
	Handler func(c context.Context,r *app.RequestContext)
}