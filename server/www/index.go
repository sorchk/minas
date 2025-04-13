// package www 包含网站前端静态资源和HTML文件
package www

import "embed" // 导入embed包，用于将静态文件嵌入到二进制文件中

//go:embed dist/static
var StaticFs embed.FS // 嵌入静态资源文件夹，包含JS、CSS等静态资源

//go:embed dist/index.html
var IndexByte []byte // 嵌入主页HTML文件，以字节数组形式存储

//go:embed dist/static/index.html
var NoUIndexByte []byte // 嵌入备用HTML文件，可能用于无UI模式或特殊场景
