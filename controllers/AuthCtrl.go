package controllers

import (
	"net/http"
	
	"vid/database"
	"vid/utils"
	. "vid/exceptions"
	. "vid/models/head"
	. "vid/models/resp"

	"github.com/gin-gonic/gin"
)

type AuthCtrl struct{}

var passUtil = new(utils.PassUtil)
var passDao = new(database.PassDao)

// POST /auth/login
func (u *AuthCtrl) Login(c *gin.Context) {

	body := reqUtil.GetBody(c.Request.Body)
	var regReq RegLogHead
	if !reqUtil.CheckJsonValid(body, &regReq) {
		c.JSON(http.StatusBadRequest, Message{
			Message: RequestBodyError.Error(),
		})
		return
	} else if !regReq.CheckFormat() {
		c.JSON(http.StatusBadRequest, Message{
			Message: LoginFormatError.Error(),
		})
		return
	}

	user, pass, ok := passDao.QueryPassRecordByUsername(regReq.Username)
	if !ok {
		c.JSON(http.StatusNotFound, Message{
			Message: UserNotExistException.Error(),
		})
		return
	}

	if !passUtil.MD5Check(regReq.Password, pass.EncryptedPass) {
		c.JSON(http.StatusUnauthorized, Message{
			Message: PasswordError.Error(),
		})
		return
	} else {
		token, err := passUtil.GenToken(pass.Uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		} else {
			c.Header("Authorization", token)
			c.JSON(http.StatusOK, user)
		}
	}
}

// POST /auth/register
func (u *AuthCtrl) Register(c *gin.Context) {

	body := reqUtil.GetBody(c.Request.Body)
	var regReq RegLogHead
	if !reqUtil.CheckJsonValid(body, &regReq) {
		c.JSON(http.StatusBadRequest, Message{
			Message: RequestBodyError.Error(),
		})
		return
	} else if !regReq.CheckFormat() {
		c.JSON(http.StatusBadRequest, Message{
			Message: RegisterFormatError.Error(),
		})
		return
	}

	encryptedPass := passUtil.MD5Encode(regReq.Password)
	query, err := passDao.InsertUserPassRecord(regReq.Username, encryptedPass)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}
