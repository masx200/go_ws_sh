package password_hashed

import (
	"testing"
)

// 示例用法
func TestHashPasswordWithSalt(t *testing.T) {
	result, err := HashPasswordWithSalt("pass")
	if err != nil {
		t.Fatal(err)
		return
	}
	println("algorithm:", result["algorithm"]) // 输出SHA-512哈希值（128字符）
	println("Hash:", result["hash"])           // 输出SHA-512哈希值（128字符）
	println("Salt:", result["salt"])           // 输出盐值（128字符，64字节）
}
