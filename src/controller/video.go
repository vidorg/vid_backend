package controller

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/constant"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/service"
)

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

// @Router              /v1/video [GET]
// @Summary             查询所有视频
// @Description         管理员权限
// @Tag                 Video
// @Tag                 Administration
// @Security            Jwt
// @Template            Order Page
// @ResponseModel 200   #Result<Page<VideoDto>>
func (v *VideoController) QueryAllVideos(c *gin.Context) {
	pageOrder := param.BindPageOrder(c, v.config)
	videos, total := v.videoService.QueryAll(pageOrder)

	ret := dto.BuildVideoDtos(videos)
	result.Ok().SetPage(pageOrder.Page, pageOrder.Limit, total, ret).JSON(c)
}

// @Router              /v1/user/{uid}/video [GET]
// @Summary             查询用户发布的所有视频
// @Tag                 Video
// @Template            Order Page
// @Param               uid path integer true "用户id"
// @ResponseModel 200   #Result<Page<VideoDto>>
func (v *VideoController) QueryVideosByUid(c *gin.Context) {
	uid, ok := param.BindRouteId(c, "uid")
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}
	pageOrder := param.BindPageOrder(c, v.config)

	videos, total, status := v.videoService.QueryByUid(uid, pageOrder)
	if status == database.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	}

	ret := dto.BuildVideoDtos(videos)
	result.Ok().SetPage(pageOrder.Page, pageOrder.Limit, total, ret).JSON(c)
}

// @Router              /v1/video/{vid} [GET]
// @Summary             查询视频
// @Tag                 Video
// @Param               vid path integer true "视频id"
// @ResponseModel 200   #Result<VideoDto>
func (v *VideoController) QueryVideoByVid(c *gin.Context) {
	vid, ok := param.BindRouteId(c, "vid")
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}

	video := v.videoService.QueryByVid(vid)
	if video == nil {
		result.Error(exception.VideoNotFoundError).JSON(c)
		return
	}

	ret := dto.BuildVideoDto(video)
	result.Ok().SetData(ret).JSON(c)
}

// @Router              /v1/video [POST]
// @Summary             新建视频
// @Tag                 Video
// @Security            Jwt
// @Param               param body #VideoParam true "请求参数"
// @ResponseModel 201   #Result<VideoDto>
func (v *VideoController) InsertVideo(c *gin.Context) {
	authUser := v.jwtService.GetContextUser(c)
	videoParam := &param.VideoParam{}
	if err := c.ShouldBind(videoParam); err != nil {
		result.Error(exception.WrapValidationError(err)).JSON(c)
		return
	}

	video := &po.Video{
		AuthorUid: authUser.Uid,
		Author:    authUser,
	}

	param.MapVideoParam(videoParam, video)
	status := v.videoService.Insert(video)
	if status == database.DbExisted {
		result.Error(exception.VideoUrlExistError).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Error(exception.VideoInsertError).JSON(c)
		return
	}

	ret := dto.BuildVideoDto(video)
	result.Created().SetData(ret).JSON(c)
}

// @Router              /v1/video/{vid} [POST]
// @Summary             更新视频
// @Description         管理员或者作者本人权限
// @Tag                 Video
// @Tag                 Administration
// @Security            Jwt
// @Param               vid   path string      true "视频id"
// @Param               param body #VideoParam true "请求参数"
// @ResponseModel 200   #Result<VideoDto>
func (v *VideoController) UpdateVideo(c *gin.Context) {
	authUser := v.jwtService.GetContextUser(c)
	vid, ok := param.BindRouteId(c, "vid")
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}
	videoParam := &param.VideoParam{}
	if err := c.ShouldBind(videoParam); err != nil {
		result.Error(exception.WrapValidationError(err)).JSON(c)
		return
	}

	video := v.videoService.QueryByVid(vid)
	if video == nil {
		result.Error(exception.VideoNotFoundError).JSON(c)
		return
	} else if authUser.Role != constant.AuthAdmin && authUser.Uid != video.AuthorUid {
		result.Error(exception.VideoNotFoundError).JSON(c)
		return
	}

	// Update
	param.MapVideoParam(videoParam, video)
	status := v.videoService.Update(video)
	if status == database.DbExisted {
		result.Error(exception.VideoUrlExistError).JSON(c)
		return
	} else if status == database.DbNotFound {
		result.Error(exception.VideoNotFoundError).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Error(exception.VideoUpdateError).JSON(c)
		return
	}

	ret := dto.BuildVideoDto(video)
	result.Ok().SetData(ret).JSON(c)
}

// @Router              /v1/video/{vid} [DELETE]
// @Summary             删除视频
// @Description         管理员或者作者本人权限
// @Tag                 Video
// @Tag                 Administration
// @Security            Jwt
// @Param               vid path string true "视频id"
// @ResponseModel 200   #Result
func (v *VideoController) DeleteVideo(c *gin.Context) {
	authUser := v.jwtService.GetContextUser(c)
	vid, ok := param.BindRouteId(c, "vid")
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}

	var status database.DbStatus
	if authUser.Role == constant.AuthAdmin {
		status = v.videoService.Delete(vid)
	} else {
		status = v.videoService.DeleteBy2Id(vid, authUser.Uid)
	}
	if status == database.DbNotFound {
		result.Error(exception.VideoNotFoundError).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Error(exception.VideoDeleteError).JSON(c)
		return
	}

	result.Ok().JSON(c)
}
