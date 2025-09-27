package go_ws_sh

import (
	"github.com/masx200/go_ws_sh/types"
)

type ClientSession struct {
	Name string `json:"name"`
}

// String method removed - cannot define new methods on non-local type types.CredentialsClient

type ConfigClient struct {
	Credentials types.CredentialsClient      `json:"credentials"`
	Sessions    ClientSession          `json:"sessions"`
	Servers     ClientConfigConnection `json:"servers"`
}

type ClientConfigConnection struct {
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
	Ca       string `json:"ca"`
	Insecure bool   `json:"insecure"`
	Host     string `json:"host"`
	IP       string `json:"ip"`
}
