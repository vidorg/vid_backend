package serializer

import (
	"github.com/vidorg/vid_backend/internal/model"
)

// User 用户序列化器
type User struct {
	ID        int64  `json:"id"`
	UserName  string `json:"username"`
	Nickname  string `json:"nickname"`
	Status    string `json:"status"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Role      string `json:"role"`
	CreatedAt int64  `json:"created_at"`
}

// BuildUserResponse 序列化用户响应
func BuildUserResponse(user *model.User) *Response {
	res := &User{
		ID:        user.ID,
		UserName:  user.UserName,
		Nickname:  user.Nickname,
		Status:    user.Status,
		Avatar:    user.Avatar,
		Email:     *user.Email,
		Role:      user.Role,
		CreatedAt: user.Created,
	}
	return &Response{
		Code: 200,
		Msg:  "success",
		Data: res,
	}
}

// BuildUsersResponse 序列化用户
func BuildUsersResponse(users []*model.User) []*User {
	res := make([]*User, len(users))
	for i, user := range users {
		res[i] = &User{
			ID:        user.ID,
			UserName:  user.UserName,
			Nickname:  user.Nickname,
			Status:    user.Status,
			Avatar:    user.Avatar,
			Email:     *user.Email,
			Role:      user.Role,
			CreatedAt: user.Created,
		}
	}
	return res
}

// Login 登录序列化器
type Login struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// BuildLoginResponse 序列化登录响应
func BuildLoginResponse(user *model.User, token string) *Response {
	res := &Login{
		User: &User{
			ID:        user.ID,
			UserName:  user.UserName,
			Nickname:  user.Nickname,
			Status:    user.Status,
			Avatar:    user.Avatar,
			Email:     *user.Email,
			Role:      user.Role,
			CreatedAt: user.Created,
		},
		Token: token,
	}
	return &Response{
		Code: 200,
		Msg:  "登录成功",
		Data: res,
	}
}
