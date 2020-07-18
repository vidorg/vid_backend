package dto

import (
	"github.com/Aoi-hosizora/ahlib/xentity"
	"github.com/vidorg/vid_backend/src/model/po"
)

// @Model         _LoginDto
// @Description   登录信息
// @Property      user   object(#_UserDto) true "用户信息"
// @Property      token  string            true "登录令牌"

// @Model         _UserDto
// @Description   用户信息
// @Property      uid           integer                          true "用户id"
// @Property      username      string                           true "用户名"
// @Property      sex           string(enum:male,female,unknown) true "用户性别"
// @Property      profile       string                           true "用户简介"
// @Property      avatar_url    string                           true "用户头像"
// @Property      birthday      string(format:date)              true "用户生日"
// @Property      role          string                           true "用户角色"
// @Property      register_time string(format:datetime)          true "用户注册时间"
type UserDto struct {
	Uid          int32  `json:"uid"`
	Username     string `json:"username"`
	Sex          string `json:"sex"`
	Profile      string `json:"profile"`
	AvatarUrl    string `json:"avatar_url"`
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

// @Model         _UserAndExtraDto
// @Description   用户与数量信息
// @Property      user  object(#_UserDto)      true "用户信息"
// @Property      extra object(#_UserExtraDto) true "用户额外信息"

// @Model         _UserExtraDto
// @Description   用户额外信息
// @Property      subscribing_cnt integer true "用户关注数量"
// @Property      subscriber_cnt  integer true "用户粉丝数量"
// @Property      video_cnt       integer true "用户视频数量"
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
