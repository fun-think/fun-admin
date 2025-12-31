package handler

import (
	"errors"

	"fun-admin/api/v1"
	"fun-admin/internal/service"
	"fun-admin/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// ProfileHandler 用户个人资料处理器
type ProfileHandler struct {
	*Handler
	profileService service.ProfileServiceInterface
}

// NewProfileHandler 创建用户个人资料处理器
func NewProfileHandler(
	handler *Handler,
	profileService service.ProfileServiceInterface,
) *ProfileHandler {
	return &ProfileHandler{
		Handler:        handler,
		profileService: profileService,
	}
}

// GetProfile 获取当前用户个人资料
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// 从上下文中获取当前用户claims
	claims, exists := c.Get("claims")
	if !exists {
		v1.HandleError(c, errors.New("无法获取用户信息"))
		return
	}

	userClaims, ok := claims.(*jwt.MyCustomClaims)
	if !ok {
		v1.HandleError(c, errors.New("用户信息格式错误"))
		return
	}

	// 获取用户详情
	profile, err := h.profileService.GetProfile(c, userClaims.UserId)
	if err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, profile)
}

// UpdateProfile 更新当前用户个人资料
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	// 从上下文中获取当前用户claims
	claims, exists := c.Get("claims")
	if !exists {
		v1.HandleError(c, errors.New("无法获取用户信息"))
		return
	}

	userClaims, ok := claims.(*jwt.MyCustomClaims)
	if !ok {
		v1.HandleError(c, errors.New("用户信息格式错误"))
		return
	}

	// 解析请求参数
	var req struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleValidationError(c, err.Error())
		return
	}

	// 更新用户资料
	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	if err := h.profileService.UpdateProfile(c, userClaims.UserId, updates); err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, nil)
}

// UpdatePassword 更新当前用户密码
func (h *ProfileHandler) UpdatePassword(c *gin.Context) {
	// 从上下文中获取当前用户claims
	claims, exists := c.Get("claims")
	if !exists {
		v1.HandleError(c, errors.New("无法获取用户信息"))
		return
	}

	userClaims, ok := claims.(*jwt.MyCustomClaims)
	if !ok {
		v1.HandleError(c, errors.New("用户信息格式错误"))
		return
	}

	// 解析请求参数
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		v1.HandleValidationError(c, err.Error())
		return
	}

	if err := h.profileService.UpdatePassword(c, userClaims.UserId, req.OldPassword, req.NewPassword); err != nil {
		v1.HandleError(c, err)
		return
	}

	v1.HandleSuccess(c, nil)
}
