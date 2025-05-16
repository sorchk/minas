// Package schtask 提供计划任务日志相关的HTTP接口
package schtask

import (
	"server/core/app/request"
	"server/core/app/response"
	"server/service/scheduled/log"

	"github.com/gin-gonic/gin"
)

// SchLogApp 计划任务日志应用结构体
type SchLogApp struct {
}

// AddRoutes 注册计划任务日志相关的路由
// parentGroup: 父路由组
func (SchLogApp) AddRoutes(parentGroup *gin.RouterGroup) {
	group := parentGroup.Group("/log")
	app := SchLogApp{}
	group.GET("/list", app.List)     // 获取日志列表
	group.GET("/load/:id", app.Load) // 加载单个日志详情
}

// List 处理获取计划任务日志列表的请求
// ctx: Gin上下文
func (SchLogApp) List(ctx *gin.Context) {
	// 获取分页查询参数
	query := request.GetPageQuery(ctx)
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

// Load 处理加载单个计划任务日志详情的请求
// ctx: Gin上下文
func (SchLogApp) Load(ctx *gin.Context) {
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
