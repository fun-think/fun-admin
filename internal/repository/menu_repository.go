package repository

import (
	"context"
	"fun-admin/internal/model"
	"fun-admin/pkg/logger"

	"gorm.io/gorm"
)

type MenuRepository interface {
	GetMenuList(ctx context.Context) ([]*model.Menu, error)
	MenuUpdate(ctx context.Context, menu *model.Menu) error
	MenuCreate(ctx context.Context, menu *model.Menu) error
	MenuDelete(ctx context.Context, id uint) error
}

func NewMenuRepository(
	logger *logger.Logger,
	db *gorm.DB,
) MenuRepository {
	return &menuRepository{
		logger: logger,
		db:     db,
	}
}

type menuRepository struct {
	logger *logger.Logger
	db     *gorm.DB
}

func (r *menuRepository) GetMenuList(ctx context.Context) ([]*model.Menu, error) {
	var list []*model.Menu
	err := r.db.Order("weight DESC").Find(&list).Error
	return list, err
}

func (r *menuRepository) MenuUpdate(ctx context.Context, menu *model.Menu) error {
	return r.db.Model(&model.Menu{}).Where("id = ?", menu.ID).Updates(menu).Error
}

func (r *menuRepository) MenuCreate(ctx context.Context, menu *model.Menu) error {
	return r.db.Create(menu).Error
}

func (r *menuRepository) MenuDelete(ctx context.Context, id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.Menu{}).Error
}
