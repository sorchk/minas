package logger

import (
	"path"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// manager 日志管理器结构体，实现了Manager接口
// 负责管理日志文件的生成和滚动
type manager struct {
	startAt time.Time      // 日志开始时间，用于生成日志文件名
	fire    chan string    // 通知通道，用于发送新日志文件名
	cr      *cron.Cron     // 定时任务调度器，用于定时触发日志滚动
	context chan int       // 上下文控制通道，用于关闭管理器
	wg      sync.WaitGroup // 等待组，确保所有goroutine完成
	lock    sync.Mutex     // 互斥锁，保护共享资源的并发访问
}

// Fire 返回通知通道，外部可通过此通道接收日志滚动事件
// 返回值: 日志文件名通知通道
func (m *manager) Fire() chan string {
	return m.fire
}

// Close 关闭日志管理器
// 关闭上下文通道并停止定时任务
func (m *manager) Close() {
	close(m.context)
	m.cr.Stop()
}

// NewManager 创建新的日志管理器
// 参数 c: 日志配置
// 返回值: 日志管理器接口和可能的错误
func NewManager(c *LoggerConfig) (Manager, error) {
	// 初始化管理器结构
	m := &manager{
		startAt: time.Now(),        // 使用当前时间作为起始时间
		cr:      cron.New(),        // 创建新的定时任务调度器
		fire:    make(chan string), // 创建通知通道
		context: make(chan int),    // 创建上下文控制通道
		wg:      sync.WaitGroup{},  // 创建等待组
	}

	// 添加定时任务，根据配置的时间模式触发日志滚动
	if _, err := m.cr.AddFunc(c.RollingTimePattern, func() {
		m.fire <- m.GenLogFileName(c) // 发送新日志文件名到通知通道
	}); err != nil {
		return nil, err
	}
	m.cr.Start() // 启动定时任务调度器

	return m, nil
}

// GenLogFileName 生成日志文件名
// 根据配置和当前时间生成唯一的日志文件名
// 参数 c: 日志配置
// 返回值: 生成的日志文件名
func (m *manager) GenLogFileName(c *LoggerConfig) (filename string) {
	m.lock.Lock() // 加锁保护共享资源
	// 生成格式为 路径/文件名-时间戳.后缀 的文件名
	filename = path.Join(c.LogPath, c.FileName+"-"+m.startAt.Format(c.TimeTagFormat)) + c.LogSuffix
	m.startAt = time.Now() // 更新开始时间为当前时间
	m.lock.Unlock()        // 解锁
	return
}
