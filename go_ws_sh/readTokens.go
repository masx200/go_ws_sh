package go_ws_sh

import (
	"encoding/json"
	"os"
)

// TokenStore 定义 Token 存储的结构体
type TokenStore []struct {
	Hash       string `json:"hash"`
	Salt       string `json:"salt"`
	Algorithm  string `json:"algorithm"`
	Identifier string `json:"identifier"`
	Username   string `json:"username"`
}

// readTokens 函数读取指定路径的 JSON 文件，并将其解析为 TokenStore 类型的数组，同时返回可能出现的错误
func readTokens(getfilepath func() (string, error)) (TokenStore, error) {
	// 获取文件路径
	filePath, err := getfilepath()
	if err != nil {
		// 返回错误，而不是使用 panic
		return nil, err
	}
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		// 返回错误，而不是使用 panic
		return nil, err
	}
	defer file.Close()

	// 创建 JSON 解码器
	var tokenStore TokenStore
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tokenStore)
	if err != nil {
		// 返回错误，而不是使用 panic
		return nil, err
	}

	return tokenStore, nil
}
