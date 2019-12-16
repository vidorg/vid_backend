package controller

import (
	"net/http"
	"strconv"
	"vid/app/database"
	"vid/app/database/dao"
	"vid/app/middleware"
	"vid/app/model/dto"
	"vid/app/model/po"
	"vid/app/util"

	"github.com/gin-gonic/gin"
)

type authCtrl struct{}

var AuthCtrl = new(authCtrl)

// POST /auth/login (Non-Auth)
func (u *authCtrl) Login(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")
	expireString := c.PostForm("expire")
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage("Request body error"))
		return
	}
	expire := util.JwtExpire
	if val, err := strconv.Atoi(expireString); err == nil {
		expire = int64(val)
	}

	passRecord := dao.PassDao.Query(username)
	if passRecord == nil {
		c.JSON(http.StatusNotFound,
			dto.Result{}.Error(http.StatusNotFound).SetMessage("User not found"))
		return
	}

	if !util.PassUtil.MD5Check(password, passRecord.EncryptedPass) {
		c.JSON(http.StatusUnauthorized,
			dto.Result{}.Error(http.StatusUnauthorized).SetMessage("Password error"))
		return
	}

	token, err := util.PassUtil.GenToken(passRecord.User.Uid, expire)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			dto.Result{}.Error(http.StatusInternalServerError).SetMessage("Login failed"))
		return
	}

	c.JSON(http.StatusOK,
		dto.Result{}.Ok().PutData("username", passRecord.User.Username).PutData("token", token))
}

// POST /auth/register (Non-Auth)
func (u *authCtrl) Register(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage("Request body error"))
		return
	}

	passRecord := &po.PassRecord{
		EncryptedPass: util.PassUtil.MD5Encode(password),
		User: &po.User{
			Username: username, RegisterIP: c.ClientIP(),
		},
	}
	status := dao.PassDao.Insert(passRecord)
	if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError,
			dto.Result{}.Error(http.StatusInternalServerError).SetMessage("Register failed"))
		return
	} else if status == database.DbExisted {
		c.JSON(http.StatusInternalServerError,
			dto.Result{}.Error(http.StatusInternalServerError).SetMessage("User has existed"))
		return
	}

	c.JSON(http.StatusOK,
		dto.Result{}.Ok().PutData("username", passRecord.User.Username))
}

// POST /auth/pass (Auth)
func (u *authCtrl) ModifyPass(c *gin.Context) {
	user := middleware.GetAuthUser(c)

	password := c.PostForm("password")
	if password == "" {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage("Request body error"))
		return
	}

	passRecord := &po.PassRecord{
		EncryptedPass: util.PassUtil.MD5Encode(password),
		User:          user,
		Uid:           user.Uid,
	}
	status := dao.PassDao.Update(passRecord)
	if status == database.DbNotFound {
		c.JSON(http.StatusInternalServerError,
			dto.Result{}.Error(http.StatusInternalServerError).SetMessage("User not found"))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError,
			dto.Result{}.Error(http.StatusInternalServerError).SetMessage("Register failed"))
		return
	}

	c.JSON(http.StatusOK,
		dto.Result{}.Ok().PutData("username", passRecord.User.Username))
}
