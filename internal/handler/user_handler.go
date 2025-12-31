package handler

import (
	v1 "fun-admin/api/v1"
	"fun-admin/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	*Handler
	userService service.UserService
}

func NewUserHandler(
	handler *Handler,
	userService service.UserService,
) *UserHandler {
	return &UserHandler{
		Handler:     handler,
		userService: userService,
	}
}

// UserUpdate godoc
// @Summary 更新管理员用户
// @Schemes
// @Description 更新管理员用户信息
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UserUpdateRequest true "参数"
// @Success 200 {object} v1.Response
// @Router /v1/admin/user [put]
func (h *UserHandler) UserUpdate(ctx *gin.Context) {
	var req v1.UserUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleValidationError(ctx, err.Error())
		return
	}
	if err := h.userService.UserUpdate(ctx, &req); err != nil {
		v1.HandleError(ctx, err)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// UserCreate godoc
// @Summary 创建管理员用户
// @Schemes
// @Description 创建新的管理员用户
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UserCreateRequest true "参数"
// @Success 200 {object} v1.Response
// @Router /v1/admin/user [post]
func (h *UserHandler) UserCreate(ctx *gin.Context) {
	var req v1.UserCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleValidationError(ctx, err.Error())
		return
	}
	if err := h.userService.UserCreate(ctx, &req); err != nil {
		v1.HandleError(ctx, err)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// UserDelete godoc
// @Summary 删除管理员用户
// @Schemes
// @Description 删除指定管理员用户
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query uint true "用户ID"
// @Success 200 {object} v1.Response
// @Router /v1/admin/user [delete]
func (h *UserHandler) UserDelete(ctx *gin.Context) {
	var req v1.UserDeleteRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleValidationError(ctx, err.Error())
		return
	}
	if err := h.userService.UserDelete(ctx, req.ID); err != nil {
		v1.HandleError(ctx, err)
		return

	}
	v1.HandleSuccess(ctx, nil)
}

// GetUsers godoc
// @Summary 获取管理员用户列表
// @Schemes
// @Description 获取管理员用户列表
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Param id query int false "用户ID"
// @Param username query string false "用户名"
// @Param nickname query string false "昵称"
// @Param phone query string false "手机号"
// @Param email query string false "邮箱"
// @Success 200 {object} v1.GetUsersResponse
// @Router /v1/admin/users [get]
func (h *UserHandler) GetUsers(ctx *gin.Context) {
	var req v1.GetUsersRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleValidationError(ctx, err.Error())
		return
	}
	data, err := h.userService.GetUsers(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, err)
		return
	}

	v1.HandleSuccess(ctx, data)
}

// GetUser godoc
// @Summary 获取管理用户信息
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.GetUserResponse
// @Router /v1/admin/user [get]
func (h *UserHandler) GetUser(ctx *gin.Context) {
	userId, err := GetUserIdFromCtx(ctx)
	if err != nil {
		v1.HandleError(ctx, err)
		return
	}
	data, err := h.userService.GetUser(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, err)
		return
	}

	v1.HandleSuccess(ctx, data)
}
