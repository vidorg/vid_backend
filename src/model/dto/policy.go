package dto

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
