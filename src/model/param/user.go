package param

import (
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/vidorg/vid_backend/src/model/constant"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("UpdateUserParam", "用户请求参数").
			Properties(
				goapidoc.NewProperty("username", "string", true, "用户名，长度在 [5, 30] 之间"),
				goapidoc.NewProperty("profile", "string", true, "用户简介，长度在 [0, 255] 之间").AllowEmpty(true),
				goapidoc.NewProperty("gender", "string", true, "性别").Enum("male", "female", "unknown"),
				goapidoc.NewProperty("birthday", "string#date", true, "生日"),
				goapidoc.NewProperty("phone_number", "string", true, "手机号码，长度为 11，仅限中国大陆手机号码"),
				goapidoc.NewProperty("avatar_url", "string", true, "头像"),
			),
	)
}

type UpdateUserParam struct {
	Username string  `json:"username"     form:"username"     binding:"required,min=5,max=30,name"` // username
	Nickname string  `json:"nickname"     form:"nickname"     binding:"required,min=5,max=30,name"` // user nickname
	Gender   int8    `json:"gender"       form:"gender"       binding:"required"`                   // user gender (0X, 1M, 2F)
	Profile  *string `json:"profile"      form:"profile"      binding:"required,min=0,max=255"`     // user profile, allowempty
	Birthday string  `json:"birthday"     form:"birthday"     binding:"required,date"`              // user birthday
	Phone    string  `json:"phone"        form:"phone"        binding:"required,phone"`             // user phone number
	Avatar   string  `json:"avatar"       form:"avatar"       binding:"required,url"`               // user avatar
}

func (u *UpdateUserParam) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"username":     u.Username,
		"nickname":     u.Nickname,
		"gender":       constant.ParseGender(u.Gender),
		"profile":      *u.Profile,
		"phone_number": u.Phone,
		"avatar":       u.Avatar,
	}

	d, err := xtime.ParseRFC3339Date(u.Birthday)
	if err == nil {
		m["birthday"] = d
	}

	return m
}
