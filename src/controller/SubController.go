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
	"net/http"
)

type subController struct {
	config  *config.ServerConfig
	userDao *dao.UserDao
	subDao  *dao.SubDao
}

func SubController(config *config.ServerConfig) *subController {
	return &subController{
		config:  config,
		userDao: dao.UserRepository(config.DatabaseConfig),
		subDao:  dao.SubRepository(config.DatabaseConfig),
	}
}

// @Router 				/user/{uid}/subscriber [GET]
// @Summary 			用户粉丝
// @Description 		查询用户所有粉丝，返回分页数据
// @Tag					Subscribe
// @Param 				uid path integer true "查询的用户id"
// @Param 				page query integer false "分页"
// @Accept 				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			404 user not found
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"count": 1,
								"page": 1,
								"data": [
									{
										"uid": 2,
										"username": "User2",
										"sex": "unknown",
										"profile": "",
										"avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",
										"birth_time": "2000-01-01",
										"authority": "normal"
									}
								]
							}
 						} */
func (s *subController) QuerySubscriberUsers(c *gin.Context) {
	uid, ok1 := param.BindRouteId(c, "uid")
	page, ok2 := param.BindQueryPage(c)
	if !ok1 || !ok2 {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}

	users, count := s.subDao.QuerySubscriberUsers(uid, page)
	if users == nil {
		common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
		return
	}

	authUser := middleware.GetAuthUser(c, s.config)
	common.Result{}.Ok().SetPage(count, page, dto.UserDto{}.FromPosThroughUser(users, authUser)).JSON(c)
}

// @Router 				/user/{uid}/subscribing [GET]
// @Summary 			用户关注
// @Description 		查询用户所有关注，返回分页数据
// @Tag					Subscribe
// @Param 				uid path integer true "查询的用户id"
// @Param 				page query integer false "分页"
// @Accept 				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			404 user not found
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"count": 1,
								"page": 1,
								"data": [
									{
										"uid": 1,
										"username": "User1",
										"sex": "male",
										"profile": "",
										"avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",
										"birth_time": "2000-01-01",
										"authority": "admin"
									}
								]
							}
 						} */
func (s *subController) QuerySubscribingUsers(c *gin.Context) {
	uid, ok1 := param.BindRouteId(c, "uid")
	page, ok2 := param.BindQueryPage(c)
	if !ok1 || !ok2 {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}

	users, count := s.subDao.QuerySubscribingUsers(uid, page)
	if users == nil {
		common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
		return
	}

	authUser := middleware.GetAuthUser(c, s.config)
	common.Result{}.Ok().SetPage(count, page, dto.UserDto{}.FromPosThroughUser(users, authUser)).JSON(c)
}

// @Router 				/user/subscribing [PUT] [Auth]
// @Summary 			关注用户
// @Description 		关注某一用户
// @Tag					Subscribe
// @Param 				to formData integer true "关注用户id"
// @Accept 				multipart/form-data
// @ErrorCode 			400 request param error
// @ErrorCode 			400 request format error
// @ErrorCode 			400 subscribe oneself invalid
// @ErrorCode 			404 user not found
// @ErrorCode 			500 subscribe failed
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"me_uid": 10,
								"to_uid": 3,
								"action": "subscribe"
							}
 						} */
func (s *subController) SubscribeUser(c *gin.Context) {
	authUser := middleware.GetAuthUser(c, s.config)
	subParam := &param.SubParam{}
	if err := c.ShouldBind(subParam); err != nil {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.WrapValidationError(err).Error()).JSON(c)
		return
	}
	if authUser.Uid == subParam.Uid {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.SubscribeSelfError.Error()).JSON(c)
		return
	}

	status := s.subDao.SubscribeUser(authUser.Uid, subParam.Uid)
	if status == database.DbNotFound {
		common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
		return
	} else if status == database.DbFailed {
		common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.SubscribeError.Error()).JSON(c)
		return
	}

	common.Result{}.Ok().PutData("me_uid", authUser.Uid).PutData("to_uid", subParam.Uid).PutData("action", "subscribe").JSON(c)
}

// @Router 				/user/subscribing [DELETE] [Auth]
// @Summary 			取消关注用户
// @Description 		取消关注某一用户
// @Tag					Subscribe
// @Param 				to formData integer true "取消关注用户id"
// @Accept 				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode 			400 request format error
// @ErrorCode			404 user not found
// @ErrorCode			500 unsubscribe failed
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"me_uid": 10,
								"to_uid": 3,
								"action": "unsubscribe"
							}
 						} */
func (s *subController) UnSubscribeUser(c *gin.Context) {
	authUser := middleware.GetAuthUser(c, s.config)
	subParam := &param.SubParam{}
	if err := c.ShouldBind(subParam); err != nil {
		common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.WrapValidationError(err).Error()).JSON(c)
		return
	}

	status := s.subDao.UnSubscribeUser(authUser.Uid, subParam.Uid)
	if status == database.DbNotFound {
		common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
		return
	} else if status == database.DbFailed {
		common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.UnSubscribeError.Error()).JSON(c)
		return
	}

	common.Result{}.Ok().PutData("me_uid", authUser.Uid).PutData("to_uid", subParam.Uid).PutData("action", "unsubscribe").JSON(c)
}
