package profile

import (
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xdatetime"
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"github.com/Aoi-hosizora/ahlib/xentity"
	"github.com/vidorg/vid_backend/src/common/constant"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/util"
)

func addDtoMappers(config *config.ServerConfig, mappers *xentity.EntityMappers) {
	// userPo -> userDto
	mappers.AddMapper(xentity.NewEntityMapper(&po.User{}, func() interface{} { return &dto.UserDto{} }, func(from interface{}, to interface{}) error {
		user := from.(*po.User)
		userDto := to.(*dto.UserDto)

		userDto.Uid = user.Uid
		userDto.Username = user.Username
		userDto.Sex = user.Sex.String()
		userDto.Profile = user.Profile
		userDto.AvatarUrl = util.CommonUtil.GetServerUrl(user.AvatarUrl, config.FileConfig.ImageUrlPrefix)
		userDto.Birthday = user.Birthday.String()
		userDto.Role = user.Role
		userDto.RegisterTime = xdatetime.NewJsonDateTime(user.CreatedAt).String()
		userDto.PhoneNumber = ""
		return nil
	}))

	// videoPo -> videoDto
	mappers.AddMapper(xentity.NewEntityMapper(&po.Video{}, func() interface{} { return &dto.VideoDto{} }, func(from interface{}, to interface{}) error {
		video := from.(*po.Video)
		videoDto := to.(*dto.VideoDto)

		videoDto.Vid = video.Vid
		videoDto.Title = video.Title
		videoDto.Description = video.Description
		videoDto.VideoUrl = video.VideoUrl
		videoDto.CoverUrl = util.CommonUtil.GetServerUrl(video.CoverUrl, config.FileConfig.ImageUrlPrefix)
		videoDto.UploadTime = xdatetime.NewJsonDateTime(video.CreatedAt).String()
		videoDto.UpdateTime = xdatetime.NewJsonDateTime(video.UpdatedAt).String()
		videoDto.Author = xcondition.First(mappers.Map(video.Author, &dto.UserDto{})).(*dto.UserDto)
		return nil
	}))

	// policyPo -> policyDto
	mappers.AddMapper(xentity.NewEntityMapper(&po.Policy{}, func() interface{} { return &dto.PolicyDto{} }, func(from interface{}, to interface{}) error {
		policy := from.(*po.Policy)
		policyDto := to.(*dto.PolicyDto)

		policyDto.Role = policy.V0
		policyDto.Path = policy.V1
		policyDto.Method = policy.V1
		return nil
	}))
}

func addParamMappers(config *config.ServerConfig, mappers *xentity.EntityMappers) {
	// userParam -> userPo
	mappers.AddMapper(xentity.NewEntityMapper(&param.UserParam{}, func() interface{} { return &po.User{} }, func(from interface{}, to interface{}) error {
		userParam := from.(*param.UserParam)
		user := to.(*po.User)

		user.Username = userParam.Username
		user.Profile = *userParam.Profile
		user.Sex = constant.ParseSexEnum(userParam.Sex)
		user.Birthday = xcondition.First(xdatetime.ParseISO8601Date(userParam.Birthday)).(xdatetime.JsonDate)
		user.PhoneNumber = userParam.PhoneNumber
		user.AvatarUrl = util.CommonUtil.GetFilenameFromUrl(userParam.AvatarUrl, config.FileConfig.ImageUrlPrefix)
		return nil
	}))

	// videoParam -> videoPo
	mappers.AddMapper(xentity.NewEntityMapper(&param.VideoParam{}, func() interface{} { return &po.Video{} }, func(from interface{}, to interface{}) error {
		videoParam := from.(*param.VideoParam)
		video := to.(*po.Video)

		video.Title = videoParam.Title
		video.Description = *videoParam.Description
		video.CoverUrl = util.CommonUtil.GetFilenameFromUrl(videoParam.CoverUrl, config.FileConfig.ImageUrlPrefix)
		video.VideoUrl = videoParam.VideoUrl
		return nil
	}))
}
