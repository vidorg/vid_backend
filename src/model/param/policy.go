package param

import (
	"github.com/Aoi-hosizora/goapidoc"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("PolicyParam", "权限策略参数").
			WithProperties(
				goapidoc.NewProperty("role", "string", true, "角色"),
				goapidoc.NewProperty("path", "string", true, "路径"),
				goapidoc.NewProperty("method", "string", true, "方法"),
			),

		goapidoc.NewDefinition("RoleParam", "修改角色请求参数").
			WithProperties(
				goapidoc.NewProperty("role", "string", true, "角色"),
			),
	)
}

// @Property      method string true "policy method"
type PolicyParam struct {
	Role   string `json:"role"   form:"role"   binding:"required"`
	Path   string `json:"path"   form:"path"   binding:"required"`
	Method string `json:"method" form:"method" binding:"required"`
}

type RoleParam struct {
	Role string `json:"role" form:"role" binding:"required"`
}
