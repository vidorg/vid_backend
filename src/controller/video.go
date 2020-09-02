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
				_adNeedAuthor, _adNeedFavoredCount, _adNeedIsFavorite, _adNeedSubscribeCount, _adNeedIsSubscribe, _adNeedIsBlock, _adNeedVideoCount, _adNeedFavoriteCount,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<VideoDto>>")),

		goapidoc.NewRoutePath("GET", "/v1/user/{uid}/video", "query videos from user").
			Tags("Video").
			Params(
				goapidoc.NewPathParam("uid", "integer#int64", true, "user id"),
				param.ADPage, param.ADLimit, param.ADOrder,
				_adNeedAuthor, _adNeedFavoredCount, _adNeedIsFavorite, _adNeedSubscribeCount, _adNeedIsSubscribe, _adNeedIsBlock, _adNeedVideoCount, _adNeedFavoriteCount,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<VideoDto>>")),

		goapidoc.NewRoutePath("GET", "/v1/video/{vid}", "query a video").
			Tags("Video").
			Params(
				goapidoc.NewPathParam("vid", "integer#int64", true, "video id"),
				_adNeedAuthor, _adNeedFavoredCount, _adNeedIsFavorite, _adNeedSubscribeCount, _adNeedIsSubscribe, _adNeedIsBlock, _adNeedVideoCount, _adNeedFavoriteCount,
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
				goapidoc.NewBodyParam("param", "InsertVideoParam", true, "update video parameter"),
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
	config       *config.Config
	jwtService   *service.JwtService
	videoService *service.VideoService
	common       *CommonController
}

func NewVideoController() *VideoController {
	return &VideoController{
		config:       xdi.GetByNameForce(sn.SConfig).(*config.Config),
		jwtService:   xdi.GetByNameForce(sn.SJwtService).(*service.JwtService),
		videoService: xdi.GetByNameForce(sn.SVideoService).(*service.VideoService),
		common:       xdi.GetByNameForce(sn.SCommonController).(*CommonController),
	}
}

// GET /v1/video
func (v *VideoController) QueryAllVideos(c *gin.Context) *result.Result {
	pp := param.BindPageOrder(c, v.config)
	videos, total, err := v.videoService.QueryAll(pp)
	if err != nil {
		return result.Error(exception.QueryVideoError).SetError(err, c)
	}

	authUser := v.jwtService.GetContextUser(c)
	authors, userExtras, err := v.common.getVideosAuthor(c, authUser, videos)
	if err != nil {
		return result.Error(exception.QueryVideoError).SetError(err, c)
	}
	videoExtras, err := v.common.getVideosExtra(c, authUser, videos)
	if err != nil {
		return result.Error(exception.QueryVideoError).SetError(err, c)
	}

	res := dto.BuildVideoDtos(videos)
	for idx, video := range res {
		video.Author = dto.BuildUserDto(authors[idx])
		if video.Author != nil {
			video.Author.Extra = userExtras[idx]
		}
		video.Extra = videoExtras[idx]
	}
	return result.Ok().SetPage(pp.Page, pp.Limit, total, res)
}

// GET /v1/user/:uid/video
func (v *VideoController) QueryVideosByUid(c *gin.Context) *result.Result {
	uid, err := param.BindRouteId(c, "uid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}
	pp := param.BindPageOrder(c, v.config)

	videos, total, err := v.videoService.QueryByUid(uid, pp)
	if err != nil {
		return result.Error(exception.QueryVideoError).SetError(err, c)
	} else if videos == nil {
		return result.Error(exception.UserNotFoundError)
	}

	authUser := v.jwtService.GetContextUser(c)
	authors, userExtras, err := v.common.getVideosAuthor(c, authUser, videos)
	if err != nil {
		return result.Error(exception.QueryVideoError).SetError(err, c)
	}
	videoExtras, err := v.common.getVideosExtra(c, authUser, videos)
	if err != nil {
		return result.Error(exception.QueryVideoError).SetError(err, c)
	}

	res := dto.BuildVideoDtos(videos)
	for idx, video := range res {
		video.Author = dto.BuildUserDto(authors[idx])
		if video.Author != nil {
			video.Author.Extra = userExtras[idx]
		}
		video.Extra = videoExtras[idx]
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

	authUser := v.jwtService.GetContextUser(c)
	authors, userExtras, err := v.common.getVideosAuthor(c, authUser, []*po.Video{video})
	if err != nil {
		return result.Error(exception.QueryVideoError).SetError(err, c)
	}
	videoExtras, err := v.common.getVideosExtra(c, authUser, []*po.Video{video})
	if err != nil {
		return result.Error(exception.QueryVideoError).SetError(err, c)
	}

	res := dto.BuildVideoDto(video)
	res.Author = dto.BuildUserDto(authors[0])
	if res.Author != nil {
		res.Author.Extra = userExtras[0]
	}
	res.Extra = videoExtras[0]
	return result.Ok().SetData(res)
}

// POST /v1/video
func (v *VideoController) InsertVideo(c *gin.Context) *result.Result {
	user := v.jwtService.GetContextUser(c)
	pa := &param.InsertVideoParam{}
	if err := c.ShouldBind(pa); err != nil {
		return result.Error(exception.WrapValidationError(err))
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
	} else if video.AuthorUid != user.Uid {
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

// DELETE /v1/video/:vid
func (v *VideoController) DeleteVideo(c *gin.Context) *result.Result {
	user := v.jwtService.GetContextUser(c)
	vid, err := param.BindRouteId(c, "vid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	video, err := v.videoService.QueryByVid(vid)
	if err != nil {
		return result.Error(exception.UpdateVideoError).SetError(err, c)
	} else if video == nil {
		return result.Error(exception.VideoNotFoundError)
	} else if video.AuthorUid != user.Uid {
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
