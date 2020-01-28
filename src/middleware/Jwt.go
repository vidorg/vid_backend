package middleware

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller/exception"
	"github.com/vidorg/vid_backend/src/database/dao"
	"github.com/vidorg/vid_backend/src/model/common"
	"github.com/vidorg/vid_backend/src/model/common/enum"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/util"
	"net/http"
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
			// AuthorizationError / TokenExpiredError
			common.Result{}.Result(http.StatusUnauthorized).SetMessage(err.Error()).JSON(c)
			c.Abort()
			return
		}
		if needAdmin && user.Authority != enum.AuthAdmin {
			common.Result{}.Result(http.StatusUnauthorized).SetMessage(exception.NeedAdminError.Error()).JSON(c)
			c.Abort()
			return
		}

		c.Set(j.UserKey, user)
		c.Next()
	}
}

func (j *JwtService) JwtCheck(token string) (*po.User, error) {
	if token == "" {
		return nil, exception.AuthorizationError
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
			return nil, exception.AuthorizationError
		}
	}

	// check redis
	ok := j.TokenDao.Query(token)
	if !ok {
		return nil, exception.AuthorizationError
	}

	// check dao & admin
	user := j.UserDao.QueryByUid(claims.UserID)
	if user == nil {
		return nil, exception.AuthorizationError
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
		common.Result{}.Result(http.StatusUnauthorized).SetMessage(exception.AuthorizationError.Error()).JSON(c)
		c.Abort()
		return nil
	}
	return user
}
