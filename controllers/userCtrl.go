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

// GET /one/:id
func (u *UserCtrl) QueryUser(c *gin.Context) {
	id, ok := reqUtil.GetIntParam(c.Params, "id")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf("Route param '%s' not found or error", "id"),
		})
		return
	}
	query, ok := userDao.QueryUser(id)
	if ok {
		c.JSON(http.StatusOK, query)
	} else {
		c.JSON(http.StatusNotFound, Message{
			Message: fmt.Sprintf("ID: %d Not Found", id),
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

	query, isExist, ok := userDao.InsertUser(user)
	if isExist {
		c.JSON(http.StatusInternalServerError, Message{
			Message: fmt.Sprintf("ID: %d already exist", user.ID),
		})
	} else if !ok {
		c.JSON(http.StatusInternalServerError, Message{
			Message: fmt.Sprintf("ID: %d insert failed", user.ID),
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

	query, isExist, ok := userDao.UpdateUser(user)
	if !isExist {
		c.JSON(http.StatusInternalServerError, Message{
			Message: fmt.Sprintf("ID: %d not exist", user.ID),
		})
	} else if !ok {
		c.JSON(http.StatusInternalServerError, Message{
			Message: fmt.Sprintf("ID: %d update failed", user.ID),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}

// DELETE /delete?id
func (u *UserCtrl) DeleteUser(c *gin.Context) {
	id, ok := reqUtil.GetIntQuery(c, "id")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf("Query param '%s' not found or error", "id"),
		})
		return
	}

	del, isExist, ok := userDao.DeleteUser(id)
	if !isExist {
		c.JSON(http.StatusInternalServerError, Message{
			Message: fmt.Sprintf("ID: %d not exist", id),
		})
	} else if !ok {
		c.JSON(http.StatusInternalServerError, Message{
			Message: fmt.Sprintf("ID: %d delete failed", id),
		})
	} else {
		c.JSON(http.StatusOK, del)
	}
}
