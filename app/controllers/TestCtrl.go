package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type testCtrl struct{}

var TestCtrl = new(testCtrl)

func (u *testCtrl) Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "Test",
	})
}

func (u *testCtrl) AuthTest(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"msg": user,
	})
}
