package serializer

import (
	"github.com/vidorg/vid_backend/pkg/logger"
	"go.uber.org/zap"
)

// Response 基础序列化器
type Response struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data,omitempty"`
	Error error       `json:"error,omitempty"`
}

// DataList 基础列表结构
type DataList struct {
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Items interface{} `json:"items"`
}

// BuildListResponse 列表构建器
func BuildListResponse(total int64, page int, limit int, items interface{}) *Response {
	return &Response{
		Code: 200,
		Msg:  "success",
		Data: &DataList{
			Total: total,
			Page:  page,
			Limit: limit,
			Items: items,
		},
	}
}

// 三位数错误编码为复用http原本含义
// 五位数错误编码为应用自定义错误
// 五开头的五位数错误编码为服务器端错误，比如数据库操作失败
// 四开头的五位数错误编码为客户端错误，有时候是客户端代码写错了，有时候是用户操作错误
const (
	CodeLoginError      = 401   // 未登录
	CodeNoRightError    = 403   // 未授权访问
	CodeParamError      = 40001 // 各种奇奇怪怪的参数错误
	CodeDBError         = 50001 // 数据库操作失败
	CodeEncryptError    = 50002 // 加密失败
	CodeServerError     = 50003 // 服务器端其他错误
	CodeUserStatusError = 50004 // 服务器端其他错误
	CodeUploadFileError = 50005 // 服务器端其他错误
)

// Err err msg
func Err(errCode int, msg string, err error) *Response {
	res := &Response{
		Code: errCode,
		Msg:  msg,
	}
	if err != nil {
		// err msg
		logger.Logger().Warn("[Error]",
			zap.Int("code", errCode),
			zap.String("msg", msg),
			zap.Error(err))
	}
	return res
}

// LoginErr 登录失败
func LoginErr() *Response {
	return Err(CodeLoginError, "未登录", nil)
}

// NoRightErr 未授权
func NoRightErr() *Response {
	return Err(CodeNoRightError, "无权限访问", nil)
}

// LoginExpiredErr 登录过期
func LoginExpiredErr() *Response {
	return Err(CodeNoRightError, "登录过期", nil)
}

// UploadFileErr 上传文件出错
func UploadFileErr(msg string, err error) *Response {
	if msg == "" {
		msg = "上传文件出错，请检查网络"
	}
	return Err(CodeUploadFileError, msg, err)
}

// ParamErr 各种参数错误
func ParamErr(msg string, err error) *Response {
	if msg == "" {
		msg = "参数错误"
	}
	return Err(CodeParamError, msg, err)
}

// UserStatusErr 账号状态
func UserStatusErr(msg string) *Response {
	if msg == "" {
		msg = "账号未被激活"
	}
	return &Response{
		Code: CodeUserStatusError,
		Msg:  msg,
	}
}

// DBErr 数据库操作失败
func DBErr(msg string, err error) *Response {
	if msg == "" {
		msg = "数据库操作失败"
	}
	return Err(CodeDBError, msg, err)
}

// EncryptErr 加密错误
func EncryptErr(msg string, err error) *Response {
	if msg == "" {
		msg = "加密错误"
	}
	return Err(CodeEncryptError, msg, err)
}

func ServerErr(msg string, err error) *Response {
	if msg == "" {
		msg = "服务器未知错误"
	}
	return Err(CodeServerError, msg, err)
}
