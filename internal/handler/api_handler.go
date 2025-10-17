package handler

import (
	v1 "fun-admin/api/v1"
	"fun-admin/internal/service"

	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	*Handler
	apiService service.ApiService
}

func NewApiHandler(
	handler *Handler,
	apiService service.ApiService,
) *ApiHandler {
	return &ApiHandler{
		Handler:    handler,
		apiService: apiService,
	}
}

// GetApis godoc
// @Summary 获取API列表
// @Schemes
// @Description 获取API列表
// @Tags API模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int true "页码"
// @Param pageSize query int true "每页数量"
// @Param group query string false "API分组"
// @Param name query string false "API名称"
// @Param path query string false "API路径"
// @Param method query string false "请求方法"
// @Success 200 {object} v1.GetApisResponse
// @Router /v1/admin/apis [get]
func (h *ApiHandler) GetApis(ctx *gin.Context) {
	var req v1.GetApisRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, err)
		return
	}
	data, err := h.apiService.GetApis(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, err)
		return
	}

	v1.HandleSuccess(ctx, data)
}

// ApiCreate godoc
// @Summary 创建API
// @Schemes
// @Description 创建新的API
// @Tags API模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.ApiCreateRequest true "参数"
// @Success 200 {object} v1.Response
// @Router /v1/admin/api [post]
func (h *ApiHandler) ApiCreate(ctx *gin.Context) {
	var req v1.ApiCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, err)
		return
	}
	if err := h.apiService.ApiCreate(ctx, &req); err != nil {
		v1.HandleError(ctx, err)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// ApiUpdate godoc
// @Summary 更新API
// @Schemes
// @Description 更新API信息
// @Tags API模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.ApiUpdateRequest true "参数"
// @Success 200 {object} v1.Response
// @Router /v1/admin/api [put]
func (h *ApiHandler) ApiUpdate(ctx *gin.Context) {
	var req v1.ApiUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, err)
		return
	}
	if err := h.apiService.ApiUpdate(ctx, &req); err != nil {
		v1.HandleError(ctx, err)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// ApiDelete godoc
// @Summary 删除API
// @Schemes
// @Description 删除指定API
// @Tags API模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query uint true "API ID"
// @Success 200 {object} v1.Response
// @Router /v1/admin/api [delete]
func (h *ApiHandler) ApiDelete(ctx *gin.Context) {
	var req v1.ApiDeleteRequest
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, err)
		return
	}
	if err := h.apiService.ApiDelete(ctx, req.ID); err != nil {
		v1.HandleError(ctx, err)
		return
	}
	v1.HandleSuccess(ctx, nil)
}
