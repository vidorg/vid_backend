package controllers

import (
	"fmt"
	"net/http"

	"vid/database"
	"vid/utils"
	. "vid/exceptions"
	. "vid/models"
	. "vid/models/resp"

	"github.com/gin-gonic/gin"
)

type UserCtrl struct{}

var reqUtil = new(utils.ReqUtil)
var userDao = new(database.UserDao)

// GET /user/all (Non-Auth)
func (u *UserCtrl) QueryAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, userDao.QueryAllUsers())
}

// GET /user/one/:uid (Non-Auth)
func (u *UserCtrl) QueryUser(c *gin.Context) {
	uid, ok := reqUtil.GetIntParam(c.Params, "uid")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "uid"),
		})
		return
	}
	query, ok := userDao.QueryUserByUid(uid)
	if ok {
		c.JSON(http.StatusOK, query)
	} else {
		c.JSON(http.StatusNotFound, Message{
			Message: UserNotExistException.Error(),
		})
	}
}

// POST /user/update (Auth)
func (u *UserCtrl) UpdateUser(c *gin.Context) {
	body := reqUtil.GetBody(c.Request.Body)
	var user User
	if !reqUtil.CheckJsonValid(body, &user) {
		c.JSON(http.StatusBadRequest, Message{
			Message: RequestBodyError.Error(),
		})
		return
	}

	authusr, _ := c.Get("user")
	user.Uid = authusr.(User).Uid

	query, err := userDao.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}

// DELETE /user/delete (Auth)
func (u *UserCtrl) DeleteUser(c *gin.Context) {
	authusr, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusUnauthorized, Message{
			Message: AuthorizationException.Error(),
		})
		return
	}
	uid := authusr.(User).Uid

	del, err := userDao.DeleteUser(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, del)
	}
}
