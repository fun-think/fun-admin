package handler

import (
	"errors"
	"fun-admin/pkg/jwt"
	"fun-admin/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger *logger.Logger
}

func NewHandler(
	logger *logger.Logger,
) *Handler {
	return &Handler{
		logger: logger,
	}
}
func GetUserIdFromCtx(ctx *gin.Context) (uint, error) {
	v, exists := ctx.Get("claims")
	if !exists {
		return 0, errors.New("user id not found")
	}
	return v.(*jwt.MyCustomClaims).UserId, nil
}
