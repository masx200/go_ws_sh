package go_ws_sh

import "gorm.io/gorm"

// MoveSession 移动会话的具体逻辑
func CopySession(sessiondb *gorm.DB, sessionName, destinationName string) (*SessionStore, error) {
	// 这里可以添加具体的数据库操作，例如更新会话的位置
	// 示例：查询会话记录
	var session SessionStore
	result := sessiondb.Where("name = ?", sessionName).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}

	// 示例：更新会话记录
	// 这里可以修改会话的相关字段，例如更新到新的目录
	session.Name = destinationName
	if err := sessiondb.Unscoped().Where("name=?", destinationName).Delete(&SessionStore{}).Error; err != nil {
		return nil, err
	}

	// if err := sessiondb.Where("name=?", sessionName).Delete(&SessionStore{}).Error; err != nil {
	// 	return nil, err
	// }
	newSession := SessionStore{
		Name: destinationName,
		Cmd:  session.Cmd,
		Args: session.Args,
		Dir:  session.Dir,
	}
	if err := sessiondb.Create(&newSession).Error; err != nil {
		return nil, err
	}

	return &newSession, nil
}
