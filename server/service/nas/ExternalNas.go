package nas

import (
	"log"
	"server/core/app/request"
	"server/core/db"
	"server/utils/config"
	"server/utils/global"
	"server/utils/rclone"
	"server/utils/xxtea"

	"github.com/duke-git/lancet/v2/maputil"
	"gorm.io/gorm"
)

// ExternalNas 外部NAS存储模型
// 用于存储和管理外部存储服务的配置信息
// 支持通过rclone连接各类存储服务如S3、FTP、WebDAV等
type ExternalNas struct {
	db.BaseModel[ExternalNas]        // 嵌入基础模型，提供ID、创建时间等公共字段
	Name                      string `gorm:"comment:'名称'" json:"name"`         // 存储服务名称，用于显示
	RcName                    string `gorm:"comment:'标识'" json:"rc_name"`      // rclone配置中的唯一标识符
	Type                      string `gorm:"comment:'类型'" json:"type"`         // 存储服务类型，如s3、ftp、webdav等
	IsAdv                     uint   `gorm:"comment:'是否开启高级配置'" json:"is_adv"` // 是否启用高级配置：0=否，1=是
	Config                    string `gorm:"comment:'配置'" json:"config"`       // 存储服务的配置信息，加密存储
	Remark                    string `gorm:"comment:'备注'" json:"remark"`       // 备注说明
	IsSync                    bool   `sql:"-" gorm:"-" json:"is_sync"`         // 是否与rclone配置同步，非数据库字段
}

// TableName 指定数据库表名
// 返回:
//   - string: 表名 "external_nas"
func (extNas ExternalNas) TableName() string {
	return "external_nas"
}

// AfterFind GORM钩子，在从数据库加载数据后执行
// 对加密存储的配置信息进行解密
// 参数:
//   - tx: GORM事务对象
//
// 返回:
//   - error: 可能的错误信息
func (u *ExternalNas) AfterFind(tx *gorm.DB) (err error) {
	u.SupperAfterFind()
	u.Config = xxtea.DecryptAuto(u.Config, config.CONF.Db.DataKey)
	return
}

// BeforeSave GORM钩子，在保存到数据库前执行
// 对配置信息进行加密存储
// 参数:
//   - tx: GORM事务对象
//
// 返回:
//   - error: 可能的错误信息
func (u *ExternalNas) BeforeSave(tx *gorm.DB) (err error) {
	u.Config = xxtea.EncryptAuto(u.Config, config.CONF.Db.DataKey)
	return
}

// 查询
func (extNas ExternalNas) List(query request.PageQuery) (list []ExternalNas, count int64, err error) {
	var entity ExternalNas
	list, err = entity.ListQuery(query)
	if err == nil {
		cMap, err := rclone.DumpConfig()
		if err == nil {
			for i, item := range list {
				list[i].IsSync = maputil.HasKey(cMap, item.RcName)
			}
		}
		count, err = entity.CountQuery(query)
		return list, count, err
	} else {
		return list, 0, err
	}
}

// 保存
func (extNas ExternalNas) Save(entity *ExternalNas, columns ...string) error {
	id := (*entity).GetID()
	if id > 0 {
		var tmp ExternalNas
		err := global.DB.Model(tmp).Take(&tmp, id).Error
		if tmp.GetID() == 0 || (err != nil && err == gorm.ErrRecordNotFound) {
			//记录不存在 新增
			return extNas.Create(entity)
		} else {
			return extNas.Update(entity, columns...)
		}
	} else {
		return extNas.Create(entity)
	}
}

// 创建
func (extNas ExternalNas) Create(entity *ExternalNas) error {
	log.Println("Create--------------------")
	err := rclone.ConfigCreate(entity.RcName, entity.Type, entity.Config)
	if err != nil {
		return err
	}
	return global.DB.Model(entity).Create(entity).Error
}

// 更新
func (extNas ExternalNas) Update(entity *ExternalNas, columns ...string) error {
	log.Println("Update--------------------")
	err := rclone.ConfigUpdate(entity.RcName, entity.Type, entity.Config)
	if err != nil {
		return err
	}
	if len(columns) > 0 {
		return global.DB.Model(entity).Select(columns).Updates(entity).Error
	} else {
		return global.DB.Model(entity).Updates(entity).Error
	}
}

// 按主键删除
func (extNas ExternalNas) Delete(id any) error {
	log.Println("Delete--------------------")
	var entity ExternalNas
	entity, err := entity.Load(id)
	if err != nil {
		return err
	}
	err = rclone.ConfigDelete(entity.RcName)
	if err != nil {
		return err
	}
	return global.DB.Model(entity).Delete(&entity, id).Error
}
