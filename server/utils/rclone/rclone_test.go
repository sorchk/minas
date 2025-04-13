package rclone

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCoreObscure 测试RClone密码加密功能
// 这个测试用例验证了：
// 1. CoreObscure函数能否成功加密字符串
// 2. 加密后的字符串是否有效（不为空且与原始字符串不同）
// 3. 能否成功获取RClone配置列表
func TestCoreObscure(t *testing.T) {
	// 定义测试用的明文密码
	plainText := "abc"

	// 调用加密函数
	result, err := CoreObscure(plainText)
	// 验证没有错误发生
	assert.Nil(t, err)

	// 输出加密结果用于调试
	log.Println("result:" + result)

	// 验证加密结果不为空
	assert.NotEmpty(t, result)

	// 验证加密后的字符串与原始字符串不同
	assert.NotEqual(t, result, plainText)

	// 获取RClone配置列表，测试API连通性
	list, err := DumpConfig()
	// 验证API调用成功
	assert.Nil(t, err)

	// 输出配置列表用于调试
	log.Printf("list:%v\n", list)
}
