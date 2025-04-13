package cache

import (
	"errors"
	"server/utils/cache/local"
	"strconv"
	"time"
)

// MemCacheSystem 内存缓存系统的实现
// 使用本地内存存储缓存数据，适用于单机部署场景
type MemCacheSystem string

// 确保 MemCacheSystem 实现了 CacheSystem 接口
var _ CacheSystem = MemCacheSystem("")

// 定义清理过期缓存的时间间隔
const cleanupInterval = 10 * time.Minute

// 默认缓存实例，用于存储所有缓存数据
var defCache = local.New(0, cleanupInterval)

// Decr 将键对应的值减1
func (cs MemCacheSystem) Decr(key string) (int64, error) {
	return defCache.DecrementInt64(key, 1, local.NoExpiration)
}

// GetInt 获取整数类型的缓存值
// 支持多种类型转换为int：int、int64、int32、uint、string等
func (cs MemCacheSystem) GetInt(key string) (int, error) {
	if x, found := defCache.Get(key); found {
		switch x.(type) {
		case int:
			return x.(int), nil
		case int64:
			return int(x.(int64)), nil
		case int32:
			return int(x.(int32)), nil
		case uint:
			return int(x.(uint)), nil
		case string:
			i, _ := strconv.Atoi(x.(string))
			return i, nil
		default:
			return 0, nil
		}
	}
	return 0, nil
}

// GetTTL 获取键的剩余生存时间
func (cs MemCacheSystem) GetTTL(key string) (time.Duration, error) {
	if item, found := defCache.GetItem(key); found {
		return time.Duration(time.Duration(item.Expiration - time.Now().UnixNano())), nil
	}
	return time.Duration(0), errors.New("Key not found")
}

// Incr 将键对应的值加1
func (cs MemCacheSystem) Incr(key string) (int64, error) {
	return defCache.IncrementInt64(key, 1, local.NoExpiration)
}

// DecrExpire 将键对应的值减1并设置过期时间
func (cs MemCacheSystem) DecrExpire(key string, ttl int64) (int64, error) {
	return defCache.DecrementInt64(key, 1, time.Duration(ttl)*time.Second)
}

// IncrExpire 将键对应的值加1并设置过期时间
func (cs MemCacheSystem) IncrExpire(key string, ttl int64) (int64, error) {
	return defCache.IncrementInt64(key, 1, time.Duration(ttl)*time.Second)
}

// Delete 删除缓存
func (cs MemCacheSystem) Delete(key string) (bool, error) {
	defCache.Delete(key)
	return true, nil
}

// Exists 检查键是否存在
func (cs MemCacheSystem) Exists(key string) (bool, error) {
	_, found := defCache.Get(key)
	return found, nil
}

// Set 设置缓存，无过期时间
func (cs MemCacheSystem) Set(key string, value interface{}) error {
	defCache.Set(key, value, local.NoExpiration)
	return nil
}

// SetExpire 设置缓存并指定过期时间（秒）
func (cs MemCacheSystem) SetExpire(key string, value interface{}, ttl int64) error {
	defCache.Set(key, value, time.Duration(ttl)*time.Second)
	return nil
}

// Get 获取字符串类型的缓存值
func (cs MemCacheSystem) Get(key string) (string, error) {
	if x, found := defCache.Get(key); found {
		result := x.(string)
		return result, nil
	}
	return "", nil
}
