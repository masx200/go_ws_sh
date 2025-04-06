package go_ws_sh

import (
	"log"

	"gorm.io/gorm"
)

// ... existing code ...

// IsUserExists 判断指定 username 的用户是否存在于 credentialdb 中
func IsUserExists(credentialdb *gorm.DB, username string) bool {
	var count int64
	err := credentialdb.Model(&CredentialStore{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		log.Printf("Error checking if user exists: %v", err)
		return false
	}
	return count > 0
}

// ... existing code ...

// IsUserExists 判断指定 username 的用户是否存在于 credentialdb 中
func IsSessionExists(sessiondb *gorm.DB, sessionname string) bool {
	var count int64
	err := sessiondb.Model(&SessionStore{}).Where("name = ?", sessionname).Count(&count).Error
	if err != nil {
		log.Printf("Error checking if user exists: %v", err)
		return false
	}
	return count > 0
}
