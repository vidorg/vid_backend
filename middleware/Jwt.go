package middleware

import (
	"net/http"
	"strings"

	"vid/database"
	. "vid/exceptions"
	. "vid/models"
	"vid/utils"

	"github.com/gin-gonic/gin"
)

var passUtil = new(utils.PassUtil)
var cmnCtrl = new(utils.CmnCtrl)
var userDao = new(database.UserDao)

func jwtAbort(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, Message{
		Message: cmnCtrl.Capitalize(err.Error()),
	})
	c.Abort()
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// No Head
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			jwtAbort(c, AuthorizationException)
			return
		}

		// No Token Magic
		parts := strings.Split(authHeader, " ")
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			jwtAbort(c, AuthorizationException)
			return
		}

		// Token Parse Err
		claims, err := passUtil.ParseToken(parts[1])
		if err != nil {
			if strings.Index(err.Error(), "token is expired by") != -1 {
				// Token Expired
				jwtAbort(c, TokenExpiredException)
			} else {
				// Other Error
				// Signature is invalid
				// illegal base64 data at input byte

				jwtAbort(c, AuthorizationException)
			}
			return
		}

		// No User
		user, ok := userDao.QueryUser(claims.UserID)
		if !ok {
			jwtAbort(c, TokenInvalidException)
		}

		c.Set("user", *user)
		c.Next()
	}
}
