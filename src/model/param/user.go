package param

import (
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/vidorg/vid_backend/src/model/constant"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("UpdateUserParam", "update user parameter").
			Properties(
				goapidoc.NewProperty("username", "string", true, "username"),
				goapidoc.NewProperty("nickname", "string", true, "user nickname"),
				goapidoc.NewProperty("gender", "integer#int32", true, "user gender, 0X | 1M | 2F").Enum(0, 1, 2),
				goapidoc.NewProperty("profile", "string", true, "user profile").AllowEmpty(true),
				goapidoc.NewProperty("birthday", "string#date", true, "user birthday").Example("2000-01-01"),
				goapidoc.NewProperty("avatar", "string", true, "user avatar").Example("https://aaa.bbb.ccc"),
			),
	)
}

type UpdateUserParam struct {
	Username string  `json:"username"     form:"username"     binding:"required,l_name,r_name"` // username
	Nickname string  `json:"nickname"     form:"nickname"     binding:"required,l_name"`        // user nickname
	Gender   int8    `json:"gender"       form:"gender"       binding:"required,o_gender"`      // user gender (0X, 1M, 2F)
	Profile  *string `json:"profile"      form:"profile"      binding:"required,l_profile"`     // user profile, allowempty
	Birthday string  `json:"birthday"     form:"birthday"     binding:"required,date"`          // user birthday
	Avatar   string  `json:"avatar"       form:"avatar"       binding:"required,url"`           // user avatar
}

func (u *UpdateUserParam) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"username": u.Username,
		"nickname": u.Nickname,
		"gender":   constant.ParseGender(u.Gender),
		"profile":  *u.Profile,
		"avatar":   u.Avatar,
	}

	d, err := xtime.ParseRFC3339Date(u.Birthday)
	if err == nil {
		m["birthday"] = d
	}

	return m
}
