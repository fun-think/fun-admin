package task

import (
	"fun-admin/internal/repository"
	"fun-admin/pkg/jwt"
	"fun-admin/pkg/logger"
	"fun-admin/pkg/sid"
)

type Task struct {
	logger *logger.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewTask(
	tm repository.Transaction,
	logger *logger.Logger,
	sid *sid.Sid,
) *Task {
	return &Task{
		logger: logger,
		sid:    sid,
		tm:     tm,
	}
}
