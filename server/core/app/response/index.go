package response

import (
	"errors"
	"net/http"
	"server/utils/logger"
	"server/utils/simple"

	"github.com/gin-gonic/gin"
)

// ResData 响应数据结构体
// 定义了API接口统一的返回格式
type ResData struct {
	Code  int         `json:"code"`  // 状态码，表示请求处理的结果
	Msg   string      `json:"msg"`   // 消息提示，对状态码的文字说明
	Total int64       `json:"total"` // 数据总数，用于分页场景
	Data  interface{} `json:"data"`  // 具体返回的数据
}

// List 返回分页列表数据
// 适用于需要分页的数据列表接口
// 参数:
//   - ctx: Gin上下文
//   - message: 响应消息
//   - total: 数据总数
//   - list: 分页后的数据列表
func List(ctx *gin.Context, message string, total int64, list interface{}) {
	ctx.JSON(http.StatusOK, ResData{Code: http.StatusOK, Msg: message, Total: total, Data: list})
	ctx.Abort()
}

// Data 返回单个数据对象
// 适用于返回单个数据对象的接口
// 参数:
//   - ctx: Gin上下文
//   - message: 响应消息
//   - data: 要返回的数据
func Data(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusOK, ResData{Code: http.StatusOK, Msg: message, Data: data})
	ctx.Abort()
}

// Success 返回成功消息
// 适用于只需要返回操作成功的接口
// 参数:
//   - ctx: Gin上下文
//   - message: 成功消息
func Success(ctx *gin.Context, message string) {
	Message(ctx, http.StatusOK, message)
}

// Error 返回错误消息
// 适用于需要返回错误信息的接口
// 参数:
//   - ctx: Gin上下文
//   - err: 错误对象
func Error(ctx *gin.Context, err error) {
	var simpleErr simple.SimpleError
	if errors.As(err, &simpleErr) {
		logger.LOG.Errorf(simpleErr.Error())
		if simpleErr.Code < 0 {
			ctx.JSON(http.StatusOK, ResData{Code: 4000, Msg: simpleErr.Msg, Data: simpleErr.Data})
			ctx.Abort()
			return
		} else {
			ctx.JSON(http.StatusOK, ResData{Code: simpleErr.Code, Msg: simpleErr.Msg, Data: simpleErr.Data})
			ctx.Abort()
			return
		}
	} else {
		logger.LOG.Errorf("Error: %s", err.Error())
		Message(ctx, 4000, err.Error())
		return
	}
}

// BadRequest 返回请求数据格式错误消息
// 适用于请求数据格式错误的场景
// 参数:
//   - ctx: Gin上下文
//   - message: 错误消息
func BadRequest(ctx *gin.Context, message string) {
	Error(ctx, simple.NewSimpleError(http.StatusBadRequest, message, nil))
}

// NotFound 返回资源不存在消息
// 适用于资源不存在的场景
// 参数:
//   - ctx: Gin上下文
//   - message: 错误消息
func NotFound(ctx *gin.Context, message string) {
	Error(ctx, simple.NewSimpleError(http.StatusNotFound, message, nil))
}

// NoContent 返回无数据消息
// 适用于无数据的场景
// 参数:
//   - ctx: Gin上下文
//   - message: 消息
func NoContent(ctx *gin.Context, message string) {
	Error(ctx, simple.NewSimpleError(http.StatusNoContent, message, nil))
}

// Message 返回自定义状态码消息
// 适用于需要自定义状态码的场景
// 参数:
//   - ctx: Gin上下文
//   - code: 自定义状态码
//   - message: 消息
func Message(ctx *gin.Context, code int, message string) {
	ctx.JSON(http.StatusOK, ResData{Code: code, Msg: message})
	ctx.Abort()
}

// Unauthorized 返回未认证消息
// 适用于需要认证的场景
// 参数:
//   - ctx: Gin上下文
//   - message: 消息
func Unauthorized(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusUnauthorized, ResData{Code: http.StatusUnauthorized, Msg: message})
	ctx.Abort()
}

// Forbidden 返回无权限消息
// 适用于需要权限的场景
// 参数:
//   - ctx: Gin上下文
//   - message: 消息
func Forbidden(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusForbidden, ResData{Code: http.StatusUnauthorized, Msg: message})
	ctx.Abort()
}
