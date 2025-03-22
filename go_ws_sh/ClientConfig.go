package go_ws_sh

import "fmt"

type ClientSession struct {
	Name string `json:"name"`
}
type CredentialsClient struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Token      string `json:"token"`
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

func (c CredentialsClient) String() string {
	return fmt.Sprintf("CredentialsClient{Username: %s, Password: %s, Token: %s, Type: %s, Identifier: %s}",
		c.Username, c.Password, c.Token, c.Type, c.Identifier)
}

type ConfigClient struct {
	Credentials CredentialsClient      `json:"credentials"`
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
