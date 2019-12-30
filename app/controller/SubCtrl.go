package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/shomali11/util/xconditions"
	"net/http"
	"strconv"
	"vid/app/controller/exception"
	"vid/app/database"
	"vid/app/database/dao"
	"vid/app/middleware"
	"vid/app/model/dto"
	"vid/app/model/dto/common"
)

type subCtrl struct{}

var SubCtrl = new(subCtrl)

// @Router 				/user/{uid}/subscriber [GET]
// @Summary 			用户粉丝
// @Description 		查询用户所有粉丝，返回分页数据
// @Param 				uid path integer true "查询的用户id"
// @Param 				page query integer false "分页"
// @Accept 				multipart/form-data
// @ErrorCode			400 request route param error
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
func (u *subCtrl) QuerySubscriberUsers(c *gin.Context) {
	uidString := c.Param("uid")
	uid, err := strconv.Atoi(uidString)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RouteParamError.Error()))
		return
	}
	pageString := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.QueryParamError.Error()))
		return
	}
	page = xconditions.IfThenElse(page == 0, 1, page).(int)

	users, count, status := dao.SubDao.QuerySubscriberUsers(uid, page)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound, common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Result{}.Ok().SetPage(count, page, dto.UserDto{}.FromPos(users)))
}

// @Router 				/user/{uid}/subscribing [GET]
// @Summary 			用户关注
// @Description 		查询用户所有关注，返回分页数据
// @Param 				uid path integer true "查询的用户id"
// @Param 				page query integer false "分页"
// @Accept 				multipart/form-data
// @ErrorCode			400 request route param error
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
func (u *subCtrl) QuerySubscribingUsers(c *gin.Context) {
	uidString := c.Param("uid")
	uid, err := strconv.Atoi(uidString)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RouteParamError.Error()))
		return
	}
	pageString := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.QueryParamError.Error()))
		return
	}
	page = xconditions.IfThenElse(page == 0, 1, page).(int)

	users, count, status := dao.SubDao.QuerySubscribingUsers(uid, page)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound, common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Result{}.Ok().SetPage(count, page, dto.UserDto{}.FromPos(users)))
}

// @Router 				/user/subscribing?to [PUT] [Auth]
// @Summary 			关注用户
// @Description 		关注某一用户
// @Param 				to query integer true "关注用户id"
// @Accept 				multipart/form-data
// @ErrorCode 			400 request query param error
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
func (u *subCtrl) SubscribeUser(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)

	toUidString := c.Query("to")
	toUid, err := strconv.Atoi(toUidString)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.QueryParamError.Error()))
		return
	}
	if authUser.Uid == toUid {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.SubscribeSelfError.Error()))
		return
	}

	status := dao.SubDao.SubscribeUser(authUser.Uid, toUid)
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
// @Param 				to query integer true "取消关注用户id"
// @Accept 				multipart/form-data
// @ErrorCode			400 request query param error
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
func (u *subCtrl) UnSubscribeUser(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)

	toUidString := c.Query("to")
	toUid, err := strconv.Atoi(toUidString)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.QueryParamError.Error()))
		return
	}

	status := dao.SubDao.UnSubscribeUser(authUser.Uid, toUid)
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
