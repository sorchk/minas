package mfa

import (
	"time"

	"github.com/xlzd/gotp"
)

// ValidCode 验证多因素认证码是否有效
// 检查当前时间和前一个时间间隔内的验证码是否匹配
// 参数 code: 需要验证的验证码
// 参数 secret: 用户的密钥
// 参数 interval: 验证码有效的时间间隔（秒）
// 返回值: 如果验证码有效返回true，否则返回false
func ValidCode(code, secret string, interval int) bool {
	totp := gotp.NewTOTP(secret, 6, interval, nil)
	now := time.Now().Unix()
	prevTime := now - int64(interval)
	return totp.Verify(code, now) || totp.Verify(code, prevTime)
}

// GetCode 根据密钥生成当前的验证码
// 参数 secret: 用户的密钥
// 参数 interval: 验证码有效的时间间隔（秒）
// 返回值: 生成的六位数验证码
func GetCode(secret string, interval int) string {
	totp := gotp.NewTOTP(secret, 6, interval, nil)
	return totp.Now()
}

// GenSecret 生成指定长度的随机密钥
// 用于创建新用户的MFA密钥
// 参数 size: 密钥的长度
// 返回值: 生成的随机密钥字符串
func GenSecret(size int) string {
	return gotp.RandomSecret(size)
}
