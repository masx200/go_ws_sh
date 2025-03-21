package go_ws_sh

type ConfigServer struct {
	SessionFile     string         `json:"session_file"`
	CredentialFile  string         `json:"credential_file"`
	InitialSessions []Session      `json:"initial_sessions"`
	Servers         []ServerConfig `json:"servers"`

	TokenFile string `json:"token_file"`

	// 添加初始化用户名和初始密码字段
	InitialUsername string `json:"initial_username"`
	InitialPassword string `json:"initial_password"`
}

type Session struct {
	Name string `json:"name"`

	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
	Dir  string   `json:"dir"`
}

type ServerConfig struct {
	Alpn     string `json:"alpn"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
	Cert     string `json:"cert"`
	Key      string `json:"key"`
}
