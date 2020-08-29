package controller

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/constant"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
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
			Params(goapidoc.NewBodyParam("param", "VideoParam", true, "视频请求参数")).
			Responses(goapidoc.NewResponse(200, "_Result<VideoDto>")),

		goapidoc.NewRoutePath("PUT", "/v1/video/{vid}", "更新视频").
			Tags("Video", "Administration").
			Securities("Jwt").
			Params(
				goapidoc.NewPathParam("vid", "integer#int32", true, "视频id"),
				goapidoc.NewBodyParam("param", "VideoParam", true, "视频请求参数"),
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
	videos, total := v.videoService.QueryAll(pp)

	ret := dto.BuildVideoDtos(videos)
	return result.Ok().SetPage(pp.Page, pp.Limit, total, ret)
}

// GET /v1/user/:uid/video
func (v *VideoController) QueryVideosByUid(c *gin.Context) *result.Result {
	uid, ok := param.BindRouteId(c, "uid")
	if !ok {
		return result.Error(exception.RequestParamError)
	}
	pp := param.BindPageOrder(c, v.config)

	videos, total, status := v.videoService.QueryByUid(uid, pp)
	if status == xstatus.DbNotFound {
		return result.Error(exception.UserNotFoundError)
	}

	ret := dto.BuildVideoDtos(videos)
	return result.Ok().SetPage(pp.Page, pp.Limit, total, ret)
}

// GET /v1/video/{vid}
func (v *VideoController) QueryVideoByVid(c *gin.Context) *result.Result {
	vid, ok := param.BindRouteId(c, "vid")
	if !ok {
		return result.Error(exception.RequestParamError)
	}

	video := v.videoService.QueryByVid(vid)
	if video == nil {
		return result.Error(exception.VideoNotFoundError)
	}

	ret := dto.BuildVideoDto(video)
	return result.Ok().SetData(ret)
}

// POST /v1/video
func (v *VideoController) InsertVideo(c *gin.Context) *result.Result {
	authUser := v.jwtService.GetContextUser(c)
	videoParam := &param.VideoParam{}
	if err := c.ShouldBind(videoParam); err != nil {
		return result.Error(exception.WrapValidationError(err))
	}

	video := &po.Video{
		AuthorUid: authUser.Uid,
		Author:    authUser,
	}

	param.MapVideoParam(videoParam, video)
	status := v.videoService.Insert(video)
	if status == xstatus.DbExisted {
		return result.Error(exception.VideoUrlExistError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.VideoInsertError)
	}

	ret := dto.BuildVideoDto(video)
	return result.Created().SetData(ret)
}

// PUT /v1/video/:vid
func (v *VideoController) UpdateVideo(c *gin.Context) *result.Result {
	authUser := v.jwtService.GetContextUser(c)
	vid, ok := param.BindRouteId(c, "vid")
	if !ok {
		return result.Error(exception.RequestParamError)
	}
	videoParam := &param.VideoParam{}
	if err := c.ShouldBind(videoParam); err != nil {
		return result.Error(exception.WrapValidationError(err))
	}

	video := v.videoService.QueryByVid(vid)
	if video == nil {
		return result.Error(exception.VideoNotFoundError)
	} else if authUser.Role != constant.AuthAdmin && authUser.Uid != video.AuthorUid {
		return result.Error(exception.VideoNotFoundError)
	}

	// Update
	param.MapVideoParam(videoParam, video)
	status := v.videoService.Update(video)
	if status == xstatus.DbExisted {
		return result.Error(exception.VideoUrlExistError)
	} else if status == xstatus.DbNotFound {
		return result.Error(exception.VideoNotFoundError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.VideoUpdateError)
	}

	ret := dto.BuildVideoDto(video)
	return result.Ok().SetData(ret)
}

// DELETE /v1/video/:vid
func (v *VideoController) DeleteVideo(c *gin.Context) *result.Result {
	authUser := v.jwtService.GetContextUser(c)
	vid, ok := param.BindRouteId(c, "vid")
	if !ok {
		return result.Error(exception.RequestParamError)
	}

	var status xstatus.DbStatus
	if authUser.Role == constant.AuthAdmin {
		status = v.videoService.Delete(vid)
	} else {
		status = v.videoService.DeleteBy2Id(vid, authUser.Uid)
	}
	if status == xstatus.DbNotFound {
		return result.Error(exception.VideoNotFoundError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.VideoDeleteError)
	}

	return result.Ok()
}
