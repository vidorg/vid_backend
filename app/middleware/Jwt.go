package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"vid/app/database/dao"
	"vid/app/model/dto/common"
	"vid/app/model/enum"
	"vid/app/model/po"
	"vid/app/util"

	"github.com/gin-gonic/gin"
	"vid/app/controller/exception"
)

var (
	JwtSecret []byte
	JwtExpire int64
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func JwtMiddleware(needAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		user, err := JwtCheck(authHeader)
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

func JwtCheck(authHeader string) (*po.User, error) {
	// No Head
	if authHeader == "" {
		return nil, exception.AuthorizationError
	}

	// No Token Magic
	parts := strings.Split(authHeader, " ")
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, exception.AuthorizationError
	}

	// Token Parse Err
	claims, err := util.AuthUtil.ParseToken(parts[1])
	if err != nil {
		if strings.Index(err.Error(), "token is expired by") != -1 {
			// Token Expired
			return nil, exception.TokenExpiredError
		} else {
			// Other Error
			// Signature is invalid
			// illegal base64 data at input byte
			return nil, exception.AuthorizationError
		}
	}

	// Check user & Admin
	user := dao.UserDao.QueryByUid(claims.UserID)
	if user == nil {
		return nil, exception.AuthorizationError
	}

	return user, nil
}

func GetAuthUser(c *gin.Context) *po.User {
	_user, exist := c.Get("user")
	if !exist { // Has not Auth
		JwtMiddleware(false)(c)
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
