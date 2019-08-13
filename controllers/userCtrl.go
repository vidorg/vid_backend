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
	defer func() {
		if err := recover(); err != nil {
			errRet(err, c)
		}
	}()
	c.JSON(http.StatusOK, userDao.QueryAllUsers())
}

// GET /one/:id
func (u *UserCtrl) QueryUser(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errRet(err, c)
		}
	}()

	id := reqUtil.GetIntParam(c.Params, "id")
	query := userDao.QueryUser(id)
	if query != nil {
		c.JSON(http.StatusOK, query)
	} else {
		c.JSON(http.StatusNotFound, Message{
			Message: fmt.Sprintf("ID: %d Not Found", id),
		})
	}
}

// PUT /insert
func (u *UserCtrl) InsertUser(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errRet(err, c)
		}
	}()

	body := reqUtil.GetBody(c.Request.Body)
	var user User
	reqUtil.CheckJsonValid(body, &user, "user")

	query := userDao.InsertUser(user)
	if query != nil {
		c.JSON(http.StatusOK, query)
	} else {
		c.JSON(http.StatusInternalServerError, Message{
			Message: fmt.Sprintf("ID: %d Insert Failed", user.ID),
		})
	}
}

// POST /update
func (u *UserCtrl) UpdateUser(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errRet(err, c)
		}
	}()

	body := reqUtil.GetBody(c.Request.Body)
	var user User
	reqUtil.CheckJsonValid(body, &user, "user")

	query := userDao.UpdateUser(user)
	if query != nil {
		c.JSON(http.StatusOK, query)
	} else {
		c.JSON(http.StatusInternalServerError, Message{
			Message: fmt.Sprintf("ID: %d Update Failed", user.ID),
		})
	}
}

// DELETE /delete
func (u *UserCtrl) DeleteUser(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errRet(err, c)
		}
	}()

	body := reqUtil.GetBody(c.Request.Body)
	var user User
	reqUtil.CheckJsonValid(body, &user, "user")
	id := user.ID

	if userDao.DeleteUser(id) {
		c.JSON(http.StatusOK, Message{
			Message: fmt.Sprintf("ID: %d Delete Success", id),
		})
	} else {
		c.JSON(http.StatusInternalServerError, Message{
			Message: fmt.Sprintf("ID: %d Delete Failed", id),
		})
	}
}
