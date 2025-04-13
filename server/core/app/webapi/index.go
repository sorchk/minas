// package webapi 提供Web API的基础接口和实现，用于处理HTTP请求
package webapi

import (
	"reflect"
	"server/core/app/request"
	"server/core/app/response"
	"server/core/db"

	"github.com/gin-gonic/gin"
)

// BaseAppApi 定义了基础应用API的接口，包含常用的CRUD操作
type BaseAppApi interface {
	// List 获取实体列表
	List(ctx *gin.Context)
	// Load 根据ID加载单个实体
	Load(ctx *gin.Context)
	// Save 保存实体（创建或更新）
	Save(ctx *gin.Context)
	// Enable 启用实体
	Enable(ctx *gin.Context)
	// Disable 禁用实体
	Disable(ctx *gin.Context)
	// Delete 删除实体
	Delete(ctx *gin.Context)
}

// BaseApp 是BaseAppApi的通用实现，使用泛型支持不同类型的数据访问对象
type BaseApp[T db.BaseDao[T]] struct {
	BaseAppApi
	// UpdateFields 指定更新时允许修改的字段列表
	UpdateFields []string
}

// AddBaseRoutes 为给定的路由组添加标准CRUD路由
func AddBaseRoutes(group *gin.RouterGroup, app BaseAppApi) {
	group.GET("/list", app.List)
	group.GET("/load/:id", app.Load)
	group.POST("/save", app.Save)
	group.POST("/enable/:id", app.Enable)
	group.POST("/disable/:id", app.Disable)
	group.POST("/delete/:id", app.Delete)
}

// List 处理列表请求，返回分页查询结果
func (BaseApp[T]) List(ctx *gin.Context) {
	// 获取分页查询参数
	query := request.GetPageQuery(ctx)
	var entity T
	// 调用实体的列表方法获取数据
	list, count, err := entity.List(query)
	if err == nil {
		// 成功时返回列表数据和总数
		response.List(ctx, "", count, list)
	} else {
		// 出错时返回无内容响应
		response.NoContent(ctx, "无数据！")
	}
}

// callReflectMethod 使用反射调用服务的SetOperator方法
func (app BaseApp[T]) callReflectMethod(service *T, ctx *gin.Context) []reflect.Value {
	getFunc := reflect.ValueOf(service).MethodByName("SetOperator")
	// 构造参数列表
	args := []reflect.Value{
		reflect.ValueOf(ctx), //sql value
	}
	// 使用 Call 方法调用函数，并获取返回值
	ret := getFunc.Call(args)
	return ret
}

// Save 处理保存请求，创建或更新实体
func (app BaseApp[T]) Save(ctx *gin.Context) {
	var entity T
	// 从请求体绑定JSON数据到实体
	if err := ctx.BindJSON(&entity); err != nil {
		response.BadRequest(ctx, "参数错误！")
		return
	}
	// 获取当前用户ID
	uid := request.GetUserID(ctx)
	// 使用反射设置操作者ID
	setOperatorFunc := reflect.ValueOf(&entity).MethodByName("SetOperatorUID")
	if setOperatorFunc.IsValid() {
		args := []reflect.Value{reflect.ValueOf(uid)}
		setOperatorFunc.Call(args)
	}
	// 保存实体，并指定允许更新的字段
	err := entity.Save(&entity, app.UpdateFields...)
	if err == nil {
		// 成功时返回保存后的实体
		response.Data(ctx, "", entity)
	} else {
		// 出错时返回错误信息
		response.Error(ctx, err)
	}
}

// Load 根据ID加载单个实体
func (BaseApp[T]) Load(ctx *gin.Context) {
	var entity T
	// 调用实体的Load方法获取数据
	entity, err := entity.Load(ctx.Param("id"))
	if err == nil {
		// 成功时返回实体数据
		response.Data(ctx, "", entity)
		return
	} else {
		// 出错时返回错误信息
		response.Error(ctx, err)
		return
	}
}

// Delete 根据ID删除实体
func (BaseApp[T]) Delete(ctx *gin.Context) {
	var entity T
	// 调用实体的Delete方法
	err := entity.Delete(ctx.Param("id"))
	if err == nil {
		// 成功时返回成功消息
		response.Success(ctx, "删除成功！")
		return
	} else {
		// 出错时返回错误信息
		response.Error(ctx, err)
		return
	}
}

// Enable 启用指定ID的实体
func (BaseApp[T]) Enable(ctx *gin.Context) {
	var entity T
	// 调用实体的Enable方法
	err := entity.Enable(ctx.Param("id"))
	if err == nil {
		// 成功时返回成功消息
		response.Success(ctx, "")
		return
	} else {
		// 出错时返回错误信息
		response.Error(ctx, err)
		return
	}
}

// Disable 禁用指定ID的实体
func (BaseApp[T]) Disable(ctx *gin.Context) {
	var entity T
	// 调用实体的Disable方法
	err := entity.Disable(ctx.Param("id"))
	if err == nil {
		// 成功时返回成功消息
		response.Success(ctx, "")
		return
	} else {
		// 出错时返回错误信息
		response.Error(ctx, err)
		return
	}
}
