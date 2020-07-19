package param

import (
	"github.com/Aoi-hosizora/ahlib/xentity"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/vidorg/vid_backend/src/model/po"
)

// https://godoc.org/github.com/go-playground/validator#hdr-Baked_In_Validators_and_Tags

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("UserParam", "用户请求参数").
			WithProperties(
				goapidoc.NewProperty("username", "string", true, "用户名，长度在 [5, 30] 之间"),
				goapidoc.NewProperty("profile", "string", true, "用户简介，长度在 [0, 255] 之间").WithAllowEmptyValue(true),
				goapidoc.NewProperty("gender", "string", true, "性别").WithEnum("male", "female", "unknown"),
				goapidoc.NewProperty("birthday", "string#date", true, "生日"),
				goapidoc.NewProperty("phone_number", "string", true, "手机号码，长度为 11，仅限中国大陆手机号码"),
				goapidoc.NewProperty("avatar_url", "string", true, "头像"),
			),
	)
}

type UserParam struct {
	Username    string  `form:"username"     json:"username"     binding:"required,min=5,max=30,name"`
	Profile     *string `form:"profile"      json:"profile"      binding:"required,min=0,max=255"`
	Gender      string  `form:"gender"       json:"gender"       binding:"required"`
	Birthday    string  `form:"birthday"     json:"birthday"     binding:"required,date"`
	PhoneNumber string  `form:"phone_number" json:"phone_number" binding:"required,phone"`
	AvatarUrl   string  `form:"avatar_url"   json:"avatar_url"   binding:"required,url"` // TODO url
}

func MapUserParam(param *UserParam, user *po.User) {
	xentity.MustMapProp(param, user)
}
