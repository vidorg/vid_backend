package controller

import (
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/database/dao"
	"github.com/vidorg/vid_backend/src/middleware"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"net/http"
)

type SubController struct {
	Config     *config.ServerConfig   `di:"~"`
	JwtService *middleware.JwtService `di:"~"`
	UserDao    *dao.UserDao           `di:"~"`
	SubDao     *dao.SubDao            `di:"~"`
	Mapper     *xmapper.EntityMapper  `di:"~"`
}

func NewSubController(dic *xdi.DiContainer) *SubController {
	ctrl := &SubController{}
	if !dic.Inject(ctrl) {
		panic("Inject failed")
	}
	return ctrl
}

// @Router              /v1/user/{uid}/subscriber?page [GET]
// @Summary             查询粉丝
// @Description         查询用户所有粉丝，返回分页数据
// @Tag                 Subscribe
// @Param               uid path integer true "查询的用户id"
// @Param               page query integer false "分页"
// @ErrorCode           400 request param error
// @ErrorCode           404 user not found
/* @Response 200        ${resp_page_users} */
func (s *SubController) QuerySubscriberUsers(c *gin.Context) {
	uid, ok1 := param.BindRouteId(c, "uid")
	page, ok2 := param.BindQueryPage(c)
	if !ok1 || !ok2 {
		result.Result{}.Result(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}

	users, count, status := s.SubDao.QuerySubscriberUsers(uid, page)
	if status == database.DbNotFound {
		result.Result{}.Result(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
		return
	}

	retDto := xcondition.First(s.Mapper.Map([]*dto.UserDto{}, users)).([]*dto.UserDto)
	result.Result{}.Ok().SetPage(count, page, retDto).JSON(c)
}

// @Router              /v1/user/{uid}/subscribing?page [GET]
// @Summary             查询关注
// @Description         查询用户所有关注，返回分页数据
// @Tag                 Subscribe
// @Param               uid path integer true "查询的用户id"
// @Param               page query integer false "分页"
// @ErrorCode           400 request param error
// @ErrorCode           404 user not found
/* @Response 200        ${resp_page_users} */
func (s *SubController) QuerySubscribingUsers(c *gin.Context) {
	uid, ok1 := param.BindRouteId(c, "uid")
	page, ok2 := param.BindQueryPage(c)
	if !ok1 || !ok2 {
		result.Result{}.Result(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}

	users, count, status := s.SubDao.QuerySubscribingUsers(uid, page)
	if status == database.DbNotFound {
		result.Result{}.Result(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
		return
	}

	retDto := xcondition.First(s.Mapper.Map([]*dto.UserDto{}, users)).([]*dto.UserDto)
	result.Result{}.Ok().SetPage(count, page, retDto).JSON(c)
}

// @Router              /v1/user/subscribing [PUT]
// @Security            Jwt
// @Template            Auth
// @Summary             关注
// @Description         关注某一用户
// @Tag                 Subscribe
// @Param               to formData integer true "关注用户id"
// @ErrorCode           400 request param error
// @ErrorCode           400 request format error
// @ErrorCode           400 subscribe oneself invalid
// @ErrorCode           404 user not found
// @ErrorCode           500 subscribe failed
/* @Response 200        ${resp_success} */
func (s *SubController) SubscribeUser(c *gin.Context) {
	authUser := s.JwtService.GetAuthUser(c)
	subParam := &param.SubParam{}
	if err := c.ShouldBind(subParam); err != nil {
		result.Result{}.Result(http.StatusBadRequest).SetMessage(exception.WrapValidationError(err).Error()).JSON(c)
		return
	}
	if authUser.Uid == subParam.Uid {
		result.Result{}.Result(http.StatusBadRequest).SetMessage(exception.SubscribeSelfError.Error()).JSON(c)
		return
	}

	status := s.SubDao.SubscribeUser(authUser.Uid, subParam.Uid)
	if status == database.DbNotFound {
		result.Result{}.Result(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Result{}.Error().SetMessage(exception.SubscribeError.Error()).JSON(c)
		return
	}

	result.Result{}.Ok().JSON(c)
}

// @Router              /v1/user/subscribing [DELETE]
// @Security            Jwt
// @Template            Auth
// @Summary             取消关注
// @Description         取消关注某一用户
// @Tag                 Subscribe
// @Param               to formData integer true "取消关注用户id"
// @ErrorCode           400 request param error
// @ErrorCode           400 request format error
// @ErrorCode           404 user not found
// @ErrorCode           500 unsubscribe failed
/* @Response 200        ${resp_success} */
func (s *SubController) UnSubscribeUser(c *gin.Context) {
	authUser := s.JwtService.GetAuthUser(c)
	subParam := &param.SubParam{}
	if err := c.ShouldBind(subParam); err != nil {
		result.Result{}.Result(http.StatusBadRequest).SetMessage(exception.WrapValidationError(err).Error()).JSON(c)
		return
	}

	status := s.SubDao.UnSubscribeUser(authUser.Uid, subParam.Uid)
	if status == database.DbNotFound {
		result.Result{}.Result(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Result{}.Error().SetMessage(exception.UnSubscribeError.Error()).JSON(c)
		return
	}

	result.Result{}.Ok().JSON(c)
}
