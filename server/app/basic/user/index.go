package user

import (
	"server/core/app/request"
	"server/core/app/response"
	"server/core/app/webapi"
	"server/middleware"
	"server/service/basic"
	"server/utils/mfa"

	"github.com/gin-gonic/gin"
)

// UserApp 用户应用结构体
// 继承了基础应用结构，用于处理用户管理相关的API请求
type UserApp struct {
	webapi.BaseApp[basic.User]
}

// AddRoutes 添加用户管理相关的路由配置
// 参数 parentGroup: 父路由组
func AddRoutes(parentGroup *gin.RouterGroup) {
	// 创建用户管理路由组
	group := parentGroup.Group("/user")
	app := UserApp{}
	// 设置用户更新操作允许修改的字段
	app.UpdateFields = []string{"id", "account", "avatar", "birthday", "email", "name", "phone", "sex", "sn", "remark", "role_keys"}
	// 添加基础CRUD路由
	webapi.AddBaseRoutes(group, &app)

	// 注释掉的保存方法
	// group.POST("/save", app.Save)

	// 添加获取用户个人资料路由
	group.GET("/profile", app.GetProfile)
	// 添加修改密码路由
	group.POST("/modify-password", app.ModifyPassword)
	// 添加修改个人资料路由
	group.POST("/modify-profile", app.ModifyProfile)
	// 添加启用多因素认证路由
	group.POST("/mfa/enable", app.EnableMFA)
	// 添加禁用多因素认证路由
	group.POST("/mfa/disable", app.DisableMFA)
}

// List 获取用户列表
// 参数 ctx: 上下文对象
func (app *UserApp) List(ctx *gin.Context) {
	query := request.GetPageQuery(ctx)
	var entity basic.User
	list, count, err := entity.List(query)
	if err == nil {
		response.List(ctx, "", count, list)
	} else {
		response.NoContent(ctx, "无数据！")
	}
}

// GetProfile 获取用户个人资料
// 参数 ctx: 上下文对象
func (app *UserApp) GetProfile(ctx *gin.Context) {
	user := basic.User{}
	user, err := user.Load(request.GetUserID(ctx))
	if err == nil {
		response.Data(ctx, "", user)
	} else {
		response.NotFound(ctx, "用户不存在！")
	}
}

// ModifyProfile 修改用户个人资料
// 参数 ctx: 上下文对象
func (app *UserApp) ModifyProfile(ctx *gin.Context) {
	entity := basic.User{}
	if err := ctx.BindJSON(&entity); err != nil {
		response.BadRequest(ctx, "参数错误！")
		return
	}
	entity.ID = request.GetUserID(ctx)
	entity.SetOperator(&entity, ctx)
	err := entity.Update(&entity, app.UpdateFields...)
	if err == nil {
		response.Data(ctx, "", entity)
	} else {
		response.Error(ctx, err)
		return
	}
}

// ModifyPasswordForm 修改密码表单结构体
type ModifyPasswordForm struct {
	Password    string `json:"password"`
	NewPassword string `json:"newpassword"`
}

// ModifyPassword 修改用户密码
// 参数 ctx: 上下文对象
func (app *UserApp) ModifyPassword(ctx *gin.Context) {
	params := ModifyPasswordForm{}
	if err := ctx.BindJSON(&params); err != nil {
		response.BadRequest(ctx, "参数错误！")
		return
	}
	user := basic.User{}
	user, err := user.Load(request.GetUserID(ctx))

	if err != nil {
		response.Error(ctx, err)
		return
	}
	if middleware.VerifyPass(user.Password, params.Password) {
		err := user.ModifyPass(user.ID, params.NewPassword)
		if err == nil {
			response.Success(ctx, "")
		} else {
			response.Error(ctx, err)
			return
		}
	} else {
		response.Error(ctx, err)
		return
	}
}

// MfaForm 多因素认证表单结构体
type MfaForm struct {
	MfaCode      string `json:"mfa_code"`
	MfaVerifCode string `json:"mfa_verify_code"`
}

// EnableMFA 启用多因素认证
// 参数 ctx: 上下文对象
func (app *UserApp) EnableMFA(ctx *gin.Context) {
	mfaForm := MfaForm{}
	if err := ctx.BindJSON(&mfaForm); err != nil {
		response.BadRequest(ctx, "开启MFA验证失败，参数错误！")
		return
	}
	if !mfa.ValidCode(mfaForm.MfaVerifCode, mfaForm.MfaCode, 30) {
		response.Message(ctx, 4000, "开启MFA验证失败，验证码错误！")
		return
	}
	user := basic.User{}
	user.ID = request.GetUserID(ctx)
	user.MFAEnable = 1
	user.MFACode = mfaForm.MfaCode
	user.SetOperator(&user, ctx)
	err := user.Update(&user, "mfa_enable", "mfa_code")
	if err != nil {
		response.Message(ctx, 4000, "开启MFA验证失败！")
	} else {
		response.Data(ctx, "开启MFA验证成功！", user)
	}
}

// DisableMFA 禁用多因素认证
// 参数 ctx: 上下文对象
func (app *UserApp) DisableMFA(ctx *gin.Context) {
	user := basic.User{}
	user.ID = request.GetUserID(ctx)
	user.MFAEnable = 0
	err := user.Update(&user, "mfa_enable")
	if err != nil {
		response.Message(ctx, 4000, "关闭MFA验证失败！请重试")
	} else {
		response.Data(ctx, "关闭MFA验证成功！", err)
	}
}
