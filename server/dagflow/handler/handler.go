// Package handler 定义DAGFlow模块的任务处理器接口和实现
package handler

import (
	"context"
	"errors"
	"server/dagflow/model"
)

// TaskHandler 任务处理器接口
type TaskHandler interface {
	// 处理节点任务
	Handle(ctx context.Context, node model.TaskNode, execCtx *model.ExecutionContext) (any, error)

	// 获取任务处理器类型
	GetType() string

	// 验证节点配置是否有效
	Validate(node model.TaskNode) error
}

// HandlerRegistry 任务处理器注册表
type HandlerRegistry struct {
	handlers map[string]TaskHandler
}

// NewHandlerRegistry 创建新的任务处理器注册表
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[string]TaskHandler),
	}
}

// Register 注册任务处理器
func (r *HandlerRegistry) Register(handler TaskHandler) {
	r.handlers[handler.GetType()] = handler
}

// Get 获取任务处理器
func (r *HandlerRegistry) Get(nodeType string) (TaskHandler, error) {
	handler, ok := r.handlers[nodeType]
	if !ok {
		return nil, errors.New("未找到类型为 " + nodeType + " 的任务处理器")
	}
	return handler, nil
}

// GetAll 获取所有任务处理器
func (r *HandlerRegistry) GetAll() map[string]TaskHandler {
	return r.handlers
}
