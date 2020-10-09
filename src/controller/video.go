package controller

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/service"
)

func init() {
	goapidoc.AddRoutePaths(
		goapidoc.NewRoutePath("GET", "/v1/video", "query all videos").
			Tags("Video", "Administration").
			Securities("Jwt").
			Params(
				param.ADPage, param.ADLimit, param.ADOrder,
				_adNeedVideoChannel, _adNeedFavoredCount, _adNeedIsFavorite,
				_adNeedChannelAuthor, _adNeedSubscriberCount, _adNeedVideoCount, _adNeedIsSubscribed,
				_adNeedFollowCount, _adNeedChannelCount, _adNeedSubscribingCount, _adNeedFavoriteCount, _adNeedIsFollow,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<VideoDto>>")),

		goapidoc.NewRoutePath("GET", "/v1/channel/{cid}/video", "query videos from user").
			Tags("Video").
			Params(
				goapidoc.NewPathParam("cid", "integer#int64", true, "channel id"),
				param.ADPage, param.ADLimit, param.ADOrder,
				_adNeedVideoChannel, _adNeedFavoredCount, _adNeedIsFavorite,
				_adNeedChannelAuthor, _adNeedSubscriberCount, _adNeedVideoCount, _adNeedIsSubscribed,
				_adNeedFollowCount, _adNeedChannelCount, _adNeedSubscribingCount, _adNeedFavoriteCount, _adNeedIsFollow,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<VideoDto>>")),

		goapidoc.NewRoutePath("GET", "/v1/video/{vid}", "query a video").
			Tags("Video").
			Params(
				goapidoc.NewPathParam("vid", "integer#int64", true, "video id"),
				_adNeedVideoChannel, _adNeedFavoredCount, _adNeedIsFavorite,
				_adNeedChannelAuthor, _adNeedSubscriberCount, _adNeedVideoCount, _adNeedIsSubscribed,
				_adNeedFollowCount, _adNeedChannelCount, _adNeedSubscribingCount, _adNeedFavoriteCount, _adNeedIsFollow,
			).
			Responses(goapidoc.NewResponse(200, "_Result<VideoDto>")),

		goapidoc.NewRoutePath("POST", "/v1/video", "create a video").
			Tags("Video").
			Securities("Jwt").
			Params(goapidoc.NewBodyParam("param", "InsertVideoParam", true, "create video parameter")).
			Responses(goapidoc.NewResponse(201, "Result")),

		goapidoc.NewRoutePath("PUT", "/v1/video/{vid}", "update a video").
			Tags("Video").
			Securities("Jwt").
			Params(
				goapidoc.NewPathParam("vid", "integer#int64", true, "video id"),
				goapidoc.NewBodyParam("param", "UpdateVideoParam", true, "update video parameter"),
			).
			Responses(goapidoc.NewResponse(200, "Result")),

		goapidoc.NewRoutePath("PUT", "/v1/video/{vid}/channel/{cid}", "update a video's channel").
			Tags("Video").
			Securities("Jwt").
			Params(
				goapidoc.NewPathParam("vid", "integer#int64", true, "video id"),
				goapidoc.NewPathParam("cid", "integer#int64", true, "channel id"),
			).
			Responses(goapidoc.NewResponse(200, "Result")),

		goapidoc.NewRoutePath("PUT", "/v1/channel/{cid}/video/channel/{cid2}", "update all video from a channel to another channel").
			Tags("Video").
			Securities("Jwt").
			Params(
				goapidoc.NewPathParam("cid", "integer#int64", true, "from channel id"),
				goapidoc.NewPathParam("cid2", "integer#int64", true, "to channel id"),
			).
			Responses(goapidoc.NewResponse(200, "Result")),

		goapidoc.NewRoutePath("DELETE", "/v1/video/{vid}", "delete a video").
			Tags("Video").
			Securities("Jwt").
			Params(goapidoc.NewPathParam("vid", "integer#int64", true, "video id")).
			Responses(goapidoc.NewResponse(200, "Result")),
	)
}

type VideoController struct {
	config         *config.Config
	common         *CommonController
	jwtService     *service.JwtService
	videoService   *service.VideoService
	channelService *service.ChannelService
}

func NewVideoController() *VideoController {
	return &VideoController{
		config:         xdi.GetByNameForce(sn.SConfig).(*config.Config),
		common:         xdi.GetByNameForce(sn.SCommonController).(*CommonController),
		jwtService:     xdi.GetByNameForce(sn.SJwtService).(*service.JwtService),
		videoService:   xdi.GetByNameForce(sn.SVideoService).(*service.VideoService),
		channelService: xdi.GetByNameForce(sn.SChannelService).(*service.ChannelService),
	}
}

// GET /v1/video
func (v *VideoController) QueryAllVideos(c *gin.Context) *result.Result {
	pp := param.BindPageOrder(c, v.config)
	videos, total, err := v.videoService.QueryAll(pp)
	if err != nil {
		return result.Error(exception.QueryVideoError).SetError(err, c)
	}

	res := dto.BuildVideoDtos(videos)
	err = v.common.PreLoadVideos(c, v.jwtService.GetContextUser(c), videos, res)
	if err != nil {
		return result.Error(exception.QueryVideoError).SetError(err, c)
	}
	return result.Ok().SetPage(pp.Page, pp.Limit, total, res)
}

// GET /v1/channel/:cid/video
func (v *VideoController) QueryVideosByCid(c *gin.Context) *result.Result {
	cid, err := param.BindRouteId(c, "cid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}
	pp := param.BindPageOrder(c, v.config)

	videos, total, err := v.videoService.QueryByCid(cid, pp)
	if err != nil {
		return result.Error(exception.QueryVideoError).SetError(err, c)
	} else if videos == nil {
		return result.Error(exception.UserNotFoundError)
	}

	res := dto.BuildVideoDtos(videos)
	err = v.common.PreLoadVideos(c, v.jwtService.GetContextUser(c), videos, res)
	if err != nil {
		return result.Error(exception.QueryVideoError).SetError(err, c)
	}
	return result.Ok().SetPage(pp.Page, pp.Limit, total, res)
}

// GET /v1/video/:vid
func (v *VideoController) QueryVideoByVid(c *gin.Context) *result.Result {
	vid, err := param.BindRouteId(c, "vid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	video, err := v.videoService.QueryByVid(vid)
	if err != nil {
		return result.Error(exception.QueryVideoError).SetError(err, c)
	} else if video == nil {
		return result.Error(exception.VideoNotFoundError)
	}

	res := dto.BuildVideoDto(video)
	err = v.common.PreLoadVideos(c, v.jwtService.GetContextUser(c), []*po.Video{video}, []*dto.VideoDto{res})
	if err != nil {
		return result.Error(exception.QueryVideoError).SetError(err, c)
	}
	return result.Ok().SetData(res)
}

// POST /v1/video
func (v *VideoController) InsertVideo(c *gin.Context) *result.Result {
	user := v.jwtService.GetContextUser(c)
	pa := &param.InsertVideoParam{}
	if err := c.ShouldBind(pa); err != nil {
		return result.Error(exception.WrapValidationError(err)).SetError(err, c)
	}

	status, err := v.videoService.Insert(pa, user.Uid)
	if status == xstatus.DbFailed {
		return result.Error(exception.InsertVideoError).SetError(err, c)
	}

	return result.Created()
}

// PUT /v1/video/:vid
func (v *VideoController) UpdateVideo(c *gin.Context) *result.Result {
	user := v.jwtService.GetContextUser(c)
	vid, err := param.BindRouteId(c, "vid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}
	pa := &param.UpdateVideoParam{}
	if err := c.ShouldBind(pa); err != nil {
		return result.Error(exception.WrapValidationError(err)).SetError(err, c)
	}

	video, err := v.videoService.QueryByVid(vid)
	if err != nil {
		return result.Error(exception.UpdateVideoError).SetError(err, c)
	} else if video == nil {
		return result.Error(exception.VideoNotFoundError)
	}
	channel, err := v.channelService.QueryByCid(video.ChannelCid)
	if err != nil {
		return result.Error(exception.UpdateVideoError).SetError(err, c)
	} else if channel == nil {
		return result.Error(exception.VideoNotFoundError)
	} else if channel.AuthorUid != user.Uid {
		return result.Error(exception.VideoPermissionError)
	}

	status, err := v.videoService.Update(vid, pa)
	if status == xstatus.DbNotFound {
		return result.Error(exception.VideoNotFoundError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.UpdateVideoError).SetError(err, c)
	}

	return result.Ok()
}

// PUT /v1/video/:vid/channel/:cid
func (v *VideoController) MoveVideoToChannel(c *gin.Context) *result.Result {
	user := v.jwtService.GetContextUser(c)
	vid, err := param.BindRouteId(c, "vid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}
	cid, err := param.BindRouteId(c, "cid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	video, err := v.videoService.QueryByVid(vid)
	if err != nil {
		return result.Error(exception.UpdateVideoChannelError).SetError(err, c)
	} else if video == nil {
		return result.Error(exception.VideoNotFoundError)
	}
	channel, err := v.channelService.QueryByCid(video.ChannelCid)
	if err != nil {
		return result.Error(exception.UpdateVideoChannelError).SetError(err, c)
	} else if channel == nil {
		return result.Error(exception.VideoNotFoundError)
	} else if channel.AuthorUid != user.Uid {
		return result.Error(exception.VideoPermissionError)
	}
	channel, err = v.channelService.QueryByCid(cid)
	if err != nil {
		return result.Error(exception.UpdateVideoChannelError).SetError(err, c)
	} else if channel == nil {
		return result.Error(exception.ChannelNotFoundError)
	} else if channel.AuthorUid != user.Uid {
		return result.Error(exception.ChannelPermissionError)
	}

	status, err := v.videoService.UpdateChannel(vid, cid)
	if status == xstatus.DbNotFound {
		return result.Error(exception.VideoNotFoundError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.UpdateVideoChannelError).SetError(err, c)
	}

	return result.Ok()
}

// PUT /v1/channel/:cid/video/channel/:cid2
func (v *VideoController) MoveAllVideosToChannel(c *gin.Context) *result.Result {
	user := v.jwtService.GetContextUser(c)
	cid1, err := param.BindRouteId(c, "cid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}
	cid2, err := param.BindRouteId(c, "cid2")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	channel1, err := v.channelService.QueryByCid(cid1)
	if err != nil {
		return result.Error(exception.UpdateVideoChannelError).SetError(err, c)
	} else if channel1 == nil {
		return result.Error(exception.ChannelNotFoundError)
	} else if channel1.AuthorUid != user.Uid {
		return result.Error(exception.ChannelPermissionError)
	}
	channel2, err := v.channelService.QueryByCid(cid2)
	if err != nil {
		return result.Error(exception.UpdateVideoChannelError).SetError(err, c)
	} else if channel2 == nil {
		return result.Error(exception.ChannelNotFoundError)
	} else if channel2.AuthorUid != user.Uid {
		return result.Error(exception.ChannelPermissionError)
	}

	status, err := v.videoService.UpdateAllToChannel(cid1, cid2)
	if status == xstatus.DbNotFound {
		return result.Error(exception.ChannelHasNoVideoError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.UpdateVideoChannelError).SetError(err, c)
	}

	return result.Ok()
}

// DELETE /v1/video/:vid
func (v *VideoController) DeleteVideo(c *gin.Context) *result.Result {
	user := v.jwtService.GetContextUser(c)
	vid, err := param.BindRouteId(c, "vid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	video, err := v.videoService.QueryByVid(vid)
	if err != nil {
		return result.Error(exception.DeleteVideoError).SetError(err, c)
	} else if video == nil {
		return result.Error(exception.VideoNotFoundError)
	}
	channel, err := v.channelService.QueryByCid(video.ChannelCid)
	if err != nil {
		return result.Error(exception.DeleteVideoError).SetError(err, c)
	} else if channel == nil {
		return result.Error(exception.VideoNotFoundError)
	} else if channel.AuthorUid != user.Uid {
		return result.Error(exception.VideoPermissionError)
	}

	status, err := v.videoService.Delete(vid)
	if status == xstatus.DbNotFound {
		return result.Error(exception.VideoNotFoundError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.DeleteVideoError).SetError(err, c)
	}

	return result.Ok()
}
