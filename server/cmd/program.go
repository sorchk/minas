package cmd

import (
	"log"
	"os"
	"runtime"
	"server/utils"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

// Program 结构体定义了服务程序
// exit 通道用于通知服务停止
type Program struct {
	exit chan struct{}
}

// Start 方法在服务启动时调用
// 参数 s 是服务实例
// 返回启动过程中的错误，如果没有错误则返回 nil
func (p *Program) Start(s service.Service) error {
	log.Println("启动服务...", service.Platform())
	p.exit = make(chan struct{}) // 初始化退出通道
	go p.run()                   // 在新的协程中运行服务主体逻辑
	log.Println("服务启动完成.")
	return nil
}

// run 方法包含服务的主要业务逻辑
func (p *Program) run() {
	log.Println("服务运行中...")
	InitConfig(nil) // 初始化配置
	StartServer()   // 启动服务器
}

// Stop 方法在服务停止时调用
// 参数 s 是服务实例
// 返回停止过程中的错误，如果没有错误则返回 nil
func (p *Program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	log.Println("停止服务.")
	close(p.exit) // 关闭退出通道，通知各个协程停止运行
	return nil
}

// init 函数在包初始化时执行，用于添加命令到根命令
func init() {
	// 添加命令
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(runCmd)
}

// runCmd 定义了运行服务的命令
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "启动服务",
	Long:  `启动服务.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// 命令执行前的准备工作
	},
	Run: func(cmd *cobra.Command, args []string) {
		ControlService("run") // 控制服务执行运行操作
	},
}

// installCmd 定义了安装服务的命令
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "安装服务",
	Long:  `安装服务.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// 命令执行前的准备工作
	},
	Run: func(cmd *cobra.Command, args []string) {
		ControlService("install") // 控制服务执行安装操作
	},
}

// uninstallCmd 定义了卸载服务的命令
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "卸载服务",
	Long:  `卸载服务.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// 命令执行前的准备工作
	},
	Run: func(cmd *cobra.Command, args []string) {
		ControlService("uninstall") // 控制服务执行卸载操作
	},
}

// startCmd 定义了启动服务的命令
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动服务",
	Long:  `启动服务.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// 命令执行前的准备工作
	},
	Run: func(cmd *cobra.Command, args []string) {
		ControlService("start") // 控制服务执行启动操作
	},
}

// stopCmd 定义了停止服务的命令
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止服务",
	Long:  `停止服务.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// 命令执行前的准备工作
	},
	Run: func(cmd *cobra.Command, args []string) {
		ControlService("stop") // 控制服务执行停止操作
	},
}

// ControlService 根据指定的操作控制服务
// 参数 action 表示要执行的操作：run, install, uninstall, start, stop
func ControlService(action string) {
	s, err := initService() // 初始化服务
	if err != nil {
		log.Println("初始化服务失败：" + err.Error())
	}
	switch action {
	case "run":
		s.Run() // 直接运行服务
	default:
		// 使用 service 包的 Control 函数执行其他操作
		err := service.Control(s, action)
		if err != nil {
			log.Printf("%s not valid, Valid actions: %q\n%v", action, service.ControlAction, err)
		}
		return
	}
}

// initService 初始化并返回服务实例
// 返回服务实例和可能的错误
func initService() (service.Service, error) {
	os.Chdir(utils.GetExeFileDirectory()) // 切换到可执行文件所在目录

	// 服务的配置信息
	options := make(service.KeyValue)
	options["LogOutput"] = true
	options["HasOutputFileSupport"] = true
	options["WorkingDirectory"] = utils.GetWorkDirectory()
	options["CurrentDirectory"] = utils.GetWorkDirectory()

	dependencies := []string{}
	if runtime.GOOS == "windows" {
		// Windows 系统特定配置
	} else {
		// 非 Windows 系统配置
		options["Restart"] = "on-failure"
		options["SuccessExitStatus"] = "1 2 8 SIGKILL"
		dependencies = append(dependencies, "Requires=network.target",
			"After=network-online.target syslog.target")
	}

	// 服务基本配置
	svcConfig := &service.Config{
		Name:         "minas",
		DisplayName:  "Minas服务",
		Description:  "支持webdav、文件备份、日志或备份清理、超级终端管理",
		Option:       options,
		Dependencies: dependencies,
		Arguments:    []string{"run"}, // 运行时的默认参数
	}

	if runtime.GOOS == "windows" {
		// Windows 系统特定服务配置
	} else {
		// 非 Windows 系统特定服务配置
		svcConfig.Dependencies = []string{
			"Requires=network.target",
			"After=network-online.target syslog.target"}
		svcConfig.UserName = "root" // 使用 root 用户运行服务
	}

	pro := &Program{}                     // 创建程序实例
	s, err := service.New(pro, svcConfig) // 创建服务
	return s, err
}
