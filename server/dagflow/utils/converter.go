// Package utils 提供DAGFlow的工具函数
package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"server/dagflow/model"
	"server/service/sflow"
	"server/utils/xxtea"
)

// FlowConverter SFlow转换器
// 负责将SFlow模型转换为DAGFlow的Flow模型
type FlowConverter struct{}

// ConvertFromSFlow 将SFlow转换为Flow
func (c *FlowConverter) ConvertFromSFlow(sflow *sflow.SFlow) (model.Flow, error) {
	if sflow == nil {
		return model.Flow{}, errors.New("SFlow为空")
	}

	if sflow.Content == "" {
		return model.Flow{}, errors.New("SFlow内容为空")
	}

	// 解析SFlow的Content字段，它应该是一个JSON字符串
	var flowData struct {
		Cells []Cell `json:"cells"`
	}

	content := xxtea.DecryptAuto(sflow.Content, "")

	err := json.Unmarshal([]byte(content), &flowData)
	if err != nil {
		return model.Flow{}, fmt.Errorf("解析SFlow内容失败: %v", err)
	}

	// 找到开始和结束节点
	startNodeID := ""
	endNodeID := ""

	// 转换为Flow格式
	nodes := make(map[string]model.TaskNode)
	edges := make(map[string]model.Edge)

	// 首先处理所有的节点
	for _, cell := range flowData.Cells {
		if cell.Shape == "dag-edge" {
			continue // 先跳过边，后面单独处理
		}

		// 处理节点
		taskNode := model.TaskNode{
			ID:         cell.ID,
			Name:       getNodeLabel(cell),
			Type:       cell.Shape,
			Properties: make(map[string]any),
		}

		// 根据节点类型设置特定属性
		if cell.Data != nil && cell.Data.Form != nil {
			// 复制表单中的所有字段到Properties
			for k, v := range cell.Data.Form {
				taskNode.Properties[k] = v
			}

			// 设置日志级别
			if logLevel, ok := cell.Data.Form["logLevel"].(string); ok {
				taskNode.LogLevel = logLevel
			}

			// 设置是否启用
			if enabled, ok := cell.Data.Form["enabled"].(bool); ok {
				taskNode.Disabled = !enabled
			}

			// 设置缓存时间
			if cacheTime, ok := cell.Data.Form["cacheTime"].(string); ok {
				taskNode.CacheTime = parseCacheTime(cacheTime)
			}

			// 设置异常处理
			if ignoreException, ok := cell.Data.Form["ignoreSimpleException"].(string); ok {
				taskNode.ExceptionHandle = ignoreException
			}

			// 设置结果名称
			if datakey, ok := cell.Data.Form["datakey"].(string); ok {
				taskNode.ResultName = datakey
			}
		}

		// 保存特定类型的节点ID
		if cell.Shape == "start" {
			if startNodeID != "" {
				return model.Flow{}, errors.New("流程图中有多个开始节点")
			}
			startNodeID = cell.ID
		} else if cell.Shape == "end" {
			if endNodeID != "" {
				return model.Flow{}, errors.New("流程图中有多个结束节点")
			}
			endNodeID = cell.ID
		}

		nodes[cell.ID] = taskNode
	}

	// 然后处理所有的边
	for _, cell := range flowData.Cells {
		if cell.Shape != "dag-edge" {
			continue
		}

		edge := model.Edge{
			ID:     cell.ID,
			Name:   getEdgeLabel(cell),
			Source: cell.Source.Cell,
			Target: cell.Target.Cell,
		}

		// 设置表达式
		if cell.Data != nil && cell.Data.Form != nil {
			if expr, ok := cell.Data.Form["expr"].(string); ok {
				edge.Expression = expr
			}
		}

		// 设置连接点
		if cell.Source.Port != "" {
			edge.SourceAnchor = cell.Source.Port
		}
		if cell.Target.Port != "" {
			edge.TargetAnchor = cell.Target.Port
		}

		edges[cell.ID] = edge
	}

	if startNodeID == "" {
		return model.Flow{}, errors.New("流程图中没有开始节点")
	}

	if endNodeID == "" {
		return model.Flow{}, errors.New("流程图中没有结束节点")
	}

	// 从结束节点获取返回结果配置
	returnResult := false
	resultType := ""
	if endNode, ok := nodes[endNodeID]; ok {
		if val, ok := endNode.Properties["returnResult"].(bool); ok {
			returnResult = val
		}
		if val, ok := endNode.Properties["resultType"].(string); ok {
			resultType = val
		}
	}

	// 创建Flow模型
	flow := model.Flow{
		ID:           sflow.ID,
		Name:         sflow.Name,
		Description:  sflow.Remark,
		Version:      "1.0",
		Disabled:     false,
		Nodes:        nodes,
		Edges:        edges,
		StartNodeID:  startNodeID,
		EndNodeID:    endNodeID,
		ReturnResult: returnResult,
		ResultType:   resultType,
	}

	return flow, nil
}

