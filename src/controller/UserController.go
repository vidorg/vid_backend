package controller

import (
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xdatetime"
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/enum"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/database/dao"
	"github.com/vidorg/vid_backend/src/middleware"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/util"
	"net/http"
)

type UserController struct {
	Config     *config.ServerConfig   `di:"~"`
	JwtService *middleware.JwtService `di:"~"`
	UserDao    *dao.UserDao           `di:"~"`
	VideoDao   *dao.VideoDao          `di:"~"`
	SubDao     *dao.SubDao            `di:"~"`
	Mapper     *xmapper.EntityMapper  `di:"~"`
}

func NewUserController(dic *xdi.DiContainer) *UserController {
	ctrl := &UserController{}
	if !dic.Inject(ctrl) {
		panic("Inject failed")
	}
	return ctrl
}

// @Router              /v1/user?page [GET]
// @Security            Jwt
// @Template            Admin Auth Page
// @Summary             查询所有用户
// @Description         管理员权限，此处可见用户手机号码
// @Tag                 User
// @Tag                 Administration
// @ResponseModel 200   #UserDtoPageResult
// @Response 200        ${resp_page_users}
func (u *UserController) QueryAllUsers(c *gin.Context) {
	page := param.BindQueryPage(c)
	users, count := u.UserDao.QueryAll(page)

	// show all user's info
	retDto := xcondition.First(u.Mapper.Map([]*dto.UserDto{}, users, dto.UserDtoAdminMapOption())).([]*dto.UserDto)
	result.Result{}.Ok().SetPage(count, page, retDto).JSON(c)
}

// @Router              /v1/user/{uid} [GET]
// @Template            ParamA
// @Summary             查询用户
// @Description         此处可见用户手机号码
// @Tag                 User
// @Param               uid path integer true false "用户id"
// @ResponseDesc 404    "user not found"
// @ResponseModel 200   #UserExtraDtoResult
// @Response 200        ${resp_user_info}
func (u *UserController) QueryUser(c *gin.Context) {
	uid, ok := param.BindRouteId(c, "uid")
	if !ok {
		result.Result{}.Result(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
		return
	}

	user := u.UserDao.QueryByUid(uid)
	if user == nil {
		result.Result{}.Result(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
		return
	}
	subscribingCnt, subscriberCnt, _ := u.SubDao.QuerySubCnt(user.Uid)
	videoCnt, _ := u.VideoDao.QueryCount(user.Uid)
	extraInfo := &dto.UserExtraDto{
		SubscribingCount: subscribingCnt,
		SubscriberCount:  subscriberCnt,
		VideoCount:       videoCnt,
	}

	// need to squeeze out whether you can see the admin info
	authUser := u.JwtService.GetAuthUser(c)
	retDto := xcondition.First(u.Mapper.Map(&dto.UserDto{}, user, dto.UserDtoExtraMapOption(authUser))).(*dto.UserDto)
	result.Result{}.Ok().PutData("user", retDto).PutData("extra", extraInfo).JSON(c)
}

// @Router              /v1/user/ [PUT]
// @Security            Jwt
// @Template            Auth Param
// @Summary             更新用户
// @Tag                 User
// @Param               param body #UserParam true false "用户请求参数"
// @ResponseDesc 400    "username has been used"
// @ResponseDesc 404    "user not found"
// @ResponseDesc 500    "user update failed"
// @ResponseModel 200   #UserDtoResult
// @Response 200        ${resp_user}
//
// @Router              /v1/user/admin/{uid} [PUT]
// @Security            Jwt
// @Template            Admin Auth Param
// @Summary             更新用户
// @Description         管理员权限
// @Tag                 User
// @Tag                 Administration
// @Param               uid   path integer    true false "用户id"
// @Param               param body #UserParam true false "用户请求参数"
// @ResponseDesc 400    "username has been used"
// @ResponseDesc 404    "user not found"
// @ResponseDesc 500    "user update failed"
// @ResponseModel 200   #UserDtoResult
// @Response 200        ${resp_user}
func (u *UserController) UpdateUser(isExact bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := &po.User{}
		if !isExact {
			user = u.JwtService.GetAuthUser(c)
		} else {
			uid, ok := param.BindRouteId(c, "uid")
			if !ok {
				result.Result{}.Result(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
				return
			}
			user = u.UserDao.QueryByUid(uid)
			if user == nil {
				result.Result{}.Result(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
				return
			}
		}
		// Update
		userParam := &param.UserParam{}
		if err := c.ShouldBind(userParam); err != nil {
			result.Result{}.Result(http.StatusBadRequest).SetMessage(exception.WrapValidationError(err).Error()).JSON(c)
			return
		}
		user.Username = userParam.Username
		user.Sex = enum.StringToSexType(userParam.Sex)
		user.Profile = *userParam.Profile
		user.BirthTime, _ = xdatetime.JsonDate{}.Parse(userParam.BirthTime, u.Config.CurrentLoc)
		user.PhoneNumber = userParam.PhoneNumber
		url, ok := util.CommonUtil.GetFilenameFromUrl(userParam.AvatarUrl, u.Config.FileConfig.ImageUrlPrefix)
		if !ok {
			result.Result{}.Result(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
			return
		}
		user.AvatarUrl = url

		status := u.UserDao.Update(user)
		if status == database.DbNotFound {
			result.Result{}.Result(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
			return
		} else if status == database.DbExisted {
			result.Result{}.Result(http.StatusBadRequest).SetMessage(exception.UsernameUsedError.Error()).JSON(c)
			return
		} else if status == database.DbFailed {
			result.Result{}.Error().SetMessage(exception.UserUpdateError.Error()).JSON(c)
			return
		}

		retDto := xcondition.First(u.Mapper.Map(&dto.UserDto{}, user)).(*dto.UserDto)
		result.Result{}.Ok().SetData(retDto).JSON(c)
	}
}

// @Router              /v1/user/ [DELETE]
// @Security            Jwt
// @Template            Auth
// @Summary             删除用户
// @Tag                 User
// @ResponseDesc 404    user not found
// @ResponseDesc 500    user delete failed
// @ResponseModel 200   #Result
// @Response 200        ${resp_success}
//
// @Router              /v1/user/admin/{uid} [DELETE]
// @Security            Jwt
// @Template            Admin Auth ParamA
// @Summary             删除用户
// @Description         管理员权限
// @Tag                 User
// @Tag                 Administration
// @Param               uid path integer true false "用户id"
// @ResponseDesc 404    "user not found"
// @ResponseDesc 500    "user delete failed"
// @ResponseModel 200   #Result
// @Response 200        ${resp_success}
func (u *UserController) DeleteUser(isExact bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		var uid int32
		if !isExact {
			uid = u.JwtService.GetAuthUser(c).Uid
		} else {
			_uid, ok := param.BindRouteId(c, "uid")
			if !ok {
				result.Result{}.Result(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c)
				return
			}
			uid = _uid
		}
		// Delete
		status := u.UserDao.Delete(uid)
		if status == database.DbNotFound {
			result.Result{}.Result(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
			return
		} else if status == database.DbFailed {
			result.Result{}.Error().SetMessage(exception.UserDeleteError.Error()).JSON(c)
			return
		}
		result.Result{}.Ok().JSON(c)
	}
}
