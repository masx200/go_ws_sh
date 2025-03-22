package go_ws_sh

import (
	"encoding/json"

	"gorm.io/gorm"
)

// ReadAllSessions 从 sessiondb 中读取所有的 SessionStore 并转换为 Session 结构体切片
func ReadAllSessions(sessiondb *gorm.DB) ([]Session, error) {
	var sessionStores []SessionStore
	// 查询 sessiondb 中的所有 SessionStore 记录
	if err := sessiondb.Find(&sessionStores).Error; err != nil {
		return nil, err
	}

	sessions := make([]Session, 0, len(sessionStores))
	for _, store := range sessionStores {
		var args []string
		// 将 Args 字段（字符串形式）反序列化为字符串切片
		if err := json.Unmarshal([]byte(store.Args), &args); err != nil {
			return nil, err
		}
		session := Session{
			Name: store.Name,
			Cmd:  store.Cmd,
			Args: args,
			Dir:  store.Dir,
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}
