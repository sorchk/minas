package system

import (
	"context"
	"errors"
	"server/dagflow/model"
)

// 节点类型常量
const (
	TypeStart = "start" // 开始节点
	TypeEnd   = "end"   // 结束节点
)

// StartNodeHandler 开始节点处理器
type StartNodeHandler struct{}

// GetType 获取处理器类型
func (h *StartNodeHandler) GetType() string {
	return TypeStart
}

// Handle 处理开始节点
func (h *StartNodeHandler) Handle(ctx context.Context, node model.TaskNode, execCtx *model.ExecutionContext) (any, error) {
	return nil, nil
}

// Validate 验证节点配置
func (h *StartNodeHandler) Validate(node model.TaskNode) error {
	// 开始节点不需要特殊配置
	return nil
}

// EndNodeHandler 结束节点处理器
type EndNodeHandler struct{}

// GetType 获取处理器类型
func (h *EndNodeHandler) GetType() string {
	return TypeEnd
}

// Handle 处理结束节点
func (h *EndNodeHandler) Handle(ctx context.Context, node model.TaskNode, execCtx *model.ExecutionContext) (any, error) {
	// 获取是否返回结果和结果类型
	returnResult, _ := node.Properties["returnResult"].(bool)
	resultType, _ := node.Properties["resultType"].(string)

	if !returnResult {
		// 不返回结果
		return nil, nil
	}

	// 根据结果类型返回不同数据
	switch resultType {
	case "all":
		// 返回所有数据
		return execCtx.Data, nil
	case "specified":
		// 返回指定数据
		specifiedKeys, ok := node.Properties["resultKeys"].([]string)
		if !ok {
			return nil, errors.New("结束节点配置错误：缺少 resultKeys 配置")
		}

		result := make(map[string]any)
		for _, key := range specifiedKeys {
			if value, exists := execCtx.GetData(key); exists {
				result[key] = value
			}
		}
		return result, nil
	default:
		// 默认返回所有数据
		return execCtx.Data, nil
	}
}

// Validate 验证节点配置
func (h *EndNodeHandler) Validate(node model.TaskNode) error {
	returnResult, _ := node.Properties["returnResult"].(bool)
	if returnResult {
		resultType, ok := node.Properties["resultType"].(string)
		if !ok {
			return errors.New("结束节点配置错误：缺少 resultType 配置")
		}

		if resultType == "specified" {
			_, ok := node.Properties["resultKeys"].([]string)
			if !ok {
				return errors.New("结束节点配置错误：缺少 resultKeys 配置")
			}
		}
	}

	return nil
}
