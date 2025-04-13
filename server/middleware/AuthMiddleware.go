// package middleware 提供系统中间件功能
// 包含身份验证、授权、登录等安全相关功能
package middleware

import (
	"errors"                           // 错误处理
	"fmt"                              // 格式化
	"net"                              // 网络相关
	"server/core/app/request"          // 请求处理
	"server/core/app/response"         // 响应处理
	"server/service/basic"             // 基础服务
	"server/utils/cache"               // 缓存工具
	toolCaptcha "server/utils/captcha" // 验证码工具
	"server/utils/config"              // 配置工具
	"server/utils/data"                // 数据工具
	"server/utils/logger"              // 日志工具
	"server/utils/mfa"                 // 多因素认证
	"server/utils/simple"              // 简单工具
	"server/utils/xxtea"               // 加密工具
	"strings"                          // 字符串处理
	"time"                             // 时间处理

	"github.com/dgrijalva/jwt-go" // JWT库，用于处理令牌
	"github.com/gin-gonic/gin"    // Web框架
	"golang.org/x/crypto/bcrypt"  // 用于密码哈希处理
)

// =====================登录相关===========================

// Claims 自定义JWT声明结构体
// 包含用户ID和标准声明
type Claims struct {
	UserID             uint // 用户唯一标识符
	jwt.StandardClaims      // JWT标准字段，如过期时间、签发者等
}

// GetCaptchaHandler 生成并返回验证码
// 用于登录验证的图形验证码处理
func GetCaptchaHandler(ctx *gin.Context) {
	captcha := toolCaptcha.GetCaptcha()
	if (captcha != toolCaptcha.Captcha{}) {
		response.Data(ctx, "验证码生成成功！", captcha)
	} else {
		response.Error(ctx, simple.NewSimpleErrorMessage("验证码生成失败！"))
	}
}

// LoginLockTime 登录错误锁定时间，单位为秒，默认为10分钟
const LoginLockTime = 60 * 10

// loginErrorCount 记录并返回登录错误次数
// 当错误次数过多时会临时锁定账号
// 参数：
//   - loginErrorCountKey: 缓存中的错误计数键名
//
// 返回：
//   - 当前错误次数
func loginErrorCount(loginErrorCountKey string) int64 {
	// 增加错误计数并设置过期时间
	errorCount, _ := cache.GetCacheSystem().IncrExpire(loginErrorCountKey, LoginLockTime)
	return errorCount
}

// getLoginErrorCount 获取当前登录错误次数
// 参数：
//   - loginErrorCountKey: 缓存中的错误计数键名
//
// 返回：
//   - 当前错误次数，如果不存在则返回0
func getLoginErrorCount(loginErrorCountKey string) int {
	// 从缓存中获取错误计数
	errorCount, _ := cache.GetCacheSystem().GetInt(loginErrorCountKey)
	return errorCount
}

// LoginUser 定义登录请求参数结构体
// LoginType 可选值: password,captcha,mfa,email,phone
type LoginUser struct {
	Username           string `json:"username"`               // 用户名
	Password           string `json:"password"`               // 密码
	Dots               string `json:"dots"`                   // 验证码点击位置
	CaptchaKey         string `json:"captcha_key"`            // 验证码键值
	MfaVerifyCode      string `json:"mfa_verify_code"`        // MFA验证码
	LoginType          string `json:"login_type"`             // 登录类型
	LoginErrorCountKey string `json:"loginErrorCountKey____"` // 错误计数键(内部使用)
}

