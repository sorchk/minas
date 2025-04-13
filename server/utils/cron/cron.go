package cron

import (
	"fmt"
	"server/utils/global"
	"server/utils/logger"

	"github.com/robfig/cron/v3"
)

// SchJob 定时任务结构体
type SchJob struct {
	Id   string   // 任务唯一标识
	Cron string   // 定时表达式，如 "0 0 * * *" 表示每天零点执行
	Job  cron.Job // 实际执行的任务
}

// init 初始化函数，在包被导入时自动执行
func init() {
	// 创建一个新的定时任务调度器
	// WithSeconds(): 支持秒级调度
	// WithChain(cron.Recover): 添加恢复机制，防止任务崩溃影响调度器
	// WithChain(cron.DelayIfStillRunning): 如果任务还在运行，则延迟执行下一次
	global.Cron = cron.New(cron.WithSeconds(), cron.WithChain(cron.Recover(cron.DefaultLogger)), cron.WithChain(cron.DelayIfStillRunning(cron.DefaultLogger)))
	global.Cron.Start() // 启动调度器
}

// mapCron 用于存储已添加的定时任务，键为任务ID，值为用户定义的任务标识
var mapCron = make(map[cron.EntryID]string)

// HashCron 检查指定ID的任务是否已存在
// 参数 id: 任务标识
// 返回值: 若任务已存在则返回true，否则返回false
func HashCron(id any) bool {
	for ei := range mapCron {
		if mapCron[ei] == fmt.Sprint(id) {
			return true
		}
	}
	return false
}

// Start 启动多个定时调度任务
// 参数 jobs: 定时任务列表
// 返回值: 错误信息
func Start(jobs []SchJob) error {
	for _, item := range jobs {
		AddJobTask(item)
	}
	return nil
}

// AddJobTask 添加单个定时调度任务
// 参数 job: 要添加的定时任务
// 返回值: 错误信息
func AddJobTask(job SchJob) error {
	//防止任务添加重复
	if !HashCron(job.Id) && job.Cron != "" {
		cid, err := global.Cron.AddJob(job.Cron, job.Job)
		if err != nil {
			return err
		}
		mapCron[cid] = job.Id
		logger.LOG.Debugf("AddJobID:%v\n", cid)
	}
	return nil
}

// RemoveJobTask 删除指定标识的定时调度任务
// 参数 id: 任务标识
func RemoveJobTask(id any) {
	keys := []cron.EntryID{}
	for ei := range mapCron {
		if mapCron[ei] == fmt.Sprint(id) {
			keys = append(keys, ei)
		}
	}
	if len(keys) > 0 {
		for _, key := range keys {
			global.Cron.Remove(key)
			delete(mapCron, key)
			logger.LOG.Debugf("RemoveJobID:%v\n", key)
		}
	}

}

// UpdateJobTask 更新定时调度任务
// 先删除原有任务，再添加新任务
// 参数 job: 要更新的定时任务
// 返回值: 错误信息
func UpdateJobTask(job SchJob) error {
	RemoveJobTask(job.Id)
	return AddJobTask(job)
}
