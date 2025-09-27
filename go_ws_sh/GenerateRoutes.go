package go_ws_sh

import (
	"gorm.io/gorm"

	"github.com/masx200/go_ws_sh/routes"
	"github.com/masx200/go_ws_sh/types"
)

// GenerateRoutesHttp 根据 openapi 文件生成路由配置
func GenerateRoutesHttp(credentialdb *gorm.DB, tokendb *gorm.DB, sessiondb *gorm.DB, initial_credentials types.InitialCredentials, initial_sessions []Session) []routes.RouteConfig {
	tokenRoutes := routes.GenerateTokenRoutes(credentialdb, tokendb, sessiondb)
	credentialRoutes := routes.GenerateCredentialRoutes(credentialdb, tokendb, sessiondb, initial_credentials)
	sessionRoutes := routes.GenerateSessionRoutes(credentialdb, tokendb, sessiondb, convertSessions(initial_sessions))

	allRoutes := append(tokenRoutes, credentialRoutes...)
	allRoutes = append(allRoutes, sessionRoutes...)

	return allRoutes
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
