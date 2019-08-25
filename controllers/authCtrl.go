package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vid/database"
	. "vid/models"
	"vid/utils"

	"github.com/gin-gonic/gin"
)

type AuthCtrl struct{}

var passUtil = new(utils.PassUtil)
var passDao = new(database.PassDao)

// POST /auth/login
func (u *AuthCtrl) Login(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"msg": "Login",
	})
}

// POST /auth/register
func (u *AuthCtrl) Register(c *gin.Context) {

	type reqHead struct {
		Username string
		Password string
	}

	body := reqUtil.GetBody(c.Request.Body)
	var regReq reqHead
	err := json.Unmarshal([]byte(body), &regReq)
	if err != nil || regReq.Username == "" || regReq.Password == "" {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf("Request body error"),
		})
		return
	}

	hashPass := passUtil.MD5Encode(regReq.Password)
	query, err := passDao.InsertUserPassRecord(regReq.Username, hashPass)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}
