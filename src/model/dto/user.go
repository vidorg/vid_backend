package dto

import (
	"github.com/Aoi-hosizora/ahlib/xentity"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/vidorg/vid_backend/src/model/po"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("UserDto", "用户信息").
			Properties(
				goapidoc.NewProperty("uid", "integer#int32", true, "用户id"),
				goapidoc.NewProperty("username", "string", true, "用户名"),
				goapidoc.NewProperty("gender", "string", true, "性别").Enum("male", "female", "unknown"),
				goapidoc.NewProperty("profile", "string", true, "简介").AllowEmpty(true),
				goapidoc.NewProperty("avatar_url", "string", true, "头像"),
				goapidoc.NewProperty("birthday", "string#date", true, "生日"),
				goapidoc.NewProperty("role", "string", true, "角色"),
				goapidoc.NewProperty("register_time", "string#date-time", true, "注册时间"),
			),

		goapidoc.NewDefinition("LoginDto", "登录信息").
			Properties(
				goapidoc.NewProperty("user", "UserDto", true, "用户信息"),
				goapidoc.NewProperty("token", "string", true, "登录令牌"),
			),

		goapidoc.NewDefinition("UserExtraDto", "用户额外信息").
			Properties(
				goapidoc.NewProperty("subscribing_cnt", "integer#int32", true, "关注数量"),
				goapidoc.NewProperty("subscriber_cnt", "integer#int32", true, "粉丝数量"),
				goapidoc.NewProperty("video_cnt", "integer#int32", true, "视频数量"),
			),

		goapidoc.NewDefinition("UserDetailDto", "用户详细信息").
			Properties(
				goapidoc.NewProperty("user", "UserDto", true, "用户信息"),
				goapidoc.NewProperty("extra", "UserExtraDto", true, "用户额外信息"),
			),
	)
}

type UserDto struct {
	Uid          int32  `json:"uid"`
	Username     string `json:"username"`
	Gender       string `json:"gender"`
	Profile      string `json:"profile"`
	AvatarUrl    string `json:"avatar_url"` // TODO url
	Birthday     string `json:"birthday"`
	Role         string `json:"role"`
	RegisterTime string `json:"register_time"`
}

func BuildUserDto(user *po.User) *UserDto {
	return xentity.MustMap(user, &UserDto{}).(*UserDto)
}

func BuildUserDtos(users []*po.User) []*UserDto {
	return xentity.MustMapSlice(users, &UserDto{}).([]*UserDto)
}

type LoginDto struct {
	User  *UserDto `json:"user"`
	Token string   `json:"token"`
}

func BuildLoginDto(user *po.User, token string) *LoginDto {
	return &LoginDto{
		User:  BuildUserDto(user),
		Token: token,
	}
}

type UserExtraDto struct {
	SubscribingCount int32 `json:"subscribing_cnt"`
	SubscriberCount  int32 `json:"subscriber_cnt"`
	VideoCount       int32 `json:"video_cnt"`
}

func BuildUserExtraDto(subscribingCnt int32, subscriberCnt int32, videoCnt int32) *UserExtraDto {
	return &UserExtraDto{
		SubscribingCount: subscribingCnt,
		SubscriberCount:  subscriberCnt,
		VideoCount:       videoCnt,
	}
}

type UserDetailDto struct {
	User  *UserDto      `json:"user"`
	Extra *UserExtraDto `json:"extra"`
}

func BuildUserDetailDto(user *po.User, extra *UserExtraDto) *UserDetailDto {
	return &UserDetailDto{
		User:  BuildUserDto(user),
		Extra: extra,
	}
}
