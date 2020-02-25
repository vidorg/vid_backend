package dto

import (
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/vidorg/vid_backend/src/common/enum"
	"github.com/vidorg/vid_backend/src/model/po"
)

// @Model         _LoginDto
// @Description   登录信息
// @Property      user   object(#_UserDto) true "用户信息"
// @Property      token  string            true "登录令牌"
// @Property      expire integer           true "登录有效期，单位为秒"

// @Model         _UserDto
// @Description   用户信息
// @Property      uid           integer                            true "用户id"
// @Property      username      string                             true "用户名"
// @Property      sex           string(enum:male,female,unknown)   true "用户性别"
// @Property      profile       string                             true "用户简介"
// @Property      avatar_url    string                             true "用户头像"
// @Property      birth_time    string(format:2000-01-01)          true "用户生日"
// @Property      authority     string(enum:normal,admin)          true "用户权限"
// @Property      phone_number  string                             true "用户手机号码，部分接口可见"
// @Property      register_time string(format:2000-01-01 00:00:00) true "用户注册时间"
type UserDto struct {
	Uid          int32  `json:"uid"`
	Username     string `json:"username"`
	Sex          string `json:"sex"`
	Profile      string `json:"profile"`
	AvatarUrl    string `json:"avatar_url"`
	BirthTime    string `json:"birth_time"`
	Authority    string `json:"authority"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	RegisterTime string `json:"register_time"`
}

// show all info
// Only used in QueryAllUsers()
func UserDtoAdminMapOption() *xmapper.MapOption {
	return xmapper.NewMapOption(&po.User{}, &UserDto{}, func(i interface{}, j interface{}) interface{} {
		user := i.(po.User)
		userDto := j.(UserDto)
		userDto.PhoneNumber = user.PhoneNumber
		return userDto
	})
}

// show info dependent on authUser
// Only used in QueryUser()
func UserDtoUserMapOption(authUser *po.User) *xmapper.MapOption {
	return xmapper.NewMapOption(&po.User{}, &UserDto{}, func(i interface{}, j interface{}) interface{} {
		if authUser == nil { // not login, nothing (default)
			return j
		}
		user := i.(po.User)
		userDto := j.(UserDto)
		if authUser.Authority == enum.AuthAdmin { // admin, all info
			userDto.PhoneNumber = user.PhoneNumber // add phone number
			return userDto
		} else { // normal, only me
			if user.Uid == authUser.Uid {
				userDto.PhoneNumber = user.PhoneNumber
			}
			return userDto
		}
	})
}
