package channel

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/internal/model"
	"github.com/vidorg/vid_backend/internal/serializer"
	"github.com/vidorg/vid_backend/pkg/orm"
)

type GetChannelListService struct {
	CategoryID int `form:"id" json:"category_id"`
	Page       int `form:"page" json:"page" query:"page"`
	Limit      int `form:"limit" json:"limit" query:"limit"`
}

func (g *GetChannelListService) GetChannelList(c *gin.Context) *serializer.Response {

	var channels []*model.Channel
	var count int64

	orm.Pagination(orm.DB().Model(&model.Channel{}), g.Page, g.Limit).
		Count(&count).
		Find(&channels)

	return &serializer.Response{
		Code: 200,
		Msg:  "success",
		Data: &serializer.DataList{
			Total: count,
			Page:  g.Page,
			Limit: g.Limit,
			Items: channels,
		},
	}
}
