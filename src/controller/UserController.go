package controller

import (
	"fmt"
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
	"github.com/vidorg/vid_backend/src/util"
	"net/http"
)

type userController struct {
	config   *config.ServerConfig
	userDao  *dao.UserDao
	videoDao *dao.VideoDao
	subDao   *dao.SubDao
}

func UserController(config *config.ServerConfig) *userController {
	return &userController{
		config:   config,
		userDao:  dao.UserRepository(config.DatabaseConfig),
		videoDao: dao.VideoRepository(config.DatabaseConfig),
		subDao:   dao.SubRepository(config.DatabaseConfig),
	}
}

// @Router				/v1/user?page [GET] [Auth]
// @Summary				查询所有用户
// @Description			管理员查询所有用户，返回分页数据，Admin
// @Tag					User
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
								"data": [
									${user}
								]
							}
 						} */
func (u *userController) QueryAllUsers(c *gin.Context) {
	page, ok := param.BindQueryPage(c)
	if !ok {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}

	users, count := u.userDao.QueryAll(page)
	common.Result{}.Ok().SetPage(count, page, dto.UserDto{}.FromPos(users, enum.DtoOptionAll)).JSON(c)
}

// @Router				/v1/user/{uid} [GET]
// @Summary				查询用户
// @Description			查询用户信息
// @Tag					User
// @Param				uid path integer true "用户id"
// @Accept				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			404 user not found
/* @Success 200			{
							"code": 200,
							"message": "success",
							"data": {
								"user": ${user},
								"extra": {
									"subscribing_cnt": 3,
									"subscriber_cnt": 2,
									"video_cnt": 3
								}
							}
 						} */
func (u *userController) QueryUser(c *gin.Context) {
	uid, ok := param.BindRouteId(c, "uid")
	if !ok {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}

	user := u.userDao.QueryByUid(uid)
	if user == nil {
		common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
		return
	}
	subscribingCnt, subscriberCnt, _ := u.subDao.QuerySubCnt(user.Uid)
	videoCnt, _ := u.videoDao.QueryCount(user.Uid)
	extraInfo := &dto.UserExtraInfo{
		SubscribingCount: subscribingCnt,
		SubscriberCount:  subscriberCnt,
		VideoCount:       videoCnt,
	}

	authUser := middleware.GetAuthUser(c, u.config)
	// Mapping from po through the authorization and administration
	common.Result{}.Ok().PutData("user", dto.UserDto{}.FromPoThroughAuth(user, authUser)).PutData("extra", extraInfo).JSON(c)
}

// @Router				/v1/user/ [PUT] [Auth]
// @Summary				更新用户
// @Description			更新用户信息
// @Tag					User
// @Param				username formData string true "用户名，长度在 [8, 30] 之间"
// @Param				sex formData string true "用户性别，允许值为 {male, female, unknown}"
// @Param				profile formData string true "用户简介，长度在 [0, 255] 之间"
// @Param				birth_time formData string true "用户生日，固定格式为 2000-01-01"
// @Param				phone_number formData string true "用户手机号码，长度为 11，仅限中国大陆手机号码"
// @Param				avatar formData file false "用户头像，默认不修改"
// @Accept				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			400 request format error
// @ErrorCode			400 request body too large
// @ErrorCode			400 username has been used
// @ErrorCode			400 image type not supported
// @ErrorCode			404 user not found
// @ErrorCode			500 image save failed
// @ErrorCode			500 user update failed
/* @Success 200			{
							"code": 200,
							"message": "success",
							"data": ${user}
 						} */
func (u *userController) UpdateUser(c *gin.Context) {
	authUser := middleware.GetAuthUser(c, u.config)
	userParam := &param.UserParam{}
	if err := c.ShouldBind(userParam); err != nil {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.WrapValidationError(err).Error()).JSON(c)
		return
	}
	// TODO Separate Api
	if avatarFile, avatarHeader, err := c.Request.FormFile("avatar"); err == nil && avatarFile != nil {
		supported, ext := util.ImageUtil.CheckImageExt(avatarHeader.Filename)
		if !supported {
			common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.ImageNotSupportedError.Error()).JSON(c)
			return
		}
		filename := fmt.Sprintf("avatar_%s.jpg", util.CommonUtil.CurrentTimeUuid())
		savePath := fmt.Sprintf("./usr/image/%s", filename)
		if err := util.ImageUtil.SaveAsJpg(avatarFile, ext, savePath); err != nil {
			common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.ImageSaveError.Error()).JSON(c)
			return
		}
		authUser.AvatarUrl = filename
	}

	authUser.Username = userParam.Username
	authUser.Sex = enum.StringToSexType(userParam.Sex)
	authUser.Profile = userParam.Profile
	authUser.BirthTime = common.JsonDate(userParam.BirthTime)
	authUser.PhoneNumber = userParam.PhoneNumber

	status := u.userDao.Update(authUser)
	if status == database.DbNotFound {
		common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
		return
	} else if status == database.DbExisted {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.UsernameUsedError.Error()).JSON(c)
		return
	} else if status == database.DbFailed {
		common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.UserUpdateError.Error()).JSON(c)
		return
	}

	common.Result{}.Ok().SetData(dto.UserDto{}.FromPo(authUser, enum.DtoOptionAll)).JSON(c)
}

// @Router				/v1/user/ [DELETE] [Auth]
// @Summary				删除用户
// @Description			删除用户所有信息
// @Tag					User
// @Accept				multipart/form-data
// @ErrorCode			404 user not found
// @ErrorCode			500 user delete failed
/* @Success 200			{
							"code": 200,
							"message": "success"
 						} */
func (u *userController) DeleteUser(c *gin.Context) {
	authUser := middleware.GetAuthUser(c, u.config)

	status := u.userDao.Delete(authUser.Uid)
	if status == database.DbNotFound {
		common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
		return
	} else if status == database.DbFailed {
		common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.UserDeleteError.Error()).JSON(c)
		return
	}

	common.Result{}.Ok().JSON(c)
}
