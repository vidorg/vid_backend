package middleware

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/enum"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database/dao"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/util"
	"log"
)

type JwtService struct {
	Config   *config.ServerConfig `di:"~"`
	TokenDao *dao.TokenDao        `di:"~"`
	UserDao  *dao.UserDao         `di:"~"`

	UserKey string `di:"-"`
}

func NewJwtService(dic *xdi.DiContainer) *JwtService {
	srv := &JwtService{UserKey: "user"}
	if !dic.Inject(srv) {
		log.Fatalln("Inject failed")
	}
	return srv
}

func (j *JwtService) JwtMiddleware(needAdmin bool) gin.HandlerFunc {
	return j.processJwt(needAdmin, true)
}

func (j *JwtService) GetAuthToken(c *gin.Context) string {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		authHeader = c.DefaultQuery("Authorization", "")
	}
	return authHeader
}

func (j *JwtService) processJwt(needAdmin bool, forMw bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := j.GetAuthToken(c)

		// check token
		user, err := j.jwtCheck(authHeader)
		if err != nil {
			if forMw { // UnAuthorizedError / TokenExpiredError / AuthorizedUserNotFoundError
				result.Error(err).JSON(c)
				c.Abort()
			}
			return
		}

		// check admin
		if needAdmin && user.Authority != enum.AuthAdmin {
			result.Error(exception.NeedAdminError).JSON(c)
			c.Abort()
			return
		}

		// Success
		c.Set(j.UserKey, user)
		c.Next()
	}
}

func (j *JwtService) jwtCheck(token string) (*po.User, *exception.ServerError) {
	if token == "" {
		return nil, exception.UnAuthorizedError
	}

	// parse
	claims, err := util.AuthUtil.ParseToken(token, j.Config.JwtConfig)
	if err != nil {
		if util.AuthUtil.IsTokenExpireError(err) {
			// Token Expired
			return nil, exception.TokenExpiredError
		} else {
			// Other Error
			// Signature is invalid
			// illegal base64 data at input byte
			return nil, exception.UnAuthorizedError
		}
	}

	// check redis
	ok := j.TokenDao.Query(token)
	if !ok {
		return nil, exception.UnAuthorizedError
	}

	// check dao
	user := j.UserDao.QueryByUid(claims.UserId)
	if user == nil {
		return nil, exception.AuthorizedUserNotFoundError
	}

	return user, nil
}

func (j *JwtService) GetContextUser(c *gin.Context) *po.User {
	_user, exist := c.Get(j.UserKey)
	if !exist { // not check token yet
		j.processJwt(false, false)(c)
		// after check token
		_user, exist = c.Get(j.UserKey)
		if !exist {
			return nil // nil directly
		}
		user, ok := _user.(*po.User)
		if !ok {
			return nil
		}
		return user
	} else { // need auth
		user, ok := _user.(*po.User)
		if !ok {
			result.Error(exception.UnAuthorizedError).JSON(c)
			c.Abort() // need abort
			return nil
		}
		return user
	}
}
