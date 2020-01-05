package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/shomali11/util/xconditions"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller/exception"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/database/dao"
	"github.com/vidorg/vid_backend/src/middleware"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/dto/common"
	"net/http"
	"strconv"
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
	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()))
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()))
		return
	}
	page = xconditions.IfThenElse(page < 1, 1, page).(int)

	users, count := s.subDao.QuerySubscriberUsers(uid, page)
	if users == nil {
		c.JSON(http.StatusNotFound, common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	}

	authUser := middleware.GetAuthUser(c, s.config)
	c.JSON(http.StatusOK, common.Result{}.Ok().SetPage(count, page, dto.UserDto{}.FromPosThroughUser(users, authUser, uid)))
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
	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()))
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()))
		return
	}
	page = xconditions.IfThenElse(page < 1, 1, page).(int)

	users, count := s.subDao.QuerySubscribingUsers(uid, page)
	if users == nil {
		c.JSON(http.StatusNotFound, common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	}

	authUser := middleware.GetAuthUser(c, s.config)
	c.JSON(http.StatusOK, common.Result{}.Ok().SetPage(count, page, dto.UserDto{}.FromPosThroughUser(users, authUser, uid)))
}

// @Router 				/user/subscribing?to [PUT] [Auth]
// @Summary 			关注用户
// @Description 		关注某一用户
// @Tag					Subscribe
// @Param 				to query integer true "关注用户id"
// @Accept 				multipart/form-data
// @ErrorCode 			400 request param error
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

	toUid, err := strconv.Atoi(c.Query("to"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()))
		return
	}
	if authUser.Uid == toUid {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.SubscribeSelfError.Error()))
		return
	}

	status := s.subDao.SubscribeUser(authUser.Uid, toUid)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound, common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError, common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.SubscribeError.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Result{}.Ok().PutData("me_uid", authUser.Uid).PutData("to_uid", toUid).PutData("action", "subscribe"))
}

// @Router 				/user/subscribing?to [DELETE] [Auth]
// @Summary 			取消关注用户
// @Description 		取消关注某一用户
// @Tag					Subscribe
// @Param 				to query integer true "取消关注用户id"
// @Accept 				multipart/form-data
// @ErrorCode			400 request param error
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

	toUid, err := strconv.Atoi(c.Query("to"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()))
		return
	}

	status := s.subDao.UnSubscribeUser(authUser.Uid, toUid)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound, common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError, common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.UnSubscribeError.Error()))
		return
	}

	c.JSON(http.StatusOK,
		common.Result{}.Ok().PutData("me_uid", authUser.Uid).PutData("to_uid", toUid).PutData("action", "unsubscribe"))
}
