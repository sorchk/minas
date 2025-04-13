package fileclean

import (
	"fmt"
	"server/service/scheduled/job"
	"server/service/scheduled/log"
	"server/utils/logger"
	"server/utils/rclone"
	"time"

	"github.com/duke-git/lancet/v2/datetime"
	"github.com/duke-git/lancet/v2/slice"
)

// FileCleanJob 文件清理任务结构体
// 定义了文件清理任务的配置和行为，用于自动清理旧文件
// 通过嵌入SchJob获得计划任务的基本属性和行为
type FileCleanJob struct {
	job.SchJob          // 嵌入基础计划任务结构体，继承其属性和方法
	WorkDir    string   // 工作目录，需要清理的目录路径
	Includes   []string // 包含的文件模式，用于筛选需要清理的文件
	Excludes   []string // 排除的文件模式，用于筛选不需要清理的文件
	KeepNum    int      // 保留文件数量，最少保留的文件数
	OffsetDay  int      // 偏移天数，保留最近多少天的文件
	IsTest     bool     // 测试模式，如果为true则只显示要删除的文件但不实际删除
	ShowDebug  bool     // 是否显示调试信息
}

// Run 执行文件清理任务
// 根据配置的规则清理文件，并记录任务执行日志
// 支持根据文件数量和日期进行清理
func (job FileCleanJob) Run() {
	schLog := log.SchLog{} // 初始化日志对象
	var logs = []string{}  // 日志内容数组

	// 创建并启动日志文件
	workDir, logPath, err := schLog.StartLogFile(job.SchJob)
	if err != nil {
		logger.LOG.Errorf("发生错误无法创建任务数据: %s\n", err.Error())
		return
	}

	// 获取文件列表
	fileList, err := rclone.List(job.WorkDir, job.Includes, job.Excludes, 0, 0)
	if err != nil {
		logs = append(logs, "ERROR: "+err.Error())
		schLog.Error(job.SchJob, logs)
		return
	}

	// 如果文件数量小于等于保留数量，则不执行清理
	if len(fileList) <= job.KeepNum {
		logs = append(logs, fmt.Sprintf("文件数小于最少保留数，不执行清理(%d/%d)", len(fileList), job.KeepNum))
		logs = append(logs, "SUCCESS!")
		schLog.Success(job.SchJob, logs)
		return
	}

	// 排序文件列表，按时间倒序和名称倒序
	less := func(a, b rclone.FileObj) bool {
		if a.ModTime == b.ModTime {
			return a.Name > b.Name
		}
		return a.ModTime > b.ModTime
	}
	slice.SortBy(fileList, less)

	// 根据保留文件数获取需要清理的日期
	startFile := fileList[int(job.KeepNum)]
	t, _ := time.Parse(time.RFC3339Nano, startFile.ModTime)
	minTime := time.Since(t)
	keepTime := time.Since(datetime.AddDay(time.Now(), -int64(job.OffsetDay)))
	if keepTime > minTime {
		minTime = keepTime
	}

	// 执行文件删除命令
	err = rclone.CmdDelete(workDir, logPath, job.WorkDir, job.Includes, job.Excludes, minTime, job.IsTest, job.ShowDebug)
	if err != nil {
		logs = append(logs, "ERROR: "+err.Error())
		schLog.Error(job.SchJob, logs)
	} else {
		logs = append(logs, "SUCCESS!")
		schLog.Success(job.SchJob, logs)
	}
}
