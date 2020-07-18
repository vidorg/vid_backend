package dto

import (
	"github.com/Aoi-hosizora/ahlib/xentity"
	"github.com/vidorg/vid_backend/src/model/po"
)

// @Model         _PolicyDto
// @Description   policy response
// @Property      role   string true "策略角色 (sub)"
// @Property      path   string true "策略路径 (obj)"
// @Property      method string true "策略方法 (act)"
type PolicyDto struct {
	Role   string `json:"role"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

func BuildPolicyDtos(policies []*po.Policy) []*PolicyDto {
	return xentity.MustMapSlice(policies, &PolicyDto{}).([]*PolicyDto)
}
