package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/util"
)

type JwtService struct {
	Config    *config.ServerConfig `di:"~"`
	UserRepo  *UserService         `di:"~"`
	TokenRepo *TokenService        `di:"~"`

	UserKey string `di:"-"`
}

func NewJwtService(dic *xdi.DiContainer) *JwtService {
	srv := &JwtService{UserKey: "user"}
	dic.MustInject(srv)
	return srv
}

func (j *JwtService) GetToken(c *gin.Context) string {
	if token := c.Request.Header.Get("Authorization"); token != "" {
		return token
	} else {
		return c.DefaultQuery("Authorization", "")
	}
}

func (j *JwtService) JwtCheck(token string) (*po.User, *exception.Error) {
	if token == "" {
		return nil, exception.UnAuthorizedError
	}

	// parse
	claims, err := util.AuthUtil.ParseToken(token, j.Config.JwtConfig)
	if err != nil {
		if util.AuthUtil.IsTokenExpireError(err) {
			return nil, exception.TokenExpiredError
		} else {
			return nil, exception.UnAuthorizedError
		}
	}

	// redis
	ok := j.TokenRepo.Query(token)
	if !ok {
		return nil, exception.UnAuthorizedError
	}

	// mysql
	user := j.UserRepo.QueryByUid(claims.UserId)
	if user != nil {
		return nil, exception.UnAuthorizedError
	}

	return user, nil
}

func (j *JwtService) GetContextUser(c *gin.Context) *po.User {
	_user, exist := c.Get(j.UserKey)
	if exist { // has jwtMw
		user, ok := _user.(*po.User)
		if !ok {
			result.Error(exception.UnAuthorizedError).JSON(c)
			c.Abort() // abort
			return nil
		}
		return user
	} else { // no jwtMw
		token := j.GetToken(c)
		user, err := j.JwtCheck(token)
		if err != nil {
			return nil // auth failed
		} else {
			c.Set(j.UserKey, user)
			return user
		}
	}
}
