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
	ID        uint     `json:"id"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
	Username  string   `json:"username"`
	Nickname  string   `json:"nickname"`
	Phone     string   `json:"phone"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
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

// HandleSuccess 处理成功响应
func HandleSuccess(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// HandleError 处理通用错误响应
func HandleError(ctx *gin.Context, err error) {
	resp := Response{
		Code:    http.StatusInternalServerError,
		Message: "服务器内部错误",
		Data:    nil,
	}
	if gin.Mode() == gin.DebugMode && err != nil {
		resp.Message = err.Error()
	}
	ctx.JSON(http.StatusInternalServerError, resp)
}

// HandleValidationError 处理校验错误
func HandleValidationError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: "参数校验失败: " + message,
		Data:    nil,
	})
}

// HandleUnauthorized 未授权响应
func HandleUnauthorized(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, Response{
		Code:    http.StatusUnauthorized,
		Message: "未授权访问: 请先登录",
		Data:    nil,
	})
}

// HandleForbidden 禁止访问响应
func HandleForbidden(ctx *gin.Context) {
	ctx.JSON(http.StatusForbidden, Response{
		Code:    http.StatusForbidden,
		Message: "禁止访问: 您没有权限执行此操作",
		Data:    nil,
	})
}

// HandleNotFound 资源未找到响应
func HandleNotFound(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, Response{
		Code:    http.StatusNotFound,
		Message: "资源未找到: 请求的资源不存在",
		Data:    nil,
	})
}

// HandleServerError 服务器内部错误
func HandleServerError(ctx *gin.Context, err error) {
	message := "服务器内部错误"
	if gin.Mode() == gin.DebugMode && err != nil {
		message = message + ": " + err.Error()
	}
	ctx.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: message,
		Data:    nil,
	})
}
