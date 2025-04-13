package global

import (
	"server/utils/cache"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// 全局变量定义，整个应用程序中可访问的共享资源
var (
	// CONF  = config.CONF  // 全局配置（当前已注释）
	DB    *gorm.DB                 // 数据库连接实例，用于与数据库交互的全局变量
	Cron  *cron.Cron               // 定时任务调度器，用于管理所有定时任务
	CACHE = cache.GetCacheSystem() // 缓存系统实例，用于存储临时数据
)
