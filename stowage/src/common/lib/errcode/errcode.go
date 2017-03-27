package errcode

import "fmt"

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (ce *CodeError) Error() string {
	return fmt.Sprintf("%d: %s", ce.Code, ce.Msg)
}

func New(code int, msg string) *CodeError {
	return &CodeError{code, msg}
}

var (
	ErrParams             = &CodeError{10000, "参数错误"}
	ErrCheckRequestFailed = &CodeError{10001, "URL请求不合法"}
	ErrRequestExpired     = &CodeError{10002, "URL请求过期"}
	ErrServerError        = &CodeError{10003, "服务器繁忙，请稍后重试"}
	ErrParamTime          = &CodeError{10004, "时间格式错误"}
	// user
	ErrGetUserInfoFailed            = &CodeError{20000, "获取用户信息失败"}
	ErrUserNotExisted               = &CodeError{20001, "用户不存在"}
	ErrUserAlreadyExisted           = &CodeError{20002, "用户已经存在"}
	ErrUserPasswordError            = &CodeError{20003, "用户密码错误~"}
	ErrBindTelFailed                = &CodeError{20004, "绑定手机号失败"}
	ErrUserUpdateFailed             = &CodeError{20005, "更新用户信息失败"}
	ErrUserCreateFailed             = &CodeError{20006, "新建用户失败"}
	ErrUserUploadPicFailed          = &CodeError{20007, "上传图片失败"}
	ErrOurUserTelHasAlreadyRegisted = &CodeError{20008, "手机号已经注册,请登录"}
	ErrGetOurUserByTelFailed        = &CodeError{20009, "未查询到用户请先联系管理员注册"}
	ErrOurUserGetAuthFailed         = &CodeError{20010, "获取签名信息失败，请重试"}
	ErrAuthCreateFailed             = &CodeError{20101, "出问题了，稍后再试吧~"}
	ErrAuthCheckFailed              = &CodeError{20102, "出问题了，稍后再试吧~"}
	ErrAuthCodeError                = &CodeError{20103, "验证码错误"}
	ErrAuthCodeExpired              = &CodeError{20104, "验证码已经失效"}
	ErrUserCodeHasAlreadyExited     = &CodeError{20106, "验证码已经发送，请60秒后重试"}
	ErrUserPremissionError          = &CodeError{20107, "您没有足够的权限查看该数据！"}
)

func ParseError(err error) (code int, msg string) {
	if e, ok := err.(*CodeError); ok {
		return e.Code, e.Msg
	}
	return ErrServerError.Code, ErrServerError.Msg
}
