package controller

import (
	"fmt"
	"net/http"
	"vid/app/controller/exception"
	"vid/app/database/dao"
	"vid/app/middleware"
	po2 "vid/app/model/po"
	"vid/app/model/resp"
	"vid/app/util"

	"github.com/gin-gonic/gin"
)

type userCtrl struct{}

var UserCtrl = new(userCtrl)

// GET /user/all (Auth) (Admin)
// @Summary QueryAllUsers
// @Description Get All User
// @Produce json
// @Router /user/all [Post]
func (u *userCtrl) QueryAllUsers(c *gin.Context) {
	authusr, _ := c.Get("user")
	if authusr.(po2.User).Authority != po2.AuthAdmin {
		c.JSON(http.StatusUnauthorized, resp.Message{
			Message: exception.NeedAdminException.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dao.UserDao.QueryAllUsers())
}

// GET /user/uid/:uid (Non-Auth)
func (u *userCtrl) QueryUser(c *gin.Context) {
	uid, ok := util.ReqUtil.GetIntParam(c.Params, "uid")
	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: fmt.Sprintf(exception.RouteParamError.Error(), "uid"),
		})
		return
	}
	query, ok := dao.UserDao.QueryUserByUid(uid)
	if ok {
		// Check Auth to include phone number

		authHeader := c.Request.Header.Get("Authorization")
		_, err := middleware.JWTCheck(authHeader)

		isAuth := err == nil
		info, _ := dao.UserDao.QueryUserExtraInfo(isAuth, query)

		c.JSON(http.StatusOK, resp.UserResp{
			User: *query,
			Info: *info,
		})
	} else {
		c.JSON(http.StatusNotFound, resp.Message{
			Message: exception.UserNotExistException.Error(),
		})
	}
}

// PUT /user/update (Auth)
func (u *userCtrl) UpdateUser(c *gin.Context) {
	body := util.ReqUtil.GetBody(c.Request.Body)
	var user po2.User
	if !user.Unmarshal(body, false) {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: exception.RequestBodyError.Error(),
		})
		return
	}

	authusr, _ := c.Get("user")
	user.Uid = authusr.(po2.User).Uid

	query, err := dao.UserDao.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}

// DELETE /user/delete (Auth)
func (u *userCtrl) DeleteUser(c *gin.Context) {
	authusr, _ := c.Get("user")
	uid := authusr.(po2.User).Uid

	del, err := dao.UserDao.DeleteUser(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, del)
	}
}
