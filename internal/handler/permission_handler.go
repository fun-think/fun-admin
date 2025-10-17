package handler

import (
	v1 "fun-admin/api/v1"
	"fun-admin/internal/service"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	*Handler
	permissionService service.PermissionService
}

func NewPermissionHandler(
	handler *Handler,
	permissionService service.PermissionService,
) *PermissionHandler {
	return &PermissionHandler{
		Handler:           handler,
		permissionService: permissionService,
	}
}

// GetUserPermissions godoc
// @Summary 获取用户权限
// @Schemes
// @Description 获取当前用户的权限列表
// @Tags 权限模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.GetUserPermissionsData
// @Router /v1/admin/user/permissions [get]
func (h *PermissionHandler) GetUserPermissions(ctx *gin.Context) {
	userId, err := GetUserIdFromCtx(ctx)
	if err != nil {
		v1.HandleError(ctx, err)
		return
	}
	data, err := h.permissionService.GetUserPermissions(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, err)
		return
	}
	// 过滤权限菜单
	v1.HandleSuccess(ctx, data)
}

// GetRolePermissions godoc
// @Summary 获取角色权限
// @Schemes
// @Description 获取指定角色的权限列表
// @Tags 权限模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param role query string true "角色名称"
// @Success 200 {object} v1.GetRolePermissionsData
// @Router /v1/admin/role/permissions [get]
func (h *PermissionHandler) GetRolePermissions(ctx *gin.Context) {
	var req v1.GetRolePermissionsRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, err)
		return
	}
	data, err := h.permissionService.GetRolePermissions(ctx, req.Role)
	if err != nil {
		v1.HandleError(ctx, err)
		return
	}
	v1.HandleSuccess(ctx, data)
}

// UpdateRolePermission godoc
// @Summary 更新角色权限
// @Schemes
// @Description 更新指定角色的权限列表
// @Tags 权限模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UpdateRolePermissionRequest true "参数"
// @Success 200 {object} v1.Response
// @Router /v1/admin/role/permissions [put]
func (h *PermissionHandler) UpdateRolePermission(ctx *gin.Context) {
	var req v1.UpdateRolePermissionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, err)
		return
	}
	err := h.permissionService.UpdateRolePermission(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, err)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// AddRoleForUser godoc
// @Summary 为用户添加角色
// @Schemes
// @Description 为指定用户添加角色
// @Tags 权限模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param user query string true "用户ID"
// @Param role query string true "角色名称"
// @Success 200 {object} v1.Response
// @Router /v1/admin/user/role [post]
func (h *PermissionHandler) AddRoleForUser(ctx *gin.Context) {
	user := ctx.Query("user")
	role := ctx.Query("role")

	if user == "" || role == "" {
		v1.HandleError(ctx, v1.ErrBadRequest)
		return
	}

	err := h.permissionService.AddRoleForUser(ctx, user, role)
	if err != nil {
		v1.HandleError(ctx, err)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// DeleteRoleForUser godoc
// @Summary 删除用户角色
// @Schemes
// @Description 删除指定用户的某个角色
// @Tags 权限模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param user query string true "用户ID"
// @Param role query string true "角色名称"
// @Success 200 {object} v1.Response
// @Router /v1/admin/user/role [delete]
func (h *PermissionHandler) DeleteRoleForUser(ctx *gin.Context) {
	user := ctx.Query("user")
	role := ctx.Query("role")

	if user == "" || role == "" {
		v1.HandleError(ctx, v1.ErrBadRequest)
		return
	}

	err := h.permissionService.DeleteRoleForUser(ctx, user, role)
	if err != nil {
		v1.HandleError(ctx, err)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// GetPermissionsForUser godoc
// @Summary 获取用户所有权限
// @Schemes
// @Description 获取指定用户的所有权限列表
// @Tags 权限模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param user query string true "用户ID"
// @Success 200 {object} v1.GetUserPermissionsData
// @Router /v1/admin/user/all-permissions [get]
func (h *PermissionHandler) GetPermissionsForUser(ctx *gin.Context) {
	user := ctx.Query("user")

	if user == "" {
		v1.HandleError(ctx, v1.ErrBadRequest)
		return
	}

	permissions, err := h.permissionService.GetPermissionsForUser(ctx, user)
	if err != nil {
		v1.HandleError(ctx, err)
		return
	}

	data := &v1.GetUserPermissionsData{
		List: permissions,
	}
	v1.HandleSuccess(ctx, data)
}

// GetAllRoles godoc
// @Summary 获取所有角色
// @Schemes
// @Description 获取系统中所有角色列表
// @Tags 权限模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} []string
// @Router /v1/admin/roles/all [get]
func (h *PermissionHandler) GetAllRoles(ctx *gin.Context) {
	roles, err := h.permissionService.GetAllRoles(ctx)
	if err != nil {
		v1.HandleError(ctx, err)
		return
	}
	v1.HandleSuccess(ctx, roles)
}
