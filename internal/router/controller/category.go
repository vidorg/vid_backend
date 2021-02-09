package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/internal/serializer"
	"github.com/vidorg/vid_backend/internal/service/category"
)

func GetCategoryList(c *gin.Context) {
	service := &category.GetCategoriesListService{}
	if err := c.ShouldBindQuery(service); err != nil {
		c.JSON(200, serializer.ParamErr("param err,", err))
	} else {
		res := service.GetCategoriesList(c)
		c.JSON(200, res)
	}
}
