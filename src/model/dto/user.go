package dto

import (
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
				goapidoc.NewProperty("followings", "integer#int32", true, "user following count"),
				goapidoc.NewProperty("followers", "integer#int32", true, "user follower count"),
				goapidoc.NewProperty("channels", "integer#int32", true, "user channel count"),
				goapidoc.NewProperty("subscribings", "integer#int32", true, "user subscribing count"),
				goapidoc.NewProperty("favorites", "integer#int32", true, "user favorite count"),
				goapidoc.NewProperty("is_following", "boolean", true, "is following this user"),
				goapidoc.NewProperty("is_followed", "boolean", true, "is followed by this user"),
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

type UserExtraDto struct {
	Followings   *int32 `json:"followings"`   // user followings count
	Followers    *int32 `json:"followers"`    // user followers count
	Channels     *int32 `json:"channels"`     // user channel count
	Subscribings *int32 `json:"subscribings"` // user subscribing count
	Favorites    *int32 `json:"favorites"`    // user favorite count
	IsFollowing  *bool  `json:"is_following"` // is following this user
	IsFollowed   *bool  `json:"is_followed"`  // is followed by this user
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
