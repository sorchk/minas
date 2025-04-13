// package middleware 提供系统中间件功能
package middleware

import (
	"log"
	"net/http"
	"server/core/app/response" // 导入响应处理模块

	"github.com/duke-git/lancet/v2/convertor" // 导入类型转换工具
	"github.com/gin-gonic/gin"                // 导入Gin框架
)

// UnifiedResponseMiddleware 是处理统一HTTP响应格式的中间件
// 该中间件确保所有API响应都遵循一致的格式，无论是正常响应还是错误响应
// 它会捕获请求处理过程中的错误并将其格式化为统一的错误响应
//
// 返回值:
//   - gin.HandlerFunc: Gin中间件处理函数
func UnifiedResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求路径和状态码到日志，方便调试和监控
		log.Println("UnifiedResponseMiddleware==============" + c.Request.URL.Path + ":" + convertor.ToString(c.Writer.Status()))

		// 继续处理后续的中间件和路由处理函数
		// c.Next()的调用会执行后续的处理函数，完成后再回到这里继续执行
		c.Next()

		// 处理完成后，检查是否在处理请求时发生了错误
		if len(c.Errors) > 0 {
			// 获取最后一个错误并以统一格式返回给客户端
			// 这确保了所有错误都以一致的格式返回
			err := c.Errors.Last()
			response.Error(c, err)
			return
		}

		// 检查是否设置了响应状态码，如果没有则设置为200 OK
		// 这确保了每个响应都有适当的HTTP状态码
		if c.Writer.Status() == 0 {
			c.Writer.WriteHeader(http.StatusOK)
		}
	}
}
