// Package dagflow 提供主服务入口
package dagflow

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/dagflow/engine"
	"server/dagflow/handler"
	"server/dagflow/handler/script"
	"server/dagflow/handler/system"
	"server/dagflow/model"
	"server/dagflow/utils"
	"server/service/sflow"
	"server/utils/logger"
)

// Service DAGFlow服务
type Service struct {
	// 流程执行引擎
	engine *engine.Engine

	// 处理器注册表
	handlerRegistry *handler.HandlerRegistry

	// SFlow转换器
	converter *utils.FlowConverter

	// 日志记录器
	logger model.LoggerInterface
}

// Logger 适配系统日志记录器
type Logger struct{}

// Debug 记录调试日志
func (l *Logger) Debug(msg string, args ...any) {
	logger.LOG.Debugf(msg, args...)
}

// Info 记录信息日志
func (l *Logger) Info(msg string, args ...any) {
	logger.LOG.Infof(msg, args...)
}

// Warn 记录警告日志
func (l *Logger) Warn(msg string, args ...any) {
	logger.LOG.Warnf(msg, args...)
}

// Error 记录错误日志
func (l *Logger) Error(msg string, args ...any) {
	logger.LOG.Errorf(msg, args...)
}

// NewService 创建新的DAGFlow服务
func NewService() *Service {
	// 创建日志记录器
	logger := &Logger{}

	// 创建处理器注册表
	registry := handler.NewHandlerRegistry()

	// 创建引擎
	eng := engine.NewEngine(registry, logger)
	// 注册基本处理器
	registry.Register(&system.StartNodeHandler{})
	registry.Register(&system.EndNodeHandler{})
	registry.Register(&system.LogNodeHandler{})
	registry.Register(&system.SleepNodeHandler{})

	// 注册JavaScript处理器
	registry.Register(script.NewJavaScriptHandler())

	// 注册休眠处理器

	// 创建SFlow转换器
	converter := &utils.FlowConverter{}

	return &Service{
		engine:          eng,
		handlerRegistry: registry,
		converter:       converter,
		logger:          logger,
	}
}

// RegisterHandler 注册自定义处理器
func (s *Service) RegisterHandler(h handler.TaskHandler) {
	s.handlerRegistry.Register(h)
}

// ExecuteFlow 执行流程
func (s *Service) ExecuteFlow(ctx context.Context, flowID string, params map[string]any, debug bool) (*model.ExecutionContext, error) {
	// 从SFlow加载流程
	sFlow := sflow.SFlow{}
	sFlow, err := sFlow.Load(flowID)
	if err != nil {
		return nil, fmt.Errorf("加载流程失败: %v", err)
	}

	// 转换为Flow模型
	flow, err := s.converter.ConvertFromSFlow(&sFlow)
	if err != nil {
		return nil, fmt.Errorf("转换流程失败: %v", err)
	}

	// 创建执行上下文
	execCtx := model.NewExecutionContext(flow.ID, params, s.logger)
	execCtx.Debug = debug // 设置debug模式

	// 调试模式下记录初始状态
	if debug {
		s.logger.Info("【调试模式】流程执行开始: %s (ID: %d)", flow.Name, flow.ID)
		s.logger.Info("【调试模式】初始数据: %v", params)
	}

	// 执行流程
	s.logger.Info("开始执行流程: %s (ID: %d)", flow.Name, flow.ID)
	execCtx, err = s.engine.Execute(ctx, flow, params)

	// 调试模式下记录完整执行结果
	if debug && execCtx != nil {
		s.logger.Info("【调试模式】流程执行完成: %s (ID: %d)", flow.Name, flow.ID)
		s.logger.Info("【调试模式】执行状态: %s", execCtx.Status)

		// 记录每个节点的执行信息
		for nodeID, status := range execCtx.NodeStatus {
			node, exists := flow.Nodes[nodeID]
			nodeName := nodeID
			if exists {
				nodeName = node.Name
			}

			s.logger.Info("【调试模式】节点 %s (%s) 状态: %s", nodeName, nodeID, status)

			// 记录节点执行时间
			if startTime, ok := execCtx.NodeStartTimes[nodeID]; ok {
				if endTime, ok := execCtx.NodeEndTimes[nodeID]; ok {
					duration := endTime.Sub(startTime)
					s.logger.Info("【调试模式】节点 %s 执行时间: %v", nodeName, duration)
				}
			}

			// 记录节点结果数据
			if result, ok := execCtx.NodeResults[nodeID]; ok {
				resultJSON, _ := json.Marshal(result)
				s.logger.Info("【调试模式】节点 %s 结果数据: %s", nodeName, string(resultJSON))
			}

			// 记录节点错误信息
			if errMsg, ok := execCtx.NodeErrors[nodeID]; ok && errMsg != "" {
				s.logger.Error("【调试模式】节点 %s 错误信息: %s", nodeName, errMsg)
			}
		}
	}

	if err != nil {
		s.logger.Error("流程执行失败: %v", err)
		return execCtx, err
	}

	s.logger.Info("流程执行完成: %s (ID: %d)", flow.Name, flow.ID)
	return execCtx, nil
}

// ValidateFlow 验证流程
func (s *Service) ValidateFlow(flowID string) (bool, []string, error) {
	// 从SFlow加载流程
	sFlow := sflow.SFlow{}
	sFlow, err := sFlow.Load(flowID)
	if err != nil {
		return false, nil, fmt.Errorf("加载流程失败: %v", err)
	}

	// 转换为Flow模型
	flow, err := s.converter.ConvertFromSFlow(&sFlow)
	if err != nil {
		return false, nil, fmt.Errorf("转换流程失败: %v", err)
	}

	// 验证节点
	errors := []string{}
	valid := true

	// 检查开始节点和结束节点
	if flow.StartNodeID == "" {
		errors = append(errors, "流程没有开始节点")
		valid = false
	}

	if flow.EndNodeID == "" {
		errors = append(errors, "流程没有结束节点")
		valid = false
	}

	// 验证每个节点的配置
	for id, node := range flow.Nodes {
		// 获取节点处理器
		handler, err := s.handlerRegistry.Get(node.Type)
		if err != nil {
			errors = append(errors, fmt.Sprintf("节点 %s (%s) 类型无效: %v", node.Name, id, err))
			valid = false
			continue
		}

		// 验证节点配置
		if err := handler.Validate(node); err != nil {
			errors = append(errors, fmt.Sprintf("节点 %s (%s) 配置无效: %v", node.Name, id, err))
			valid = false
		}
	}

	return valid, errors, nil
}

// GetRegisteredHandlers 获取已注册的处理器类型
func (s *Service) GetRegisteredHandlers() []string {
	handlers := s.handlerRegistry.GetAll()
	types := make([]string, 0, len(handlers))

	for handlerType := range handlers {
		types = append(types, handlerType)
	}

	return types
}

// 全局服务实例
var defaultService *Service

// Init 初始化DAGFlow服务
func Init() {
	if defaultService == nil {
		defaultService = NewService()
		log.Println("DAGFlow服务已初始化")
	}
}

// GetService 获取DAGFlow服务实例
func GetService() *Service {
	if defaultService == nil {
		Init()
	}
	return defaultService
}
