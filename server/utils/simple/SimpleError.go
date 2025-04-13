package simple

import (
	"fmt"
)

// 预定义的错误状态码常量
const (
	StatusMFAAuth      = 4401 // 需要MFA验证
	StatusExistData    = 4001 // 数据已存在 一般用于新建数据时
	StatusExistRefData = 4002 // 存在关联数据 一般用于删除数据时
)

// SimpleError 简单错误结构体
// 用于统一的错误返回格式，包含错误码、错误消息和相关数据
type SimpleError struct {
	Code int         `json:"code"` // 错误码，用于前端识别错误类型
	Msg  string      `json:"msg"`  // 错误消息，描述错误的具体原因
	Data interface{} `json:"data"` // 相关数据，可携带额外的错误上下文信息
}

// NewSimpleError 创建带有自定义错误码、消息和数据的错误
// 参数 code: 自定义错误码
// 参数 msg: 错误消息
// 参数 data: 相关数据
// 返回值: SimpleError实例
func NewSimpleError(code int, msg string, data interface{}) SimpleError {
	return SimpleError{Code: code, Msg: msg, Data: data}
}

// NewSimpleErrorData 创建带有默认错误码(4000)、自定义消息和数据的错误
// 参数 msg: 错误消息
// 参数 data: 相关数据
// 返回值: SimpleError实例
func NewSimpleErrorData(msg string, data interface{}) SimpleError {
	return SimpleError{Code: 4000, Msg: msg, Data: data}
}

// NewSimpleErrorMessage 创建只带有默认错误码(4000)和消息的错误
// 参数 msg: 错误消息
// 返回值: SimpleError实例
func NewSimpleErrorMessage(msg string) SimpleError {
	return SimpleError{Code: 4000, Msg: msg, Data: nil}
}

// Error 实现error接口的Error方法
// 将SimpleError格式化为字符串
// 返回值: 格式化后的错误字符串
func (ee SimpleError) Error() string {
	return fmt.Sprintf("Error Code: %d, Message: %s, Data: %v", ee.Code, ee.Msg, ee.Data)
}
