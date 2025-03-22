package go_ws_sh

import (
	"os"

	"gorm.io/gorm"
)

// EnsureSessions 函数用于确保会话信息存在于数据库中，如果不存在则从配置中初始化
func EnsureSessions(config ConfigServer, sessiondb *gorm.DB) error {
	// 获取 SessionFile，如果为空则使用默认值
	SessionFile := config.SessionFile
	if SessionFile == "" {
		SessionFile = "session_store.db"
	}

	// 如果初始会话列表为空，则直接返回
	if len(config.InitialSessions) == 0 {
		return nil
	}

	// 检查数据库中是否存在记录
	var count int64
	sessiondb.Model(&SessionStore{}).Count(&count)

	// 检查文件是否存在
	if _, err := os.Stat(SessionFile); os.IsNotExist(err) || count == 0 {
		// 遍历初始会话列表
		for _, initialSession := range config.InitialSessions {
			args, err := StringSlice(initialSession.Args).Value()
			if err != nil {
				return err
			}

			// 创建 SessionStore 结构体实例
			session := SessionStore{
				Name: initialSession.Name,
				Cmd:  initialSession.Cmd,
				Args: string(args),
				Dir:  initialSession.Dir,
			}
			// 将会话信息保存到数据库中
			if err := sessiondb.Create(&session).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
