// Package sflow
package sflow

import (
	"log"
	"server/core/db"

	// 加密解密工具
	// 生成唯一ID
	"gorm.io/gorm"
)

// type: api,job,crawler,
// SFlow 定义SFlow 模型
// 用于存储SFlow
type SFlow struct {
	db.BaseModel[SFlow]        // 继承基础模型，提供通用字段和方法
	Name                string `gorm:"comment:'名称' size:128" json:"name"`                  // 名称
	Content             string `gorm:"comment:'脚本' size:102400 default:''" json:"content"` // 脚本（JSON格式）
	Type                string `gorm:"comment:'类型' size:20 default:''" json:"type"`        // 类型
	LogLevel            int    `gorm:"comment:'日志级别' default:0" json:"log_level"`          // 日志级别
	LastStatus          int    `gorm:"default:0;comment:'最近状态'" json:"last_status"`        // 最近执行状态
	LastRunTime         string `gorm:"comment:'上次执行时间'" json:"last_run_time"`              // 上次执行时间
	ProjectDirID        string `gorm:"comment:'项目目录ID';default:0" json:"project_dir_id"`   // 项目目录ID，关联到项目目录
	Remark              string `gorm:"comment:'备注'" json:"remark"`                         // 备注说明
}

// TableName 返回SFlow表名
// 实现gorm的Tabler接口
func (SFlow) TableName() string {
	return "sflow" // 数据库表名
}

// BeforeCreate 创建记录前的钩子函数
func (u *SFlow) BeforeCreate(tx *gorm.DB) (err error) {
	// 调用父类的创建前函数
	u.SupperBeforeCreate()
	return
}

// AfterFind 查询记录后的钩子函数
func (u *SFlow) AfterFind(tx *gorm.DB) (err error) {
	// 调用父类的查询后函数
	u.SupperAfterFind()
	return
}

// BeforeSave 保存记录前的钩子函数
func (u *SFlow) BeforeSave(tx *gorm.DB) (err error) {
	log.Println("SFlow BeforeSave")
	return
}
