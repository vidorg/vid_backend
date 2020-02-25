package profile

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xdatetime"
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/po"
	"strings"
)

func loadDtoProfile(config *config.ServerConfig, mapper *xmapper.EntityMapper) *xmapper.EntityMapper {
	user := func(i interface{}) po.User { return i.(po.User) }
	video := func(i interface{}) po.Video { return i.(po.Video) }

	mapper = mapper.
		CreateMapper(&po.User{}, &dto.UserDto{}).
		ForMember("Sex", func(i interface{}) interface{} {
			return user(i).Sex.String() // SexType
		}).
		ForMember("BirthTime", func(i interface{}) interface{} {
			return user(i).BirthTime.String() // JsonDate
		}).
		ForMember("Authority", func(i interface{}) interface{} {
			return user(i).Authority.String() // AuthType
		}).
		ForMember("RegisterTime", func(i interface{}) interface{} {
			return xdatetime.NewJsonDateTime(user(i).CreatedAt).String() // time.Time
		}).
		ForMember("AvatarUrl", func(i interface{}) interface{} {
			avatar := user(i).AvatarUrl
			if !strings.HasPrefix(avatar, "http") {
				if avatar == "" {
					avatar = "avatar.jpg"
				}
				avatar = fmt.Sprintf("%s%s", config.FileConfig.ImageUrlPrefix, avatar)
			}
			return avatar
		}).
		ForMember("PhoneNumber", func(i interface{}) interface{} {
			return ""
		}).
		Build()

	mapper = mapper.
		CreateMapper(&po.Video{}, &dto.VideoDto{}).
		ForMember("UploadTime", func(i interface{}) interface{} {
			return xdatetime.NewJsonDateTime(video(i).CreatedAt).String() // time.Time
		}).
		ForMember("UpdateTime", func(i interface{}) interface{} {
			return xdatetime.NewJsonDateTime(video(i).UpdatedAt).String() // time.Time
		}).
		ForMember("CoverUrl", func(i interface{}) interface{} {
			cover := video(i).CoverUrl
			if !strings.HasPrefix(cover, "http") {
				if cover == "" {
					cover = "cover.jpg"
				}
				cover = fmt.Sprintf("%s%s", config.FileConfig.ImageUrlPrefix, cover)
			}
			return cover
		}).
		ForNest("Author", "Author").
		Build()

	return mapper
}
