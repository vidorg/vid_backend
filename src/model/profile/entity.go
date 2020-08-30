package profile

import (
	"github.com/Aoi-hosizora/ahlib/xentity"
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
)

func addDtoMappers() {
	// videoPo -> videoDto
	xentity.AddMapper(xentity.NewMapper(&po.Video{}, func() interface{} { return &dto.VideoDto{} }, func(from interface{}, to interface{}) error {
		video := from.(*po.Video)
		videoDto := to.(*dto.VideoDto)

		videoDto.Vid = video.Vid
		videoDto.Title = video.Title
		videoDto.Description = video.Description
		videoDto.CoverUrl = video.CoverUrl // TODO
		videoDto.VideoUrl = video.VideoUrl
		videoDto.UploadTime = xtime.NewJsonDateTime(video.CreatedAt).String()
		videoDto.UpdateTime = xtime.NewJsonDateTime(video.UpdatedAt).String()
		videoDto.Author = xentity.MustMap(video.Author, &dto.UserDto{}).(*dto.UserDto)
		return nil
	}))

	// policyPo -> policyDto
	xentity.AddMapper(xentity.NewMapper(&po.RbacRule{}, func() interface{} { return &dto.PolicyDto{} }, func(from interface{}, to interface{}) error {
		policy := from.(*po.RbacRule)
		policyDto := to.(*dto.PolicyDto)

		policyDto.Role = policy.V0
		policyDto.Path = policy.V1
		policyDto.Method = policy.V1
		return nil
	}))
}

func addParamMappers() {
	// videoParam -> videoPo
	xentity.AddMapper(xentity.NewMapper(&param.VideoParam{}, func() interface{} { return &po.Video{} }, func(from interface{}, to interface{}) error {
		videoParam := from.(*param.VideoParam)
		video := to.(*po.Video)

		video.Title = videoParam.Title
		video.Description = *videoParam.Description
		video.CoverUrl = videoParam.CoverUrl // TODO
		video.VideoUrl = videoParam.VideoUrl // TODO
		return nil
	}))
}
