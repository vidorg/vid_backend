package controllers

import (
	"net/http"
	"vid/exceptions"

	"github.com/gin-gonic/gin"
)

type TestCtrl struct{}

func (u *TestCtrl) Test(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errRet(err, c)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"msg": "Test",
	})
}

func (u *TestCtrl) Error(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errRet(err, c)
		}
	}()
	panic(exceptions.NewTestError("Test error", "Error detail"))
}
