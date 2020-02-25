package profile

import (
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xdatetime"
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/vidorg/vid_backend/src/common/enum"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/util"
)

func loadParamProfile(config *config.ServerConfig, mapper *xmapper.EntityMapper) *xmapper.EntityMapper {
	userParam := func(i interface{}) param.UserParam { return i.(param.UserParam) }
	videoParam := func(i interface{}) param.VideoParam { return i.(param.VideoParam) }

	mapper = mapper.
		CreateMapper(&param.UserParam{}, &po.User{}). // update user
		ForMember("Sex", func(i interface{}) interface{} {
			return enum.ParseSexType(userParam(i).Sex) // SexType
		}).
		ForMember("Profile", func(i interface{}) interface{} {
			return *userParam(i).Profile // string
		}).
		ForMember("BirthTime", func(i interface{}) interface{} {
			return xcondition.First(xdatetime.JsonDate{}.Parse(userParam(i).BirthTime, config.MetaConfig.CurrentLoc)) // JsonDate
		}).
		ForMember("AvatarUrl", func(i interface{}) interface{} {
			return util.CommonUtil.GetFilenameFromUrl(userParam(i).AvatarUrl, config.FileConfig.ImageUrlPrefix) // string
		}).
		Build()

	mapper = mapper.
		CreateMapper(&param.VideoParam{}, &po.Video{}). // create / update video
		ForMember("Description", func(i interface{}) interface{} {
			return *videoParam(i).Description // string
		}).
		ForMember("CoverUrl", func(i interface{}) interface{} {
			return util.CommonUtil.GetFilenameFromUrl(videoParam(i).CoverUrl, config.FileConfig.ImageUrlPrefix) // string
		}).
		Build()

	return mapper
}
