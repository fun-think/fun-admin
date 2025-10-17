package job

import (
	"fun-admin/internal/repository"
	"fun-admin/pkg/jwt"
	"fun-admin/pkg/logger"
	"fun-admin/pkg/sid"
)

type Job struct {
	logger *logger.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewJob(
	tm repository.Transaction,
	logger *logger.Logger,
	sid *sid.Sid,
) *Job {
	return &Job{
		logger: logger,
		sid:    sid,
		tm:     tm,
	}
}
