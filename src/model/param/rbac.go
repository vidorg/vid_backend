package param

import (
	"github.com/Aoi-hosizora/goapidoc"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("ChangeUserRoleParam", "change role param").
			Properties(
				goapidoc.NewProperty("sub", "string", true, "new subject"),
			),

		goapidoc.NewDefinition("RbacSubjectParam", "insert / delete rbac subject param").
			Properties(
				goapidoc.NewProperty("new", "string", true, "new subject"),
				goapidoc.NewProperty("from", "string", true, "subject inherited from"),
			),

		goapidoc.NewDefinition("RbacPolicyParam", "insert / delete rbac policy param").
			Properties(
				goapidoc.NewProperty("sub", "string", true, "new subject"),
				goapidoc.NewProperty("obj", "string", true, "new object"),
				goapidoc.NewProperty("act", "string", true, "new action"),
			),
	)
}

type ChangeUserRoleParam struct {
	Sub string `json:"sub" form:"sub" binding:"required"` // new subject
}

type RbacSubjectParam struct {
	New  string `json:"new"  form:"new"  binding:"required"` // new subject
	From string `json:"from" form:"from" binding:"required"` // subject inherited from
}

type RbacPolicyParam struct {
	Sub string `json:"sub" form:"sub" binding:"required"` // sub
	Obj string `json:"obj" form:"obj" binding:"required"` // obj
	Act string `json:"act" form:"act" binding:"required"` // act
}
