package logger

import (
	"fmt"
	"io"
	"os"
	"server/utils/config"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// 日志格式相关常量定义
const (
	TimeFormat         = "2006-01-02 15:04:05" // 标准时间格式，用于日志输出中的时间戳
	FileTImeFormat     = "2006-01-02"          // 文件时间格式，用于日志文件名中的日期
	RollingTimePattern = "0 0  * * *"          // 日志滚动时间模式，定义为每天0点
)

// LOG 全局日志记录器实例，可在整个应用程序中使用
var LOG *logrus.Logger

// Init 初始化日志系统
// 创建新的日志记录器并配置输出
func Init() {
	LOG = logrus.New()                    // 创建新的logrus日志实例
	setOutput(LOG, config.CONF.LogConfig) // 根据配置设置日志输出
	LOG.Info("init logger successfully")  // 记录初始化成功的信息
}

// setOutput 配置日志输出
// 参数 logger: 要配置的日志记录器
// 参数 conf: 日志配置信息
func setOutput(logger *logrus.Logger, conf config.LogConfig) {
	// 根据配置创建日志写入器
	writer, err := NewWriterFromConfig(&LoggerConfig{
		LogPath:            conf.LogPath,       // 日志文件存储路径
		FileName:           conf.LogName,       // 日志文件名称
		TimeTagFormat:      FileTImeFormat,     // 时间标签格式
		MaxRemain:          conf.MaxBackup,     // 最大保留的日志文件数量
		RollingTimePattern: RollingTimePattern, // 日志滚动时间模式
		LogSuffix:          conf.LogSuffix,     // 日志文件后缀
	})
	if err != nil {
		panic(err) // 如果创建写入器失败，则触发panic
	}

	// 解析日志级别
	level, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		panic(err) // 如果解析级别失败，则触发panic
	}

	// 创建多重写入器，同时输出到文件和标准输出
	fileAndStdoutWriter := io.MultiWriter(writer, os.Stdout)

	// 配置日志记录器
	logger.SetOutput(fileAndStdoutWriter)   // 设置输出目标
	logger.SetLevel(level)                  // 设置日志级别
	logger.SetFormatter(new(MineFormatter)) // 设置自定义格式化器
}

// MineFormatter 自定义日志格式化器，实现logrus.Formatter接口
type MineFormatter struct{}

// Format 实现了Formatter接口的Format方法，定义日志条目的格式化逻辑
// 参数 entry: 需要格式化的日志条目
// 返回值: 格式化后的字节数组和可能的错误
func (s *MineFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	detailInfo := ""
	if len(entry.Data) == 0 {
		// 如果没有额外的数据字段，则使用简单格式
		msg := fmt.Sprintf("[%s] [%s] %s %s \n", time.Now().Format(TimeFormat), strings.ToUpper(entry.Level.String()), entry.Message, detailInfo)
		return []byte(msg), nil
	}
	// 如果有额外的数据字段，将其添加到日志消息中
	msg := fmt.Sprintf("[%s] [%s] %s %s {%v} \n", time.Now().Format(TimeFormat), strings.ToUpper(entry.Level.String()), entry.Message, detailInfo, entry.Data)
	return []byte(msg), nil
}
