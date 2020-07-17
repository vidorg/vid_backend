package param

// @Model         PolicyParam
// @Description   权限策略参数
// @Property      role   string true "policy role"
// @Property      path   string true "policy path"
// @Property      method string true "policy method"
type PolicyParam struct {
	Role   string `json:"role"   form:"role"   binding:"required"`
	Path   string `json:"path"   form:"path"   binding:"required"`
	Method string `json:"method" form:"method" binding:"required"`
}

// @Model         RoleParam
// @Description   修改权限参数
// @Property      role   string true "policy role"
type RoleParam struct {
	Role string `json:"role" form:"role" binding:"required"`
}
