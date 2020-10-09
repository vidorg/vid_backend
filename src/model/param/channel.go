package param

import (
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/vidorg/vid_backend/src/model/po"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("InsertChannelParam", "insert channel parameter").
			Properties(
				goapidoc.NewProperty("name", "string", true, "channel name"),
				goapidoc.NewProperty("description", "string", true, "channel description"),
				goapidoc.NewProperty("cover_url", "string", true, "channel cover url").Example("https://aaa.bbb.ccc"),
			),

		goapidoc.NewDefinition("UpdateChannelParam", "update channel parameter").
			Properties(
				goapidoc.NewProperty("name", "string", true, "channel name"),
				goapidoc.NewProperty("description", "string", true, "channel description"),
				goapidoc.NewProperty("cover_url", "string", true, "channel cover url").Example("https://aaa.bbb.ccc"),
			),
	)
}

type InsertChannelParam struct {
	Name        string `json:"name"        form:"name"        binding:"required,l_title"`       // channel name
	Description string `json:"description" form:"description" binding:"required,l_description"` // channel description
	CoverUrl    string `json:"cover_url"   form:"cover_url"   binding:"required,url"`           // channel cover url (oss)
}

func (i *InsertChannelParam) ToChannelPo() *po.Channel {
	return &po.Channel{
		Name:        i.Name,
		Description: i.Description,
		CoverUrl:    i.CoverUrl,
	}
}

type UpdateChannelParam struct {
	Name        string `json:"name"        form:"name"        binding:"required,l_title"`       // channel name
	Description string `json:"description" form:"description" binding:"required,l_description"` // channel description
	CoverUrl    string `json:"cover_url"   form:"cover_url"   binding:"required,url"`           // channel cover url (oss)
}

func (u *UpdateChannelParam) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":        u.Name,
		"description": u.Description,
		"cover_url":   u.CoverUrl,
	}
}
