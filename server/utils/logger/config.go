package logger

import (
	"errors"
	"io"
	"os"
	"path"
)

// 日志系统全局变量定义
var (
	BufferSize         = 0x100000                              // 日志缓冲区大小，1MB
	DefaultFileMode    = os.FileMode(0644)                     // 日志文件默认权限模式，可读可写
	DefaultFileFlag    = os.O_RDWR | os.O_CREATE | os.O_APPEND // 文件打开标志，读写、创建、追加模式
	ErrInvalidArgument = errors.New("error argument invalid")  // 无效参数错误
	QueueSize          = 1024                                  // 日志队列大小
	ErrClosed          = errors.New("error write on close")    // 关闭后写入错误
)

// LoggerConfig 日志配置结构体，包含日志记录器的所有配置选项
type LoggerConfig struct {
	TimeTagFormat      string // 时间标签格式，如 "2006-01-02 15:04:05"
	LogPath            string // 日志文件存储路径
	FileName           string // 日志文件名称
	LogSuffix          string // 日志文件后缀，如 ".log"
	MaxRemain          int    // 最大保留的日志文件数量
	RollingTimePattern string // 日志滚动时间模式，如 "0 0 0 * * ?" 表示每天零点
}

// Manager 日志管理器接口，定义了日志管理的基本操作
type Manager interface {
	Fire() chan string // 获取日志消息通道
	Close()            // 关闭日志管理器
}

// RollingWriter 滚动日志写入器接口，实现日志文件的自动滚动
type RollingWriter interface {
	io.Writer     // 嵌入io.Writer接口，提供Write方法
	Close() error // 关闭写入器
}

// FilePath 根据配置生成日志文件的完整路径
// 参数 c: 日志配置
// 返回值: 完整的日志文件路径
func FilePath(c *LoggerConfig) (filepath string) {
	filepath = path.Join(c.LogPath, c.FileName) + c.LogSuffix
	return
}
