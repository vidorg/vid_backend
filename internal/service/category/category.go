package category

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/internal/model"
	"github.com/vidorg/vid_backend/internal/serializer"
	"github.com/vidorg/vid_backend/pkg/orm"
)

type GetCategoriesListService struct {
	CategoryID int `form:"id" json:"category_id"`
}

func (g *GetCategoriesListService) GetCategoriesList(c *gin.Context) *serializer.Response {
	var categories []*model.Category

	// 查指定分类ID
	if g.CategoryID != 0 {
		orm.DB().Where("id = ?", g.CategoryID).Preload("Categories").
			Preload("Categories.Categories").Find(&categories)
		return &serializer.Response{
			Code: 200,
			Msg:  "success",
			Data: categories,
		}
	}

	// 查所有
	orm.DB().Where("category_id = ?", g.CategoryID).Preload("Categories").
		Preload("Categories.Categories").Find(&categories)

	return &serializer.Response{
		Code: 200,
		Msg:  "success",
		Data: categories,
	}
}
