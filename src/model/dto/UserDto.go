package dto

import (
	"fmt"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/enum"
	"github.com/vidorg/vid_backend/src/model/po"
	"strings"
)

type UserDto struct {
	Uid         int    `json:"uid"`
	Username    string `json:"username"`
	Sex         string `json:"sex"`
	Profile     string `json:"profile"`
	AvatarUrl   string `json:"avatar_url"`
	BirthTime   string `json:"birth_time"`
	Authority   string `json:"authority"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

func (UserDto) FromPo(user *po.User, config *config.ServerConfig, option enum.DtoOption, otherOption ...interface{}) *UserDto {
	if !strings.HasPrefix(user.AvatarUrl, "http") {
		if user.AvatarUrl == "" {
			user.AvatarUrl = fmt.Sprintf("%savatar.jpg", config.FileConfig.ImageUrlPrefix)
		} else {
			user.AvatarUrl = fmt.Sprintf("%s%s", config.FileConfig.ImageUrlPrefix, user.AvatarUrl)
		}
	}
	dto := &UserDto{
		Uid:       user.Uid,
		Username:  user.Username,
		Sex:       user.Sex.String(),
		Profile:   user.Profile,
		AvatarUrl: user.AvatarUrl,
		BirthTime: user.BirthTime.String(),
		Authority: user.Authority.String(),
	}

	// All: Return All Number
	// Self + uid: Filter Me
	// None: No Alloc
	if option == enum.DtoOptionAll ||
		(option == enum.DtoOptionOnlySelf && len(otherOption) > 0 && otherOption[0].(int) == user.Uid) {

		dto.PhoneNumber = user.PhoneNumber
	}

	return dto
}

func (UserDto) FromPos(users []*po.User, config *config.ServerConfig, option enum.DtoOption, otherOption ...interface{}) []*UserDto {
	dtos := make([]*UserDto, len(users))
	for idx, user := range users {
		dtos[idx] = UserDto{}.FromPo(user, config, option, otherOption...)
	}
	return dtos
}

// 返回单个的情况
// 根据是否认证，是否为管理员，是否为本人判断
//
// 没有认证: DtoOptionNone
// 已经认证 && (管理员 || 本人): DtoOptionAll
// 已经认证 && 非管理员 && 非本人: DtoOptionNone
func (UserDto) FromPoThroughAuth(retUser *po.User, authUser *po.User, config *config.ServerConfig) *UserDto {
	if authUser != nil && (authUser.Authority == enum.AuthAdmin || authUser.Uid == retUser.Uid) { // IsSelfOrAdmin
		return UserDto{}.FromPo(retUser, config, enum.DtoOptionAll)
	} else {
		return UserDto{}.FromPo(retUser, config, enum.DtoOptionNone)
	}
}

// 返回数组的情况
//
// 没有认证: DtoOptionNone
// 已经认证 && 管理员: DtoOptionAll
// 已经认证 && 非管理员: DtoOptionOnlySelf
func (UserDto) FromPosThroughUser(users []*po.User, authUser *po.User, config *config.ServerConfig) []*UserDto {
	if authUser == nil { // None
		return UserDto{}.FromPos(users, config, enum.DtoOptionNone)
	}
	if authUser.Authority == enum.AuthAdmin { // Admin
		return UserDto{}.FromPos(users, config, enum.DtoOptionAll)
	}
	return UserDto{}.FromPos(users, config, enum.DtoOptionOnlySelf, authUser.Uid)
}
