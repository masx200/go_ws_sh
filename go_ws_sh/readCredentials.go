package go_ws_sh

import (
	"encoding/json"
	"os"
)

type CredentialsStore []CredentialStore

// readTokens 函数读取指定路径的 JSON 文件，并将其解析为 CredentialsStore 类型的数组，同时返回可能出现的错误
func readCredentials(getfilepath func() (string, error)) (CredentialsStore, error) {
	// 获取文件路径
	filePath, err := getfilepath()
	if err != nil {
		// 返回错误，而不是使用 panic
		return nil, err
	}

	if !FileExists(filePath) {
		return nil, nil
	}
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		// 返回错误，而不是使用 panic
		return nil, err
	}
	defer file.Close()

	// 创建 JSON 解码器
	var CredentialsStore CredentialsStore
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&CredentialsStore)
	if err != nil {
		// 返回错误，而不是使用 panic
		return nil, err
	}

	return CredentialsStore, nil
}
