package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/controller/exception"
	"github.com/vidorg/vid_backend/src/database/dao"
	"github.com/vidorg/vid_backend/src/model/dto/common"
	"github.com/vidorg/vid_backend/src/model/enum"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/util"
	"net/http"
)

func JwtMiddleware(needAdmin bool, config *config.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		user, err := JwtCheck(authHeader, config)
		if err != nil {
			c.JSON(http.StatusUnauthorized, common.Result{}.Error(http.StatusUnauthorized).SetMessage(err.Error()))
			c.Abort()
			return
		}
		if needAdmin && user.Authority != enum.AuthAdmin {
			c.JSON(http.StatusUnauthorized, common.Result{}.Error(http.StatusUnauthorized).SetMessage(exception.NeedAdminError.Error()))
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func JwtCheck(authHeader string, config *config.ServerConfig) (*po.User, error) {
	if authHeader == "" {
		return nil, exception.AuthorizationError
	}

	// parse
	claims, err := util.AuthUtil.ParseToken(authHeader, config.JwtConfig)
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

	// check dao & admin
	user := dao.UserRepository(config.DatabaseConfig).QueryByUid(claims.UserID)
	if user == nil {
		return nil, exception.AuthorizationError
	}

	return user, nil
}

func GetAuthUser(c *gin.Context, config *config.ServerConfig) *po.User {
	_user, exist := c.Get("user")
	if !exist { // Has not Auth
		JwtMiddleware(false, config)(c)
		_user, exist = c.Get("user")
		if !exist { // Non-Auth
			return nil
		}
	}
	user, ok := _user.(*po.User)
	if !ok { // Auth Failed
		c.JSON(http.StatusUnauthorized, common.Result{}.Error(http.StatusUnauthorized).SetMessage(exception.AuthorizationError.Error()))
		c.Abort()
		return nil
	}
	return user
}
