package middleware

import (
	"net/http"
	"strings"
	"vid/app/database/dao"
	"vid/app/model/dto"
	"vid/app/model/po"
	"vid/app/util"

	"github.com/gin-gonic/gin"
	. "vid/app/controller/exception"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		user, err := JWTCheck(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.Result{}.Error(http.StatusUnauthorized).SetMessage(err.Error()))
			c.Abort()
			return
		} else {
			c.Set("user", user)
			c.Next()
		}
	}
}

func JWTCheck(authHeader string) (*po.User, error) {

	// No Head
	if authHeader == "" {
		return nil, AuthorizationException
	}

	// No Token Magic
	parts := strings.Split(authHeader, " ")
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, AuthorizationException
	}

	// Token Parse Err
	claims, err := util.PassUtil.ParseToken(parts[1])
	if err != nil {
		if strings.Index(err.Error(), "token is expired by") != -1 {
			// Token Expired
			return nil, TokenExpiredException
		} else {
			// Other Error
			// Signature is invalid
			// illegal base64 data at input byte

			return nil, AuthorizationException
		}
	}

	// No User
	user, ok := dao.UserDao.QueryUserByUid(claims.UserID)
	if !ok {
		return nil, TokenInvalidException
	}

	return user, nil
}

func GetAuthUser(c *gin.Context) *po.User {
	_user, exist := c.Get("user")
	if !exist {
		return nil
	}
	user, ok := _user.(*po.User)
	if !ok {
		return nil
	}
	return user
}
