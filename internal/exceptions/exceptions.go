package exceptions

import (
	"net/http"

	"funnel/pkg/log"
)

// Error 表示自定义错误，包括状态码、消息和日志级别。
type Error struct {
	Code  int
	Msg   string
	Level log.Level
}

// ApiExceptions 表示自定义的错误类型
var (
	ServerError  = NewError(200500, log.LevelError, "系统异常，请稍后重试!")
	ParamsError  = NewError(200501, log.LevelInfo, "参数错误")
	ZFError      = NewError(200502, log.LevelError, "正方异常，请稍后重试!")
	ZFLoginError = NewError(200503, log.LevelError, "正方登录失败，请稍后重试!")
	OauthError   = NewError(200599, log.LevelError, "统一异常，请稍后重试!")

	NotFound = NewError(200404, log.LevelWarn, http.StatusText(http.StatusNotFound))
)

// Error 方法实现了 error 接口，返回错误的消息内容
func (e *Error) Error() string {
	return e.Msg
}

// NewError 创建并返回一个新的自定义错误实例
func NewError(code int, level log.Level, msg string) *Error {
	return &Error{
		Code:  code,
		Msg:   msg,
		Level: level,
	}
}
