package system

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"server/dagflow/model"
	"strconv"
	"time"
)

// 休眠节点类型常量
const (
	TypeSleep = "Sleep" // 休眠节点
)

// SleepNodeHandler 休眠节点处理器
// 支持在min和max时间范围内随机休眠
type SleepNodeHandler struct{}

// GetType 获取处理器类型
func (h *SleepNodeHandler) GetType() string {
	return TypeSleep
}

// Handle 处理休眠节点
func (h *SleepNodeHandler) Handle(ctx context.Context, node model.TaskNode, execCtx *model.ExecutionContext) (any, error) {
	// 获取最小休眠时间(毫秒)
	minTimeStr, _ := node.Properties["min"].(string)
	if minTimeStr == "" {
		minTimeStr = "0" // 默认最小时间0毫秒
	}

	// 获取最大休眠时间(毫秒)
	maxTimeStr, _ := node.Properties["max"].(string)
	if maxTimeStr == "" {
		return nil, errors.New("休眠节点配置错误：缺少或为空的max配置")
	}

	// 解析最小和最大休眠时间
	minTime, err := strconv.ParseInt(minTimeStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("休眠节点配置错误：min值无效 - %v", err)
	}

	maxTime, err := strconv.ParseInt(maxTimeStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("休眠节点配置错误：max值无效 - %v", err)
	}

	// 验证时间范围
	if minTime < 0 {
		minTime = 0 // 最小值不能小于0
	}

	if maxTime < minTime {
		return nil, fmt.Errorf("休眠节点配置错误：max值(%d)小于min值(%d)", maxTime, minTime)
	}

	// 确定实际休眠时间
	var sleepTime int64
	if minTime == maxTime {
		sleepTime = minTime // 如果最小值等于最大值，则固定休眠
	} else {
		// 生成随机休眠时间
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		sleepTime = minTime + rng.Int63n(maxTime-minTime+1)
	}

	// 记录休眠信息
	execCtx.Log("info", "休眠节点开始执行，休眠时间: %d 毫秒", sleepTime)

	// 创建一个可取消的计时器
	timer := time.NewTimer(time.Duration(sleepTime) * time.Millisecond)
	defer timer.Stop()

	// 等待休眠时间结束或上下文取消
	select {
	case <-ctx.Done():
		// 上下文被取消，可能是流程被强制停止
		return nil, ctx.Err()
	case <-timer.C:
		// 休眠时间到，继续执行
	}

	// 记录执行结果
	result := map[string]any{
		"sleepTime": sleepTime,
		"unit":      "ms",
	}

	execCtx.Log("info", "休眠节点执行完成，实际休眠: %d 毫秒", sleepTime)
	return result, nil
}

// Validate 验证节点配置
func (h *SleepNodeHandler) Validate(node model.TaskNode) error {
	// 检查最小休眠时间
	minTimeStr, _ := node.Properties["min"].(string)
	// min是可选的，默认为0

	// 检查最大休眠时间(必须)
	maxTimeStr, ok := node.Properties["max"].(string)
	if !ok || maxTimeStr == "" {
		return errors.New("休眠节点配置错误：缺少或为空的max配置")
	}

	// 验证最小时间格式
	if minTimeStr != "" {
		minTime, err := strconv.ParseInt(minTimeStr, 10, 64)
		if err != nil {
			return fmt.Errorf("休眠节点配置错误：min值无效 - %v", err)
		}

		if minTime < 0 {
			return fmt.Errorf("休眠节点配置错误：min值不能小于0")
		}
	}

	// 验证最大时间格式
	maxTime, err := strconv.ParseInt(maxTimeStr, 10, 64)
	if err != nil {
		return fmt.Errorf("休眠节点配置错误：max值无效 - %v", err)
	}

	if maxTime <= 0 {
		return fmt.Errorf("休眠节点配置错误：max值必须大于0")
	}

	// 如果min存在，验证范围关系
	if minTimeStr != "" {
		minTime, _ := strconv.ParseInt(minTimeStr, 10, 64)
		if maxTime < minTime {
			return fmt.Errorf("休眠节点配置错误：max值必须大于或等于min值")
		}
	}

	return nil
}
