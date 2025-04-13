package job

// SchJob 计划任务基础结构体
// 定义了所有计划任务共有的基础属性
// 被各种具体任务类型结构体嵌入，提供通用字段
type SchJob struct {
	TaskId   uint   `gorm:"comment:'任务编号'" json:"task_id"`  // 任务唯一标识ID
	TaskName string `gorm:"comment:'名称'" json:"task_name"`  // 任务名称，用于显示
	TaskType string `gorm:"comment:'类型' " json:"task_type"` // 任务类型，用于区分不同种类的任务
}
