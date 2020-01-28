package controller

import (
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller/exception"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/database/dao"
	"github.com/vidorg/vid_backend/src/middleware"
	"github.com/vidorg/vid_backend/src/model/common"
	"github.com/vidorg/vid_backend/src/model/common/enum"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/dto/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/util"
	"net/http"
	"time"
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
		panic("Inject failed")
	}
	return ctrl
}

// @Router              /v1/video?page [GET]
// @Template            Auth, Admin
// @Summary             查询所有视频
// @Description         管理员查询所有视频，返回分页数据，管理员权限
// @Tag                 Video
// @Tag                 Administration
// @Param               page query integer false "分页"
// @Accept              multipart/form-data
// @ErrorCode           400 request param error
/* @Response 200		{
							"code": 200,
							"message": "success",
							"data": {
								"count": 1,
								"page": 1,
								"data": [ ${video} ]
							}
 						} */
func (v *VideoController) QueryAllVideos(c *gin.Context) {
	page, ok := param.BindQueryPage(c)
	if !ok {
		common.Result{}.Result(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}

	videos, count := v.VideoDao.QueryAll(page)

	retDto := xcondition.First(v.Mapper.Map([]*dto.VideoDto{}, videos)).([]*dto.VideoDto)
	common.Result{}.Ok().SetPage(count, page, retDto).JSON(c)
}

// @Router              /v1/user/{uid}/video?page [GET]
// @Summary             查询用户视频
// @Description         查询作者为指定用户的所有视频，返回分页数据
// @Tag                 Video
// @Param               uid path integer true "用户id"
// @Param               page query integer false "分页"
// @Accept              multipart/form-data
// @ErrorCode           400 request param error
// @ErrorCode           404 user not found
/* @Response 200		{
							"code": 200,
							"message": "success",
							"data": {
								"count": 1,
								"page": 1,
								"data": [ ${video} ]
							}
 						} */
func (v *VideoController) QueryVideosByUid(c *gin.Context) {
	uid, ok1 := param.BindRouteId(c, "uid")
	page, ok2 := param.BindQueryPage(c)
	if !ok1 || !ok2 {
		common.Result{}.Result(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}

	videos, count, status := v.VideoDao.QueryByUid(uid, page)
	if status == database.DbNotFound {
		common.Result{}.Result(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
		return
	}

	retDto := xcondition.First(v.Mapper.Map([]*dto.VideoDto{}, videos)).([]*dto.VideoDto)
	common.Result{}.Ok().SetPage(count, page, retDto).JSON(c)
}

// @Router              /v1/video/{vid} [GET]
// @Summary             查询视频
// @Description         查询视频信息，作者id为-1表示已删除的用户
// @Tag                 Video
// @Param               vid path integer true "视频id"
// @Accept              multipart/form-data
// @ErrorCode           400 request param error
// @ErrorCode           404 video not found
/* @Response 200		{
							"code": 200,
							"message": "success",
							"data": ${video}
 						} */
func (v *VideoController) QueryVideoByVid(c *gin.Context) {
	vid, ok := param.BindRouteId(c, "vid")
	if !ok {
		common.Result{}.Result(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}

	video := v.VideoDao.QueryByVid(vid)
	if video == nil {
		common.Result{}.Result(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()).JSON(c)
		return
	}

	retDto := xcondition.First(v.Mapper.Map(&dto.VideoDto{}, video)).(*dto.VideoDto)
	common.Result{}.Ok().SetData(retDto).JSON(c)
}

// @Router              /v1/video/ [POST]
// @Template            Auth
// @Summary             新建视频
// @Description         新建用户视频
// @Tag                 Video
// @Param               title formData string true "视频标题，长度在 [1, 100] 之间"
// @Param               description formData string true "视频简介，长度在 [0, 1024] 之间"
// @Param               cover_url formData string false "视频封面链接"
// @Param               video_url formData string true "视频资源链接"
// @Accept              multipart/form-data
// @ErrorCode           400 request param error
// @ErrorCode           400 request format error
// @ErrorCode           400 video has been updated
// @ErrorCode           500 video insert failed
/* @Response 200		{
							"code": 201,
							"message": "created",
							"data": ${video}
 						} */
func (v *VideoController) InsertVideo(c *gin.Context) {
	authUser := v.JwtService.GetAuthUser(c)
	videoParam := &param.VideoParam{}
	if err := c.ShouldBind(videoParam); err != nil {
		common.Result{}.Result(http.StatusBadRequest).SetMessage(exception.WrapValidationError(err).Error()).JSON(c)
		return
	}
	coverUrl, ok := util.CommonUtil.GetFilenameFromUrl(videoParam.CoverUrl, v.Config.FileConfig.ImageUrlPrefix)
	if !ok {
		common.Result{}.Result(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}

	video := &po.Video{
		Title:       videoParam.Title,
		Description: *videoParam.Description,
		CoverUrl:    coverUrl,
		VideoUrl:    videoParam.VideoUrl, // TODO
		UploadTime:  common.JsonDateTime(time.Now()),
		AuthorUid:   authUser.Uid,
		Author:      authUser,
	}
	status := v.VideoDao.Insert(video)
	if status == database.DbExisted {
		common.Result{}.Result(http.StatusBadRequest).SetMessage(exception.VideoExistError.Error()).JSON(c)
		return
	} else if status == database.DbFailed {
		common.Result{}.Error().SetMessage(exception.VideoInsertError.Error()).JSON(c)
		return
	}

	retDto := xcondition.First(v.Mapper.Map(&dto.VideoDto{}, video)).(*dto.VideoDto)
	common.Result{}.Result(http.StatusCreated).SetData(retDto).JSON(c)
}

// @Router              /v1/video/{vid} [POST]
// @Template            Auth, Admin
// @Summary             更新视频
// @Description         更新用户视频信息，管理员或者作者本人可以操作
// @Tag                 Video
// @Tag                 Administration
// @Param               vid path string true "更新视频id"
// @Param               title formData string true "视频标题，长度在 [1, 100] 之间"
// @Param               description formData string true "视频简介，长度在 [0, 1024] 之间"
// @Param               cover_url formData string true "视频封面链接"
// @Param               video_url formData string true "视频资源链接"
// @Accept              multipart/form-data
// @ErrorCode           400 request param error
// @ErrorCode           400 request format error
// @ErrorCode           400 video has been updated
// @ErrorCode           404 video not found
// @ErrorCode           500 video update failed
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": ${video}
 						} */
func (v *VideoController) UpdateVideo(c *gin.Context) {
	authUser := v.JwtService.GetAuthUser(c)
	vid, ok := param.BindRouteId(c, "vid")
	if !ok {
		common.Result{}.Result(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}
	videoParam := &param.VideoParam{}
	if err := c.ShouldBind(videoParam); err != nil {
		common.Result{}.Result(http.StatusBadRequest).SetMessage(exception.WrapValidationError(err).Error()).JSON(c)
		return
	}

	video := v.VideoDao.QueryByVid(vid)
	if video == nil {
		common.Result{}.Result(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()).JSON(c)
		return
	}
	if authUser.Authority != enum.AuthAdmin && authUser.Uid != video.AuthorUid {
		common.Result{}.Result(http.StatusUnauthorized).SetMessage(exception.NeedAdminError.Error()).JSON(c)
		return
	}
	// Update
	video.Title = videoParam.Title
	video.Description = *videoParam.Description
	video.VideoUrl = videoParam.VideoUrl // TODO
	coverUrl, ok := util.CommonUtil.GetFilenameFromUrl(videoParam.CoverUrl, v.Config.FileConfig.ImageUrlPrefix)
	if !ok {
		common.Result{}.Result(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}
	video.CoverUrl = coverUrl

	status := v.VideoDao.Update(video)
	if status == database.DbExisted {
		common.Result{}.Result(http.StatusBadRequest).SetMessage(exception.VideoExistError.Error()).JSON(c)
		return
	} else if status == database.DbNotFound {
		common.Result{}.Result(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()).JSON(c)
		return
	} else if status == database.DbFailed {
		common.Result{}.Error().SetMessage(exception.VideoUpdateError.Error()).JSON(c)
		return
	}

	retDto := xcondition.First(v.Mapper.Map(&dto.VideoDto{}, video)).(*dto.VideoDto)
	common.Result{}.Ok().SetData(retDto).JSON(c)
}

// @Router              /v1/video/{vid} [DELETE]
// @Template            Auth, Admin
// @Summary             删除视频
// @Description         删除用户视频，管理员或者作者本人可以操作
// @Tag                 Video
// @Tag                 Administration
// @Param               vid path string true "删除视频id"
// @Accept              multipart/form-data
// @ErrorCode           400 request param error
// @ErrorCode           404 video not found
// @ErrorCode           500 video delete failed
/* @Response 200		{
							"code": 200,
							"message": "success"
 						} */
func (v *VideoController) DeleteVideo(c *gin.Context) {
	authUser := v.JwtService.GetAuthUser(c)
	vid, ok := param.BindRouteId(c, "vid")
	if !ok {
		common.Result{}.Result(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}
	// Check author and authorization
	video := v.VideoDao.QueryByVid(vid)
	if video == nil {
		common.Result{}.Result(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()).JSON(c)
		return
	}
	if authUser.Authority != enum.AuthAdmin && authUser.Uid != video.AuthorUid {
		common.Result{}.Result(http.StatusUnauthorized).SetMessage(exception.NeedAdminError.Error()).JSON(c)
		return
	}
	// Delete
	status := v.VideoDao.Delete(vid)
	if status == database.DbNotFound {
		common.Result{}.Result(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()).JSON(c)
		return
	} else if status == database.DbFailed {
		common.Result{}.Error().SetMessage(exception.VideoDeleteError.Error()).JSON(c)
		return
	}

	common.Result{}.Ok().JSON(c)
}
