// Package nas 提供网络存储相关的服务和模型
package nas

import (
	"log"
	"server/core/db"
	"server/utils/config"
	"server/utils/global"
	"server/utils/xxtea" // 加密解密工具

	"github.com/segmentio/ksuid" // 生成唯一ID
	"gorm.io/gorm"
)

// Webdav 定义WebDAV账号模型
// 用于存储WebDAV账号信息，包括账号、令牌、权限等
type Webdav struct {
	db.BaseModel[Webdav] // 继承基础模型，提供通用字段和方法

	// 基本信息
	Name string `gorm:"comment:'名称'" json:"name"`     // WebDAV账号名称
	Home string `gorm:"comment:'主目录'" json:"home"`    // WebDAV根目录路径
	Size uint   `gorm:"comment:'配额(MB)'" json:"size"` // 存储配额，单位MB

	// 认证信息
	Account string `gorm:"comment:'访问账号'" json:"account"` // 登录账号
	Token   string `gorm:"comment:'访问令牌'" json:"token"`   // 访问令牌，用于认证

	// 其他信息
	Remark string `gorm:"comment:'备注'" json:"remark"` // 备注信息
	Perms  string `gorm:"comment:'权限'" json:"perms"`  // 权限字符串，如"CLGUD"
}

// TableName 返回WebDAV表名
// 实现gorm的Tabler接口
func (Webdav) TableName() string {
	return "webdavs" // 数据库表名
}

// BeforeCreate 创建记录前的钩子函数
// 在创建新WebDAV账号前自动生成令牌
func (u *Webdav) BeforeCreate(tx *gorm.DB) (err error) {
	// 调用父类的创建前函数
	u.SupperBeforeCreate()
	// 生成唯一的令牌
	u.Token = ksuid.New().String()
	return
}

// AfterFind 查询记录后的钩子函数
// 在从数据库读取数据后自动解密令牌
func (u *Webdav) AfterFind(tx *gorm.DB) (err error) {
	// 调用父类的查询后函数
	u.SupperAfterFind()
	// 解密令牌
	u.Token = xxtea.DecryptAuto(u.Token, config.CONF.Db.DataKey)
	return
}

// BeforeSave 保存记录前的钩子函数
// 在保存数据前自动加密令牌
func (u *Webdav) BeforeSave(tx *gorm.DB) (err error) {
	log.Println("Webdav BeforeSave")
	// 加密令牌
	u.Token = xxtea.EncryptAuto(u.Token, config.CONF.Db.DataKey)
	return
}

// GetByToken 根据账号和令牌获取WebDAV账号信息
// 参数:
//   - 无需额外参数，使用对象自身的Account和Token属性
//
// 返回:
//   - Webdav: 查询到的WebDAV账号信息
//   - error: 错误信息，如果没有错误则为nil
func (e *Webdav) GetByToken() (Webdav, error) {
	// 初始化结果对象
	webdav := Webdav{}

	// 根据账号和令牌查询数据库
	if err := global.DB.Where("account = ?", e.Account).Where("token = ?", e.Token).First(&webdav).Error; err != nil {
		// 查询失败，返回错误
		return webdav, err
	}

	// 查询成功，返回结果
	return webdav, nil
}
