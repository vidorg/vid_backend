package controller

import (
	"github.com/Aoi-hosizora/ahlib-web/xstatus"
	"github.com/Aoi-hosizora/ahlib/xdi"
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
	goapidoc.AddPaths(
		goapidoc.NewPath("GET", "/v1/policy", "查询所有策略").
			WithTags("Policy", "Administration").
			WithSecurities("Jwt").
			WithResponses(goapidoc.NewResponse(200).WithType("_Result<_Page<PolicyDto>>")),

		goapidoc.NewPath("PUT", "/v1/policy/{uid}/role", "修改用户角色").
			WithTags("Policy", "Administration").
			WithSecurities("Jwt").
			WithParams(
				goapidoc.NewPathParam("uid", "integer#int32", true, "用户id"),
				goapidoc.NewBodyParam("param", "RoleParam", true, "修改角色请求参数"),
			).
			WithResponses(goapidoc.NewResponse(200).WithType("Result")),

		goapidoc.NewPath("POST", "/v1/policy", "新建策略").
			WithTags("Policy", "Administration").
			WithSecurities("Jwt").
			WithParams(goapidoc.NewBodyParam("param", "PolicyParam", true, "策略请求参数")).
			WithResponses(goapidoc.NewResponse(200).WithType("Result")),

		goapidoc.NewPath("DELETE", "/v1/policy", "删除策略").
			WithTags("Policy", "Administration").
			WithSecurities("Jwt").
			WithParams(goapidoc.NewBodyParam("param", "PolicyParam", true, "策略请求参数")).
			WithResponses(goapidoc.NewResponse(200).WithType("Result")),
	)
}

type PolicyController struct {
	config        *config.Config
	userService   *service.UserService
	casbinService *service.CasbinService
}

func NewPolicyController() *PolicyController {
	return &PolicyController{
		config:        xdi.GetByNameForce(sn.SConfig).(*config.Config),
		userService:   xdi.GetByNameForce(sn.SUserService).(*service.UserService),
		casbinService: xdi.GetByNameForce(sn.SCasbinService).(*service.CasbinService),
	}
}

// GET /v1/policy
func (r *PolicyController) Query(c *gin.Context) {
	pp := param.BindPage(c, r.config)
	total, policies := r.casbinService.GetPolicies(pp.Limit, pp.Page)

	ret := dto.BuildPolicyDtos(policies)
	result.Ok().SetPage(pp.Page, pp.Limit, total, ret).JSON(c)
}

// PUT /v1/policy/:uid/role
func (r *PolicyController) SetRole(c *gin.Context) {
	uid, ok := param.BindRouteId(c, "uid")
	roleParam := &param.RoleParam{}
	if err := c.ShouldBind(roleParam); err != nil || !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}

	user := r.userService.QueryByUid(uid)
	if user == nil {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	}
	if user.Role == "root" {
		result.Error(exception.PolicySetRoleError).JSON(c)
		return
	}
	user.Role = roleParam.Role

	status := r.userService.Update(user)
	if status == xstatus.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	} else if status == xstatus.DbFailed {
		result.Error(exception.UserUpdateError).JSON(c)
		return
	}

	result.Ok().JSON(c)
}

// POST /v1/policy
func (r *PolicyController) Insert(c *gin.Context) {
	policyParam := &param.PolicyParam{}
	if err := c.ShouldBind(policyParam); err != nil {
		result.Error(exception.RequestParamError).JSON(c)
	}

	status := r.casbinService.AddPolicy(policyParam.Role, policyParam.Path, policyParam.Method)
	if status == xstatus.DbExisted {
		result.Error(exception.PolicyExistedError).JSON(c)
		return
	} else if status == xstatus.DbFailed {
		result.Error(exception.PolicyInsertError).JSON(c)
		return
	}

	result.Ok().JSON(c)
}

// DELETE /v1/policy
func (r *PolicyController) Delete(c *gin.Context) {
	policyParam := &param.PolicyParam{}
	if err := c.ShouldBind(policyParam); err != nil {
		result.Error(exception.RequestParamError).JSON(c)
	}

	status := r.casbinService.DeletePolicy(policyParam.Role, policyParam.Path, policyParam.Method)
	if status == xstatus.DbNotFound {
		result.Error(exception.PolicyNotFountError).JSON(c)
		return
	} else if status == xstatus.DbFailed {
		result.Error(exception.PolicyDeleteError).JSON(c)
		return
	}

	result.Ok().JSON(c)
}
