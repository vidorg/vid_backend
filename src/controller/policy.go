package controller

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xentity"
	"github.com/Aoi-hosizora/ahlib/xslice"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/service"
)

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

// @Router              /v1/policy [GET]
// @Summary             查询所有策略
// @Tag                 Policy
// @Security            Jwt
// @Template            Page
// @ResponseModel 200   #Result<Page<PolicyDto>>
func (r *PolicyController) Query(c *gin.Context) {
	page := param.BindPage(c, r.config)
	total, policies := r.casbinService.GetPolicies(page.Limit, page.Page)

	policiesDto := xentity.MustMapSlice(xslice.Sti(policies), &dto.PolicyDto{}).([]*dto.PolicyDto)
	result.Ok().SetPage(total, page.Page, page.Limit, policiesDto).JSON(c)
}

// @Router              /v1/policy/role/{uid} [PUT]
// @Summary             修改用户权限
// @Tag                 Policy
// @Security            Jwt
// @Param               uid   path integer    true "用户id"
// @Param               param body #RoleParam true "请求参数"
// @ResponseModel 200   #Result
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
	if status == database.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Error(exception.UserUpdateError).JSON(c)
		return
	}

	result.Ok().JSON(c)
}

// @Router              /v1/policy [POST]
// @Summary             插入策略
// @Tag                 Policy
// @Security            Jwt
// @Param               param body #PolicyParam true "请求参数"
// @ResponseModel 200   #Result
func (r *PolicyController) Insert(c *gin.Context) {
	policyParam := &param.PolicyParam{}
	if err := c.ShouldBind(policyParam); err != nil {
		result.Error(exception.RequestParamError).JSON(c)
	}

	status := r.casbinService.AddPolicy(policyParam.Role, policyParam.Path, policyParam.Method)
	if status == database.DbExisted {
		result.Error(exception.PolicyExistedError).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Error(exception.PolicyInsertError).JSON(c)
		return
	}

	result.Ok().JSON(c)
}

// @Router              /v1/policy [DELETE]
// @Summary             删除策略
// @Security            Jwt
// @Tag                 Policy
// @Param               param body #PolicyParam true "请求参数"
// @ResponseModel 200   #Result
func (r *PolicyController) Delete(c *gin.Context) {
	policyParam := &param.PolicyParam{}
	if err := c.ShouldBind(policyParam); err != nil {
		result.Error(exception.RequestParamError).JSON(c)
	}

	status := r.casbinService.DeletePolicy(policyParam.Role, policyParam.Path, policyParam.Method)
	if status == database.DbNotFound {
		result.Error(exception.PolicyNotFountError).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Error(exception.PolicyDeleteError).JSON(c)
		return
	}

	result.Ok().JSON(c)
}
