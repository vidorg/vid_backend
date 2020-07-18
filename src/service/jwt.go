package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/util"
)

type JwtService struct {
	config       *config.Config
	userService  *UserService
	tokenService *TokenService

	UserKey string
}

func NewJwtService() *JwtService {
	return &JwtService{
		config:       xdi.GetByNameForce(sn.SConfig).(*config.Config),
		userService:  xdi.GetByNameForce(sn.SUserService).(*UserService),
		tokenService: xdi.GetByNameForce(sn.STokenService).(*TokenService),
		UserKey:      "user",
	}
}

func (j *JwtService) GetToken(c *gin.Context) string {
	token := c.Request.Header.Get("Authorization")
	if token != "" {
		return token
	}
	return c.DefaultQuery("Authorization", "")
}

func (j *JwtService) JwtCheck(token string) (*po.User, *exception.Error) {
	if token == "" {
		return nil, exception.UnAuthorizedError
	}

	// parse
	claims, err := util.AuthUtil.ParseToken(token, j.config.Jwt)
	if err != nil {
		if util.AuthUtil.IsTokenExpireError(err) {
			return nil, exception.TokenExpiredError
		} else {
			return nil, exception.UnAuthorizedError
		}
	}

	// redis
	ok := j.tokenService.Query(token)
	if !ok {
		return nil, exception.UnAuthorizedError
	}

	// mysql
	user := j.userService.QueryByUid(claims.UserId)
	if user != nil {
		return nil, exception.UnAuthorizedError
	}

	return user, nil
}

func (j *JwtService) GetContextUser(c *gin.Context) *po.User {
	_user, exist := c.Get(j.UserKey)
	if exist { // jwt middleware has checked (or this method will not be invoked)
		user, ok := _user.(*po.User)
		if !ok {
			result.Error(exception.UnAuthorizedError).JSON(c)
			c.Abort() // abort
			return nil
		}
		return user
	}

	// otherwise, jwt is not required in this http method
	token := j.GetToken(c)
	user, err := j.JwtCheck(token)
	if err != nil {
		return nil // auth failed
	}
	c.Set(j.UserKey, user)
	return user
}
