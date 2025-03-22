package go_ws_sh

// import (
// 	"encoding/json"
// 	"os"
// 	"path/filepath"
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
// )

// // WriteCredentials 函数将 CredentialsStore 类型的数据写入到指定的文件中
// func WriteCredentials(getfilepath func() (string, error), CredentialsStore CredentialsStore) error {
// 	// 获取文件路径
// 	filePath, err := getfilepath()
// 	if err != nil {
// 		return err
// 	}
// 	//创建出这个文件所在的文件夹
// 	dir := filepath.Dir(filePath)
// 	err = os.MkdirAll(dir, os.ModePerm)
// 	if err != nil {
// 		return err
// 	}
// 	// 创建或打开文件
// 	file, err := os.Create(filePath)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// 创建 JSON 编码器
// 	encoder := json.NewEncoder(file)
// 	encoder.SetIndent("", "  ")

// 	// 将 CredentialsStore 编码为 JSON 并写入文件
// 	err = encoder.Encode(CredentialsStore)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
