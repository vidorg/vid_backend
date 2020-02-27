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
	"github.com/vidorg/vid_backend/src/model/po"
	"log"
)

type UserController struct {
	Config     *config.ServerConfig   `di:"~"`
	JwtService *middleware.JwtService `di:"~"`
	UserDao    *dao.UserDao           `di:"~"`
	VideoDao   *dao.VideoDao          `di:"~"`
	SubDao     *dao.SubDao            `di:"~"`
	SearchDao  *dao.SearchDao         `di:"~"`
	Mapper     *xmapper.EntityMapper  `di:"~"`
}

func NewUserController(dic *xdi.DiContainer) *UserController {
	ctrl := &UserController{}
	if !dic.Inject(ctrl) {
		log.Fatalln("Inject failed")
	}
	return ctrl
}

// @Router              /v1/user [GET]
// @Security            Jwt
// @Template            Admin Auth Order Page
// @Summary             查询所有用户
// @Description         管理员权限，此处可见用户手机号码
// @Tag                 User
// @Tag                 Administration
// @ResponseModel 200   #Result<Page<UserDto>>
// @ResponseEx 200      ${resp_page_users}
func (u *UserController) QueryAllUsers(c *gin.Context) {
	page := param.BindQueryPage(c)
	order := param.BindQueryOrder(c)
	users, count := u.UserDao.QueryAll(page, order)

	retDto := xcondition.First(u.Mapper.Map([]*dto.UserDto{}, users, dto.UserDtoAdminMapOption())).([]*dto.UserDto)
	result.Ok().SetPage(count, page, retDto).JSON(c)
}

// @Router              /v1/user/{uid} [GET]
// @Template            ParamA
// @Summary             查询用户
// @Description         此处用户本人可见手机号码，管理员不受限制
// @Tag                 User
// @Param               uid path integer true "用户id"
// @ResponseDesc 404    "user not found"
// @ResponseModel 200   #Result<UserExtraDto>
// @ResponseEx 200      ${resp_user_info}
func (u *UserController) QueryUser(c *gin.Context) {
	uid, ok := param.BindRouteId(c, "uid")
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}

	user := u.UserDao.QueryByUid(uid)
	if user == nil {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	}
	subscribingCnt, subscriberCnt, _ := u.SubDao.QueryCountByUid(user.Uid)
	videoCnt, _ := u.VideoDao.QueryCountByUid(user.Uid)
	extraInfo := &dto.UserExtraDto{
		SubscribingCount: subscribingCnt,
		SubscriberCount:  subscriberCnt,
		VideoCount:       videoCnt,
	}

	authUser := u.JwtService.GetContextUser(c)
	retDto := xcondition.First(u.Mapper.Map(&dto.UserDto{}, user, dto.UserDtoUserMapOption(authUser))).(*dto.UserDto)
	result.Ok().
		PutData("user", retDto).
		PutData("extra", extraInfo).JSON(c)
}

// @Router              /v1/user [PUT]
// @Security            Jwt
// @Template            Auth Param
// @Summary             更新用户
// @Tag                 User
// @Param               param body #UserParam true "请求参数"
// @ResponseDesc 400    "username has been used"
// @ResponseDesc 404    "user not found"
// @ResponseDesc 500    "user update failed"
// @ResponseModel 200   #Result<UserDto>
// @ResponseEx 200      ${resp_user}
//
// @Router              /v1/user/admin/{uid} [PUT]
// @Security            Jwt
// @Template            Admin Auth Param
// @Summary             更新用户
// @Description         管理员权限
// @Tag                 User
// @Tag                 Administration
// @Param               uid   path integer    true "用户id"
// @Param               param body #UserParam true "请求参数"
// @ResponseDesc 400    "username has been used"
// @ResponseDesc 404    "user not found"
// @ResponseDesc 500    "user update failed"
// @ResponseModel 200   #Result<UserDto>
// @ResponseEx 200      ${resp_user}
func (u *UserController) UpdateUser(isSpec bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := &po.User{}
		if !isSpec {
			user = u.JwtService.GetContextUser(c)
		} else {
			uid, ok := param.BindRouteId(c, "uid")
			if !ok {
				result.Error(exception.RequestParamError).JSON(c)
				return
			}
			user = u.UserDao.QueryByUid(uid)
			if user == nil {
				result.Error(exception.UserNotFoundError).JSON(c)
				return
			}
		}
		// Update
		userParam := &param.UserParam{}
		if err := c.ShouldBind(userParam); err != nil {
			result.Error(exception.WrapValidationError(err)).PutData("error", err.Error()).JSON(c)
			return
		}

		_ = u.Mapper.MapProp(userParam, user)
		status := u.UserDao.Update(user)
		if status == database.DbNotFound {
			result.Error(exception.UserNotFoundError).JSON(c)
			return
		} else if status == database.DbExisted {
			result.Error(exception.UsernameUsedError).JSON(c)
			return
		} else if status == database.DbFailed {
			result.Error(exception.UserUpdateError).JSON(c)
			return
		}

		retDto := xcondition.First(u.Mapper.Map(&dto.UserDto{}, user, dto.UserDtoAdminMapOption())).(*dto.UserDto)
		result.Ok().SetData(retDto).JSON(c)
	}
}

// @Router              /v1/user [DELETE]
// @Security            Jwt
// @Template            Auth
// @Summary             删除用户
// @Tag                 User
// @ResponseDesc 404    "user not found"
// @ResponseDesc 500    "user delete failed"
// @ResponseModel 200   #Result
// @ResponseEx 200      ${resp_success}
//
// @Router              /v1/user/admin/{uid} [DELETE]
// @Security            Jwt
// @Template            Admin Auth ParamA
// @Summary             删除用户
// @Description         管理员权限
// @Tag                 User
// @Tag                 Administration
// @Param               uid path integer true "用户id"
// @ResponseDesc 404    "user not found"
// @ResponseDesc 500    "user delete failed"
// @ResponseModel 200   #Result
// @ResponseEx 200      ${resp_success}
func (u *UserController) DeleteUser(isSpec bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		var uid int32
		if !isSpec {
			uid = u.JwtService.GetContextUser(c).Uid
		} else {
			var ok bool
			uid, ok = param.BindRouteId(c, "uid")
			if !ok {
				result.Error(exception.RequestParamError).JSON(c)
				return
			}
		}
		// Delete
		status := u.UserDao.Delete(uid)
		if status == database.DbNotFound {
			result.Error(exception.UserNotFoundError).JSON(c)
			return
		} else if status == database.DbFailed {
			result.Error(exception.UserDeleteError).JSON(c)
			return
		}
		result.Ok().JSON(c)
	}
}
