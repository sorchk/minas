package route

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// INDEX 定义默认的索引文件名称
const INDEX = "index.html"

// ServeFileSystem 定义了一个文件系统服务接口
// 该接口扩展了标准http.FileSystem接口，并添加了文件存在性检查功能
type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

// Serve 返回一个中间件处理函数，用于提供静态文件服务
// 参数:
//   - urlPrefix: URL前缀，用于从请求路径中剥离
//   - fs: 实现了ServeFileSystem接口的文件系统对象
//
// 返回:
//   - gin.HandlerFunc: Gin中间件函数
func Serve(urlPrefix string, fs ServeFileSystem) gin.HandlerFunc {
	// 创建标准的HTTP文件服务器
	fileserver := http.FileServer(fs)
	// 如果指定了URL前缀，则在处理请求时剥离该前缀
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		// 检查请求的文件是否存在
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			// 文件存在，使用文件服务器处理请求
			fileserver.ServeHTTP(c.Writer, c.Request)
			// 中止后续中间件的执行
			c.Abort()
			return
		}
		// 文件不存在，继续执行下一个中间件
		c.Next()
	}
}

// embedFileSystem 是一个包装了http.FileSystem的结构体
// 用于处理嵌入式文件系统的文件访问
type embedFileSystem struct {
	http.FileSystem
}

// Exists 检查指定路径的文件是否存在于嵌入式文件系统中
// 参数:
//   - prefix: URL前缀(在此实现中未使用，但保留以满足接口要求)
//   - path: 要检查的文件路径
//
// 返回:
//   - bool: 如果文件存在返回true，否则返回false
func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	return err == nil
}

// EmbedFolder 将嵌入式文件系统转换为ServeFileSystem接口
// 参数:
//   - fsEmbed: 嵌入式文件系统
//   - targetPath: 目标路径，指定要使用的子目录
//
// 返回:
//   - ServeFileSystem: 可用于静态文件服务的文件系统接口实现
func EmbedFolder(fsEmbed embed.FS, targetPath string) ServeFileSystem {
	// 获取嵌入式文件系统中的子文件系统
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	// 包装为embedFileSystem并返回
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}
