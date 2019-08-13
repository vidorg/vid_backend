package controllers

import (
	"fmt"
	"net/http"
	. "vid/exceptions"

	"github.com/gin-gonic/gin"
)

func NoRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"msg": fmt.Sprintf("Route %s %s is not found.", c.Request.Method, c.Request.URL.Path),
	})
	return
}

// func NoMethod(c *gin.Context) {
// 	c.JSON(http.StatusMethodNotAllowed, gin.H{
// 		"msg": fmt.Sprintf("Method %s is not allowed.", c.Request.Method),
// 	})
// 	return
// }

func errRet(err interface{}, c *gin.Context) {
	switch err.(type) {

	case *TestError:
		c.JSON(http.StatusBadRequest, err.(*TestError).Info())
		break

	case error:
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.(error).Error(),
		})
		break
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		break
	}
}
