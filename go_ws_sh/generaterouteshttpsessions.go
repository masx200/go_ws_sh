package go_ws_sh

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"

	"github.com/masx200/go_ws_sh/types"
)


func GenerateRoutesHttpsessions(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, initial_sessions []Session) []types.RouteConfig {
	routes_copy_and_move := []types.RouteConfig{

		{
			Path:   "/sessions",
			Method: "COPY",
			MiddleWare: HertzCompose(AuthorizationMiddleware(credentialdb, tokendb, sessiondb), func(c context.Context, r *app.RequestContext, next types.HertzNext) {
				CopyMiddleware(credentialdb, tokendb, sessiondb, c, r, next)
			}),
		},
		{
			Path:   "/sessions",
			Method: "MOVE",
			MiddleWare: HertzCompose(AuthorizationMiddleware(credentialdb, tokendb, sessiondb), func(c context.Context, r *app.RequestContext, next types.HertzNext) {
				MoveMiddleware(initial_sessions, credentialdb, tokendb, sessiondb, c, r, next)
			}),
		},
	}
	return routes_copy_and_move
}
