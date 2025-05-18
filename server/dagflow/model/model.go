// Package model 定义DAGFlow模块的数据模型
package model

import (
	"encoding/json"
	"time"
)

// Node 节点基础接口
type Node interface {
	GetID() string
	GetName() string
	IsDisabled() bool
}

// TaskNode 任务节点模型
type TaskNode struct {
	ID              string         `json:"id"`              // 节点ID
	Name            string         `json:"name"`            // 节点名称
	Type            string         `json:"type"`            // 节点类型，决定如何处理数据
	ResultName      string         `json:"resultName"`      // 结果命名，为空则使用task_id作为数据key
	CacheTime       int            `json:"cacheTime"`       // 缓存时间(秒)，默认-1不缓存
	ExceptionHandle string         `json:"exceptionHandle"` // 异常处理方式，根据el表达式判断是否继续运行，为空则发生异常时终止运行
	Disabled        bool           `json:"disabled"`        // 是否禁用，默认启用
	LogLevel        string         `json:"logLevel"`        // 日志级别，默认无
	Properties      map[string]any `json:"properties"`      // 节点属性，根据节点类型不同包含不同的属性
}

// GetID 获取节点ID
func (t TaskNode) GetID() string {
	return t.ID
}

// GetName 获取节点名称
func (t TaskNode) GetName() string {
	return t.Name
}

// IsDisabled 检查节点是否被禁用
func (t TaskNode) IsDisabled() bool {
	return t.Disabled
}

// GetResultKey 获取结果数据的键名
func (t TaskNode) GetResultKey() string {
	if t.ResultName == "" {
		return t.ID
	}
	return t.ResultName
}

// Edge 连线模型
type Edge struct {
	ID           string `json:"id"`           // 连线ID
	Name         string `json:"name"`         // 连线名称
	Expression   string `json:"expression"`   // 表达式，el表达式判断是否执行后续任务
	Source       string `json:"source"`       // 源节点ID
	Target       string `json:"target"`       // 目标节点ID
	SourceAnchor string `json:"sourceAnchor"` // 源连接点
	TargetAnchor string `json:"targetAnchor"` // 目标连接点
}

// GetID 获取连线ID
func (e Edge) GetID() string {
	return e.ID
}

// GetName 获取连线名称
func (e Edge) GetName() string {
	return e.Name
}

// IsDisabled 连线不支持禁用
func (e Edge) IsDisabled() bool {
	return false
}

// Flow 流程图模型
type Flow struct {
	ID           uint                `json:"id"`           // 流程ID
	Name         string              `json:"name"`         // 流程名称
	Description  string              `json:"description"`  // 流程描述
	Version      string              `json:"version"`      // 版本号
	Disabled     bool                `json:"disabled"`     // 是否禁用
	Nodes        map[string]TaskNode `json:"nodes"`        // 节点集合
	Edges        map[string]Edge     `json:"edges"`        // 连线集合
	StartNodeID  string              `json:"startNodeId"`  // 开始节点ID
	EndNodeID    string              `json:"endNodeId"`    // 结束节点ID
	ReturnResult bool                `json:"returnResult"` // 是否返回结果
	ResultType   string              `json:"resultType"`   // 返回结果类型
}

// NodeExecutionStatus 节点执行状态
type NodeExecutionStatus string

const (
	Pending   NodeExecutionStatus = "pending"   // 待执行
	Running   NodeExecutionStatus = "running"   // 正在执行
	Completed NodeExecutionStatus = "completed" // 已完成
	Failed    NodeExecutionStatus = "failed"    // 执行失败
	Skipped   NodeExecutionStatus = "skipped"   // 已跳过
)

