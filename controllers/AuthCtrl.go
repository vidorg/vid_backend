package controllers

import (
	"net/http"

	. "vid/database"
	. "vid/exceptions"
	. "vid/models"
	. "vid/models/req"
	. "vid/models/resp"
	. "vid/utils"

	"github.com/gin-gonic/gin"
)

type authCtrl struct{}

var AuthCtrl = new(authCtrl)

// POST /auth/login (Non-Auth)
func (u *authCtrl) Login(c *gin.Context) {

	body := ReqUtil.GetBody(c.Request.Body)
	var regReq RegLogReq
	if !regReq.Unmarshal(body) {
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

	user, pass, ok := PassDao.QueryPassRecordByUsername(regReq.Username)
	if !ok {
		c.JSON(http.StatusNotFound, Message{
			Message: UserNotExistException.Error(),
		})
		return
	}

	if !PassUtil.MD5Check(regReq.Password, pass.EncryptedPass) {
		c.JSON(http.StatusUnauthorized, Message{
			Message: PasswordException.Error(),
		})
		return
	} else {
		token, err := PassUtil.GenToken(pass.Uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		} else {
			c.Header("Authorization", token)
			c.JSON(http.StatusOK, user)
		}
	}
}

// POST /auth/register (Non-Auth)
func (u *authCtrl) Register(c *gin.Context) {

	body := ReqUtil.GetBody(c.Request.Body)
	var regReq RegLogReq
	if !regReq.Unmarshal(body) {
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

	encryptedPass := PassUtil.MD5Encode(regReq.Password)
	query, err := PassDao.InsertUserPassRecord(regReq.Username, encryptedPass, c.ClientIP())

	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}

// POST /auth/modifypass (Auth)
func (u *authCtrl) ModifyPass(c *gin.Context) {
	body := ReqUtil.GetBody(c.Request.Body)
	var passReq PassReq
	if !passReq.Unmarshal(body) {
		c.JSON(http.StatusBadRequest, Message{
			Message: RequestBodyError.Error(),
		})
		return
	} else if !passReq.CheckFormat() {
		c.JSON(http.StatusBadRequest, Message{
			Message: RegisterFormatError.Error(),
		})
		return
	}
	authusrtmp, _ := c.Get("user")
	authusr := authusrtmp.(User)

	encryptedPass := PassUtil.MD5Encode(passReq.Password)
	passRecord := PassRecord{
		Uid:           authusr.Uid,
		EncryptedPass: encryptedPass,
	}
	_, err := PassDao.UpdatePass(passRecord)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, authusr)
	}
}
