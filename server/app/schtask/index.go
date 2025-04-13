// package schtask 定义了计划任务的API接口和路由处理
package schtask

import (
	"log"
	"server/core/app/request"
	"server/core/app/response"
	"server/core/app/webapi"
	"server/service/scheduled"

	"github.com/gin-gonic/gin"
)

// SchTaskApp 计划任务应用结构体
// 继承了BaseApp基类，用于处理计划任务的API请求
type SchTaskApp struct {
	webapi.BaseApp[scheduled.SchTask] // 内嵌基类，指定模型类型为scheduled.SchTask
}

// AddRoutes 注册计划任务相关的路由
// 在父路由组下创建计划任务路由组，并注册相关的路由处理函数
func AddRoutes(parentGroup *gin.RouterGroup) {
	// 创建计划任务路由组
	group := parentGroup.Group("/schtask")
	// 创建应用实例
	app := SchTaskApp{}
	// 注册基础路由（CRUD操作）
	webapi.AddBaseRoutes(group, &app)
	// 设置允许更新的字段列表
	app.UpdateFields = []string{"name", "type", "cron", "source", "log_keep_num", "script", "remark"}

	// 注册手动执行任务的路由
	group.POST("/exec/:id", app.Exec)
}

// Exec 手动执行指定的计划任务
// 根据ID加载任务并执行，返回执行结果
func (app SchTaskApp) Exec(ctx *gin.Context) {
	log.Println("----------Exec------")
	// 创建任务实例
	schTask := scheduled.SchTask{}
	// 从路径参数中获取任务ID
	id := ctx.Param("id")
	// 根据ID加载任务信息
	schTask, err := schTask.Load(id)
	if err != nil {
		// 加载失败时返回错误
		response.Error(ctx, err)
		return
	}
	// 执行任务
	err = schTask.Exec()
	if err != nil {
		// 执行失败时返回错误
		response.Error(ctx, err)
		return
	} else {
		// 执行成功时返回成功响应
		response.Success(ctx, "")
		return
	}
}

// List 获取计划任务列表
// 支持按名称模糊查询和分页
func (app SchTaskApp) List(ctx *gin.Context) {
	log.Println("----------List------")
	log.Println("List externalNas")
	// 创建任务实例
	externalNas := scheduled.SchTask{}
	// 从查询参数中获取名称过滤条件
	name := ctx.Query("name")
	// 获取分页查询参数
	query := request.GetPageQuery(ctx)
	// 添加名称模糊查询条件
	query.AddFilter(request.NewLikeFilter("name", name))
	// 执行分页查询
	list, count, err := externalNas.List(query)
	if err == nil {
		// 查询成功时返回列表数据
		response.List(ctx, "", count, list)
	} else {
		// 查询失败或无数据时返回无内容响应
		response.NoContent(ctx, "无数据！")
	}
}
