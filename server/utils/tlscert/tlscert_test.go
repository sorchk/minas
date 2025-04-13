package tlscert

import (
	"testing"
)

// TestTls TLS证书生成功能的测试函数
// 用于测试证书生成功能是否正常工作
// 生成的测试证书将保存在/tmp/test目录下
// 参数 t: 测试上下文
func TestTls(t *testing.T) {
	// 调用证书生成函数，参数说明：
	// 证书文件路径: /tmp/test/1.crt
	// 私钥文件路径: /tmp/test/1.key
	// 通用名称: test
	// 组织名称: aide
	// 附加域名: test.com
	GenerateSelfSignedCert("/tmp/test/1.crt", "/tmp/test/1.key", "test", "aide", "test.com")
}
