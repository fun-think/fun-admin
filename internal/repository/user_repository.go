package repository

import (
	"context"
	v1 "fun-admin/api/v1"
	"fun-admin/internal/model"
	"fun-admin/pkg/logger"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUser(ctx context.Context, id uint) (*model.User, error)
	GetUsers(ctx context.Context, req *v1.GetUsersRequest) ([]*model.User, int64, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByMobile(ctx context.Context, mobile string) (*model.User, error)
	UserUpdate(ctx context.Context, user *model.User) error
	UserCreate(ctx context.Context, user *model.User) error
	UserDelete(ctx context.Context, id uint) error
}

func NewUserRepository(
	logger *logger.Logger,
	db *gorm.DB,
	enforcer *casbin.SyncedEnforcer,
) UserRepository {
	return &userRepository{
		logger:   logger,
		db:       db,
		enforcer: enforcer,
	}
}

type userRepository struct {
	logger   *logger.Logger
	db       *gorm.DB
	enforcer *casbin.SyncedEnforcer
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *userRepository) GetUserByMobile(ctx context.Context, mobile string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ?", mobile).First(&user).Error
	return &user, err
}

func (r *userRepository) GetUsers(ctx context.Context, req *v1.GetUsersRequest) ([]*model.User, int64, error) {
	var list []*model.User
	db := r.db.Model(&model.User{}).Order("id DESC")

	if req.Username != "" {
		db = db.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Nickname != "" {
		db = db.Where("nickname LIKE ?", "%"+req.Nickname+"%")
	}
	if req.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+req.Phone+"%")
	}
	if req.Email != "" {
		db = db.Where("email LIKE ?", "%"+req.Email+"%")
	}

	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	err = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize)).Find(&list).Error
	return list, total, err
}

func (r *userRepository) GetUser(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *userRepository) UserUpdate(ctx context.Context, user *model.User) error {
	if user.Password == "" {
		return r.db.Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
			"username": user.Username,
			"nickname": user.Nickname,
			"phone":    user.Phone,
			"email":    user.Email,
		}).Error
	}
	return r.db.Model(&model.User{}).Where("id = ?", user.ID).Updates(user).Error
}

func (r *userRepository) UserCreate(ctx context.Context, user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) UserDelete(ctx context.Context, id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.User{}).Error
}
