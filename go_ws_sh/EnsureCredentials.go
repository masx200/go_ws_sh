package go_ws_sh

import (
	"os"

	"gorm.io/gorm"

	password_hashed "github.com/masx200/go_ws_sh/password-hashed"
)

// EnsureCredentials 函数用于确保认证信息文件存在，如果不存在则生成初始化认证信息
func EnsureCredentials(config ConfigServer, credentialdb *gorm.DB) error {
	// 获取 CredentialFile，如果为空则使用默认值
	credentialFile := config.CredentialFile
	if credentialFile == "" {
		credentialFile = "credential_store.db"
	}
	// 检查数据库中是否存在记录
	var count int64
	credentialdb.Model(&CredentialStore{}).Count(&count)
	// 检查文件是否存在
	if _, err := os.Stat(credentialFile); os.IsNotExist(err) || count == 0 {
		// 获取 InitialUsername 和 InitialPassword，如果为空则使用默认值
		username := config.InitialUsername
		if username == "" {
			username = "admin"
		}
		password := config.InitialPassword
		if password == "" {
			password = "pass"
		}

		// 生成认证信息
		//salt := "random_salt" // 这里需要生成真正的随机盐
		hashresult, err := password_hashed.HashPasswordWithSalt(password /* , salt */, password_hashed.Options{
			Algorithm: "SHA-512",
		})
		if err != nil {
			return err
		}
		// 创建 Credentials 结构体
		credentials := CredentialStore{

			Username:  username,
			Hash:      hashresult.Hash,
			Salt:      hashresult.Salt,
			Algorithm: "SHA-512", // 假设使用 SHA-512 算法

		}

		// 将认证信息保存到文件中
		return credentialdb.Create(&credentials).Error
	}
	return nil
}
