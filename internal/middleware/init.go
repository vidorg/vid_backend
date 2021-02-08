package middleware

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	r.Use(Recover(), Cors(), Logger())
}
