package param

import (
	"github.com/Aoi-hosizora/goapidoc"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("SubscribeParam", "关注请求参数").
			WithProperties(
				goapidoc.NewProperty("to", "integer#int32", true, "用户id"),
			),
	)
}

type SubscribeParam struct {
	To int32 `form:"to" json:"to" binding:"required,min=1"`
}
