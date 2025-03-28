package go_ws_sh

import "time"

type InitialCredentials []struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type ConfigServer struct {
	SessionFile     string         `json:"session_file"`
	CredentialFile  string         `json:"credential_file"`
	InitialSessions []Session      `json:"initial_sessions"`
	Servers         []ServerConfig `json:"servers"`

	TokenFile          string             `json:"token_file"`
	InitialCredentials InitialCredentials `json:"initial_credentials"`
	// 添加初始化用户名和初始密码字段

}

type Session struct {
	Name string `json:"name"`

	Cmd       string    `json:"cmd"`
	Args      []string  `json:"args"`
	Dir       string    `json:"dir"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ServerConfig struct {
	Alpn     string `json:"alpn"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
	Cert     string `json:"cert"`
	Key      string `json:"key"`
}
