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
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/service"
)

func init() {
	goapidoc.AddRoutePaths(
		goapidoc.NewRoutePath("GET", "/v1/video", "管理员查询所有视频").
			Tags("Video", "Administration").
			Securities("Jwt").
			Params(param.ADPage, param.ADLimit, param.ADOrder).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<VideoDto>>")),

		goapidoc.NewRoutePath("GET", "/v1/user/{uid}/video", "查询用户发布的所有视频").
			Tags("Video").
			Params(
				goapidoc.NewPathParam("uid", "integer#int32", true, "用户id"),
				param.ADPage, param.ADLimit, param.ADOrder,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<VideoDto>>")),

		goapidoc.NewRoutePath("GET", "/v1/video/{vid}", "查询视频").
			Tags("Video").
			Params(goapidoc.NewPathParam("vid", "integer#int32", true, "视频id")).
			Responses(goapidoc.NewResponse(200, "_Result<VideoDto>")),

		goapidoc.NewRoutePath("POST", "/v1/video/", "新建视频").
			Tags("Video").
			Securities("Jwt").
			Params(goapidoc.NewBodyParam("param", "InsertVideoParam", true, "视频请求参数")).
			Responses(goapidoc.NewResponse(200, "_Result<VideoDto>")),

		goapidoc.NewRoutePath("PUT", "/v1/video/{vid}", "更新视频").
			Tags("Video", "Administration").
			Securities("Jwt").
			Params(
				goapidoc.NewPathParam("vid", "integer#int32", true, "视频id"),
				goapidoc.NewBodyParam("param", "InsertVideoParam", true, "视频请求参数"),
			).
			Responses(goapidoc.NewResponse(200, "_Result<VideoDto>")),

		goapidoc.NewRoutePath("DELETE", "/v1/video/{vid}", "删除视频").
			Tags("Video", "Administration").
			Securities("Jwt").
			Params(goapidoc.NewPathParam("vid", "integer#int32", true, "视频id")).
			Responses(goapidoc.NewResponse(200, "Result")),
	)
}

type VideoController struct {
	config       *config.Config
	jwtService   *service.JwtService
	videoService *service.VideoService
}

func NewVideoController() *VideoController {
	return &VideoController{
		config:       xdi.GetByNameForce(sn.SConfig).(*config.Config),
		jwtService:   xdi.GetByNameForce(sn.SJwtService).(*service.JwtService),
		videoService: xdi.GetByNameForce(sn.SVideoService).(*service.VideoService),
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

	res := dto.BuildVideoDtos(videos)
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
