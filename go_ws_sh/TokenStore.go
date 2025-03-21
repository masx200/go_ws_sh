package go_ws_sh

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type TokenStore struct {
	Identifier string `json:"identifier" gorm:"primarykey;unique;index;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	Hash      string `json:"hash" gorm:"index;not null"`
	Salt      string `json:"salt" gorm:"index;not null"`
	Algorithm string `json:"algorithm" gorm:"index;not null"`

	Username    string `json:"username" gorm:"index;not null"`
	Description string `json:"description" gorm:"index;not null"`
}

// Token 定义 Token 结构体

func (t TokenStore) String() string {
	return fmt.Sprintf("TokenStore{ CreatedAt: %v, UpdatedAt: %v, DeletedAt: %v, Hash: %s, Salt: %s, Algorithm: %s, Identifier: %s, Username: %s}",
		t.CreatedAt, t.UpdatedAt, t.DeletedAt, t.Hash, t.Salt, t.Algorithm, t.Identifier, t.Username)
}
func (TokenStore) TableName() string {
	return strings.ToLower("TokenStore")
}
