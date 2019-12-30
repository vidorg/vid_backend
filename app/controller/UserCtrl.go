package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shomali11/util/xconditions"
	"net/http"
	"strconv"
	"vid/app/controller/exception"
	"vid/app/database"
	"vid/app/database/dao"
	"vid/app/middleware"
	"vid/app/model"
	"vid/app/model/dto"
	"vid/app/model/dto/common"
	"vid/app/model/enum"
	"vid/app/util"
)

type userCtrl struct{}

var UserCtrl = new(userCtrl)

// @Router 				/user?page [GET] [Auth]
// @Summary 			查询所有用户
// @Description 		管理员查询所有用户，返回分页数据，Admin
// @Tag					User
// @Tag					Administration
// @Param 				page query integer false "分页"
// @Accept 				multipart/form-data
// @ErrorCode			400 request query param error
// @ErrorCode			401 need admin authority
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"count": 1,
								"page": 1,
								"data": [
									{
										"uid": 10,
										"username": "aoihosizora",
										"sex": "male",
										"profile": "Demo Profile",
										"avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",
										"birth_time": "2019-12-26",
										"authority": "admin"
									}
								]
							}
 						} */
func (u *userCtrl) QueryAllUsers(c *gin.Context) {
	pageString := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.QueryParamError.Error()))
		return
	}
	page = xconditions.IfThenElse(page == 0, 1, page).(int)

	users, count := dao.UserDao.QueryAll(page)
	c.JSON(http.StatusOK, common.Result{}.Ok().SetPage(count, page, dto.UserDto{}.FromPos(users, enum.DtoOptionAll)))
}

// @Router 				/user/{uid} [GET]
// @Summary 			查询用户
// @Description 		查询用户信息
// @Tag					User
// @Param 				uid path integer true "用户id"
// @Accept 				multipart/form-data
// @ErrorCode			400 request route param error
// @ErrorCode			404 user not found
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"user": {
									"uid": 10,
									"username": "aoihosizora",
									"sex": "male",
									"profile": "Demo Profile",
									"avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",
									"birth_time": "2019-12-26",
									"authority": "admin"
									"phone_number": "13512345678"
								},
								"extra": {
									"subscribing_cnt": 3,
									"subscriber_cnt": 2,
									"video_cnt": 3,
									"playlist_cnt": 0
								}
							}
 						} */
func (u *userCtrl) QueryUser(c *gin.Context) {
	uidString := c.Param("uid")
	uid, err := strconv.Atoi(uidString)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RouteParamError.Error()))
		return
	}

	user := dao.UserDao.QueryByUid(uid)
	if user == nil {
		c.JSON(http.StatusNotFound, common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	}
	authUser := middleware.GetAuthUser(c)

	subscribingCnt, subscriberCnt, _ := dao.SubDao.QuerySubCnt(user.Uid)
	videoCnt, _ := dao.VideoDao.QueryCount(user.Uid)
	extraInfo := &dto.UserExtraInfo{
		SubscribingCount: subscribingCnt,
		SubscriberCount:  subscriberCnt,
		VideoCount:       videoCnt,
		PlaylistCount:    0, // TODO
	}

	c.JSON(http.StatusOK, common.Result{}.Ok().PutData("user", dto.UserDto{}.FromPoThroughUser(user, authUser, uid)).PutData("extra", extraInfo))
}

// @Router 				/user/ [PUT] [Auth]
// @Summary 			更新用户
// @Description 		更新用户信息
// @Tag					User
// @Param 				username formData string true "用户名" minLength(8) maxLength(30)
// @Param 				sex formData string true "用户性别" enum(male, female, unknown)
// @Param 				profile formData string true "用户简介" minLength(0) maxLength(255)
// @Param 				birth_time formData string true "用户生日，固定格式为2000-01-01"
// @Param 				phone_number formData string true "用户手机号码"
// @Param 				avatar formData file false "用户头像，默认不修改"
// @Accept 				multipart/form-data
// @ErrorCode 			400 request form data error
// @ErrorCode 			400 request format error
// @ErrorCode 			400 request body too large
// @ErrorCode 			400 username has been used
// @ErrorCode 			400 image type not supported
// @ErrorCode 			404 user not found
// @ErrorCode 			500 image save failed
// @ErrorCode 			500 user update failed
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"uid": 10,
								"username": "aoihosizora",
								"sex": "male",
								"profile": "Demo Profile",
								"avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",
								"birth_time": "2019-12-26",
								"authority": "admin",
								"phone_number": "13512345678"
							}
 						} */
func (u *userCtrl) UpdateUser(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)

	username, exist1 := c.GetPostForm("username")
	profile, exist2 := c.GetPostForm("profile")
	sex, exist3 := c.GetPostForm("sex")
	birthTimeStr, exist4 := c.GetPostForm("birth_time")
	birthTime, err1 := common.JsonDate{}.Parse(birthTimeStr)
	phoneNumber, exist5 := c.GetPostForm("phone_number")
	if !exist1 || !exist2 || !exist3 || !exist4 || !exist5 {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.FormParamError.Error()))
		return
	}
	if !model.FormatCheck.Username(username) || !model.FormatCheck.UserProfile(profile) || !model.FormatCheck.PhoneNumber(phoneNumber) || err1 != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.FormatError.Error()))
		return
	}
	avatarFile, avatarHeader, err2 := c.Request.FormFile("avatar")
	if err2 == nil && avatarFile != nil {
		supported, ext := util.ImageUtil.CheckImageExt(avatarHeader.Filename)
		if !supported {
			c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.ImageNotSupportedError.Error()))
			return
		}
		filename := fmt.Sprintf("avatar_%s.jpg", util.CommonUtil.CurrentTimeUuid())
		savePath := fmt.Sprintf("./usr/image/%d/%s", authUser.Uid, filename)
		if err := util.ImageUtil.SaveAsJpg(avatarFile, ext, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.ImageSaveError.Error()))
			return
		}
		authUser.AvatarUrl = filename
	}

	authUser.Username = username
	authUser.Sex = enum.SexType("").FromString(sex)
	authUser.Profile = profile
	authUser.BirthTime = birthTime
	authUser.PhoneNumber = phoneNumber

	status := dao.UserDao.Update(authUser)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound, common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	} else if status == database.DbExisted {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.UserNameUsedError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError, common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.UserUpdateError.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Result{}.Ok().SetData(dto.UserDto{}.FromPo(authUser, enum.DtoOptionAll)))
}

// @Router 				/user/ [DELETE] [Auth]
// @Summary 			删除用户
// @Description 		删除用户所有信息
// @Tag					User
// @Accept 				multipart/form-data
// @ErrorCode 			404 user not found
// @ErrorCode 			500 user delete failed
/* @Success 200 		{
							"code": 200,
							"message": "success"
 						} */
func (u *userCtrl) DeleteUser(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)

	status := dao.UserDao.Delete(authUser.Uid)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound, common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError, common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.UserDeleteError.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Result{}.Ok())
}
