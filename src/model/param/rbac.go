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
				goapidoc.NewProperty("sub", "string", true, "new subject"),
				goapidoc.NewProperty("sub2", "string", true, "subject inherited from"),
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
	Sub  string `json:"sub"  form:"sub"  binding:"required"` // sub
	Sub2 string `json:"sub2" form:"sub2" binding:"required"` // sub2
}

type RbacPolicyParam struct {
	Sub string `json:"sub" form:"sub" binding:"required"` // sub
	Obj string `json:"obj" form:"obj" binding:"required"` // obj
	Act string `json:"act" form:"act" binding:"required"` // act
}
