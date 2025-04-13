package tools

// AppClean 应用程序清理接口
// 定义了需要在应用程序关闭前执行清理工作的组件应实现的方法
type AppClean interface {
	Clean() // 清理方法，实现该方法的组件应在此进行资源释放和清理工作
}

// appCleans 已注册的所有需要执行清理的组件列表
var appCleans []AppClean

// RegisterAppClean 注册需要执行清理操作的组件
// 参数 appClean: 实现了AppClean接口的组件实例
func RegisterAppClean(appClean AppClean) {
	appCleans = append(appCleans, appClean)
}

// ExecuteClean 执行所有注册组件的清理操作
// 通常在应用程序退出前调用，确保资源被正确释放
func ExecuteClean() {
	for i := 0; i < len(appCleans); i++ {
		appCleans[i].Clean()
	}
}
