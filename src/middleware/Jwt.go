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
		panic("Inject failed")
	}
	return srv
}

func (j *JwtService) JwtMiddleware(needAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		user, err := j.JwtCheck(authHeader)
		if err != nil {
			// UnAuthorizedError / TokenExpiredError
			result.Error(err).JSON(c)
			c.Abort()
			return
		}
		if needAdmin && user.Authority != enum.AuthAdmin {
			result.Error(exception.NeedAdminError).JSON(c)
			c.Abort()
			return
		}

		c.Set(j.UserKey, user)
		c.Next()
	}
}

func (j *JwtService) JwtCheck(token string) (*po.User, *exception.ServerError) {
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

	// check dao & admin
	user := j.UserDao.QueryByUid(claims.UserId)
	if user == nil {
		return nil, exception.AuthorizedUserNotFoundError
	}

	return user, nil
}

func (j *JwtService) GetAuthUser(c *gin.Context) *po.User {
	_user, exist := c.Get(j.UserKey)
	if !exist { // not jet check auth
		j.JwtMiddleware(false)(c)
		_user, exist = c.Get(j.UserKey)
		if !exist { // Non-Auth
			return nil
		}
	}
	user, ok := _user.(*po.User)
	if !ok { // auth failed
		result.Error(exception.UnAuthorizedError).JSON(c)
		c.Abort()
		return nil
	}
	return user
}
