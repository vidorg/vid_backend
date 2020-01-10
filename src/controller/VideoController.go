package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller/exception"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/database/dao"
	"github.com/vidorg/vid_backend/src/middleware"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/dto/common"
	"github.com/vidorg/vid_backend/src/model/dto/param"
	"github.com/vidorg/vid_backend/src/model/enum"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/util"
	"log"
	"net/http"
	"time"
)

type videoController struct {
	config   *config.ServerConfig
	videoDao *dao.VideoDao
}

func VideoController(config *config.ServerConfig) *videoController {
	return &videoController{
		config:   config,
		videoDao: dao.VideoRepository(config.MySqlConfig),
	}
}

// @Router				/v1/video?page [GET] [Auth]
// @Summary				查询所有视频
// @Description			管理员查询所有视频，返回分页数据，Admin
// @Tag					Video
// @Tag					Administration
// @Param				page query integer false "分页"
// @Accept				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			401 need admin authority
/* @Success 200			{
							"code": 200,
							"message": "success",
							"data": {
								"count": 1,
								"page": 1,
								"data": [ ${video} ]
							}
 						} */
func (v *videoController) QueryAllVideos(c *gin.Context) {
	page, ok := param.BindQueryPage(c)
	if !ok {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}

	videos, count := v.videoDao.QueryAll(page)
	common.Result{}.Ok().SetPage(count, page, dto.VideoDto{}.FromPos(videos, v.config)).JSON(c)
}

// @Router				/v1/user/{uid}/video?page [GET]
// @Summary				查询用户视频
// @Description			查询作者为用户的所有视频，返回分页数据
// @Tag					Video
// @Param				uid path integer true "用户id"
// @Param				page query integer false "分页"
// @Accept				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			404 user not found
/* @Success 200			{
							"code": 200,
							"message": "success",
							"data": {
								"count": 1,
								"page": 1,
								"data": [ ${video} ]
							}
 						} */
func (v *videoController) QueryVideosByUid(c *gin.Context) {
	uid, ok1 := param.BindRouteId(c, "uid")
	page, ok2 := param.BindQueryPage(c)
	if !ok1 || !ok2 {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}

	videos, count, status := v.videoDao.QueryByUid(uid, page)
	if status == database.DbNotFound {
		common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
		return
	}

	common.Result{}.Ok().SetPage(count, page, dto.VideoDto{}.FromPos(videos, v.config)).JSON(c)
}

// @Router				/v1/video/{vid} [GET]
// @Summary				查询视频
// @Description			查询视频信息，作者id为-1表示已删除的用户
// @Tag					Video
// @Param				vid path integer true "视频id"
// @Accept				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			404 video not found
/* @Success 200			{
							"code": 200,
							"message": "success",
							"data": ${video}
 						} */
func (v *videoController) QueryVideoByVid(c *gin.Context) {
	vid, ok := param.BindRouteId(c, "vid")
	if !ok {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}

	video := v.videoDao.QueryByVid(vid)
	if video == nil {
		common.Result{}.Error(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()).JSON(c)
		return
	}

	common.Result{}.Ok().SetData(dto.VideoDto{}.FromPo(video, v.config)).JSON(c)
}

// @Router				/v1/video/ [POST] [Auth]
// @Summary				新建视频
// @Description			新建用户视频
// @Tag					Video
// @Param				title formData string true "视频标题，长度在 [1, 100] 之间"
// @Param				description formData string true "视频简介，长度在 [0, 1024] 之间"
// @Param				cover_url formData string false "视频封面链接"
// @Param				video_url formData string true "视频资源链接"
// @Accept				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			400 request format error
// @ErrorCode			400 request body too large
// @ErrorCode			400 image type not supported
// @ErrorCode			400 video resource has been used
// @ErrorCode			500 image save failed
// @ErrorCode			500 video insert failed
/* @Success 200			{
							"code": 201,
							"message": "created",
							"data": ${video}
 						} */
