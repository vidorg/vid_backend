package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/internal/serializer"
	"github.com/vidorg/vid_backend/internal/service/channel"
)

func GetChannelList(c *gin.Context) {
	service := &channel.GetChannelListService{}
	if err := c.ShouldBindQuery(service); err != nil {
		c.JSON(200, serializer.ParamErr("param err,", err))
	} else {
		res := service.GetChannelList(c)
		c.JSON(200, res)
	}
}
