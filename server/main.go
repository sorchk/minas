// package main 定义了程序的入口点
package main

import (
	"log"
	"os"
	"os/signal"
	"server/cmd"         // 命令行处理模块
	"server/utils/tools" // 工具函数模块
	"syscall"
	"time"
)

// main 程序的主函数，负责初始化和启动应用
func main() {

	// 创建一个信号通道，用于接收系统信号
	sigs := make(chan os.Signal, 1)
	// 注册要监听的系统信号：中断、终止和退出
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 启动一个goroutine来处理信号
	go func() {
		// 阻塞等待信号
		sig := <-sigs
		log.Println()
		log.Println(sig)
		log.Println("程序正在清理资源并退出...")
		time.Sleep(1 * time.Second) // 模拟清理资源的过程
		// 执行清理操作
		tools.ExecuteClean()
		// 正常退出程序
		os.Exit(0)
	}()

	// 执行命令行处理
	cmd.Execute()
}
