package dto

import (
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/vidorg/vid_backend/src/model/common/enum"
	"github.com/vidorg/vid_backend/src/model/po"
)

type UserDto struct {
	Uid         int32  `json:"uid"`
	Username    string `json:"username"`
	Sex         string `json:"sex"`
	Profile     string `json:"profile"`
	AvatarUrl   string `json:"avatar_url"`
	BirthTime   string `json:"birth_time"`
	Authority   string `json:"authority"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

// show all info
// Only used in QueryAllUsers()
func UserDtoAdminMapper(mapper *xmapper.EntityMapper) *xmapper.EntityMapper {
	return mapper.CreateMapper(&po.User{}, &UserDto{}).ForExtra(func(i interface{}, j interface{}) interface{} {
		user := i.(po.User)
		userDto := j.(UserDto)
		userDto.PhoneNumber = user.PhoneNumber
		return userDto
	}).Build()
}

// show info dependent on authUser
// Only used in QueryUser()
func UserDtoExtraMapper(mapper *xmapper.EntityMapper, authUser *po.User) *xmapper.EntityMapper {
	extra := func(i interface{}, j interface{}) interface{} {
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
	}
	return mapper.CreateMapper(&po.User{}, &UserDto{}).ForExtra(extra).Build()
}
