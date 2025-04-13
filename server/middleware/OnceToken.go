package middleware

import (
	"fmt"
	"server/core/app/request"
	"server/core/app/response"
	"server/utils/cache"
	"server/utils/data"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AUTH_TIME_OFFSET_SECOND 定义一次性令牌的有效期，单位为秒
const AUTH_TIME_OFFSET_SECOND = 30

// GetOnceToken 生成一次性令牌的中间件
// 该函数会先进行API密钥认证，然后为认证通过的用户生成一个一次性令牌
// 令牌将被存储在缓存系统中，并设置过期时间
func GetOnceToken(ctx *gin.Context) {
	// 首先验证API密钥
	AuthApiKeyMiddleware(ctx)
	// 获取当前用户ID
	id := request.GetUserID(ctx)
	if id > 0 {
		// 将用户ID转换为字符串
		userId := fmt.Sprint(id)
		// 生成唯一的UUID作为令牌
		token := uuid.New().String()
		// 构造缓存键名
		key := "auth:once:token:" + token
		// 将令牌与用户ID关联并设置过期时间
		cache.GetCacheSystem().SetExpire(key, userId, AUTH_TIME_OFFSET_SECOND)
		// 返回令牌和用户ID给客户端
		response.Data(ctx, "", data.Map{"token": token, "userId": userId})
	}
}

// VerificationOnceToken 验证一次性令牌的有效性
// 传入用户ID和令牌，验证令牌是否有效且属于该用户
// 返回布尔值表示验证结果
func VerificationOnceToken(userId string, token string) bool {
	// 构造缓存键名
	key := "auth:once:token:" + token
	// 从缓存中获取令牌对应的用户ID
	val, err := cache.GetCacheSystem().Get(key)
	if err != nil {
		return false
	}
	// 比较存储的用户ID是否与提供的用户ID一致
	return val == userId
}
