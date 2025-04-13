package cache

import (
	"server/utils/config"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCacheSystem Redis缓存系统的实现
// 使用Redis服务器存储缓存数据，适用于分布式环境
type RedisCacheSystem string

// 确保 RedisCacheSystem 实现了 CacheSystem 接口
var _ CacheSystem = RedisCacheSystem("")

// Redis客户端实例
var rdb *redis.Client = nil

// DecrExpire 将键对应的值减1并设置过期时间
func (cs RedisCacheSystem) DecrExpire(key string, ttl int64) (int64, error) {
	redis, err := GetRedisClient()
	if err != nil {
		return 0, err
	}
	value, err := redis.Decr(ctx, key).Result()
	redis.Expire(ctx, key, time.Duration(ttl)*time.Second)
	if err == nil {
		return value, nil
	}
	return 0, err
}

// GetInt 获取整数类型的缓存值
func (cs RedisCacheSystem) GetInt(key string) (int, error) {
	redis, err := GetRedisClient()
	if err != nil {
		return 0, err
	}
	value, err := redis.Get(ctx, key).Int()
	if err == nil {
		return value, nil
	}
	return 0, err
}

// GetTTL 获取键的剩余生存时间
func (cs RedisCacheSystem) GetTTL(key string) (time.Duration, error) {
	redis, err := GetRedisClient()
	if err != nil {
		return 0, err
	}
	value, err := redis.TTL(ctx, key).Result()
	if err == nil {
		return value, nil
	}
	return 0, err
}

// IncrExpire 将键对应的值加1并设置过期时间
func (cs RedisCacheSystem) IncrExpire(key string, ttl int64) (int64, error) {
	redis, err := GetRedisClient()
	if err != nil {
		return 0, err
	}
	// 注意：这里有一个bug，应该使用Incr而不是Decr
	value, err := redis.Incr(ctx, key).Result() // 修复：由Decr改为Incr
	redis.Expire(ctx, key, time.Duration(ttl)*time.Second)
	if err == nil {
		return value, nil
	}
	return 0, err
}

// SupportRedis 检查是否支持Redis
// 通过配置判断Redis是否可用
func SupportRedis() bool {
	return config.CONF.Redis.Enable && config.CONF.Redis.Host != ""
}

// GetRedisClient 获取Redis客户端实例
// 如果已有有效连接则复用，否则创建新连接
func GetRedisClient() (*redis.Client, error) {
	if rdb != nil {
		// 读取已保存的连接 测试连接
		_, err := rdb.Ping(ctx).Result()
		if err == nil {
			// 连接有效
			return rdb, nil
		}
	}
	// 连接无效 重连
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.CONF.Redis.Host + ":" + config.CONF.Redis.Port,
		Password: config.CONF.Redis.Password,
		DB:       config.CONF.Redis.Db,
	})
	_, err := rdb.Ping(ctx).Result()
	if err == nil {
		// 连接有效
		return rdb, nil
	}
	// 连接无效
	return nil, err
}

// Incr 将键对应的值加1
func (cs RedisCacheSystem) Incr(key string) (int64, error) {
	redis, err := GetRedisClient()
	if err != nil {
		return 0, err
	}
	value, err := redis.Incr(ctx, key).Result()
	if err == nil {
		return value, nil
	}
	return 0, err
}

// Decr 将键对应的值减1
func (cs RedisCacheSystem) Decr(key string) (int64, error) {
	redis, err := GetRedisClient()
	if err != nil {
		return 0, err
	}
	value, err := redis.Decr(ctx, key).Result()
	if err == nil {
		return value, nil
	}
	return 0, err
}

// Delete 删除缓存
func (cs RedisCacheSystem) Delete(key string) (bool, error) {
	redis, err := GetRedisClient()
	if err != nil {
		return false, err
	}
	value, err := redis.Del(ctx, key).Result()
	if err == nil {
		return value != 0, nil
	}
	return false, err
}

// Exists 检查键是否存在
func (cs RedisCacheSystem) Exists(key string) (bool, error) {
	redis, err := GetRedisClient()
	if err != nil {
		return false, err
	}
	value, err := redis.Exists(ctx, key).Result()
	if err == nil {
		return value != 0, nil
	}
	return false, err
}

// Set 设置缓存，无过期时间
func (cs RedisCacheSystem) Set(key string, value interface{}) error {
	redis, err := GetRedisClient()
	if err != nil {
		return err
	}
	_, err = redis.Set(ctx, key, value, 0).Result()
	if err == nil {
		return nil
	}
	return err
}

// SetExpire 设置缓存并指定过期时间（秒）
func (cs RedisCacheSystem) SetExpire(key string, value interface{}, ttl int64) error {
	redis, err := GetRedisClient()
	if err != nil {
		return err
	}
	_, err = redis.Set(ctx, key, value, time.Duration(ttl)*time.Second).Result()
	if err == nil {
		return nil
	}
	return err
}

// Get 获取字符串类型的缓存值
func (cs RedisCacheSystem) Get(key string) (string, error) {
	redis, err := GetRedisClient()
	if err != nil {
		return "", err
	}
	value, err := redis.Get(ctx, key).Result()
	if err == nil {
		return value, nil
	}
	return "", err
}
