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
}

func JwtMiddleware(needAdmin bool, dic xdi.DiContainer) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		user, err := JwtCheck(authHeader, dic)
		if err != nil {
			// AuthorizationError / TokenExpiredError
			common.Result{}.Error(http.StatusUnauthorized).SetMessage(err.Error()).JSON(c)
			c.Abort()
			return
		}
		if needAdmin && user.Authority != enum.AuthAdmin {
			common.Result{}.Error(http.StatusUnauthorized).SetMessage(exception.NeedAdminError.Error()).JSON(c)
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func JwtCheck(token string, dic xdi.DiContainer) (*po.User, error) {
	if token == "" {
		return nil, exception.AuthorizationError
	}

	jwtCtrl := &struct {
		Config   *config.ServerConfig `di:"~"`
		TokenDao *dao.TokenDao        `di:"~"`
		UserDao  *dao.UserDao         `di:"~"`
	}{}
	dic.Inject(jwtCtrl)

	// parse
	claims, err := util.AuthUtil.ParseToken(token, jwtCtrl.Config.JwtConfig)
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
	ok := jwtCtrl.TokenDao.Query(token)
	if !ok {
		return nil, exception.AuthorizationError
	}

	// check dao & admin
	user := jwtCtrl.UserDao.QueryByUid(claims.UserID)
	if user == nil {
		return nil, exception.AuthorizationError
	}

	return user, nil
}

func GetAuthUser(c *gin.Context, dic xdi.DiContainer) *po.User {
	_user, exist := c.Get("user")
	if !exist { // not jet check auth
		JwtMiddleware(false, dic)(c)
		_user, exist = c.Get("user")
		if !exist { // Non-Auth
			return nil
		}
	}
	user, ok := _user.(*po.User)
	if !ok { // auth failed
		common.Result{}.Error(http.StatusUnauthorized).SetMessage(exception.AuthorizationError.Error()).JSON(c)
		c.Abort()
		return nil
	}
	return user
}
