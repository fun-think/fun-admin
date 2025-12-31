package handler

import (
	v1 "fun-admin/api/v1"
	"fun-admin/internal/service"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	*Handler
	loginService service.LoginService
}

func NewLoginHandler(
	handler *Handler,
	loginService service.LoginService,
) *LoginHandler {
	return &LoginHandler{
		Handler:      handler,
		loginService: loginService,
	}
}

// Login godoc
// @Summary 账号登录
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.LoginRequest true "params"
// @Success 200 {object} v1.LoginResponse
// @Router /v1/login [post]
func (h *LoginHandler) Login(ctx *gin.Context) {
	var req v1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleValidationError(ctx, err.Error())
		return
	}

	response, err := h.loginService.Login(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, err)
		return
	}
	v1.HandleSuccess(ctx, response)
}

// SendSMSCode godoc
// @Summary 发送短信验证码
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body object{mobile:string} true "手机号"
// @Success 200 {object} v1.Response
// @Router /v1/send-sms-code [post]
func (h *LoginHandler) SendSMSCode(ctx *gin.Context) {
	var req struct {
		Mobile string `json:"mobile" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleValidationError(ctx, err.Error())
		return
	}

	if err := h.loginService.SendSMSCode(ctx, req.Mobile); err != nil {
		v1.HandleError(ctx, err)
		return
	}

	v1.HandleSuccess(ctx, nil)
}
