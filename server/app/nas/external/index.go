package external

import (
	"log"
	"net/http"
	"os"
	"server/core/app/request"
	"server/core/app/response"
	"server/core/app/webapi"
	"server/service/nas"
	"server/utils/rclone"
	"strings"

	"github.com/gin-gonic/gin"
)

// ExternalNasApp 外部网络存储应用结构体
// 继承了基础应用结构，用于处理外部NAS的API请求
type ExternalNasApp struct {
	webapi.BaseApp[nas.ExternalNas]
}

// RCloneProxyApiBase RClone API代理的基础路径
const RCloneProxyApiBase = "/rclone/api/"

// AddRoutes 添加路由配置
// 设置外部NAS相关的API路由
// 参数 parentGroup: 父路由组
func AddRoutes(parentGroup *gin.RouterGroup) {
	// 创建路由组
	group := parentGroup.Group("/external")
	app := ExternalNasApp{}
	// 添加基础CRUD路由
	webapi.AddBaseRoutes(group, &app)
	// 设置更新操作时允许更新的字段
	app.UpdateFields = []string{"name", "type", "rc_name", "is_adv", "config", "remark"}

	// 添加列出目录内容的路由
	group.GET("/list-dir", app.ListDir)
	// 添加RClone API代理路由，支持任意URI
	group.POST(RCloneProxyApiBase+"*uri", app.RCloneApi)
}

// RCloneApi RClone API代理方法
// 参数 ctx: 请求上下文
func (app ExternalNasApp) RCloneApi(ctx *gin.Context) {
	uri := ctx.Param("uri")
	code, body, err := rclone.ApiProxy(uri, ctx.Request.Body)
	if err != nil {
		response.Message(ctx, code, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, response.ResData{Code: code, Msg: "", Data: body})
	ctx.Abort()
}

// List 列出外部NAS
// 参数 ctx: 请求上下文
func (app ExternalNasApp) List(ctx *gin.Context) {
	log.Println("List externalNas")
	externalNas := nas.ExternalNas{}
	name := ctx.Query("name")
	rc_name := ctx.Query("rc_name")
	query := request.GetPageQuery(ctx)
	query.AddFilter(request.NewLikeFilter("name", name), request.NewLikeFilter("rc_name", rc_name))
	list, count, err := externalNas.List(query)
	if err == nil {
		response.List(ctx, "", count, list)
	} else {
		response.NoContent(ctx, "无数据！")
	}
}

// ListDir 列出目录内容
// 参数 ctx: 请求上下文
func (app ExternalNasApp) ListDir(ctx *gin.Context) {
	PthSep := string(os.PathSeparator)
	path := ctx.Query("path")
	if !strings.HasSuffix(path, PthSep) {
		path += PthSep
	}
	nas := ctx.Query("nas")
	dir := path
	if nas != "0" && nas != "" {
		dir = nas + ":" + path
	}
	fileList, err := rclone.ListDir(dir, strings.Split("*", "\n"), []string{}, 0, 0)
	if err != nil {
		response.Error(ctx, err)
	} else {
		response.Data(ctx, "", fileList)
	}
}
