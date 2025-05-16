// Package projectdir 提供NAS相关的API控制器
package projectdir

import (
	"server/core/app/request"
	"server/core/app/response"
	"server/core/app/webapi"
	"server/service/basic"

	"github.com/gin-gonic/gin"
)

// ProjectDirApp 项目目录接口控制器
// 提供项目目录的CRUD操作和树形结构查询
type ProjectDirApp struct {
	webapi.BaseApp[basic.ProjectDir] // 继承通用接口实现
}

// AddRoutes 添加项目目录相关路由
// 设置项目目录API的所有路由和对应的处理函数
// 参数:
//   - group: 路由组，路由将被添加到此组中

// 参数 parentGroup: 父路由组

func AddRoutes(parentGroup *gin.RouterGroup) {

	// 创建项目目录路由组
	group := parentGroup.Group("/project-dir")
	// 创建项目目录应用
	app := ProjectDirApp{}

	// 添加基础CRUD路由
	webapi.AddBaseRoutes(group, &app)

	// 设置更新操作时允许更新的字段
	app.UpdateFields = []string{"name", "parent_id", "remark"}
	// 获取目录树结构
	group.GET("/tree", app.GetTree)
	// 获取指定父目录下的子目录
	group.GET("/children", app.GetChildren)
}

// 参数 ctx: 请求上下文
func (app *ProjectDirApp) List(ctx *gin.Context) {
	var entity basic.ProjectDir
	query := request.GetPageQuery(ctx)

	// 获取父目录ID
	parentID := ctx.Query("parent_id")

	if parentID != "" {
		query.AddFilter(request.NewEqualFilter("parent_id", parentID))
	}

	// 获取目录名称
	name := ctx.Query("name")
	if name != "" {
		query.AddFilter(request.NewLikeFilter("name", name))
	}
	list, count, err := entity.List(query)
	if err == nil {
		response.List(ctx, "", count, list)
	} else {
		response.NoContent(ctx, "无数据！")
	}
}

// GetTree 获取目录树结构
func (a *ProjectDirApp) GetTree(ctx *gin.Context) {

	var entity basic.ProjectDir
	// 调用服务方法获取目录树
	tree, err := entity.GetTree()
	if err != nil {
		// 出错时返回错误信息
		response.Error(ctx, err)
		return
	}

	// 返回目录树数据
	response.Data(ctx, "", tree)
}

// GetChildren 获取子目录列表
func (a *ProjectDirApp) GetChildren(ctx *gin.Context) {
	// 获取父目录ID
	parentID := uint(ctx.GetInt("parent_id"))

	var entity basic.ProjectDir
	// 调用服务方法获取子目录
	children, err := entity.GetByParent(parentID)
	if err != nil {
		// 出错时返回错误信息
		response.Error(ctx, err)
		return
	}

	// 返回子目录数据
	response.Data(ctx, "", children)
}
