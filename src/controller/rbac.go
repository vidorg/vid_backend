package controller

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/service"
)

func init() {
	goapidoc.AddRoutePaths(
		goapidoc.NewRoutePath("GET", "/v1/rbac/rule", "Get rbac rules").
			Tags("Rbac", "Administration").
			Securities("Jwt").
			Params(param.ADPage, param.ADLimit).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<RbacRuleDto>>")),

		goapidoc.NewRoutePath("PUT", "/v1/rbac/user/{uid}", "Change user role").
			Tags("Rbac", "Administration").
			Securities("Jwt").
			Params(
				goapidoc.NewPathParam("uid", "integer#int64", true, "user id"),
				goapidoc.NewBodyParam("body", "ChangeUserRoleParam", true, "change role param"),
			).
			Responses(goapidoc.NewResponse(200, "Result")),

		goapidoc.NewRoutePath("POST", "/v1/rbac/subject", "Insert subject").
			Tags("Rbac", "Administration").
			Securities("Jwt").
			Params(goapidoc.NewBodyParam("body", "RbacSubjectParam", true, "insert subject param")).
			Responses(goapidoc.NewResponse(201, "Result")),

		goapidoc.NewRoutePath("POST", "/v1/rbac/policy", "Insert policy").
			Tags("Rbac", "Administration").
			Securities("Jwt").
			Params(goapidoc.NewBodyParam("body", "RbacPolicyParam", true, "insert policy param")).
			Responses(goapidoc.NewResponse(201, "Result")),

		goapidoc.NewRoutePath("DELETE", "/v1/rbac/subject", "Delete subject").
			Tags("Rbac", "Administration").
			Securities("Jwt").
			Params(goapidoc.NewBodyParam("body", "RbacSubjectParam", true, "delete subject param")).
			Responses(goapidoc.NewResponse(200, "Result")),

		goapidoc.NewRoutePath("DELETE", "/v1/rbac/policy", "Delete policy").
			Tags("Rbac", "Administration").
			Securities("Jwt").
			Params(goapidoc.NewBodyParam("body", "RbacPolicyParam", true, "delete policy param")).
			Responses(goapidoc.NewResponse(200, "Result")),
	)
}

type RbacController struct {
	config        *config.Config
	userService   *service.UserService
	jwtService    *service.JwtService
	casbinService *service.CasbinService
}

func NewRbacController() *RbacController {
	return &RbacController{
		config:        xdi.GetByNameForce(sn.SConfig).(*config.Config),
		userService:   xdi.GetByNameForce(sn.SUserService).(*service.UserService),
		jwtService:    xdi.GetByNameForce(sn.SJwtService).(*service.JwtService),
		casbinService: xdi.GetByNameForce(sn.SCasbinService).(*service.CasbinService),
	}
}

// GET /v1/rbac/rule
func (r *RbacController) GetRules(c *gin.Context) *result.Result {
	pp := param.BindPage(c, r.config)
	rules, total, err := r.casbinService.GetAllRules(pp)
	if err != nil {
		return result.Error(exception.QueryRbacRuleError).SetError(err, c)
	}

	res := dto.BuildRbacRuleDtos(rules)
	return result.Ok().SetPage(pp.Page, pp.Limit, total, res)
}

// PUT /v1/rbac/user/:uid
func (r *RbacController) ChangeRole(c *gin.Context) *result.Result {
	uid, err := param.BindRouteId(c, "uid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	user := r.jwtService.GetContextUser(c)
	if user == nil {
		return nil
	}

	pa := &param.ChangeUserRoleParam{}
	if err := c.ShouldBind(pa); err != nil {
		return result.Error(exception.WrapValidationError(err)).SetError(err, c)
	}

	if user.Uid == uid {
		return result.Error(exception.ChangeSelfRoleError)
	}

	status, err := r.userService.UpdateRole(uid, pa.Sub)
	if status == xstatus.DbNotFound {
		return result.Error(exception.UserNotFoundError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.ChangeRoleError).SetError(err, c)
	}

	return result.Ok()
}

// POST /v1/rbac/subject
func (r *RbacController) InsertSubject(c *gin.Context) *result.Result {
	pa := &param.RbacSubjectParam{}
	if err := c.ShouldBind(pa); err != nil {
		return result.Error(exception.WrapValidationError(err)).SetError(err, c)
	}

	status, err := r.casbinService.AddSubject(pa.New, pa.From)
	if status == xstatus.DbExisted {
		return result.Error(exception.RbacSubjectExistedError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.RbacSubjectInsertError).SetError(err, c)
	}

	return result.Created()
}

// POST /v1/rbac/policy
func (r *RbacController) InsertPolicy(c *gin.Context) *result.Result {
	pa := &param.RbacPolicyParam{}
	if err := c.ShouldBind(pa); err != nil {
		return result.Error(exception.WrapValidationError(err)).SetError(err, c)
	}

	status, err := r.casbinService.AddPolicy(pa.Sub, pa.Obj, pa.Act)
	if status == xstatus.DbExisted {
		return result.Error(exception.RbacPolicyExistedError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.RbacPolicyInsertError).SetError(err, c)
	}

	return result.Created()

}

// DELETE /v1/rbac/subject
func (r *RbacController) DeleteSubject(c *gin.Context) *result.Result {
	pa := &param.RbacSubjectParam{}
	if err := c.ShouldBind(pa); err != nil {
		return result.Error(exception.WrapValidationError(err)).SetError(err, c)
	}

	status, err := r.casbinService.RemoveSubject(pa.New, pa.From)
	if status == xstatus.DbNotFound {
		return result.Error(exception.RbacSubjectNotFoundError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.RbacSubjectDeleteError).SetError(err, c)
	}

	return result.Ok()
}

// DELETE /v1/rbac/policy
func (r *RbacController) DeletePolicy(c *gin.Context) *result.Result {
	pa := &param.RbacPolicyParam{}
	if err := c.ShouldBind(pa); err != nil {
		return result.Error(exception.WrapValidationError(err)).SetError(err, c)
	}

	status, err := r.casbinService.RemovePolicy(pa.Sub, pa.Obj, pa.Act)
	if status == xstatus.DbNotFound {
		return result.Error(exception.RbacPolicyNotFoundError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.RbacPolicyDeleteError).SetError(err, c)
	}

	return result.Ok()
}
