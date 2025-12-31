package repository

import (
	"context"
	"fun-admin/internal/model"
	"fun-admin/pkg/logger"
	"time"

	"gorm.io/gorm"
)

type LoginRepository interface {
	// 保存短信验证码
	SaveSMSCode(ctx context.Context, mobile, code string, expiredAt time.Time) error
	// 获取有效的短信验证码
	GetValidSMSCode(ctx context.Context, mobile, code string) (*model.SMSCode, error)
	// 标记验证码为已使用
	MarkSMSCodeAsUsed(ctx context.Context, id uint) error
}

func NewLoginRepository(
	logger *logger.Logger,
	db *gorm.DB,
) LoginRepository {
	return &loginRepository{
		logger: logger,
		db:     db,
	}
}

type loginRepository struct {
	logger *logger.Logger
	db     *gorm.DB
}

// SaveSMSCode 保存短信验证码
func (r *loginRepository) SaveSMSCode(ctx context.Context, mobile, code string, expiredAt time.Time) error {
	smsCode := &model.SMSCode{
		Mobile:    mobile,
		Code:      code,
		ExpiredAt: expiredAt,
		IsUsed:    false,
	}
	return r.db.WithContext(ctx).Create(smsCode).Error
}

// GetValidSMSCode 获取有效的短信验证码
func (r *loginRepository) GetValidSMSCode(ctx context.Context, mobile, code string) (*model.SMSCode, error) {
	var smsCode model.SMSCode
	err := r.db.WithContext(ctx).Where(
		"mobile = ? AND code = ? AND expired_at > ? AND is_used = ?",
		mobile, code, time.Now(), false,
	).First(&smsCode).Error
	if err != nil {
		return nil, err
	}
	return &smsCode, nil
}

// MarkSMSCodeAsUsed 标记验证码为已使用
func (r *loginRepository) MarkSMSCodeAsUsed(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.SMSCode{}).Where("id = ?", id).Update("is_used", true).Error
}