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
				goapidoc.NewProperty("birthday", "string#date", true, "user birthday"),
				goapidoc.NewProperty("phone", "string", true, "user phone number"),
				goapidoc.NewProperty("avatar", "string", true, "user avatar"),
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
