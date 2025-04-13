// Package webdav 提供WebDAV协议支持，允许客户端通过HTTP协议进行文件管理操作
package webdav

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"server/service/nas"
	"server/utils/config"
	"strings"

	"golang.org/x/net/webdav"
)

// WebDAV访问路径常量
// DAV_ROOT 表示WebDAV的根路径
const DAV_ROOT = "/dav"

// DAV_PREFIX 表示WebDAV的路径前缀
const DAV_PREFIX = "/dav/"

// IsWebDav 判断请求URL是否是WebDAV请求
// 参数:
//   - url: 请求的URL路径
//
// 返回:
//   - bool: 如果是WebDAV请求返回true，否则返回false
func IsWebDav(url string) bool {
	return strings.HasPrefix(url, config.CONF.App.ContextPath+DAV_PREFIX) || url == config.CONF.App.ContextPath+DAV_ROOT
}

// GetUser 根据账号和令牌获取WebDAV用户信息
// 参数:
//   - account: 用户账号
//   - token: 访问令牌
//
// 返回:
//   - User: WebDAV用户信息，包含目录、权限等
func GetUser(account string, token string) User {
	// 创建WebDAV数据库对象
	var webdavDb = nas.Webdav{}
	// 初始化用户对象，默认不允许访问
	user := User{Directory: "", Account: account, Perms: []string{}, Allow: false}

	// 设置查询条件
	webdavDb.Account = account
	webdavDb.Token = token

	// 从数据库获取用户信息
	webdave, err := webdavDb.GetByToken()
	if err != nil {
		// 如果查询出错，返回默认用户对象（无权限）
		return user
	}

	// 设置用户主目录
	user.Directory = webdave.Home

	// 确保权限字符串不为空再进行分割
	if webdave.Perms != "" {
		// 将权限字符串分割为字符数组，每个字符代表一种权限
		user.Perms = strings.Split(webdave.Perms, "")
	}

	// 设置用户是否允许访问（IsDisable为0表示启用）
	user.Allow = webdave.IsDisable == 0
	return user
}

// acls 定义HTTP方法与权限码的映射关系
// 权限码说明:
// - C: 上传权限
// - U: 修改权限(复制、移动、重命名)
// - D: 删除权限
// - L: 列表权限
// - G: 下载权限
var acls = map[string]string{
	// 上传相关的HTTP方法
	"PUT":   "C", // 上传文件
	"MKCOL": "C", // 创建目录

	// 修改相关的HTTP方法
	"COPY":      "U", // 复制文件/目录
	"MOVE":      "U", // 移动或重命名文件/目录
	"PROPPATCH": "U", // 修改属性

	// 列表相关的HTTP方法
	"PROPFIND": "L", // 获取目录列表或文件信息

	// 删除相关的HTTP方法
	"DELETE": "D", // 删除文件/目录

	// 下载相关的HTTP方法
	"GET":  "G", // 获取文件内容
	"HEAD": "G", // 获取文件头信息
	"POST": "G", // 某些客户端使用POST下载文件
}

// matchAcl 检查用户是否有权限执行指定HTTP方法
// 参数:
//   - userPerms: 用户权限列表，包含权限码字符
//   - method: HTTP请求方法
//
// 返回:
//   - bool: 如果用户有权限返回true，否则返回false
func matchAcl(userPerms []string, method string) bool {
	// 获取请求方法对应的权限码
	perm := acls[strings.ToUpper(method)]

	// 如果方法没有对应的权限码，表示不需要特殊权限
	if perm == "" {
		// 不管控的HTTP方法，允许访问
		return true
	}

	// 遍历用户权限列表，检查是否包含所需权限
	for _, item := range userPerms {
		if item == perm {
			return true
		}
	}

	// 没有找到匹配的权限，拒绝访问
	return false
}

// ServeHTTP 处理WebDAV请求，包括认证、授权和请求处理
// 参数:
//   - w: HTTP响应写入器
//   - r: HTTP请求
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 设置基本认证头
	w.Header().Set("WWW-Authenticate", `Basic realm="Webdav realm"`)

	// 获取基本认证信息
	username, token, ok := r.BasicAuth()
	if !ok {
		// 没有提供认证信息，返回未授权错误
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	// 根据账号和令牌获取用户信息
	user := GetUser(username, token)

	// 检查用户是否被禁用
	if !user.Allow {
		http.Error(w, "User disabled", http.StatusUnauthorized)
		return
	}

	// 检查用户是否有权限执行请求的方法
	if !matchAcl(user.Perms, r.Method) {
		http.Error(w, "No permission to access this request method:"+r.Method, http.StatusForbidden)
		return
	}

	// 创建WebDAV处理器
	var handler = webdav.Handler{
		Prefix:     config.CONF.App.ContextPath + DAV_ROOT,
		FileSystem: webdav.Dir(user.Directory),
		LockSystem: webdav.NewMemLS(), // 使用内存锁系统
	}

	// 根据请求方法自动创建父目录，提高用户体验
	switch r.Method {
	case "MKCOL": // 创建目录时
		MkdirAllParentDir(r.URL.Path, user.Directory)
	case "PUT": // 上传文件时
		MkdirAllParentDir(r.URL.Path, user.Directory)
	case "COPY", "MOVE": // 复制或移动时
		// 获取目标路径
		hdr := r.Header.Get("Destination")
		if hdr != "" {
			u, err := url.Parse(hdr)
			if err == nil {
				// 为目标路径创建父目录
				MkdirAllParentDir(u.Path, user.Directory)
			}
		}
	}

	// 执行WebDAV请求处理
	handler.ServeHTTP(w, r)

	// 如果是OPTIONS请求，返回200 OK状态
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	}
}

// MkdirAllParentDir 创建指定路径的父目录
// 参数:
//   - urlPath: WebDAV请求的URL路径
//   - userDirectory: 用户的根目录
func MkdirAllParentDir(urlPath string, userDirectory string) {
	// 从请求路径中移除上下文路径和WebDAV前缀
	path := strings.TrimPrefix(urlPath, config.CONF.App.ContextPath+DAV_PREFIX)

	// 将相对路径转换为用户目录下的绝对路径
	path = filepath.Join(userDirectory, path)

	// 使用更安全的权限设置创建父目录
	// 0755权限表示所有者有读写执行权限，其他用户有读和执行权限
	os.MkdirAll(filepath.Dir(path), 0755)
}
