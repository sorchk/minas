// package db 提供数据库模型和操作的核心功能
package db

import (
	"server/core/app/request" // 导入请求处理模块

	"gorm.io/gorm" // 导入GORM ORM库
)

// 定义常量
const (
	LOCAL_TIME_FORMAT = "2006-01-02 15:04:05" // 本地时间格式
	IS_DISABLE_FIELD  = "is_disable"          // 禁用字段名
	ID_FIELD          = "id"                  // ID字段名
)

// BaseDao 定义了基础数据访问对象接口
// 泛型接口，适用于任何类型
type BaseDao[T any] interface {
	GetID() uint                                       // 获取ID
	List(query request.PageQuery) ([]T, int64, error)  // 分页列表查询
	Load(id any) (T, error)                            // 根据ID加载单个实体
	Save(data *T, columns ...string) error             // 保存实体（创建或更新）
	Create(data *T) error                              // 创建新实体
	Update(data *T, columns ...string) error           // 更新现有实体
	Delete(id any) error                               // 删除实体
	Enable(id any) error                               // 启用实体
	Disable(id any) error                              // 禁用实体
	Count() (int64, error)                             // 计算总数
	CountQuery(query request.PageQuery) (int64, error) // 根据查询条件计算数量
	ListQuery(query request.PageQuery) ([]T, error)    // 根据查询条件获取列表
}

// GetID 获取模型ID的方法实现
func (m BaseModel[T]) GetID() uint {
	return m.ID
}

// BaseModel 基础模型结构体
// 包含了所有模型共有的字段和方法
type BaseModel[T BaseDao[T]] struct {
	ID            uint           `gorm:"primary_key" json:"id" mapstructure:"id"`    // 主键ID
	IsDisable     uint           `gorm:"default:0;comment:'是否禁用'" json:"is_disable"` // 是否禁用：0-启用，1-禁用
	CreatedAt     LocalTime      `gorm:"comment:'创建时间'" json:"created_at"`           // 创建时间
	CreatedBy     uint           `gorm:"default:0;comment:'创建人'" json:"created_by"`  // 创建人ID
	UpdatedAt     LocalTime      `gorm:"comment:'更新时间'" json:"updated_at"`           // 更新时间
	UpdatedBy     uint           `gorm:"default:0;comment:'更新人'" json:"updated_by"`  // 更新人ID
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`                    // 删除时间（软删除）
	CreatedByName string         `sql:"-" gorm:"-" json:"created_by_name"`           // 创建人姓名（非数据库字段）
	UpdatedByName string         `sql:"-" gorm:"-" json:"updated_by_name"`           // 更新人姓名（非数据库字段）
}

// SupperAfterFind 在查询后填充创建人和更新人的姓名
func (baseModel *BaseModel[T]) SupperAfterFind() {
	if baseModel.CreatedBy > 0 {
		baseModel.CreatedByName = GetUserNameCache(baseModel.CreatedBy)
	}
	if baseModel.UpdatedBy > 0 {
		baseModel.UpdatedByName = GetUserNameCache(baseModel.UpdatedBy)
	}
}

// AfterFind GORM的钩子函数，在记录查询后自动调用
func (baseModel *BaseModel[T]) AfterFind(tx *gorm.DB) (err error) {
	baseModel.SupperAfterFind()
	return
}

// SupperBeforeCreate 在创建记录前设置创建人
func (baseModel *BaseModel[T]) SupperBeforeCreate() {
	baseModel.CreatedBy = baseModel.UpdatedBy
}

// BeforeCreate GORM的钩子函数，在记录创建前自动调用
func (baseModel *BaseModel[T]) BeforeCreate(tx *gorm.DB) (err error) {
	baseModel.SupperBeforeCreate()
	return
}
