package go_ws_sh

import "gorm.io/gorm"

type SessionStore struct {
	gorm.Model

	Name string   `json:"name" gorm:"index;unique;not null"`
	Cmd  string   `json:"cmd" gorm:"index;not null"`
	Args []string `json:"args" gorm:"index;not null"`
	Dir  string   `json:"dir" gorm:"index;not null"`
}
