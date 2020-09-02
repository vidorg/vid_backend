package dto

import (
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/vidorg/vid_backend/src/model/po"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("UserDto", "user response").
			Properties(
				goapidoc.NewProperty("uid", "integer#int64", true, "user id"),
				goapidoc.NewProperty("username", "string", true, "username"),
				goapidoc.NewProperty("email", "string", true, "user email"),
				goapidoc.NewProperty("nickname", "string", true, "user nickname"),
				goapidoc.NewProperty("gender", "string", true, "user gender").Enum("secret", "male", "female"),
				goapidoc.NewProperty("profile", "string", true, "user profile").AllowEmpty(true),
				goapidoc.NewProperty("avatar", "string", true, "user avatar"),
				goapidoc.NewProperty("birthday", "string#date", true, "user birthday"),
				goapidoc.NewProperty("role", "string", true, "user role"),
				goapidoc.NewProperty("state", "string", true, "user state"),
				goapidoc.NewProperty("register_time", "string#date-time", true, "user register time"),
				goapidoc.NewProperty("extra", "UserExtraDto", true, "user extra information"),
			),

		goapidoc.NewDefinition("UserExtraDto", "user extra response").
			Properties(
				goapidoc.NewProperty("subscribings", "integer#int32", true, "user subscribing count"),
				goapidoc.NewProperty("subscribers", "integer#int32", true, "user subscriber count"),
				goapidoc.NewProperty("is_subscribing", "boolean", true, "authorized user is subscribing this user"),
				goapidoc.NewProperty("is_subscribed", "boolean", true, "authorized user is subscribed by this user"),
				goapidoc.NewProperty("is_blocking", "boolean", true, "authorized user is blocking this user"),
				goapidoc.NewProperty("videos", "integer#int32", true, "user video count"),
				goapidoc.NewProperty("favorites", "integer#int32", true, "user favorite count"),
			),

		goapidoc.NewDefinition("LoginDto", "login response").
			Properties(
				goapidoc.NewProperty("user", "UserDto", true, "authorized user"),
				goapidoc.NewProperty("token", "string", true, "access token"),
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
	State        string        `json:"state"`         // user state
	RegisterTime string        `json:"register_time"` // user register time
	Extra        *UserExtraDto `json:"extra"`         // user extra information
}

func BuildUserDto(user *po.User) *UserDto {
	if user == nil {
		return nil
	}
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
		State:        user.State.String(),
		RegisterTime: xtime.NewJsonDateTime(user.CreatedAt).String(),
		Extra:        &UserExtraDto{},
	}
}

func BuildUserDtos(users []*po.User) []*UserDto {
	out := make([]*UserDto, len(users))
	for idx, user := range users {
		out[idx] = BuildUserDto(user)
	}
	return out
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
		"state":         xproperty.NewValue(false, "state"),
		"register_time": xproperty.NewValue(false, "created_at"),
	}
}

type UserExtraDto struct {
	Subscribings  *int32 `json:"subscribings"`   // user subscribing count
	Subscribers   *int32 `json:"subscribers"`    // user subscriber count
	IsSubscribing *bool  `json:"is_subscribing"` // authorized user is subscribing this user
	IsSubscribed  *bool  `json:"is_subscribed"`  // authorized user is subscribed by this user
	IsBlocking    *bool  `json:"is_blocking"`    // authorized user is blocking this user
	Videos        *int32 `json:"videos"`         // user video count
	Favorites     *int32 `json:"favorites"`      // user favorite count
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
