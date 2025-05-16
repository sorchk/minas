package sflow

import (
	"log"
	"reflect"
	"server/core/app/request"
	"server/core/app/response"
	"server/core/app/webapi"
	"server/service/sflow"
	"server/utils/logger"

	"github.com/gin-gonic/gin"
)

// SFlowApp 作业流程应用结构体
// 继承了BaseApp基类，用于处理作业流程的API请求
type SFlowApp struct {
	webapi.BaseApp[sflow.SFlow] // 内嵌基类，指定模型类型为sflow.SFlow
}

// AddRoutes 注册作业流程相关的路由
// 在父路由组下创建作业流程路由组，并注册相关的路由处理函数
func (SFlowApp) AddRoutes(parentGroup *gin.RouterGroup) {
	// 创建作业流程路由组
	group := parentGroup.Group("/sflow")
	// 创建应用实例
	app := SFlowApp{}
	// 注册基础路由（CRUD操作）
	webapi.AddBaseRoutes(group, &app)
	// 设置允许更新的字段列表
	app.UpdateFields = []string{"name", "type", "log_level", "uri", "project_dir_id", "remark"}

	// 注册手动执行任务的路由
	group.POST("/exec/:id", app.Exec)
	group.POST("/test/:id", app.Test)
	group.POST("/saveContent", app.saveContent)
}
func (app SFlowApp) Test(ctx *gin.Context) {

	log.Println("----------Test------")

}
func (app SFlowApp) saveContent(ctx *gin.Context) {
	log.Println("----------saveContent------")
	var entity sflow.SFlow
	// 从请求体绑定JSON数据到实体
	if err := ctx.BindJSON(&entity); err != nil {
		response.BadRequest(ctx, "参数错误！")
		logger.LOG.Errorf("参数错误:%s", err.Error())
		return
	}
	// 获取当前用户ID
	uid := request.GetUserID(ctx)
	// 使用反射设置操作者ID
	setOperatorFunc := reflect.ValueOf(&entity).MethodByName("SetOperatorUID")
	if setOperatorFunc.IsValid() {
		args := []reflect.Value{reflect.ValueOf(uid)}
		setOperatorFunc.Call(args)
	}
	// 保存实体，并指定允许更新的字段
	err := entity.Save(&entity, "content")
	if err == nil {
		// 成功时返回保存后的实体
		response.Data(ctx, "", entity)
	} else {
		// 出错时返回错误信息
		response.Error(ctx, err)
	}
}

// Exec 手动执行指定的作业流程
// 根据ID加载任务并执行，返回执行结果
func (app SFlowApp) Exec(ctx *gin.Context) {
	log.Println("----------Exec------")
	// 创建任务实例
	sFlow := sflow.SFlow{}
	// 从路径参数中获取任务ID
	id := ctx.Param("id")
	// 根据ID加载任务信息
	sFlow, err := sFlow.Load(id)
	if err != nil {
		// 加载失败时返回错误
		response.Error(ctx, err)
		return
	}
	// 执行任务
	// err = sFlow.Exec()
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

// List 获取作业流程列表
// 支持按名称模糊查询和分页
func (app SFlowApp) List(ctx *gin.Context) {
	// 创建任务实例
	sFlow := sflow.SFlow{}
	// 获取分页查询参数
	query := request.GetPageQuery(ctx)
	// 添加名称模糊查询条件
	name := ctx.Query("name")
	if name != "" {
		query.AddFilter(request.NewLikeFilter("name", name))
	}
	// 添加类型过滤条件
	flowType := ctx.Query("type")
	query.AddFilter(request.NewEqualFilter("type", flowType))

	project_dir_id := ctx.Query("project_dir_id")
	if project_dir_id != "" {
		// 如果提供了项目目录ID，则添加过滤条件
		query.AddFilter(request.NewEqualFilter("project_dir_id", project_dir_id))
	}
	// 执行分页查询
	list, count, err := sFlow.List(query)
	if err == nil {
		// 查询成功时返回列表数据
		response.List(ctx, "", count, list)
	} else {
		// 查询失败或无数据时返回无内容响应
		response.NoContent(ctx, "无数据！")
	}
}
