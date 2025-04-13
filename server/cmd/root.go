// package cmd 包含命令行相关功能的实现
package cmd

import (
	"log"                 // 日志包
	"server/utils/config" // 系统配置包

	"github.com/spf13/cobra" // 命令行工具库
)

var (
	// VersionShow 是否显示版本信息的标志，默认为false
	VersionShow bool = false
)

// init 初始化函数，在包被导入时自动执行
func init() {
	// 添加全局参数 --version/-v 用于显示版本信息
	rootCmd.PersistentFlags().BoolVarP(&VersionShow, "version", "v", false, "查看版本")
}

// rootCmd 表示在没有指定子命令时的基本命令
// 这是程序的根命令，所有其他命令都是它的子命令
var rootCmd = &cobra.Command{
	Use:   "",                                    // 命令的使用方式
	Short: "Minas" + config.Version,              // 简短描述，包含版本号
	Long:  `Minas 系统  Version:` + config.Version, // 详细描述
	Run: func(cmd *cobra.Command, args []string) {
		// 当直接运行根命令时的行为
		if VersionShow {
			// 如果指定了--version标志，显示版本信息
			log.Printf("版本号: %s\n", config.Version)
		} else {
			// 否则显示帮助信息
			cmd.Help()
		}
	},
}

// Execute 执行根命令
// 这个函数由 main.main() 调用，用于启动命令行应用
// 它会将所有子命令添加到root命令并适当设置标志
// 在程序生命周期内只需要调用一次
func Execute() {
	// 执行根命令并检查错误
	cobra.CheckErr(rootCmd.Execute())
}
