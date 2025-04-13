package basic

import (
	"server/core/app/request"
	"server/core/db"
	"server/utils/config"
	"server/utils/global"
	"server/utils/xxtea"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// User 用户模型，存储系统用户信息
// 继承自基础模型，并包含用户特有的字段
type User struct {
	db.BaseModel[User]        // 嵌入基础模型，提供ID、创建时间等公共字段
	Account            string `gorm:"comment:'账号';type:varchar(50);unique;not null;" json:"account"`                                         // 用户账号，唯一且必填
	Password           string `gorm:"comment:'密码'" json:"password"`                                                                          // 用户密码
	Phone              string `gorm:"comment:'手机号码'" json:"phone"`                                                                           // 用户手机号码
	Name               string `gorm:"comment:'姓名';type:varchar(50);column:name" json:"name" structs:"name" validate:"required,min=1,max=50"` // 用户姓名，必填字段
	Avatar             string `gorm:"force;comment:'头像URL'" json:"avatar"`                                                                   // 用户头像URL
	Birthday           string `gorm:"comment:'出生日期'" json:"birthday"`                                                                        // 用户生日
	Sex                int    `gorm:"comment:'性别：0 未填写 1男 2女';default:0" json:"sex"`                                                         // 用户性别：0=未填写，1=男，2=女
	Email              string `gorm:"comment:'邮箱';type:varchar(100);unique_index" json:"email"`                                              // 用户邮箱，唯一索引
	Sn                 int    `gorm:"comment:'排序';default:0" json:"sn"`                                                                      // 排序号
	IsAdmin            int    `gorm:"comment:'是超级管理';default:0" json:"is_admin"`                                                             // 是否为超级管理员：0=否，1=是
	MFAEnable          int    `gorm:"comment:'MFA认证是否启用';default:0" json:"mfa_enable"`                                                       // 多因素认证是否启用：0=否，1=是
	MFACode            string `gorm:"comment:'MFA认证密钥'" json:"mfa_code"`                                                                     // MFA认证密钥
	Remark             string `gorm:"comment:'备注';default:''" json:"remark"`                                                                 // 用户备注信息
	LastLoginIp        string `gorm:"comment:'最后登录IP'" json:"last_login_ip"`                                                                 // 用户最后登录IP地址
	LastLoginTime      int    `gorm:"comment:'最后登录时间';default:0" json:"last_login_time"`                                                     // 用户最后登录时间戳
}

// SetOperator 设置操作者信息
// 在用户模型更新时，记录更新操作的执行者ID
// 参数:
//   - entity: 用户实体指针
//   - ctx: Gin上下文，用于获取当前请求的用户ID
func (*User) SetOperator(entity *User, ctx *gin.Context) {
	entity.UpdatedBy = request.GetUserID(ctx)
}

// TableName 指定数据库表名
// 返回:
//   - string: 表名 "users"
func (*User) TableName() string {
	return "users" // 这里返回正确的表名
}

// AfterFind GORM钩子，在从数据库加载数据后执行
// 处理需要解密的字段
// 参数:
//   - tx: GORM事务对象
//
// 返回:
//   - error: 可能的错误信息
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	u.SupperAfterFind()
	u.MFACode = xxtea.DecryptAuto(u.MFACode, config.CONF.Db.DataKey)
	return
}

// BeforeSave GORM钩子，在保存到数据库前执行
// 处理需要加密的字段
// 参数:
//   - tx: GORM事务对象
//
// 返回:
//   - error: 可能的错误信息
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	u.MFACode = xxtea.DecryptAuto(u.MFACode, config.CONF.Db.DataKey)
	return
}

// ModifyPass 修改用户密码
// 参数:
//   - id: 用户ID
//   - password: 新密码
//
// 返回:
//   - error: 更新过程中可能发生的错误
func (e *User) ModifyPass(id uint, password string) error {
	user := User{}
	user.ID = id
	user.Password = password
	err := user.Update(&user, "password")
	if err != nil {
		return err
	}
	return nil
}

// GetByAccount 根据账号获取用户信息
// 参数:
//   - account: 用户账号
//
// 返回:
//   - User: 查询到的用户对象
//   - error: 查询过程中可能发生的错误
func (e *User) GetByAccount(account string) (User, error) {
	user := User{}
	err := global.DB.Model(user).Where("account =?", account).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