// getLoginUser 根据不同登录类型验证用户身份
// 支持密码、验证码、MFA等多种登录方式
func getLoginUser(loginUser LoginUser) (basic.User, error) {
	var user basic.User = basic.User{}

	// 普通密码登录方式
	if loginUser.LoginType == "password" || loginUser.LoginType == "" {
		// 检查是否被锁定
		errorCount := getLoginErrorCount(loginUser.LoginErrorCountKey)
		if errorCount > 5 {
			return user, errors.New("登录错误次数过多，账号暂时锁定！")
		}

		// 验证用户账号
		user, err := user.GetByAccount(loginUser.Username)
		if err != nil {
			logger.LOG.Errorf("登录失败：%s", err.Error())
			return user, errors.New("账号密码有误或账号冻结！")
		}

		// 验证密码
		if VerifyPass(user.Password, loginUser.Password) {
			return user, nil
		} else {
			return user, errors.New("账号密码有误或账号冻结！")
		}
	} else if loginUser.LoginType == "captcha" {
		// 验证码登录方式
		// 先验证图形验证码
		_, errCapt := toolCaptcha.Check(loginUser.CaptchaKey, loginUser.Dots)
		if errCapt != nil {
			return user, errors.New("图形验证未通过！")
		}

		// 再验证用户名和密码
		user, err := user.GetByAccount(loginUser.Username)
		if err != nil {
			logger.LOG.Errorf("登录失败：%s", err.Error())
			return user, errors.New("账号密码有误或账号冻结！")
		}

		if VerifyPass(user.Password, loginUser.Password) {
			return user, nil
		} else {
			return user, errors.New("账号密码有误或账号冻结！")
		}
	} else if loginUser.LoginType == "mfa" {
		// MFA双因素验证登录方式
		user, err := user.GetByAccount(loginUser.Username)
		if err != nil {
			logger.LOG.Errorf("登录失败：%s", err.Error())
			return user, errors.New("账号密码有误或账号冻结！")
		}

		// 先验证密码
		if VerifyPass(user.Password, loginUser.Password) {
			// 再验证MFA码
			if mfa.ValidCode(loginUser.MfaVerifyCode, user.MFACode, 30) {
				return user, nil
			} else {
				return user, errors.New("MFA验证码错误！")
			}
		} else {
			return user, errors.New("账号密码有误或账号冻结！")
		}

	} else if loginUser.LoginType == "email" {
		// 邮箱登录方式（尚未实现）
		return user, errors.New("登录类型错误：" + loginUser.LoginType)
	} else {
		// 其他未知登录方式
		return user, errors.New("登录类型错误：" + loginUser.LoginType)
	}
}

// LoginHandler 用户登录处理函数
// 处理用户登录请求，验证身份并生成令牌
func LoginHandler(ctx *gin.Context) {
	var loginUser LoginUser
	// 解析请求JSON
	if err := ctx.BindJSON(&loginUser); err != nil {
		response.BadRequest(ctx, "格式错误")
		return
	}

	// 获取用户IP，用于锁定特定IP+用户名组合
	ip := GetUserIP(ctx)
	loginUser.LoginErrorCountKey = "login_error_count:" + loginUser.Username + ":" + ip

	// 验证用户身份
	user, err := getLoginUser(loginUser)
	if err == nil {
		// 判断是否需要MFA验证
		if user.MFAEnable == 1 && user.MFACode != "" && loginUser.LoginType != "mfa" {
			response.Message(ctx, 4401, "需要MFA验证！")
			return
		} else {
			// 生成JWT令牌
			config := config.CONF
			expire := config.Jwt.Expire
			expireTime := time.Now().Add(time.Duration(expire) * time.Second)
			claims := &Claims{
				UserID: user.ID,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expireTime.Unix(),  // 过期时间
					IssuedAt:  time.Now().Unix(),  // 签发时间
					Issuer:    config.Jwt.Issuer,  // 签名颁发者
					Subject:   config.Jwt.Subject, // 签名主题
				},
			}

			// 使用HS256算法创建令牌
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenStr, errToken := token.SignedString([]byte(config.Jwt.AppKey))

			// 用户权限（目前为通配符，表示拥有所有权限）
			var permissions = []string{"*"}

			if errToken != nil {
				// 令牌生成失败，增加登录错误计数
				errorCount := loginErrorCount(loginUser.LoginErrorCountKey)
				response.Error(ctx, simple.NewSimpleErrorData("获取用户信息失败！", errorCount))
				return
			}

			// 将令牌存入缓存并清除错误计数
			cache.GetCacheSystem().Set(fmt.Sprintf("%d", user.ID), tokenStr)
			cache.GetCacheSystem().Delete(loginUser.LoginErrorCountKey)

			// 返回登录成功及用户信息
			response.Data(ctx, "登录成功！",
				data.Map{"token": tokenStr, "id": user.ID, "username": user.Account, "avatar": user.Avatar, "name": user.Name, "perms": permissions, "ip": ip})
			return
		}
	} else {
		// 登录失败，增加错误计数
		errorCount := loginErrorCount(loginUser.LoginErrorCountKey)
		response.Error(ctx, simple.NewSimpleError(4000, err.Error(), errorCount))
		return
	}
}