// ExecutionContext 执行上下文，保存流程执行过程中的数据
type ExecutionContext struct {
	FlowID         uint                           `json:"flowId"`         // 流程ID
	ExecutionID    string                         `json:"executionId"`    // 执行ID
	StartTime      time.Time                      `json:"startTime"`      // 开始时间
	EndTime        time.Time                      `json:"endTime"`        // 结束时间
	Status         NodeExecutionStatus            `json:"status"`         // 执行状态
	Params         map[string]any                 `json:"params"`         // 参数
	Data           map[string]any                 `json:"data"`           // 数据存储
	NodeStatus     map[string]NodeExecutionStatus `json:"nodeStatus"`     // 节点状态
	NodeStartTimes map[string]time.Time           `json:"nodeStartTimes"` // 节点开始时间
	NodeEndTimes   map[string]time.Time           `json:"nodeEndTimes"`   // 节点结束时间
	NodeResults    map[string]any                 `json:"nodeResults"`    // 节点执行结果
	NodeErrors     map[string]string              `json:"nodeErrors"`     // 节点执行错误
	Debug          bool                           `json:"debug"`          // 是否为调试模式
	ParentContext  *ExecutionContext              `json:"-"`              // 父执行上下文(用于子流程)
	SubContexts    map[string][]*ExecutionContext `json:"-"`              // 子执行上下文(用于迭代节点)
	logger         LoggerInterface                `json:"-"`              // 日志记录器
}

// LoggerInterface 日志接口
type LoggerInterface interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// NewExecutionContext 创建新的执行上下文
func NewExecutionContext(flowID uint, params map[string]any, logger LoggerInterface) *ExecutionContext {
	now := time.Now()
	return &ExecutionContext{
		FlowID:         flowID,
		ExecutionID:    generateExecutionID(),
		StartTime:      now,
		Status:         Pending,
		Params:         params,
		Data:           make(map[string]any),
		NodeStatus:     make(map[string]NodeExecutionStatus),
		NodeStartTimes: make(map[string]time.Time),
		NodeEndTimes:   make(map[string]time.Time),
		NodeResults:    make(map[string]any),
		NodeErrors:     make(map[string]string),
		SubContexts:    make(map[string][]*ExecutionContext),
		Debug:          false,
		logger:         logger,
	}
}

// SetData 设置数据
func (ctx *ExecutionContext) SetData(key string, value any) {
	ctx.Data[key] = value
}

// GetData 获取数据
func (ctx *ExecutionContext) GetData(key string) (any, bool) {
	v, ok := ctx.Data[key]
	return v, ok
}

// SetNodeResult 设置节点执行结果
func (ctx *ExecutionContext) SetNodeResult(nodeID string, result any) {
	ctx.NodeResults[nodeID] = result
}

// GetNodeResult 获取节点执行结果
func (ctx *ExecutionContext) GetNodeResult(nodeID string) (any, bool) {
	v, ok := ctx.NodeResults[nodeID]
	return v, ok
}

// SetNodeStatus 设置节点执行状态
func (ctx *ExecutionContext) SetNodeStatus(nodeID string, status NodeExecutionStatus) {
	ctx.NodeStatus[nodeID] = status
}

// GetNodeStatus 获取节点执行状态
func (ctx *ExecutionContext) GetNodeStatus(nodeID string) NodeExecutionStatus {
	if status, ok := ctx.NodeStatus[nodeID]; ok {
		return status
	}
	return Pending
}

// Log 记录日志
func (ctx *ExecutionContext) Log(level string, message string, args ...any) {
	if ctx.logger == nil {
		return
	}

	switch level {
	case "debug":
		ctx.logger.Debug(message, args...)
	case "info":
		ctx.logger.Info(message, args...)
	case "warn":
		ctx.logger.Warn(message, args...)
	case "error":
		ctx.logger.Error(message, args...)
	default:
		ctx.logger.Info(message, args...)
	}
}

// Clone 克隆执行上下文(用于迭代节点)
func (ctx *ExecutionContext) Clone() *ExecutionContext {
	dataCopy := make(map[string]any, len(ctx.Data))
	for k, v := range ctx.Data {
		dataCopy[k] = v
	}

	clone := NewExecutionContext(ctx.FlowID, dataCopy, ctx.logger)
	clone.ParentContext = ctx
	return clone
}

// ToJSON 将ExecutionContext转为JSON字符串
func (ctx *ExecutionContext) ToJSON() (string, error) {
	bytes, err := json.Marshal(ctx)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON 从JSON字符串恢复ExecutionContext
func (ctx *ExecutionContext) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), ctx)
}

// generateExecutionID 生成执行ID
func generateExecutionID() string {
	// 使用当前时间戳和随机数生成唯一ID
	return time.Now().Format("20060102150405") + "_" + RandomString(8)
}

// RandomString 生成指定长度的随机字符串
func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}
