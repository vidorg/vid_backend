package dto

import (
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/Aoi-hosizora/ahlib/xtime"
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
	Uid          uint64        `json:"uid"`           // user uid
	Username     string        `json:"username"`      // username
	Email        string        `json:"email"`         // user email
	Nickname     string        `json:"nickname"`      // user nickname
	Gender       string        `json:"gender"`        // user gender
	Profile      string        `json:"profile"`       // user profile, allowempty
	Avatar       string        `json:"avatar"`        // user avatar url
	Birthday     string        `json:"birthday"`      // user birthday
	Role         string        `json:"role"`          // user role
	RegisterTime string        `json:"register_time"` // user register time
	Extra        *UserExtraDto `json:"extra"`         // user extra information
}

func BuildUserDto(user *po.User) *UserDto {
	return &UserDto{
		Uid:          user.Uid,
		Username:     user.Username,
		Email:        user.Email,
		Nickname:     user.Nickname,
		Gender:       user.Gender.String(),
		Profile:      user.Profile,
		Avatar:       user.Avatar,
		Birthday:     user.Birthday.String(),
		Role:         user.Role,
		RegisterTime: xtime.NewJsonDateTime(user.CreatedAt).String(),
		Extra:        nil,
	}
}

func BuildUserDtos(users []*po.User) []*UserDto {
	out := make([]*UserDto, len(users))
	for idx, user := range users {
		out[idx] = BuildUserDto(user)
	}
	return out
}

type LoginDto struct {
	User  *UserDto `json:"user"`  // authorized user
	Token string   `json:"token"` // access token
}

func BuildLoginDto(user *po.User, token string) *LoginDto {
	return &LoginDto{
		User:  BuildUserDto(user),
		Token: token,
	}
}

type UserExtraDto struct {
	Subscribings *int32 `json:"subscribings"`
	Subscribers  *int32 `json:"subscribers"`
	Videos       *int32 `json:"videos"`
}

func BuildUserExtraDto(dto *UserExtraDto) *UserExtraDto {
	if dto.Subscribings == nil && dto.Subscribers == nil && dto.Videos == nil {
		return nil
	}
	return dto
}

func BuildUserPropertyMapper() xproperty.PropertyDict {
	return xproperty.PropertyDict{
		"uid":           xproperty.NewValue(false, "uid"),
		"username":      xproperty.NewValue(false, "username"),
		"email":         xproperty.NewValue(false, "email"),
		"nickname":      xproperty.NewValue(false, "nickname"),
		"gender":        xproperty.NewValue(false, "gender"),
		"profile":       xproperty.NewValue(false, "profile"),
		"avatar":        xproperty.NewValue(false, "avatar"),
		"birthday":      xproperty.NewValue(false, "birthday"),
		"age":           xproperty.NewValue(true, "birthday"),
		"role":          xproperty.NewValue(false, "role"),
		"register_time": xproperty.NewValue(false, "created_at"),
	}
}
