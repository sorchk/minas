package config

// LogConfig 日志配置结构体，用于定义和管理应用程序日志的各种参数
type LogConfig struct {
	LogPath   string `mapstructure:"log-path" json:"logPath" yaml:"log-path"`       // 日志文件存储路径
	Level     string `mapstructure:"level" json:"level" yaml:"level"`               // 日志记录级别(如:debug, info, warn, error)
	TimeZone  string `mapstructure:"time-zone" json:"timeZone" yaml:"time-zone"`    // 日志时间戳的时区设置
	LogName   string `mapstructure:"log-name" json:"logName" yaml:"log-name"`       // 日志文件基础名称
	LogSuffix string `mapstructure:"log-suffix" json:"logSuffix" yaml:"log-suffix"` // 日志文件扩展名后缀
	MaxBackup int    `mapstructure:"max-backup" json:"maxBackup" yaml:"max-backup"` // 保留的最大日志文件数量
}
