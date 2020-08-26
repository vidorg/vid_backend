package dto

import (
	"github.com/Aoi-hosizora/ahlib/xentity"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/vidorg/vid_backend/src/model/po"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("PolicyDto", "权限策略信息").
			Properties(
				goapidoc.NewProperty("role", "string", true, "角色 (sub)"),
				goapidoc.NewProperty("path", "string", true, "路径 (obj)"),
				goapidoc.NewProperty("method", "string", true, "方法 (act)"),
			),
	)
}

type PolicyDto struct {
	Role   string `json:"role"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

func BuildPolicyDtos(policies []*po.Policy) []*PolicyDto {
	return xentity.MustMapSlice(policies, &PolicyDto{}).([]*PolicyDto)
}
