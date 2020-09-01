package param

import (
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/vidorg/vid_backend/src/model/po"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("InsertVideoParam", "insert video parameter").
			Properties(
				goapidoc.NewProperty("title", "string", true, "video title"),
				goapidoc.NewProperty("description", "string", true, "video description"),
				goapidoc.NewProperty("video_url", "string", true, "video source url"),
				goapidoc.NewProperty("cover_url", "string", true, "video cover url"),
			),
	)
}

type InsertVideoParam struct {
	Title       string `json:"title"       form:"title"       binding:"required,l_title"`       // video title
	Description string `json:"description" form:"description" binding:"required,l_description"` // video description
	VideoUrl    string `json:"video_url"   form:"video_url"   binding:"required,url"`           // video source url (oss)
	CoverUrl    string `json:"cover_url"   form:"cover_url"   binding:"required,url"`           // video cover url (oss)
}

func (i *InsertVideoParam) ToPo() *po.Video {
	return &po.Video{
		Title:       i.Title,
		Description: i.Description,
		VideoUrl:    i.VideoUrl,
		CoverUrl:    i.CoverUrl,
	}
}

type UpdateVideoParam struct {
	Title       string `json:"title"       form:"title"       binding:"required,l_title"`       // video title
	Description string `json:"description" form:"description" binding:"required,l_description"` // video description
	VideoUrl    string `json:"video_url"   form:"video_url"   binding:"required,url"`           // video source url (oss)
	CoverUrl    string `json:"cover_url"   form:"cover_url"   binding:"required,url"`           // video cover url (oss)
}

func (u *UpdateVideoParam) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"title":       u.Title,
		"description": u.Description,
		"video_url":   u.VideoUrl,
		"cover_url":   u.CoverUrl,
	}
}
