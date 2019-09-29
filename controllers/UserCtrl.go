package controllers

import (
	"fmt"
	"net/http"

	. "vid/database"
	. "vid/exceptions"
	"vid/middleware"
	. "vid/models"
	. "vid/models/resp"
	. "vid/utils"

	"github.com/gin-gonic/gin"
)

type userCtrl struct{}

var UserCtrl = new(userCtrl)

// GET /user/all (Auth) (Admin)
func (u *userCtrl) QueryAllUsers(c *gin.Context) {
	authusr, _ := c.Get("user")
	if authusr.(User).Authority != AuthAdmin {
		c.JSON(http.StatusUnauthorized, Message{
			Message: NeedAdminException.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UserDao.QueryAllUsers())
}

// GET /user/uid/:uid (Non-Auth)
func (u *userCtrl) QueryUser(c *gin.Context) {
	uid, ok := ReqUtil.GetIntParam(c.Params, "uid")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "uid"),
		})
		return
	}
	query, ok := UserDao.QueryUserByUid(uid)
	if ok {
		// Check Auth to include phone number

		authHeader := c.Request.Header.Get("Authorization")
		_, err := middleware.JWTCheck(authHeader)

		isAuth := err == nil
		info, _ := UserDao.QueryUserExtraInfo(isAuth, query)

		c.JSON(http.StatusOK, UserResp{
			User: *query,
			Info: *info,
		})
	} else {
		c.JSON(http.StatusNotFound, Message{
			Message: UserNotExistException.Error(),
		})
	}
}

// POST /user/update (Auth)
func (u *userCtrl) UpdateUser(c *gin.Context) {
	body := ReqUtil.GetBody(c.Request.Body)
	var user User
	if !user.Unmarshal(body, false) {
		c.JSON(http.StatusBadRequest, Message{
			Message: RequestBodyError.Error(),
		})
		return
	}

	authusr, _ := c.Get("user")
	user.Uid = authusr.(User).Uid

	query, err := UserDao.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}

// DELETE /user/delete (Auth)
func (u *userCtrl) DeleteUser(c *gin.Context) {
	authusr, _ := c.Get("user")
	uid := authusr.(User).Uid

	del, err := UserDao.DeleteUser(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, del)
	}
}
