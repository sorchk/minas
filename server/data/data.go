package data

import (
	_ "embed" // 导入embed包，用于嵌入静态资源
	"log"
	"path/filepath"
	"server/utils"
)

// ConfigFilePath 配置文件路径
var ConfigFilePath = "data/config.yaml"

// WorkDir 工作目录路径
var WorkDir = "data/workdir"

// DbPath 数据库文件路径
var DbPath = "data/minas.db"

// init 初始化函数，在包被导入时自动执行
func init() {
	// 打印可执行文件所在目录
	log.Println("BinPath:" + utils.GetExeFileDirectory())
	// 打印当前工作目录
	log.Println("WorkDirectory:" + utils.GetWorkDirectory())
	// 获取配置文件的绝对路径并打印
	path, _ := filepath.Abs(ConfigFilePath)
	log.Println("ConfigFilePath:" + path)
}

// DefaultYamlData 嵌入的默认配置文件内容
// 通过go:embed指令将config.yaml文件内容嵌入到二进制文件中
//
//go:embed config.yaml
var DefaultYamlData []byte
