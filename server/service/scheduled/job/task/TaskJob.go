package task

import (
	"server/service/scheduled/job"
	"server/service/scheduled/log"
)

// TaskJob 任务作业结构体
// 定义了任务作业的配置和行为，用于执行自定义任务
// 通过嵌入SchJob获得计划任务的基本属性和行为
type TaskJob struct {
	job.SchJob        // 嵌入基础计划任务结构体，继承其属性和方法
	Script     string `gorm:"comment:'脚本'" json:"script"` // 要执行的脚本内容
}

// Run 执行任务作业
// 记录任务开始和完成状态
// 目前为占位实现，实际执行逻辑待实现
func (job TaskJob) Run() {
	schLog := log.SchLog{}   // 初始化日志对象
	var logs = []string{}    // 日志内容数组
	schLog.Start(job.SchJob) // 记录任务开始状态
	//TODO 这里需要实现具体的任务执行逻辑
	schLog.Success(job.SchJob, logs) // 记录任务成功完成
}
