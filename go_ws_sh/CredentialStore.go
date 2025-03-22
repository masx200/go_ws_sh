package go_ws_sh

import (
	"fmt"

	"gorm.io/gorm"
)

// 定义结构体以匹配JSON结构
type CredentialStore struct {
	gorm.Model
	Username  string `json:"username" gorm:"index;unique;not null"`
	Hash      string `json:"hash" gorm:"index;not null"`
	Salt      string `json:"salt" gorm:"index;not null"`
	Algorithm string `json:"algorithm" gorm:"index;not null"`
}

func (c CredentialStore) String() string {
	return fmt.Sprintf("Credentials{ID: %d, CreatedAt: %v, UpdatedAt: %v, DeletedAt: %v, Username: %s, Hash: %s, Salt: %s, Algorithm: %s}",
		c.ID, c.CreatedAt, c.UpdatedAt, c.DeletedAt, c.Username, c.Hash, c.Salt, c.Algorithm)
}
func (CredentialStore) TableName() string {
	return "credentials"
}
