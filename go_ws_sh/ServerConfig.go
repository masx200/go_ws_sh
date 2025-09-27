package go_ws_sh

import (
	"github.com/masx200/go_ws_sh/types"
)

type ConfigServer struct {
	SessionFile     string         `json:"session_file"`
	CredentialFile  string         `json:"credential_file"`
	InitialSessions []types.Session `json:"initial_sessions"`
	Servers         []ServerConfig `json:"servers"`

	TokenFile          string               `json:"token_file"`
	InitialCredentials types.InitialCredentials `json:"initial_credentials"`
	// 添加初始化用户名和初始密码字段

}

type ServerConfig struct {
	Alpn     string `json:"alpn"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
	Cert     string `json:"cert"`
	Key      string `json:"key"`
}
