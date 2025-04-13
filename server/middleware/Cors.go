package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件
// 返回一个处理跨域请求的Gin中间件函数，允许来自不同源的HTTP请求访问API
func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		// 设置允许所有来源访问
		context.Header("Access-Control-Allow-Origin", "*")
		// 设置允许的HTTP请求头
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		// 设置允许的HTTP方法
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		// 设置可以被客户端访问的响应头
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		// 设置是否允许发送Cookie
		context.Header("Access-Control-Allow-Credentials", "true")
		// 对于OPTIONS预检请求，直接返回204状态码(无内容)
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		// 继续处理请求
		context.Next()
	}
}
