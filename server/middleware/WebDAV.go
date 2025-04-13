// Package middleware 提供各种中间件功能
package middleware

import (
	"server/utils/webdav"

	"github.com/gin-gonic/gin"
)

// WebDavHandler 创建WebDAV处理中间件
// 该中间件将拦截WebDAV请求并进行处理，非WebDAV请求则传递给下一个中间件
// 返回:
//   - gin.HandlerFunc: Gin中间件函数
func WebDavHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 判断请求是否是WebDAV请求
		if webdav.IsWebDav(ctx.Request.URL.Path) {
			// 如果是WebDAV请求，交给WebDAV处理器处理
			webdav.ServeHTTP(ctx.Writer, ctx.Request)
			// 中止后续中间件处理
			ctx.Abort()
			return
		} else {
			// 如果不是WebDAV请求，继续执行下一个中间件
			ctx.Next()
		}
	}
}
