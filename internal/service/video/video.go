package video

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/internal/model"
	"github.com/vidorg/vid_backend/internal/serializer"
	"github.com/vidorg/vid_backend/pkg/orm"
)

type GetVideoListService struct {
	CategoryID *int   `form:"category_id" json:"category_id" query:"category_id"`
	ChannelID  *int64 `form:"channel_id" json:"channel_id" query:"channel_id"`
	Page       int    `form:"page" json:"page" query:"page"`
	Limit      int    `form:"limit" json:"limit" query:"limit"`
}

func (g *GetVideoListService) GetVideoList(c *gin.Context) *serializer.Response {

	var videos []*model.Video
	var total int64

	tx := orm.Pagination(orm.DB().Model(&model.Video{}), g.Page, g.Limit).
		Preload("Author")

	if g.CategoryID == nil {
		tx.Count(&total).
			Find(&videos)
		return &serializer.Response{
			Code: 200,
			Msg:  "success",
			Data: &serializer.DataList{
				Total: total,
				Page:  g.Page,
				Limit: g.Limit,
				Items: videos,
			},
		}
	}
	if g.ChannelID != nil {
		tx.Where("channel_id = ?", g.ChannelID).
			Count(&total).
			Find(&videos)
		return &serializer.Response{
			Code: 200,
			Msg:  "success",
			Data: &serializer.DataList{
				Total: total,
				Page:  g.Page,
				Limit: g.Limit,
				Items: videos,
			},
		}
	}
	tx.Where("category_id = ?", g.CategoryID).
		Count(&total).
		Find(&videos)
	return &serializer.Response{
		Code: 200,
		Msg:  "success",
		Data: &serializer.DataList{
			Total: total,
			Page:  g.Page,
			Limit: g.Limit,
			Items: videos,
		},
	}
}
