package profile

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xdatetime"
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/po"
	"strings"
)

func CreateMapperProfile(config *config.ServerConfig) *xmapper.EntityMapper {
	mapper := xmapper.NewEntityMapper()

	mapper = mapper.
		CreateMapper(&po.User{}, &dto.UserDto{}).
		ForMember("Sex", func(i interface{}) interface{} { return i.(po.User).Sex.String() }).
		ForMember("BirthTime", func(i interface{}) interface{} { return i.(po.User).BirthTime.String() }).
		ForMember("Authority", func(i interface{}) interface{} { return i.(po.User).Authority.String() }).
		ForMember("RegisterTime", func(i interface{}) interface{} { return xdatetime.NewJsonDateTime(i.(po.User).CreatedAt).String() }).
		ForMember("AvatarUrl", func(i interface{}) interface{} {
			avatar := i.(po.User).AvatarUrl
			if !strings.HasPrefix(avatar, "http") {
				avatar = xcondition.IfThenElse(avatar == "",
					fmt.Sprintf("%savatar.jpg", config.FileConfig.ImageUrlPrefix),
					fmt.Sprintf("%s%s", config.FileConfig.ImageUrlPrefix, avatar)).(string)
			}
			return avatar
		}).
		ForMember("PhoneNumber", func(i interface{}) interface{} { return "" }).
		Build()

	mapper = mapper.
		CreateMapper(&po.Video{}, &dto.VideoDto{}).
		ForMember("UploadTime", func(i interface{}) interface{} { return i.(po.Video).UploadTime.String() }).
		ForMember("UpdateTime", func(i interface{}) interface{} { return xdatetime.NewJsonDateTime(i.(po.Video).UpdatedAt).String() }).
		ForMember("CoverUrl", func(i interface{}) interface{} {
			cover := i.(po.Video).CoverUrl
			if !strings.HasPrefix(cover, "http") {
				cover = xcondition.IfThenElse(cover == "",
					fmt.Sprintf("%scover.jpg", config.FileConfig.ImageUrlPrefix),
					fmt.Sprintf("%s%s", config.FileConfig.ImageUrlPrefix, cover)).(string)
			}
			return cover
		}).
		ForNest("Author", "Author").
		Build()

	return mapper
}
