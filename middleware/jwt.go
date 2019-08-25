package middleware

import (
	"net/http"
	"strings"
	"time"
	"vid/database"
	. "vid/models"
	"vid/utils"

	"github.com/gin-gonic/gin"
)

var passUtil = new(utils.PassUtil)
var userDao = new(database.UserDao)

func jwtAbort(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, Message{
		Message: msg,
	})
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// No Head
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			jwtAbort(c, "Authorization fail")
			return
		}

		// No Eoken Magic
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			jwtAbort(c, "Authorization fail")
			return
		}

		// Token Parse Err
		claims, err := passUtil.ParseToken(parts[1])
		if err != nil {
			jwtAbort(c, err.Error())
			return
		}

		// Token Expire
		if time.Now().Unix() > claims.ExpiresAt {
			jwtAbort(c, "Token has expired")
			return
		}

		// No User
		user, ok := userDao.QueryUser(claims.UserID)
		if !ok {
			jwtAbort(c, "Token invalid")
		}

		c.Set("user", *user)
		c.Next()
	}
}