// PasswordLoginOut PC后台退出处理
// 清除用户缓存中的令牌信息
func PasswordLoginOut(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	// 检查令牌是否存在
	if tokenString == "" {
		response.Unauthorized(ctx, "未提供授权令牌")
		return
	}

	// 解析令牌获取用户信息
	claims, err := ParseToken(tokenString)
	if err != nil {
		logger.LOG.Warnf("令牌解析失败: %v", err)
		response.Unauthorized(ctx, "无效的令牌")
		return
	}

	// 从缓存中删除用户令牌
	userID := fmt.Sprintf("%d", claims.UserID)
	ok, err := cache.GetCacheSystem().Delete(userID)

	// 首先处理删除操作中的错误
	if err != nil {
		logger.LOG.Errorf("从缓存中删除用户令牌失败: %v", err)
		response.Error(ctx, simple.NewSimpleErrorMessage("退出处理失败"))
		return
	}

	// 处理令牌未找到的情况
	if !ok {
		// 这种情况可能是令牌已经过期或者已被删除
		// 仍然返回成功，因为用户的目标是退出登录
		logger.LOG.Infof("用户 %s 的令牌已不存在于缓存中", userID)
		response.Success(ctx, "已退出登录")
		return
	}

	// 令牌成功删除
	logger.LOG.Infof("用户 %s 成功退出登录", userID)
	response.Success(ctx, "安全退出成功")
}

// VerifyPass 验证密码是否正确
// 使用bcrypt比较哈希密码和用户输入密码
// 参数：
//   - password: 数据库中存储的哈希密码
//   - loginPass: 用户输入的密码（可能是加密的）
//
// 返回：
//   - 如果密码匹配返回true，否则返回false
func VerifyPass(password string, loginPass string) bool {
	// 检查密码是否为空
	if password == "" || loginPass == "" {
		return false
	}
	// 调试日志
	logger.LOG.Debugf("password:%s,loginPass:%s", password, loginPass)
	// 解密用户输入的密码并与数据库中的哈希值比较
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(xxtea.DecryptAuto(loginPass, "")))
	return err == nil
}

// =====================登录相关===========================

// =====================Token获取用户信息===========================

// GetUserInfo 根据令牌获取用户信息
// 从令牌中提取用户ID并查询用户详细信息
func GetUserInfo(ctx *gin.Context) {
	id := request.GetUserID(ctx)
	if id > 0 {
		user := basic.User{}
		user, err := user.Load(id)
		var permissions = []string{"*"}
		// 获取用户权限信息
		if err == nil {
			response.Data(ctx, "获取用户信息成功！", data.Map{"user": user, "perms": permissions})
		} else {
			response.Error(ctx, simple.NewSimpleErrorMessage("获取用户信息失败！"))
		}
	} else {
		response.Unauthorized(ctx, "Authorization Failed！")
		return
	}
}

// AuthApiKeyMiddleware API密钥认证中间件
// 支持通过API密钥+MFA进行身份验证，也支持回退到标准JWT认证
func AuthApiKeyMiddleware(ctx *gin.Context) {
	apiKey := ctx.GetHeader("Api-Key")
	apiAllow := false

	// 检查API密钥是否有效
	if apiKey != "" && apiKey == config.CONF.App.ApiKey && config.CONF.App.ApiUser > 0 {
		// 通过API密钥获取用户
		user := basic.User{}
		user, _ = user.Load(uint(config.CONF.App.ApiUser))

		// 如果用户启用了MFA，需要验证MFA码
		if user.MFAEnable == 1 {
			mfaKey := ctx.GetHeader("Mfa-Key")
			if mfa.ValidCode(mfaKey, user.MFACode, 30) {
				// MFA验证通过
				apiAllow = true
			}
		} else {
			// 未启用MFA，直接通过API密钥验证
			apiAllow = true
		}
	}

	if apiAllow {
		// 设置用户ID
		ctx.Set("UserID", config.CONF.App.ApiUser)
	} else {
		// API密钥验证失败，回退到标准JWT认证
		AuthMiddleware(ctx)
	}
}