// Cell 表示DAGFlow中的一个单元格（节点或连接线）
type Cell struct {
	ID       string    `json:"id"`
	Shape    string    `json:"shape"`
	Position *Position `json:"position,omitempty"`
	Size     *Size     `json:"size,omitempty"`
	View     string    `json:"view,omitempty"`
	Data     *CellData `json:"data,omitempty"`
	Ports    *Ports    `json:"ports,omitempty"`
	ZIndex   int       `json:"zIndex,omitempty"`
	Source   *Endpoint `json:"source,omitempty"`
	Target   *Endpoint `json:"target,omitempty"`
	Attrs    *Attrs    `json:"attrs,omitempty"`
	Labels   []string  `json:"labels,omitempty"`
}

// Position 表示节点的位置
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Size 表示节点的大小
type Size struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// CellData 表示节点的数据
type CellData struct {
	Parent bool                   `json:"parent"`
	Form   map[string]interface{} `json:"form"`
}

// Ports 表示节点的连接点
type Ports struct {
	Groups map[string]PortGroup `json:"groups"`
	Items  []PortItem           `json:"items"`
}

// PortGroup 表示连接点组
type PortGroup struct {
	Position string              `json:"position"`
	Attrs    map[string]PortAttr `json:"attrs"`
}

// PortAttr 表示连接点属性
type PortAttr struct {
	R           int    `json:"r,omitempty"`
	Magnet      bool   `json:"magnet,omitempty"`
	Fill        string `json:"fill,omitempty"`
	FillOpacity string `json:"fillOpacity,omitempty"`
	Stroke      string `json:"stroke,omitempty"`
	StrokeWidth int    `json:"strokeWidth,omitempty"`
	Style       *Style `json:"style,omitempty"`
}

// Style 表示样式
type Style struct {
	Visibility string `json:"visibility"`
}

// PortItem 表示连接点项
type PortItem struct {
	Group string `json:"group"`
	ID    string `json:"id"`
}

// Endpoint 表示连接线的端点
type Endpoint struct {
	Cell string `json:"cell"`
	Port string `json:"port"`
}

// Attrs 表示连接线的属性
type Attrs struct {
	Line *Line `json:"line"`
}

// Line 表示线的属性
type Line struct {
	Stroke string `json:"stroke"`
}

// 获取节点标签
func getNodeLabel(cell Cell) string {
	if cell.Data != nil && cell.Data.Form != nil {
		if label, ok := cell.Data.Form["label"].(string); ok {
			return label
		}
	}
	return cell.Shape
}

// 获取边标签
func getEdgeLabel(cell Cell) string {
	if len(cell.Labels) > 0 {
		return cell.Labels[0]
	}
	return ""
}

// 解析缓存时间
func parseCacheTime(cacheTimeStr string) int {
	var cacheTime int
	_, err := fmt.Sscanf(cacheTimeStr, "%d", &cacheTime)
	if err != nil {
		return -1 // 默认不缓存
	}
	return cacheTime
}
