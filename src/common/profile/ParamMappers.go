package profile

import (
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xdatetime"
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"github.com/Aoi-hosizora/ahlib/xentity"
	"github.com/vidorg/vid_backend/src/common/enum"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/util"
)

func addParamMappers(config *config.ServerConfig, mappers *xentity.EntityMappers) {
	mappers.AddMapper(xentity.NewEntityMapper(&param.UserParam{}, &po.User{}, func(from interface{}, to interface{}) error {
		userParam := from.(*param.UserParam)
		user := to.(*po.User)

		user.Username = userParam.Username
		user.Profile = *userParam.Profile
		user.Sex = enum.ParseSexType(userParam.Sex)
		user.Birthday = xcondition.First(xdatetime.ParseISO8601Date(userParam.Birthday)).(xdatetime.JsonDate)
		user.PhoneNumber = userParam.PhoneNumber
		user.AvatarUrl = util.CommonUtil.GetFilenameFromUrl(userParam.AvatarUrl, config.FileConfig.ImageUrlPrefix)
		return nil
	}))

	mappers.AddMapper(xentity.NewEntityMapper(&param.VideoParam{}, &po.Video{}, func(from interface{}, to interface{}) error {
		videoParam := from.(*param.VideoParam)
		video := to.(*po.Video)

		video.Title = videoParam.Title
		video.Description = *videoParam.Description
		video.CoverUrl = util.CommonUtil.GetFilenameFromUrl(videoParam.CoverUrl, config.FileConfig.ImageUrlPrefix)
		video.VideoUrl = videoParam.VideoUrl
		return nil
	}))
}