// =====================Token获取用户信息===========================

// =====================授权相关===========================

// AuthMiddleware 标准JWT身份认证中间件
// 验证请求中的Authorization头中的JWT令牌
func AuthMiddleware(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		response.Unauthorized(ctx, "Authorization Failed！")
		return
	}

	// 解析令牌
	claims, err := ParseToken(tokenString)
	if err != nil {
		response.Unauthorized(ctx, err.Error())
		return
	} else {
		// 验证令牌是否在缓存中（是否有效）
		token, errToken := cache.GetCacheSystem().Get(fmt.Sprintf("%d", claims.UserID))
		if errToken != nil && len(token) == 0 {
			response.Unauthorized(ctx, errToken.Error())
			return
		}
		// 将用户ID设置到上下文中
		ctx.Set("UserID", claims.UserID)
	}
}

// ParseToken 解析JWT令牌
// 支持Bearer格式的令牌
func ParseToken(tokenString string) (*Claims, error) {
	// 处理Bearer前缀
	if strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
		tokenString = strings.Trim(tokenString[7:], " ")
	}

	// 使用应用密钥解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{},
		func(token *jwt.Token) (i interface{}, err error) {
			return []byte(config.CONF.Jwt.AppKey), nil
		})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

// =====================授权相关===========================

// GetUserIP 获取用户真实IP地址
// 处理代理服务器、负载均衡器等场景下的IP获取
func GetUserIP(c *gin.Context) string {
	ip := c.ClientIP()
	ipx := net.ParseIP(ip)

	// 如果是内部IP，尝试从代理头部获取真实IP
	if ipx == nil || isPrivateIP(ipx) {
		proxyHeaders := c.Request.Header.Get("X-Real-IP,X-Forwarded-For")
		// 分割代理头部，取第一个有效的公网IP作为真实IP
		ips := strings.Split(proxyHeaders, ",")
		for _, v := range ips {
			if isValidIP(v) && !isPrivateIP(net.ParseIP(v)) {
				return v
			}
		}
	}
	return ip
}

// isPrivateIP 检查IP是否为内部私有IP
// 检查常见的私有IP范围
// 参数：
//   - ip: 要检查的IP地址
//
// 返回：
//   - 如果是私有IP返回true，否则返回false
func isPrivateIP(ip net.IP) bool {
	// 如果IP为空，返回false
	if ip == nil {
		return false
	}
	// 检查是否为本地回环 IP（127.0.0.1等）
	if ip.IsLoopback() {
		return true
	}
	// 转换为IPv4格式
	ip4 := ip.To4()
	if ip4 == nil {
		// 如果不是IPv4地址，返回false
		return false
	}
	// 检查常见的私有IP范围
	privateRanges := []string{
		"10.0.0.0/8",     // 私有网络 A 类
		"172.16.0.0/12",  // 私有网络 B 类
		"192.168.0.0/16", // 私有网络 C 类
	}
	// 遍历所有私有网段范围进行检查
	for _, r := range privateRanges {
		_, network, _ := net.ParseCIDR(r)
		if network.Contains(ip4) {
			return true
		}
	}
	return false
}

// isValidIP 检查字符串是否为有效的IP地址
// 参数：
//   - ip: 要检查的IP地址字符串
//
// 返回：
//   - 如果是有效IP返回true，否则返回false
func isValidIP(ip string) bool {
	// 尝试解析IP地址
	ipx := net.ParseIP(ip)
	return ipx != nil
}
