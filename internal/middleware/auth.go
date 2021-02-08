package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/internal/model"
	"github.com/vidorg/vid_backend/internal/serializer"
	"github.com/vidorg/vid_backend/pkg/jwt"
	"net/http"
)

// Auth
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//token,err := c.Cookie("token")
		//if errors.Is(err,http.ErrNoCookie){
		//	c.AbortWithStatusJSON(403,&serializer.Response{
		//		Code: 403,
		//		Msg:  "未登录",
		//		Data: nil,
		//	})
		//}
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, &serializer.Response{
				Code: 403,
				Msg:  "未登录",
				Data: nil,
			})
			return
		}
		userClaims, err := jwt.ParseToken([]byte(token))
		if err != nil || userClaims.UID == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, &serializer.Response{
				Code: 403,
				Msg:  "token失效",
				Data: err,
			})
			return
		}
		c.Set("user_id", userClaims.UID)
		user, err := model.GetUser(userClaims.UID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, &serializer.Response{
				Code: 403,
				Msg:  "没有找到该用户",
				Data: err,
			})
		}
		c.Set("user", user)
		c.Next()
	}
}
