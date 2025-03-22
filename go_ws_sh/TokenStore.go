package go_ws_sh

import (
	"fmt"

	"gorm.io/gorm"
)

type TokenStore struct {
	gorm.Model
	Hash       string `json:"hash" gorm:"index;not null"`
	Salt       string `json:"salt" gorm:"index;not null"`
	Algorithm  string `json:"algorithm" gorm:"index;not null"`
	Identifier string `json:"identifier" gorm:"unique;index;not null"`
	Username   string `json:"username" gorm:"index;not null"`
}

// Token 定义 Token 结构体

func (t TokenStore) String() string {
	return fmt.Sprintf("TokenStore{ID: %d, CreatedAt: %v, UpdatedAt: %v, DeletedAt: %v, Hash: %s, Salt: %s, Algorithm: %s, Identifier: %s, Username: %s}",
		t.ID, t.CreatedAt, t.UpdatedAt, t.DeletedAt, t.Hash, t.Salt, t.Algorithm, t.Identifier, t.Username)
}
func (TokenStore) TableName() string {
	return "TokenStore"
}
