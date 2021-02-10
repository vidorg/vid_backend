package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/internal/serializer"
	"github.com/vidorg/vid_backend/internal/service/video"
)

func GetVideoList(c *gin.Context) {
	service := &video.GetVideoListService{}
	if err := c.ShouldBindQuery(service); err != nil {
		c.JSON(200, serializer.ParamErr("param err,", err))
	} else {
		res := service.GetVideoList(c)
		c.JSON(200, res)
	}
}
