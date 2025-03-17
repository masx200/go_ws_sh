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
	println(result.String())           // 输出盐值（128字符，64字节）
}
