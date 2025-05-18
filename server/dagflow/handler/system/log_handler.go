package system

import (
	"context"
	"errors"
	"fmt"
	"server/dagflow/model"
	"server/dagflow/utils"
)

// 日志相关的节点类型常量
const (
	// ...existing code...
	TypeLog = "Log" // 增强日志节点
)

// LogNodeHandler 增强日志节点处理器
// 支持使用EL表达式处理日志内容
type LogNodeHandler struct{}

// GetType 获取处理器类型
func (h *LogNodeHandler) GetType() string {
	return TypeLog
}

// Handle 处理增强日志节点
func (h *LogNodeHandler) Handle(ctx context.Context, node model.TaskNode, execCtx *model.ExecutionContext) (any, error) {
	// 获取日志级别
	logLevel, _ := node.Properties["logLevel"].(string)
	if logLevel == "" {
		logLevel = "info" // 默认级别
	}

	// 计算日志内容
	content := utils.GetStr(node, "logInfo", "", execCtx.Data)

	// 记录日志
	execCtx.Log(logLevel, "[%s] %s", node.Name, content)

	// 设置节点调试信息
	if execCtx.Debug {
		nodeInfo := map[string]any{
			"content": content,
			"level":   logLevel,
		}
		execCtx.SetNodeResult(node.ID, nodeInfo)
	}

	// 返回日志内容
	return content, nil
}

// Validate 验证节点配置
func (h *LogNodeHandler) Validate(node model.TaskNode) error {
	// 检查日志模板
	logInfo, ok := node.Properties["logInfo"].(string)
	if !ok || logInfo == "" {
		return errors.New("日志节点配置错误：缺少或为空的logInfo配置")
	}

	// 验证日志级别
	if logLevel, ok := node.Properties["logLevel"].(string); ok && logLevel != "" {
		// 检查日志级别是否有效
		validLevels := map[string]bool{
			"debug": true,
			"info":  true,
			"warn":  true,
			"error": true,
		}

		if !validLevels[logLevel] {
			return fmt.Errorf("增强日志节点配置错误：无效的日志级别 %s", logLevel)
		}
	}

	return nil
}
