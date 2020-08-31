package param

import (
	"github.com/Aoi-hosizora/goapidoc"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("LoginParam", "登录请求参数").
			Properties(
				goapidoc.NewProperty("username", "string", true, "用户名"),
				goapidoc.NewProperty("password", "string", true, "密码"),
			),

		goapidoc.NewDefinition("RegisterParam", "注册请求参数").
			Properties(
				goapidoc.NewProperty("username", "string", true, "用户名，长度在 [5, 30] 之间"),
				goapidoc.NewProperty("password", "string", true, "密码，长度在 [8, 30] 之间"),
			),

		goapidoc.NewDefinition("PasswordParam", "修改密码请求参数").
			Properties(
				goapidoc.NewProperty("password", "string", true, "密码，长度在 [8, 30] 之间"),
			),
	)
}

type RegisterParam struct {
	Email    string `json:"email"    form:"email"    binding:"required,min=5,max=30,email"` // register email
	Password string `json:"password" form:"password" binding:"required,min=8,max=30,pwd"`   // register password
}

type LoginParam struct {
	Parameter string `json:"parameter" form:"parameter" binding:"required"` // login parameter
	Password  string `json:"password"  form:"password"  binding:"required"` // login password
}

type UpdatePasswordParam struct {
	Old string `json:"old" form:"old" binding:"required,pwd"` // old password
	New string `json:"new" form:"new" binding:"required,pwd"` // new password
}
