package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// LoginResponseData 登录响应数据
type LoginResponseData struct {
	AccessToken string   `json:"accessToken"`
	UserInfo    UserInfo `json:"userInfo"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// ResponseData 响应数据
type ResponseData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

// GetUsersResponseData 获取用户列表响应数据
type GetUsersResponseData struct {
	List  []UserDataItem `json:"list"`
	Total int64          `json:"total"`
}

// UserDataItem 用户数据项
type UserDataItem struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

// GetUserResponseData 获取用户详情响应数据
type GetUserResponseData struct {
	ID        uint     `json:"id"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
	Username  string   `json:"username"`
	Nickname  string   `json:"nickname"`
	Phone     string   `json:"phone"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
}

// ErrorResponse 错误响应结构体
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// HandleSuccess 成功响应
func HandleSuccess(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// HandleError 处理错误响应
func HandleError(ctx *gin.Context, err error) {
	// 这里可以根据错误类型返回不同的错误码和消息
	ctx.JSON(http.StatusOK, Response{
		Code:    1,
		Message: "error",
		Data:    err.Error(),
	})
}

// HandleValidationError 处理验证错误
func HandleValidationError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, Response{
		Code:    400,
		Message: "参数验证失败",
		Data:    message,
	})
}

// HandleUnauthorized 未授权错误
func HandleUnauthorized(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code:    401,
		Message: "未授权访问",
		Data:    "请先登录",
	})
}

// HandleForbidden 禁止访问错误
func HandleForbidden(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code:    403,
		Message: "禁止访问",
		Data:    "您没有权限执行此操作",
	})
}

// HandleNotFound 资源未找到错误
func HandleNotFound(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code:    404,
		Message: "资源未找到",
		Data:    "请求的资源不存在",
	})
}

// HandleServerError 服务器内部错误
func HandleServerError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, Response{
		Code:    500,
		Message: "服务器内部错误",
		Data:    err.Error(),
	})
}
