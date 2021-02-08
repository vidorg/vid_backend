package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/internal/serializer"
	"github.com/vidorg/vid_backend/internal/service/user"
)

func UserLogin(c *gin.Context) {
	service := &user.LoginService{}
	if err := c.ShouldBind(service); err != nil {
		c.JSON(200, serializer.ParamErr("param err,", err))
	} else {
		res := service.Login()
		c.JSON(200, res)
	}
}
