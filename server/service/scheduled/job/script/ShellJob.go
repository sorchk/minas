package shell

import (
	"server/service/scheduled/job"
	"server/service/scheduled/log"
	"server/utils/cmd"
	"server/utils/logger"
	"time"
)

// ShellJob Shell脚本任务结构体
// 定义了Shell脚本任务的配置和行为，用于执行自定义Shell命令或脚本
// 通过嵌入SchJob获得计划任务的基本属性和行为
type ShellJob struct {
	job.SchJob        // 嵌入基础计划任务结构体，继承其属性和方法
	Script     string `gorm:"comment:'脚本'" json:"script"` // 要执行的Shell脚本内容
}

// Run 执行Shell脚本任务
// 在指定的工作目录下执行配置的Shell脚本
// 设置超时时间并记录执行日志和错误信息
func (job ShellJob) Run() {
	schLog := log.SchLog{} // 初始化日志对象
	var logs = []string{}  // 日志内容数组

	// 创建并启动日志文件
	workDir, logPath, err := schLog.StartLogFile(job.SchJob)
	if err != nil {
		logger.LOG.Errorf("发生错误无法创建任务数据: %s\n", err.Error())
		return
	}

	// 执行Shell脚本，设置5小时超时
	err = cmd.ExecCronjobWithTimeOut(job.Script, workDir, logPath, time.Hour*5)
	if err != nil {
		logs = append(logs, "ERROR: "+err.Error())
		schLog.Error(job.SchJob, logs)
	} else {
		logs = append(logs, "SUCCESS!")
		schLog.Success(job.SchJob, logs)
	}
}
