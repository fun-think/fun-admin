package service

import (
	"fun-admin/internal/repository"
	"fun-admin/pkg/jwt"
	"fun-admin/pkg/logger"
	"fun-admin/pkg/sid"
)

type Service struct {
	logger *logger.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewService(
	tm repository.Transaction,
	logger *logger.Logger,
	sid *sid.Sid,
	jwt *jwt.JWT,
) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
		jwt:    jwt,
		tm:     tm,
	}
}
