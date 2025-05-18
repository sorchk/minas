// Package engine 实现DAGFlow的流程执行引擎
package engine

import (
	"context"
	"errors"
	"fmt"
	"server/dagflow/core/el"
	"server/dagflow/handler"
	"server/dagflow/model"
	"sync"
	"time"
)

// Engine 流程执行引擎
type Engine struct {
	handlerRegistry *handler.HandlerRegistry
	logger          model.LoggerInterface
}

// NewEngine 创建新的流程执行引擎
func NewEngine(registry *handler.HandlerRegistry, logger model.LoggerInterface) *Engine {
	return &Engine{
		handlerRegistry: registry,
		logger:          logger,
	}
}

// Execute 执行流程
func (e *Engine) Execute(ctx context.Context, flow model.Flow, initialData map[string]any) (*model.ExecutionContext, error) {
	if flow.StartNodeID == "" {
		return nil, errors.New("流程图没有指定开始节点")
	}

	if flow.EndNodeID == "" {
		return nil, errors.New("流程图没有指定结束节点")
	}

	// 初始化执行上下文
	execCtx := model.NewExecutionContext(flow.ID, initialData, e.logger)

	// 检查流程是否被禁用
	if flow.Disabled {
		execCtx.Status = model.Skipped
		execCtx.Log("warn", "流程已被禁用，跳过执行")
		return execCtx, nil
	}

	// 开始执行流程
	execCtx.Status = model.Running
	execCtx.Log("info", "开始执行流程: %s", flow.Name)

	// 从开始节点开始执行
	err := e.executeNode(ctx, flow, flow.StartNodeID, execCtx)

	// 设置执行结束时间
	execCtx.EndTime = time.Now()

	// 根据执行结果设置状态
	if err != nil {
		execCtx.Status = model.Failed
		execCtx.Log("error", "流程执行失败: %v", err)
		return execCtx, err
	}

	execCtx.Status = model.Completed
	execCtx.Log("info", "流程执行完成: %s", flow.Name)
	return execCtx, nil
}

// ExecuteSubFlow 执行子流程，用于迭代节点
func (e *Engine) ExecuteSubFlow(ctx context.Context, flow model.Flow, initialData map[string]any) (*model.ExecutionContext, error) {
	// 创建子流程的执行上下文
	execCtx := model.NewExecutionContext(flow.ID, initialData, e.logger)

	// 执行子流程，重用Execute方法的逻辑
	if flow.StartNodeID == "" {
		return nil, errors.New("子流程没有指定开始节点")
	}

	if flow.EndNodeID == "" {
		return nil, errors.New("子流程没有指定结束节点")
	}

	// 检查流程是否被禁用
	if flow.Disabled {
		execCtx.Status = model.Skipped
		execCtx.Log("warn", "子流程已被禁用，跳过执行")
		return execCtx, nil
	}

	// 开始执行流程
	execCtx.Status = model.Running
	execCtx.Log("info", "开始执行子流程: %s", flow.Name)

	// 从开始节点开始执行
	err := e.executeNode(ctx, flow, flow.StartNodeID, execCtx)

	// 设置执行结束时间
	execCtx.EndTime = time.Now()

	// 根据执行结果设置状态
	if err != nil {
		execCtx.Status = model.Failed
		execCtx.Log("error", "子流程执行失败: %v", err)
		return execCtx, err
	}

	execCtx.Status = model.Completed
	execCtx.Log("info", "子流程执行完成: %s", flow.Name)
	return execCtx, nil
}

// executeNode 执行单个节点
func (e *Engine) executeNode(ctx context.Context, flow model.Flow, nodeID string, execCtx *model.ExecutionContext) error {
	// 检查节点是否已经执行过
	status := execCtx.GetNodeStatus(nodeID)
	if status == model.Completed || status == model.Skipped {
		execCtx.Log("debug", "节点 %s 已经执行过，跳过", nodeID)
		return nil
	}

	// 获取节点
	node, ok := flow.Nodes[nodeID]
	if !ok {
		return fmt.Errorf("未找到节点: %s", nodeID)
	}

	// 检查节点是否被禁用
	if node.IsDisabled() {
		execCtx.SetNodeStatus(nodeID, model.Skipped)
		execCtx.Log("info", "节点 %s 已被禁用，跳过执行", node.Name)
		return e.executeNextNodes(ctx, flow, nodeID, execCtx)
	}

	// 标记节点为运行中
	execCtx.SetNodeStatus(nodeID, model.Running)
	execCtx.NodeStartTimes[nodeID] = time.Now()

	execCtx.Log("info", "开始执行节点: %s (%s)", node.Name, node.Type)

	// 如果是结束节点，直接返回
	if nodeID == flow.EndNodeID {
		execCtx.SetNodeStatus(nodeID, model.Completed)
		execCtx.NodeEndTimes[nodeID] = time.Now()
		execCtx.Log("info", "结束节点执行完成")
		return nil
	}

	// 获取处理器
	taskHandler, err := e.handlerRegistry.Get(node.Type)
	if err != nil {
		execCtx.SetNodeStatus(nodeID, model.Failed)
		execCtx.NodeErrors[nodeID] = err.Error()
		execCtx.Log("error", "获取节点处理器失败: %v", err)
		return err
	}

	// 执行节点处理器
	result, err := taskHandler.Handle(ctx, node, execCtx)

	// 记录执行结束时间
	execCtx.NodeEndTimes[nodeID] = time.Now()

	// 处理执行结果
	if err != nil {
		execCtx.SetNodeStatus(nodeID, model.Failed)
		execCtx.NodeErrors[nodeID] = err.Error()
		execCtx.Log("error", "节点执行失败: %v", err)

		// 检查异常处理方式
		if node.ExceptionHandle != "" {
			// TODO: 根据el表达式判断是否继续执行
			continueExecution := true
			if continueExecution {
				execCtx.Log("warn", "根据异常处理配置，继续执行后续节点")
				return e.executeNextNodes(ctx, flow, nodeID, execCtx)
			}
		}

		return err
	}

	// 保存结果到执行上下文
	resultKey := node.GetResultKey()
	execCtx.SetData(resultKey, result)
	execCtx.Log("debug", "节点 %s 执行结果: %v 数据类型：%T", node.Name, result, result)
	execCtx.SetNodeResult(nodeID, result)
	execCtx.SetNodeStatus(nodeID, model.Completed)

	execCtx.Log("info", "节点 %s 执行完成", node.Name)

	// 执行后续节点
	return e.executeNextNodes(ctx, flow, nodeID, execCtx)
}

