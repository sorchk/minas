package filebackup

import (
	"fmt"
	"server/service/scheduled/job"
	"server/service/scheduled/log"
	"server/utils/logger"
	"server/utils/rclone"
	"time"
)

// FileBackupJob 文件备份任务结构体
// 定义了文件备份任务的配置和行为，支持多种备份同步策略
// 通过嵌入SchJob获得计划任务的基本属性和行为
type FileBackupJob struct {
	job.SchJob                  // 嵌入基础计划任务结构体，继承其属性和方法
	Type               uint     `gorm:"comment:'同步类型1备份2镜像3双向同步'  default:1 " json:"type"` // 备份类型：1=仅备份，2=镜像，3=双向同步，4=完整备份
	Src                string   // 源路径，备份的来源目录
	Dst                string   // 目标路径，备份的目标目录
	CreateEmptySrcDirs bool     // 是否创建空的源目录
	RemoveEmptyDirs    bool     // 是否删除空目录
	ShowDebug          bool     // 是否显示调试信息
	Includes           []string // 包含的文件模式，用于筛选需要备份的文件
	Excludes           []string // 排除的文件模式，用于筛选不需要备份的文件
}

// Run 执行文件备份任务
// 根据配置的备份类型选择不同的备份策略
// 记录任务执行日志和错误信息
func (job FileBackupJob) Run() {
	schLog := log.SchLog{} // 初始化日志对象
	var logs = []string{}  // 日志内容数组

	// 创建并启动日志文件
	workDir, logPath, err := schLog.StartLogFile(job.SchJob)
	if err != nil {
		// 创建日志文件失败时记录错误并退出
		logger.LOG.Errorf("发生错误无法创建任务数据: %s\n", err.Error())
		return
	}

	// 根据配置的备份类型执行不同的备份操作
	switch job.Type {
	case 1:
		// 仅复制文件，不删除目标位置已有文件
		err = rclone.CmdCopy(workDir, logPath, job.Src, job.Dst, job.CreateEmptySrcDirs, job.ShowDebug, job.Includes, job.Excludes)
	case 2:
		// 镜像复制，会删除目标位置有而源位置没有的文件
		err = rclone.CmdSync(workDir, logPath, job.Src, job.Dst, job.CreateEmptySrcDirs, job.ShowDebug, job.Includes, job.Excludes)
	case 3:
		// 双向同步，使源和目标位置的文件保持一致
		err = rclone.CmdBisync(workDir, logPath, job.Src, job.Dst, job.CreateEmptySrcDirs, job.RemoveEmptyDirs, job.ShowDebug, job.Includes, job.Excludes)
	case 4:
		// 完整备份，创建带有时间戳的新目录进行完整复制
		err = rclone.CmdCopy(workDir, logPath, job.Src, fmt.Sprintf("%s/%s", job.Dst, time.Now().Format("20060102_150405")), job.CreateEmptySrcDirs, job.ShowDebug, job.Includes, job.Excludes)
	default:
		// 未知的备份类型，不执行任何操作
	}

	// 根据执行结果记录日志
	if err != nil {
		logs = append(logs, "ERROR: "+err.Error()) // 记录错误信息
		schLog.Error(job.SchJob, logs)             // 将错误日志写入日志文件
	} else {
		logs = append(logs, "SUCCESS!")  // 记录成功信息
		schLog.Success(job.SchJob, logs) // 将成功日志写入日志文件
	}
}
