package dto

import (
	"fmt"
	"strings"
	"vid/app/model/po"
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

func (UserDto) FromPo(user *po.User, hasPhone bool) *UserDto {
	if !strings.HasPrefix(user.AvatarUrl, "http") {
		if user.AvatarUrl == "" {
			user.AvatarUrl = "http://localhost:3344/raw/image/default/avatar.jpg"
		} else {
			user.AvatarUrl = fmt.Sprintf("http://localhost:3344/raw/image/%d/%s", user.Uid, user.AvatarUrl)
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
	if hasPhone {
		dto.PhoneNumber = user.PhoneNumber
	}
	return dto
}

func (UserDto) FromPos(users []*po.User) []*UserDto {
	dtos := make([]*UserDto, len(users))
	for idx, user := range users {
		dtos[idx] = UserDto{}.FromPo(user, false)
	}
	return dtos
}
