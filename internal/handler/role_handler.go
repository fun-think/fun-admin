package handler

import (
	v1 "fun-admin/api/v1"
	"fun-admin/internal/service"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	*Handler
	roleService service.RoleService
}

func NewRoleHandler(
	handler *Handler,
	roleService service.RoleService,
) *RoleHandler {
	return &RoleHandler{
		Handler:     handler,
		roleService: roleService,
	}
}

// GetRoles godoc
// @Summary 获取角色列表
// @Schemes
// @Description 获取角色列表
// @Tags 角色模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Param sid query string false "角色ID"
// @Param name query string false "角色名称"
// @Success 200 {object} v1.GetRolesResponse
// @Router /v1/admin/roles [get]
func (h *RoleHandler) GetRoles(ctx *gin.Context) {
	var req v1.GetRoleListRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, err)
		return
	}
	data, err := h.roleService.GetRoles(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, err)
		return
	}

	v1.HandleSuccess(ctx, data)
}

// RoleCreate godoc
// @Summary 创建角色
// @Schemes
// @Description 创建新的角色
// @Tags 角色模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.RoleCreateRequest true "参数"
// @Success 200 {object} v1.Response
// @Router /v1/admin/role [post]
func (h *RoleHandler) RoleCreate(ctx *gin.Context) {
	var req v1.RoleCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, err)
		return
	}
	if err := h.roleService.RoleCreate(ctx, &req); err != nil {
		v1.HandleError(ctx, err)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// RoleUpdate godoc
// @Summary 更新角色
// @Schemes
// @Description 更新角色信息
// @Tags 角色模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.RoleUpdateRequest true "参数"
// @Success 200 {object} v1.Response
// @Router /v1/admin/role [put]
func (h *RoleHandler) RoleUpdate(ctx *gin.Context) {
	var req v1.RoleUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, err)
		return
	}
	if err := h.roleService.RoleUpdate(ctx, &req); err != nil {
		v1.HandleError(ctx, err)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// RoleDelete godoc
// @Summary 删除角色
// @Schemes
// @Description 删除指定角色
// @Tags 角色模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query uint true "角色ID"
// @Success 200 {object} v1.Response
// @Router /v1/admin/role [delete]
func (h *RoleHandler) RoleDelete(ctx *gin.Context) {
	var req v1.RoleDeleteRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, err)
		return
	}
	if err := h.roleService.RoleDelete(ctx, req.ID); err != nil {
		v1.HandleError(ctx, err)
		return
	}
	v1.HandleSuccess(ctx, nil)
}