// executeNextNodes 执行后续节点
func (e *Engine) executeNextNodes(ctx context.Context, flow model.Flow, nodeID string, execCtx *model.ExecutionContext) error {
	// 如果是结束节点，直接返回
	if nodeID == flow.EndNodeID {
		return nil
	}

	// 查找后续连线和节点
	nextEdges := findOutgoingEdges(flow, nodeID)
	if len(nextEdges) == 0 {
		execCtx.Log("warn", "节点 %s 没有后续节点", nodeID)
		return nil
	}

	// 过滤有表达式条件的连线
	availableEdges := make([]model.Edge, 0)
	for _, edge := range nextEdges {
		if edge.Expression != "" {
			// 使用EL表达式评估器计算表达式结果
			result, err := el.Evaluate(edge.Expression, execCtx.Data)
			if err != nil {
				execCtx.Log("error", "计算连线 %s 的表达式失败: %v, 跳过该连线", edge.Name, err)
				return err
			}
			execCtx.Log("info", "连线 %s 的表达式计算结果: %v, 数据类型为：%T", edge.Name, result, result)
			// 将结果转换为布尔值
			var expressionResult bool
			switch v := result.(type) {
			case bool:
				expressionResult = v
			case string:
				// 尝试将字符串解析为布尔值
				expressionResult = v == "true" || v == "True" || v == "TRUE" || v == "1"
			case int, int64, float64:
				// 非零值视为true
				expressionResult = v != 0
			default:
				// 其他类型视为存在即为true
				expressionResult = result != nil
			}

			if !expressionResult {
				execCtx.Log("info", "连线 %s 表达式条件不满足，跳过", edge.Name)

				// 标记目标节点为已跳过
				targetNode, ok := flow.Nodes[edge.Target]
				if ok {
					execCtx.SetNodeStatus(edge.Target, model.Skipped)
					execCtx.Log("info", "节点 %s 被跳过(表达式结果为假)", targetNode.Name)

					// 判断该节点是否为单入口节点(只有一条入边)
					isOnlyIncoming := isSingleIncomingNode(flow, edge.Target)
					if isOnlyIncoming {
						// 继续标记后续分支的节点为已跳过
						e.markSkippedBranch(flow, edge.Target, execCtx)
					}
				}
				continue
			}
		}

		// 如果没有表达式或表达式为真，则将边添加到可用列表
		availableEdges = append(availableEdges, edge)
	}

	// 如果没有可用的连线，返回成功
	if len(availableEdges) == 0 {
		return nil
	}

	// 执行所有可用连线的目标节点
	if len(availableEdges) == 1 {
		// 只有一个后续节点，直接执行
		return e.executeNode(ctx, flow, availableEdges[0].Target, execCtx)
	} else {
		// 多个后续节点，并行执行
		return e.executeParallelNodes(ctx, flow, availableEdges, execCtx)
	}
}

// isSingleIncomingNode 判断节点是否只有一条入边（单入口节点）
func isSingleIncomingNode(flow model.Flow, nodeID string) bool {
	count := 0
	for _, edge := range flow.Edges {
		if edge.Target == nodeID {
			count++
			if count > 1 {
				return false
			}
		}
	}
	return count == 1
}

// executeParallelNodes 并行执行多个节点
func (e *Engine) executeParallelNodes(ctx context.Context, flow model.Flow, edges []model.Edge, execCtx *model.ExecutionContext) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(edges))

	for _, edge := range edges {
		wg.Add(1)
		go func(edge model.Edge) {
			defer wg.Done()
			err := e.executeNode(ctx, flow, edge.Target, execCtx)
			if err != nil {
				errChan <- err
			}
		}(edge)
	}

	// 等待所有分支执行完成
	wg.Wait()
	close(errChan)

	// 检查是否有错误
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// markSkippedBranch 标记被跳过的分支上的所有节点为已跳过
func (e *Engine) markSkippedBranch(flow model.Flow, startNodeID string, execCtx *model.ExecutionContext) {
	// 已经标记过的节点不需要重复标记
	if execCtx.GetNodeStatus(startNodeID) == model.Skipped {
		return
	}

	execCtx.SetNodeStatus(startNodeID, model.Skipped)

	// 查找所有后续节点并标记
	nextEdges := findOutgoingEdges(flow, startNodeID)
	for _, edge := range nextEdges {
		e.markSkippedBranch(flow, edge.Target, execCtx)
	}
}

// findOutgoingEdges 查找节点的所有出边
func findOutgoingEdges(flow model.Flow, nodeID string) []model.Edge {
	var edges []model.Edge
	for _, edge := range flow.Edges {
		if edge.Source == nodeID {
			edges = append(edges, edge)
		}
	}
	return edges
}
