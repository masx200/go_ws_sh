package go_ws_sh

import (
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// StringSlice 自定义类型，用于存储字符串切片
type StringSlice []string

// Value 实现 driver.Valuer 接口，用于将 StringSlice 转换为可存储的值
func (s StringSlice) Value() ([]byte, error) {
	return json.Marshal(s)
}

// Scan 实现 sql.Scanner 接口，用于将存储的值转换为 StringSlice
func (s *StringSlice) Scan(value []byte) error {
	if value == nil {
		return nil
	}
	var bytes []byte = value

	return json.Unmarshal(bytes, s)
}

type SessionStore struct {
	gorm.Model

	Name string `json:"name" gorm:"index;unique;not null"`
	Cmd  string `json:"cmd" gorm:"index;not null"`
	Args string `json:"args" gorm:"index;not null"`
	Dir  string `json:"dir" gorm:"index;not null"`
}

// String 方法用于将 SessionStore 结构体转换为字符串表示
func (s SessionStore) String() string {
	return fmt.Sprintf("SessionStore{ID: %d, CreatedAt: %v, UpdatedAt: %v, DeletedAt: %v, Name: %s, Cmd: %s, Args: %v, Dir: %s}",
		s.ID, s.CreatedAt, s.UpdatedAt, s.DeletedAt, s.Name, s.Cmd, s.Args, s.Dir)
}

// TableName 方法用于指定 SessionStore 结构体对应的数据库表名
func (SessionStore) TableName() string {
	return strings.ToLower("SessionStore")
}
