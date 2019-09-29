package middleware

import (
	"net/http"
	"strings"

	. "vid/database"
	. "vid/exceptions"
	. "vid/models"
	. "vid/models/resp"
	. "vid/utils"

	"github.com/gin-gonic/gin"
)

func jwtAbort(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, Message{
		Message: CmnUtil.Capitalize(err.Error()),
	})
	c.Abort()
}

func JWTCheck(authHeader string) (*User, error) {

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
	claims, err := PassUtil.ParseToken(parts[1])
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
	user, ok := UserDao.QueryUserByUid(claims.UserID)
	if !ok {
		return nil, TokenInvalidException
	}

	return user, nil
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		user, err := JWTCheck(authHeader)
		if err != nil {
			jwtAbort(c, err)
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
