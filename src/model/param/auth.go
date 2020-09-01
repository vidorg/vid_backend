package param

import (
	"github.com/Aoi-hosizora/goapidoc"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("RegisterParam", "register parameter").
			Properties(
				goapidoc.NewProperty("email", "string", true, "register email").Example("aaa@bbb.ccc"),
				goapidoc.NewProperty("password", "string", true, "register password"),
			),

		goapidoc.NewDefinition("LoginParam", "login parameter").
			Properties(
				goapidoc.NewProperty("parameter", "string", true, "login parameter, support uid | username | email"),
				goapidoc.NewProperty("password", "string", true, "login password"),
			),

		goapidoc.NewDefinition("UpdatePasswordParam", "update password parameter").
			Properties(
				goapidoc.NewProperty("old", "string", true, "old password"),
				goapidoc.NewProperty("new", "string", true, "new password"),
			),
	)
}

type RegisterParam struct {
	Email    string `json:"email"    form:"email"    binding:"required,l_email,email"` // register email
	Password string `json:"password" form:"password" binding:"required,l_pwd,r_pwd"`   // register password
}

type LoginParam struct {
	Parameter string `json:"parameter" form:"parameter" binding:"required"` // login parameter
	Password  string `json:"password"  form:"password"  binding:"required"` // login password
}

type UpdatePasswordParam struct {
	Old string `json:"old" form:"old" binding:"required,l_pwd,r_pwd"` // old password
	New string `json:"new" form:"new" binding:"required,l_pwd,r_pwd"` // new password
}
