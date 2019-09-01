package controllers

import (
	"fmt"
	"net/http"

	"vid/database"
	"vid/utils"
	. "vid/exceptions"
	. "vid/models"

	"github.com/gin-gonic/gin"
)

type UserCtrl struct{}

var reqUtil = new(utils.ReqUtil)
var userDao = new(database.UserDao)

// GET /user/all
func (u *UserCtrl) QueryAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, userDao.QueryAllUsers())
}

// GET /user/one/:uid
func (u *UserCtrl) QueryUser(c *gin.Context) {
	uid, ok := reqUtil.GetIntParam(c.Params, "uid")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "uid"),
		})
		return
	}
	query, ok := userDao.QueryUser(uid)
	if ok {
		c.JSON(http.StatusOK, query)
	} else {
		c.JSON(http.StatusNotFound, Message{
			Message: UserNotExistException.Error(),
		})
	}
}

// POST /user/update
func (u *UserCtrl) UpdateUser(c *gin.Context) {
	body := reqUtil.GetBody(c.Request.Body)
	var user User
	if !reqUtil.CheckJsonValid(body, &user) {
		c.JSON(http.StatusBadRequest, Message{
			Message: RequestBodyError.Error(),
		})
		return
	}

	query, err := userDao.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}

// DELETE /user/delete?uid
func (u *UserCtrl) DeleteUser(c *gin.Context) {
	uid, ok := reqUtil.GetIntQuery(c, "uid")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(QueryParamError.Error(), "uid"),
		})
		return
	}

	del, err := userDao.DeleteUser(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, del)
	}
}
