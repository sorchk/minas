// Package sflow 提供作业流程日志相关的HTTP接口
package sflow

import (
	"server/core/app/request"
	"server/core/app/response"
	"server/service/scheduled/log"

	"github.com/gin-gonic/gin"
)

// SFlowLogApp 作业流程日志应用结构体
type SFlowLogApp struct {
}

// AddRoutes 注册作业流程日志相关的路由
// parentGroup: 父路由组
func (SFlowLogApp) AddRoutes(parentGroup *gin.RouterGroup) {
	group := parentGroup.Group("/log")
	app := SFlowLogApp{}
	group.GET("/list", app.List)     // 获取日志列表
	group.GET("/load/:id", app.Load) // 加载单个日志详情
}

// List 处理获取作业流程日志列表的请求
// ctx: Gin上下文
func (SFlowLogApp) List(ctx *gin.Context) {
	// 获取分页查询参数
	query := request.GetPageQuery(ctx)
	flow_id := ctx.Query("flow_id")
	query.AddFilter(request.NewEqualFilter("flow_id", flow_id))
	var entity log.SchLog
	// 调用服务层获取日志列表
	list, count, err := entity.List(query)
	if err == nil {
		// 成功时返回日志列表数据
		response.List(ctx, "", count, list)
	} else {
		// 失败时返回无内容响应
		response.NoContent(ctx, "无数据！")
	}
}

// Load 处理加载单个作业流程日志详情的请求
// ctx: Gin上下文
func (SFlowLogApp) Load(ctx *gin.Context) {
	var entity log.SchLog
	// 根据ID加载日志详情
	entity, err := entity.Load(ctx.Param("id"))
	if err == nil {
		// 成功时返回日志详情
		response.Data(ctx, "", entity)
		return
	} else {
		// 失败时返回错误信息
		response.Error(ctx, err)
		return
	}
}
