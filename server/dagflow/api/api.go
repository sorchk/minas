// Package dagflow/api 提供DAGFlow的API接口
package api

import (
	"errors"
	"server/core/app/response"
	"server/dagflow"

	"github.com/gin-gonic/gin"
)

// DAGFlowAPI DAGFlow API控制器
type DAGFlowAPI struct{}

// AddRoutes 注册DAGFlow相关的路由
func (api *DAGFlowAPI) AddRoutes(group *gin.RouterGroup) {

	// 注册路由
	group.POST("/execute/:id", api.ExecuteFlow)
	group.POST("/validate/:id", api.ValidateFlow)
	group.GET("/handlers", api.GetHandlers)
	group.POST("/debug/:id", api.DebugFlow)
}

// ExecuteRequest 执行请求参数
type ExecuteRequest struct {
	Data map[string]any `json:"data"` // 初始数据
}

// ExecuteFlow 执行流程
func (api *DAGFlowAPI) ExecuteFlow(ctx *gin.Context) {
	// 获取流程ID
	id := ctx.Param("id")
	if id == "" {
		response.BadRequest(ctx, "流程ID不能为空")
		return
	}

	// 解析请求体
	var req ExecuteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	// 确保初始数据不为空
	if req.Data == nil {
		req.Data = make(map[string]any)
	}

	// 获取服务实例
	service := dagflow.GetService()
	if service == nil {
		response.Error(ctx, errors.New("DAGFlow服务未初始化"))
		return
	}

	// 执行流程
	execCtx, err := service.ExecuteFlow(ctx, id, req.Data, false)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	// 返回执行结果
	response.Data(ctx, "流程执行成功", execCtx)
}

// ValidateFlow 验证流程
func (api *DAGFlowAPI) ValidateFlow(ctx *gin.Context) {
	// 获取流程ID
	id := ctx.Param("id")
	if id == "" {
		response.BadRequest(ctx, "流程ID不能为空")
		return
	}

	// 获取服务实例
	service := dagflow.GetService()
	if service == nil {
		response.Error(ctx, errors.New("DAGFlow服务未初始化"))
		return
	}

	// 验证流程
	valid, errors, err := service.ValidateFlow(id)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	// 返回验证结果
	response.Data(ctx, "流程验证完成", gin.H{
		"valid":  valid,
		"errors": errors,
	})
}

// GetHandlers 获取所有注册的处理器类型
func (api *DAGFlowAPI) GetHandlers(ctx *gin.Context) {
	// 获取服务实例
	service := dagflow.GetService()
	if service == nil {
		response.Error(ctx, errors.New("DAGFlow服务未初始化"))
		return
	}

	// 获取处理器类型
	handlers := service.GetRegisteredHandlers()

	// 返回处理器类型列表
	response.Data(ctx, "获取处理器类型成功", handlers)
}

// DebugFlow 调试流程
func (api *DAGFlowAPI) DebugFlow(ctx *gin.Context) {
	// 获取流程ID
	id := ctx.Param("id")
	if id == "" {
		response.BadRequest(ctx, "流程ID不能为空")
		return
	}

	// 解析请求体
	var req ExecuteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	// 确保初始数据不为空
	if req.Data == nil {
		req.Data = make(map[string]any)
	}

	// 获取服务实例
	service := dagflow.GetService()
	if service == nil {
		response.Error(ctx, errors.New("DAGFlow服务未初始化"))
		return
	}

	// 执行流程，开启调试模式
	execCtx, err := service.ExecuteFlow(ctx, id, req.Data, true)
	if err != nil {
		response.DataCode(ctx, 500, "流程调试失败", execCtx)
		return
	}
	// 返回执行结果和调试信息
	response.Data(ctx, "流程调试完成", execCtx)
}
