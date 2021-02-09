package user

import (
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/internal/model"
	"github.com/vidorg/vid_backend/internal/serializer"
	"github.com/vidorg/vid_backend/pkg/jwt"
	"github.com/vidorg/vid_backend/pkg/orm"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// LoginService 管理用户登录的服务
type LoginService struct {
	UserName string `form:"username" json:"username" binding:"required,min=3,max=12"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
}

// NoParamsService 无参数的服务
type NoParamsService struct{}

// QueryUsersService 用户分页查询的服务
type QueryUsersService struct {
	Page  int `form:"page" json:"page" query:"page"`
	Limit int `form:"limit" json:"limit" query:"limit"`
}

// RegisterService 管理用户注册的服务
type RegisterService struct {
	UserName string `form:"username" json:"username" binding:"required,min=3,max=12"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
	NickName string `form:"nickname" json:"nickname" binding:"required,min=3,max=20"`
	Email    string `form:"email" json:"email" binding:"required,min=3"`
}

// ResetPasswordService 管理用户重新设置密码的服务
type ResetPasswordService struct {
	UserName    string `form:"username" json:"username" binding:"required,min=3,max=12"`
	Password    string `form:"password" json:"password" binding:"required,min=3,max=20"`
	NewPassword string `form:"new_password" json:"new_password" binding:"required,min=3,max=20"`
}

// Login 用户登录
func (u *LoginService) Login() *serializer.Response {
	user := &model.User{}

	// 查找用户
	if rdb := orm.DB().Where("username = ?", u.UserName).First(user); rdb.RowsAffected == 0 {
		return serializer.ParamErr("账号或密码错误", nil)
	} else if err := rdb.Error; err != nil {
		return serializer.DBErr("查找用户错误", err)
	}

	// 检查密码
	if ok, err := user.MatchPassword(u.Password); err != nil {
		return serializer.EncryptErr("密码校验失败", err)
	} else if !ok {
		return serializer.ParamErr("账号或密码错误", nil)
	}

	// JWT
	token, err := jwt.GenerateToken(user.ID, time.Hour)
	if err != nil {
		return serializer.EncryptErr("令牌生成失败", err)
	}

	return serializer.BuildLoginResponse(user, string(token))
}

// Logout 用户登出
func (s *NoParamsService) Logout(c *gin.Context) *serializer.Response {

	return &serializer.Response{
		Code: 200,
		Msg:  "注销成功",
	}
}

func (s *NoParamsService) Auth(c *gin.Context) *serializer.Response {
	user, exists := c.Get("user")
	if !exists {
		return serializer.LoginExpiredErr()
	}
	return serializer.BuildUserResponse(user.(*model.User))
}

// Register 用户注册
func (u *RegisterService) Register() *serializer.Response {
	user := &model.User{
		UserName: u.UserName,
		Nickname: u.NickName,
		Status:   model.UserActive,
		Email:    &u.Email,
		Role:     "normal",
	}

	// 表单验证
	var count int64
	orm.DB().Model(&model.User{}).Where("username = ?", u.UserName).Count(&count)
	if count > 0 {
		count = 0
		return serializer.ParamErr("用户名已经注册", nil)
	}

	// 加密密码
	if err := user.SetPassword(u.Password); err != nil {
		return serializer.EncryptErr("密码加密失败", err)
	}

	// 创建用户
	if err := orm.DB().Create(&user).Error; err != nil {
		return serializer.ParamErr("注册失败", err)
	}

	return serializer.BuildUserResponse(user)
}

// QueryUsers 分页查询用户
func (q *QueryUsersService) QueryUsers() *serializer.Response {
	page := q.Page
	limit := q.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 1
	}

	users := make([]*model.User, 0)
	var total int64
	orm.DB().Model(&model.User{}).Count(&total)
	orm.DB().Offset((page - 1) * limit).Limit(limit).Model(&model.User{}).Find(&users)

	res := serializer.BuildListResponse(total, page, limit, serializer.BuildUsersResponse(users))
	return res
}

// ResetPassword 更新密码
func (u *ResetPasswordService) ResetPassword() *serializer.Response {
	user := &model.User{}

	// 查找用户
	if rdb := orm.DB().Where("username = ?", u.UserName).First(user); rdb.RowsAffected == 0 {
		return serializer.ParamErr("账号或密码错误", nil)
	} else if err := rdb.Error; err != nil {
		return serializer.DBErr("查找用户错误", err)
	}

	// 检查密码
	if ok, err := user.MatchPassword(u.Password); err != nil {
		return serializer.EncryptErr("密码校验失败", err)
	} else if !ok {
		return serializer.ParamErr("账号或密码错误", nil)
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(u.NewPassword), model.PasswordCost)
	if err != nil {
		return serializer.ParamErr("参数错误", err)
	}

	update := orm.DB().Model(&model.User{}).Update("password", string(bytes))
	if update == nil {
		return serializer.ParamErr("密码加密失败", err)
	}
	return &serializer.Response{
		Code: 200,
		Msg:  "更新成功",
	}
}
