package cache

import (
	"context"
	"time"
)

// 创建一个全局上下文对象，用于Redis操作
var ctx = context.Background()

// CacheSystem 定义缓存系统的通用接口
// 所有的缓存实现（Redis或内存）都必须实现这个接口
type CacheSystem interface {
	// Get 获取字符串类型的缓存值
	Get(key string) (string, error)

	// Set 设置缓存，无过期时间
	Set(key string, value interface{}) error

	// GetInt 获取整数类型的缓存值
	GetInt(key string) (int, error)

	// Exists 检查键是否存在
	Exists(key string) (bool, error)

	// Incr 将键对应的值加1
	Incr(key string) (int64, error)

	// Decr 将键对应的值减1
	Decr(key string) (int64, error)

	// Delete 删除缓存
	Delete(key string) (bool, error)

	// IncrExpire 将键对应的值加1并设置过期时间
	IncrExpire(key string, ttl int64) (int64, error)

	// DecrExpire 将键对应的值减1并设置过期时间
	DecrExpire(key string, ttl int64) (int64, error)

	// SetExpire 设置缓存并指定过期时间（秒）
	SetExpire(key string, value interface{}, ttl int64) error

	// GetTTL 获取键的剩余生存时间
	GetTTL(key string) (time.Duration, error)
}

// GetCacheSystem 根据配置返回适当的缓存系统实现
// 如果Redis可用则使用Redis，否则使用内存缓存
func GetCacheSystem() CacheSystem {
	if SupportRedis() {
		return RedisCacheSystem("")
	} else {
		return MemCacheSystem("")
	}
}
