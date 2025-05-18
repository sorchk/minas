package script

import (
	"context"
	"errors"
	"fmt"
	"server/dagflow/model"
	"server/dagflow/utils"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

// JavaScript脚本节点类型常量
const (
	// ...existing code...
	TypeJavaScript = "JavaScript" // JavaScript脚本节点
)

// JavaScriptHandler JavaScript脚本处理器
type JavaScriptHandler struct {
	registry *require.Registry // 模块加载器
}

// NewJavaScriptHandler 创建JavaScript脚本处理器
func NewJavaScriptHandler() *JavaScriptHandler {
	registry := require.NewRegistry()
	registry.RegisterNativeModule("console", console.RequireWithPrinter(newJSPrinter()))

	return &JavaScriptHandler{
		registry: registry,
	}
}

// GetType 获取处理器类型
func (h *JavaScriptHandler) GetType() string {
	return TypeJavaScript
}

// Handle 处理JavaScript脚本节点
func (h *JavaScriptHandler) Handle(ctx context.Context, node model.TaskNode, execCtx *model.ExecutionContext) (any, error) {
	// 获取脚本文本
	scriptText, ok := node.Properties["scriptText"].(string)
	if !ok || scriptText == "" {
		return nil, errors.New("JavaScript节点配置错误：缺少或为空的scriptText配置")
	}
	var data = execCtx.Data
	// 获取可选的变量配置
	scriptVars, err := utils.GetMap(node, "scriptVars", data)
	if err != nil {
		return nil, fmt.Errorf("获取脚本变量失败：%v", err)
	}
	// 创建JS运行时
	vm := goja.New()
	// 启用模块支持
	h.registry.Enable(vm)

	// 注入控制台
	console.Enable(vm)

	// 注入上下文数据
	dataObj := vm.NewObject()
	// 将执行上下文数据转换为JS对象
	for k, v := range execCtx.Data {
		if err := dataObj.Set(k, v); err != nil {
			return nil, fmt.Errorf("设置JS环境变量失败：%v", err)
		}
	}

	// 设置全局变量
	if err := vm.Set("ctx", dataObj); err != nil {
		return nil, fmt.Errorf("设置JS环境变量失败：%v", err)
	}

	// 针对配置的变量注入上下文
	for key := range scriptVars {
		value := scriptVars[key]
		if err := vm.Set(key, value); err != nil {
			return nil, fmt.Errorf("设置JS自定义环境变量失败：%v", err)
		}
	}

	// 注入日志函数
	if err := vm.Set("log", func(level string, message string, args ...interface{}) {
		logLevel := "info" // 默认日志级别
		if level != "" {
			logLevel = level
		}
		execCtx.Log(logLevel, "[JS] %s", fmt.Sprintf(message, args...))
	}); err != nil {
		return nil, fmt.Errorf("设置JS日志函数失败：%v", err)
	}

	// 设置执行超时(默认5秒)
	var timeoutMS int64 = 5000
	if timeout, ok := node.Properties["timeout"].(int64); ok && timeout > 0 {
		timeoutMS = timeout
	}

	// 创建带超时的上下文
	execTimeout := time.Duration(timeoutMS) * time.Millisecond
	timeoutCtx, cancel := context.WithTimeout(ctx, execTimeout)
	defer cancel()

	// 在后台监控超时
	done := make(chan struct{})
	var execErr error
	var result goja.Value

	go func() {
		defer close(done)
		// 执行脚本
		result, execErr = vm.RunString(scriptText)
	}()

	// 等待执行完成或超时
	select {
	case <-timeoutCtx.Done():
		// 如果是超时导致的取消
		if timeoutCtx.Err() == context.DeadlineExceeded {
			return nil, errors.New("JavaScript脚本执行超时")
		}
		return nil, timeoutCtx.Err()
	case <-done:
		// 执行完成
	}

	// 检查执行错误
	if execErr != nil {
		return nil, fmt.Errorf("JavaScript脚本执行错误: %v", execErr)
	}

	// 如果没有返回值或返回undefined
	if result == nil || goja.IsUndefined(result) || goja.IsNull(result) {
		return nil, nil
	}

	// 将结果导出为Go值
	goValue := result.Export()
	fmt.Printf("JavaScript脚本执行结果: %v 类型：%T\n", goValue, goValue)

	return goValue, nil
}

// Validate 验证节点配置
func (h *JavaScriptHandler) Validate(node model.TaskNode) error {
	// 检查脚本文本
	scriptText, ok := node.Properties["scriptText"].(string)
	if !ok || scriptText == "" {
		return errors.New("JavaScript节点配置错误：缺少或为空的scriptText配置")
	}

	// 检查脚本是否可以编译
	compilable, _ := node.Properties["compilable"].(bool)
	if compilable {
		// 预编译脚本检查语法错误
		_, err := goja.Compile("", scriptText, false)
		if err != nil {
			return fmt.Errorf("JavaScript脚本编译错误: %v", err)
		}
	}

	return nil
}

// JSConsole 自定义JS控制台实现
type JSConsole struct {
	logger model.LoggerInterface
}

// 创建新的JS打印器
func newJSPrinter() console.Printer {
	return &JSPrinter{}
}

// JSPrinter 实现Printer接口
type JSPrinter struct{}

// Error implements console.Printer.
func (p *JSPrinter) Error(string) {
	panic("unimplemented")
}

// Log implements console.Printer.
func (p *JSPrinter) Log(string) {
	panic("unimplemented")
}

// Warn implements console.Printer.
func (p *JSPrinter) Warn(string) {
	panic("unimplemented")
}

// Print 打印消息
func (p *JSPrinter) Print(msg string) {
	fmt.Print(msg)
}

// Printf 格式化打印消息
func (p *JSPrinter) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
