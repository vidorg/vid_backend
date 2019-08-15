package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestCtrl struct{}

func (u *TestCtrl) Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "Test",
	})
}
