package go_ws_sh

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"

	"github.com/masx200/go_ws_sh/routes"
)

func GenerateRoutesHttpsessions(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, initial_sessions []Session) []RouteConfig {
	routes_copy_and_move := []RouteConfig{

		{
			Path:   "/sessions",
			Method: "COPY",
			MiddleWare: HertzCompose(AuthorizationMiddleware(credentialdb, tokendb, sessiondb), func(c context.Context, r *app.RequestContext, next HertzNext) {
				routes.CopyMiddleware(credentialdb, tokendb, sessiondb, c, r, routes.HertzNext(next))
			}),
		},
		{
			Path:   "/sessions",
			Method: "MOVE",
			MiddleWare: HertzCompose(AuthorizationMiddleware(credentialdb, tokendb, sessiondb), func(c context.Context, r *app.RequestContext, next HertzNext) {
				routes.MoveMiddleware(convertSessions(initial_sessions), credentialdb, tokendb, sessiondb, c, r, routes.HertzNext(next))
			}),
		},
	}
	return routes_copy_and_move
}

func convertSessions(sessions []Session) []routes.Session {
	result := make([]routes.Session, len(sessions))
	for i, session := range sessions {
		result[i] = routes.Session{
			Name:      session.Name,
			Cmd:       session.Cmd,
			Args:      session.Args,
			Dir:       session.Dir,
			CreatedAt: session.CreatedAt,
			UpdatedAt: session.UpdatedAt,
		}
	}
	return result
}
