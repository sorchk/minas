// Package webdav 提供WebDAV协议支持
package webdav

// User 表示WebDAV用户信息
// 包含用户账号、主目录、权限和访问状态
type User struct {
	// Account 用户账号
	Account string
	// Directory 用户的WebDAV根目录
	Directory string
	// Perms 用户权限列表，每个元素代表一种权限
	// 包含: C(上传), U(修改), D(删除), L(列表), G(下载)
	Perms []string
	// Allow 用户是否允许访问，true表示启用，false表示禁用
	Allow bool
}
