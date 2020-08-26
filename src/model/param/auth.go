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

type LoginParam struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type RegisterParam struct {
	Username string `form:"username" json:"username" binding:"required,min=5,max=30,name"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=30,pwd"`
}

type PassParam struct {
	Password string `form:"password" json:"password" binding:"required,min=8,max=30,pwd"`
}
