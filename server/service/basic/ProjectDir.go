// Package nas 提供网络存储相关的服务和模型
package basic

import (
	"server/core/db"
	"server/utils/global"

	"gorm.io/gorm"
)

// ProjectDir 定义项目目录实体模型
// 用于存储项目的目录结构信息
type ProjectDir struct {
	db.BaseModel[ProjectDir] // 继承基础模型，提供通用字段和方法

	// 基本信息
	Name     string `gorm:"comment:'名称';not null" json:"name"`         // 目录名称
	ParentID uint   `gorm:"comment:'上级ID';default:0" json:"parent_id"` // 上级目录ID，0表示顶级目录
	Remark   string `gorm:"comment:'备注'" json:"remark"`                // 备注说明

	// 非数据库字段
	Children []ProjectDir `gorm:"-" json:"children,omitempty"` // 子目录列表，不存储到数据库
	HasChild bool         `gorm:"-" json:"has_child"`          // 是否有子目录
}

// TableName 返回项目目录表名
// 实现gorm的Tabler接口
func (ProjectDir) TableName() string {
	return "project_dirs"
}

// GetTree 获取指定模块类型的目录树
// 返回:
//   - []ProjectDir: 目录树结构
//   - error: 错误信息，如果没有错误则为nil
func (e *ProjectDir) GetTree() ([]ProjectDir, error) {
	// 查询所有目录
	var allDirs []ProjectDir
	query := global.DB.Model(&ProjectDir{})

	// 按ID排序获取所有目录
	if err := query.Order("id").Find(&allDirs).Error; err != nil {
		return nil, err
	}

	// 构建目录树
	return buildTree(allDirs, 0), nil
}

// GetByParent 获取指定父目录下的所有子目录
// 参数:
//   - parentID: 父目录ID，0表示顶级目录
//
// 返回:
//   - []ProjectDir: 子目录列表
//   - error: 错误信息，如果没有错误则为nil
func (e *ProjectDir) GetByParent(parentID uint) ([]ProjectDir, error) {
	var dirs []ProjectDir

	// 根据父目录ID查询子目录
	if err := global.DB.Where("parent_id = ?", parentID).Order("id").Find(&dirs).Error; err != nil {
		return nil, err
	}

	return dirs, nil
}

// HasChildren 检查是否有子目录
// 参数:
//   - id: 目录ID
//
// 返回:
//   - bool: 是否有子目录
//   - error: 错误信息，如果没有错误则为nil
func (e *ProjectDir) HasChildren(id uint) (bool, error) {
	var count int64

	// 计算子目录数量
	if err := global.DB.Model(&ProjectDir{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// AfterFind 查询后的钩子函数
// 在从数据库查询记录后自动检查是否有子目录
func (e *ProjectDir) AfterFind(tx *gorm.DB) (err error) {
	// 调用父类的查询后函数
	e.SupperAfterFind()

	// 检查是否有子目录
	hasChild, _ := e.HasChildren(e.ID)
	e.HasChild = hasChild

	return
}

// buildTree 辅助函数，用于递归构建目录树
// 参数:
//   - allDirs: 所有目录列表
//   - parentID: 父目录ID
//
// 返回:
//   - []ProjectDir: 构建好的目录树
func buildTree(allDirs []ProjectDir, parentID uint) []ProjectDir {
	var result []ProjectDir

	// 遍历所有目录，找出当前层级的目录
	for _, dir := range allDirs {
		if dir.ParentID == parentID {
			// 递归构建子目录
			dir.Children = buildTree(allDirs, dir.ID)
			// 添加到结果中
			result = append(result, dir)
		}
	}

	return result
}
