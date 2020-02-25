package controller

import (
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/enum"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/database/dao"
	"github.com/vidorg/vid_backend/src/middleware"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"log"
	"net/http"
)

type VideoController struct {
	Config     *config.ServerConfig   `di:"~"`
	JwtService *middleware.JwtService `di:"~"`
	VideoDao   *dao.VideoDao          `di:"~"`
	Mapper     *xmapper.EntityMapper  `di:"~"`
}

func NewVideoController(dic *xdi.DiContainer) *VideoController {
	ctrl := &VideoController{}
	if !dic.Inject(ctrl) {
		log.Fatalln("Inject failed")
	}
	return ctrl
}

// @Router              /v1/video [GET]
// @Security            Jwt
// @Template            Admin Auth Order Page
// @Summary             查询所有视频
// @Description         管理员权限
// @Tag                 Video
// @Tag                 Administration
// @ResponseModel 200   #Result<Page<VideoDto>>
// @ResponseEx 200      ${resp_page_videos}
func (v *VideoController) QueryAllVideos(c *gin.Context) {
	page := param.BindQueryPage(c)
	order := param.BindQueryOrder(c)
	videos, count := v.VideoDao.QueryAll(page, order)

	retDto := xcondition.First(v.Mapper.Map([]*dto.VideoDto{}, videos)).([]*dto.VideoDto)
	result.Ok().SetPage(count, page, retDto).JSON(c)
}

// @Router              /v1/user/{uid}/video [GET]
// @Template            ParamA Order Page
// @Summary             查询用户发布的所有视频
// @Tag                 Video
// @Param               uid path integer true "用户id"
// @ResponseDesc 404    "user not found"
// @ResponseModel 200   #Result<Page<VideoDto>>
// @ResponseEx 200      ${resp_page_videos}
func (v *VideoController) QueryVideosByUid(c *gin.Context) {
	uid, ok := param.BindRouteId(c, "uid")
	page := param.BindQueryPage(c)
	order := param.BindQueryOrder(c)
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}

	videos, count, status := v.VideoDao.QueryByUid(uid, page, order)
	if status == database.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	}

	retDto := xcondition.First(v.Mapper.Map([]*dto.VideoDto{}, videos)).([]*dto.VideoDto)
	result.Ok().SetPage(count, page, retDto).JSON(c)
}

// @Router              /v1/video/{vid} [GET]
// @Template            ParamA
// @Summary             查询视频
// @Description         作者为 null 表示用户已删除
// @Tag                 Video
// @Param               vid path integer true "视频id"
// @ResponseDesc 404    "video not found"
// @ResponseModel 200   #Result<VideoDto>
// @ResponseEx 200      ${resp_video}
func (v *VideoController) QueryVideoByVid(c *gin.Context) {
	vid, ok := param.BindRouteId(c, "vid")
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}

	video := v.VideoDao.QueryByVid(vid)
	if video == nil {
		result.Error(exception.VideoNotFoundError).JSON(c)
		return
	}

	retDto := xcondition.First(v.Mapper.Map(&dto.VideoDto{}, video)).(*dto.VideoDto)
	result.Ok().SetData(retDto).JSON(c)
}

// @Router              /v1/video [POST]
// @Security            Jwt
// @Template            Auth Param
// @Summary             新建视频
// @Tag                 Video
// @Param               param body #VideoParam true "请求参数"
// @ResponseDesc 400    "video url has been used"
// @ResponseDesc 500    "video insert failed"
// @ResponseModel 201   #Result<VideoDto>
// @Response 201        ${resp_new_video}
func (v *VideoController) InsertVideo(c *gin.Context) {
	authUser := v.JwtService.GetContextUser(c)
	videoParam := &param.VideoParam{}
	if err := c.ShouldBind(videoParam); err != nil {
		result.Error(exception.WrapValidationError(err)).JSON(c)
		return
	}

	video := &po.Video{
		AuthorUid:   authUser.Uid,
		Author:      authUser,
	}

	_ = v.Mapper.MapProp(videoParam, video)
	status := v.VideoDao.Insert(video)
	if status == database.DbExisted {
		result.Error(exception.VideoUrlExistError).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Error(exception.VideoInsertError).JSON(c)
		return
	}

	retDto := xcondition.First(v.Mapper.Map(&dto.VideoDto{}, video)).(*dto.VideoDto)
	result.Status(http.StatusCreated).SetData(retDto).JSON(c)
}

// @Router              /v1/video/{vid} [POST]
// @Security            Jwt
// @Template            Auth Param
// @Summary             更新视频
// @Description         管理员或者作者本人权限
// @Tag                 Video
// @Tag                 Administration
// @Param               vid   path string      true "视频id"
// @Param               param body #VideoParam true "请求参数"
// @ResponseDesc 400    "video url has been used"
// @ResponseDesc 404    "video not found"
// @ResponseDesc 500    "video update failed"
// @ResponseModel 200   #Result<VideoDto>
// @ResponseEx 200      ${resp_video}
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

	video := v.VideoDao.QueryByVid(vid)
	if video == nil {
		result.Error(exception.VideoNotFoundError).JSON(c)
		return
	} else if authUser.Authority != enum.AuthAdmin && authUser.Uid != video.AuthorUid {
		result.Error(exception.VideoNotFoundError).JSON(c)
		return
	}
	// Update
	_ = v.Mapper.MapProp(videoParam, video)
	status := v.VideoDao.Update(video)
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

	retDto := xcondition.First(v.Mapper.Map(&dto.VideoDto{}, video)).(*dto.VideoDto)
	result.Ok().SetData(retDto).JSON(c)
}

// @Router              /v1/video/{vid} [DELETE]
// @Security            Jwt
// @Template            Auth ParamA
// @Summary             删除视频
// @Description         管理员或者作者本人权限
// @Tag                 Video
// @Tag                 Administration
// @Param               vid path string true "视频id"
// @ResponseDesc 404    "video not found"
// @ResponseDesc 500    "video delete failed"
// @ResponseModel 200   #Result
// @ResponseEx 200      ${resp_success}
func (v *VideoController) DeleteVideo(c *gin.Context) {
	authUser := v.JwtService.GetContextUser(c)
	vid, ok := param.BindRouteId(c, "vid")
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}

	var status database.DbStatus
	if authUser.Authority == enum.AuthAdmin {
		status = v.VideoDao.Delete(vid)
	} else {
		status = v.VideoDao.DeleteBy2Id(vid, authUser.Uid)
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