func (v *videoController) InsertVideo(c *gin.Context) {
	authUser := middleware.GetAuthUser(c, v.config)
	videoParam := &param.VideoParam{}
	if err := c.ShouldBind(videoParam); err != nil {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.WrapValidationError(err).Error()).JSON(c)
		return
	}
	coverUrl, ok := util.CommonUtil.GetFilenameFromUrl(videoParam.CoverUrl, v.config.FileConfig.ImageUrlPrefix)
	if !ok {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
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
	status := v.videoDao.Insert(video)
	if status == database.DbExisted {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.VideoExistError.Error()).JSON(c)
		return
	} else if status == database.DbFailed {
		common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.VideoInsertError.Error()).JSON(c)
		return
	}

	common.Result{}.Created().SetData(dto.VideoDto{}.FromPo(video, v.config)).JSON(c)
}

// @Router				/v1/video/{vid} [POST] [Auth]
// @Summary				更新视频
// @Description			更新用户视频信息，管理员或者作者本人可以操作
// @Tag					Video
// @Tag					Administration
// @Param				vid path string true "更新视频id"
// @Param				title formData string true "视频标题，长度在 [1, 100] 之间"
// @Param				description formData string true "视频简介，长度在 [0, 1024] 之间"
// @Param				cover_url formData string true "视频封面链接"
// @Param				video_url formData string true "视频资源链接"
// @Accept				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			400 request format error
// @ErrorCode			401 need admin authority
// @ErrorCode			404 video not found
// @ErrorCode			500 video update failed
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": ${video}
 						} */
func (v *videoController) UpdateVideo(c *gin.Context) {
	authUser := middleware.GetAuthUser(c, v.config)
	vid, ok := param.BindRouteId(c, "vid")
	if !ok {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}
	videoParam := &param.VideoParam{}
	if err := c.ShouldBind(videoParam); err != nil {
		log.Println(err)
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.WrapValidationError(err).Error()).JSON(c)
		return
	}

	video := v.videoDao.QueryByVid(vid)
	if video == nil {
		common.Result{}.Error(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()).JSON(c)
		return
	}
	if authUser.Authority != enum.AuthAdmin && authUser.Uid != video.AuthorUid {
		common.Result{}.Error(http.StatusUnauthorized).SetMessage(exception.NeedAdminError.Error()).JSON(c)
		return
	}
	// Update
	video.Title = videoParam.Title
	video.Description = *videoParam.Description
	video.VideoUrl = videoParam.VideoUrl // TODO
	coverUrl, ok := util.CommonUtil.GetFilenameFromUrl(videoParam.CoverUrl, v.config.FileConfig.ImageUrlPrefix)
	if !ok {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}
	video.CoverUrl = coverUrl

	status := v.videoDao.Update(video)
	if status == database.DbNotFound {
		common.Result{}.Error(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()).JSON(c)
		return
	} else if status == database.DbFailed {
		common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.VideoUpdateError.Error()).JSON(c)
		return
	}

	common.Result{}.Ok().SetData(dto.VideoDto{}.FromPo(video, v.config)).JSON(c)
}

// @Router				/v1/video/{vid} [DELETE] [Auth]
// @Summary				删除视频
// @Description			删除用户视频，管理员或者作者本人可以操作
// @Tag					Video
// @Tag					Administration
// @Param				vid path string true "删除视频id"
// @Accept				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			401 need admin authority
// @ErrorCode			404 video not found
// @ErrorCode			500 video delete failed
/* @Success 200			{
							"code": 200,
							"message": "success"
 						} */
func (v *videoController) DeleteVideo(c *gin.Context) {
	authUser := middleware.GetAuthUser(c, v.config)
	vid, ok := param.BindRouteId(c, "vid")
	if !ok {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}
	// Check author and authorization
	video := v.videoDao.QueryByVid(vid)
	if video == nil {
		common.Result{}.Error(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()).JSON(c)
		return
	}
	if authUser.Authority != enum.AuthAdmin && authUser.Uid != video.AuthorUid {
		common.Result{}.Error(http.StatusUnauthorized).SetMessage(exception.NeedAdminError.Error()).JSON(c)
		return
	}
	// Delete
	status := v.videoDao.Delete(vid)
	if status == database.DbNotFound {
		common.Result{}.Error(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()).JSON(c)
		return
	} else if status == database.DbFailed {
		common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.VideoDeleteError.Error()).JSON(c)
		return
	}

	common.Result{}.Ok().JSON(c)
}
