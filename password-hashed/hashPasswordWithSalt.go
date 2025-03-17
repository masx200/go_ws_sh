package password_hashed

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"hash"
)

type Options struct {
	algorithm  string
	saltLength int
}

// HashPasswordWithSalt 生成加盐哈希，支持自定义算法和盐值长度
func HashPasswordWithSalt(password string, options ...Options) (map[string]string, error) {
	var option Options
	// 可选参数,判断参数长度
	if len(options) > 0 {
		option = options[0]
	}
	var saltLength = option.saltLength
	var algorithm = option.algorithm
	if saltLength == 0 {
		saltLength = 64 // 默认盐值长度为16字节
	}
	if algorithm == "" {
		algorithm = "SHA-512" // 默认使用SHA-256算法
	}
	// 1. 生成随机盐值 [[7]]
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil { // 使用crypto/rand生成安全随机数
		return nil, err
	}

	// 2. 合并盐值和密码（盐值在前，密码在后）[[4]]
	passwordBytes := []byte(password) // Go原生字符串即UTF-8编码 [[3]]
	merged := append(salt, passwordBytes...)

	// 3. 选择哈希算法
	var hasher hash.Hash
	switch algorithm {
	case "SHA-384":
		hasher = sha512.New384()
	case "SHA-256":
		hasher = sha256.New()
	case "SHA-512":
		hasher = sha512.New()
	default:
		return nil, errors.New("unsupported algorithm")
	}

	// 4. 计算哈希值 [[8]]
	if _, err := hasher.Write(merged); err != nil {
		return nil, err
	}
	hashBytes := hasher.Sum(nil)

	// 5. 转换为十六进制字符串 [[9]]
	return map[string]string{
		"hash":      hex.EncodeToString(hashBytes),
		"salt":      hex.EncodeToString(salt),
		"algorithm": algorithm,
	}, nil
}
