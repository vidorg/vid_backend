package controllers

import (
	"fmt"
	"net/http"
	"vid/database"
	. "vid/models"
	"vid/utils"

	"github.com/gin-gonic/gin"
)

type UserCtrl struct{}

var reqUtil = new(utils.ReqUtil)
var userDao = new(database.UserDao)

// GET /all
func (u *UserCtrl) QueryAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, userDao.QueryAllUsers())
}

// GET /one/:uid
func (u *UserCtrl) QueryUser(c *gin.Context) {
	uid, ok := reqUtil.GetIntParam(c.Params, "uid")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf("Route param '%s' not found or error", "uid"),
		})
		return
	}
	query, ok := userDao.QueryUser(uid)
	if ok {
		c.JSON(http.StatusOK, query)
	} else {
		c.JSON(http.StatusNotFound, Message{
			Message: fmt.Sprintf("Uid: %d Not Found", uid),
		})
	}
}

// PUT /insert
func (u *UserCtrl) InsertUser(c *gin.Context) {
	body := reqUtil.GetBody(c.Request.Body)
	var user User
	if !reqUtil.CheckJsonValid(body, &user) {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf("Request body error"),
		})
		return
	}

	query, err := userDao.InsertUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}

// POST /update
func (u *UserCtrl) UpdateUser(c *gin.Context) {
	body := reqUtil.GetBody(c.Request.Body)
	var user User
	if !reqUtil.CheckJsonValid(body, &user) {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf("Request body error"),
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

// DELETE /delete?uid
func (u *UserCtrl) DeleteUser(c *gin.Context) {
	uid, ok := reqUtil.GetIntQuery(c, "uid")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf("Query param '%s' not found or error", "uid"),
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
