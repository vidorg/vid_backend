package controller

import (
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xentity"
	"github.com/Aoi-hosizora/ahlib/xslice"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/common/constant"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/service"
	"net/http"
)

type VideoController struct {
	Config       *config.Config         `di:"~"`
	Logger       *logrus.Logger         `di:"~"`
	Mappers      *xentity.EntityMappers `di:"~"`
	JwtService   *service.JwtService    `di:"~"`
	VideoService *service.VideoService  `di:"~"`
}

func NewVideoController(dic *xdi.DiContainer) *VideoController {
	ctrl := &VideoController{}
	dic.MustInject(ctrl)
	return ctrl
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
	pageOrder := param.BindPageOrder(c, v.Config)
	videos, count := v.VideoService.QueryAll(pageOrder)

	retDto := xcondition.First(v.Mappers.MapSlice(xslice.Sti(videos), &dto.VideoDto{})).([]*dto.VideoDto)
	result.Ok().SetPage(count, pageOrder.Page, pageOrder.Limit, retDto).JSON(c)
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
	pageOrder := param.BindPageOrder(c, v.Config)

	videos, count, status := v.VideoService.QueryByUid(uid, pageOrder)
	if status == database.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	}

	retDto := xcondition.First(v.Mappers.MapSlice(xslice.Sti(videos), &dto.VideoDto{})).([]*dto.VideoDto)
	result.Ok().SetPage(count, pageOrder.Page, pageOrder.Limit, retDto).JSON(c)
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

	video := v.VideoService.QueryByVid(vid)
	if video == nil {
		result.Error(exception.VideoNotFoundError).JSON(c)
		return
	}

	retDto := xcondition.First(v.Mappers.Map(video, &dto.VideoDto{})).(*dto.VideoDto)
	result.Ok().SetData(retDto).JSON(c)
}

// @Router              /v1/video [POST]
// @Summary             新建视频
// @Tag                 Video
// @Security            Jwt
// @Param               param body #VideoParam true "请求参数"
// @ResponseModel 201   #Result<VideoDto>
func (v *VideoController) InsertVideo(c *gin.Context) {
	authUser := v.JwtService.GetContextUser(c)
	videoParam := &param.VideoParam{}
	if err := c.ShouldBind(videoParam); err != nil {
		result.Error(exception.WrapValidationError(err)).JSON(c)
		return
	}

	video := &po.Video{
		AuthorUid: authUser.Uid,
		Author:    authUser,
	}

	_ = v.Mappers.MapProp(videoParam, video)
	status := v.VideoService.Insert(video)
	if status == database.DbExisted {
		result.Error(exception.VideoUrlExistError).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Error(exception.VideoInsertError).JSON(c)
		return
	}

	retDto := xcondition.First(v.Mappers.Map(video, &dto.VideoDto{})).(*dto.VideoDto)
	result.Status(http.StatusCreated).SetData(retDto).JSON(c)
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
	authUser := v.JwtService.GetContextUser(c)
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

	video := v.VideoService.QueryByVid(vid)
	if video == nil {
		result.Error(exception.VideoNotFoundError).JSON(c)
		return
	} else if authUser.Role != constant.AuthAdmin && authUser.Uid != video.AuthorUid {
		result.Error(exception.VideoNotFoundError).JSON(c)
		return
	}
	// Update
	_ = v.Mappers.MapProp(videoParam, video)
	status := v.VideoService.Update(video)
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

	retDto := xcondition.First(v.Mappers.Map(video, &dto.VideoDto{})).(*dto.VideoDto)
	result.Ok().SetData(retDto).JSON(c)
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
	authUser := v.JwtService.GetContextUser(c)
	vid, ok := param.BindRouteId(c, "vid")
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}

	var status database.DbStatus
	if authUser.Role == constant.AuthAdmin {
		status = v.VideoService.Delete(vid)
	} else {
		status = v.VideoService.DeleteBy2Id(vid, authUser.Uid)
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
