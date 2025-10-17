package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorCode 错误码类型
type ErrorCode int

// 预定义错误码
const (
	// 成功
	CodeSuccess ErrorCode = 0

	// 客户端错误 4xx
	CodeBadRequest      ErrorCode = 400
	CodeUnauthorized    ErrorCode = 401
	CodeForbidden       ErrorCode = 403
	CodeNotFound        ErrorCode = 404
	CodeConflict        ErrorCode = 409
	CodeTooManyRequests ErrorCode = 429

	// 服务端错误 5xx
	CodeInternalError      ErrorCode = 500
	CodeNotImplemented     ErrorCode = 501
	CodeBadGateway         ErrorCode = 502
	CodeServiceUnavailable ErrorCode = 503

	// 业务错误 1xxx
	CodeValidationFailed       ErrorCode = 1001
	CodeDuplicateEntry         ErrorCode = 1002
	CodeRecordNotFound         ErrorCode = 1003
	CodeInvalidCredentials     ErrorCode = 1004
	CodeTokenExpired           ErrorCode = 1005
	CodeInsufficientPermission ErrorCode = 1006
)

// AppError 应用错误结构
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
	Cause   error     `json:"-"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 返回底层错误
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithDetails 添加错误详情
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithCause 添加底层错误
func (e *AppError) WithCause(cause error) *AppError {
	e.Cause = cause
	return e
}

// GetHTTPStatus 获取对应的HTTP状态码
func (e *AppError) GetHTTPStatus() int {
	switch {
	case e.Code >= 400 && e.Code < 500:
		return int(e.Code)
	case e.Code >= 500 && e.Code < 600:
		return int(e.Code)
	case e.Code == CodeValidationFailed:
		return http.StatusBadRequest
	case e.Code == CodeDuplicateEntry:
		return http.StatusConflict
	case e.Code == CodeRecordNotFound:
		return http.StatusNotFound
	case e.Code == CodeInvalidCredentials:
		return http.StatusUnauthorized
	case e.Code == CodeTokenExpired:
		return http.StatusUnauthorized
	case e.Code == CodeInsufficientPermission:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// 错误创建函数
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func Newf(code ErrorCode, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   err,
	}
}

func Wrapf(err error, code ErrorCode, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Cause:   err,
	}
}

// 预定义错误实例
var (
	// 成功
	ErrSuccess = New(CodeSuccess, "success")

	// 客户端错误
	ErrBadRequest   = New(CodeBadRequest, "请求参数错误")
	ErrUnauthorized = New(CodeUnauthorized, "未授权访问")
	ErrForbidden    = New(CodeForbidden, "权限不足")
	ErrNotFound     = New(CodeNotFound, "资源不存在")

	// 服务端错误
	ErrInternalError = New(CodeInternalError, "服务器内部错误")

	// 业务错误
	ErrValidationFailed       = New(CodeValidationFailed, "数据验证失败")
	ErrDuplicateEntry         = New(CodeDuplicateEntry, "数据重复")
	ErrRecordNotFound         = New(CodeRecordNotFound, "记录不存在")
	ErrInvalidCredentials     = New(CodeInvalidCredentials, "用户名或密码错误")
	ErrTokenExpired           = New(CodeTokenExpired, "令牌已过期")
	ErrInsufficientPermission = New(CodeInsufficientPermission, "权限不足")
)

// IsAppError 检查是否为应用错误
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// AsAppError 转换为应用错误
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	ok := errors.As(err, &appErr)
	return appErr, ok
}

// GetCode 获取错误码
func GetCode(err error) ErrorCode {
	if appErr, ok := AsAppError(err); ok {
		return appErr.Code
	}
	return CodeInternalError
}

// GetMessage 获取错误消息
func GetMessage(err error) string {
	if appErr, ok := AsAppError(err); ok {
		return appErr.Message
	}
	return err.Error()
}

// GetHTTPStatus 获取HTTP状态码
func GetHTTPStatus(err error) int {
	if appErr, ok := AsAppError(err); ok {
		return appErr.GetHTTPStatus()
	}
	return http.StatusInternalServerError
}
