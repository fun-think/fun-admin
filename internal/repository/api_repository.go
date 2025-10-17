package repository

import (
	"context"
	v1 "fun-admin/api/v1"
	"fun-admin/internal/model"
	"fun-admin/pkg/logger"

	"gorm.io/gorm"
)

type ApiRepository interface {
	GetApis(ctx context.Context, req *v1.GetApisRequest) ([]*model.Api, int64, error)
	GetApiGroups(ctx context.Context) ([]string, error)
	ApiUpdate(ctx context.Context, api *model.Api) error
	ApiCreate(ctx context.Context, api *model.Api) error
	ApiDelete(ctx context.Context, id uint) error
}

func NewApiRepository(
	logger *logger.Logger,
	db *gorm.DB,
) ApiRepository {
	return &apiRepository{
		logger: logger,
		db:     db,
	}
}

type apiRepository struct {
	logger *logger.Logger
	db     *gorm.DB
}

func (r *apiRepository) GetApis(ctx context.Context, req *v1.GetApisRequest) ([]*model.Api, int64, error) {
	var list []*model.Api
	db := r.db.Model(&model.Api{}).Order("id DESC")

	if req.Group != "" {
		db = db.Where("group = ?", req.Group)
	}
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Path != "" {
		db = db.Where("path LIKE ?", "%"+req.Path+"%")
	}
	if req.Method != "" {
		db = db.Where("method = ?", req.Method)
	}

	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	err = db.Offset(int((req.Page - 1) * req.PageSize)).Limit(int(req.PageSize)).Find(&list).Error
	return list, total, err
}

func (r *apiRepository) GetApiGroups(ctx context.Context) ([]string, error) {
	var groups []string
	err := r.db.Model(&model.Api{}).Distinct().Pluck("group", &groups).Error
	return groups, err
}

func (r *apiRepository) ApiUpdate(ctx context.Context, api *model.Api) error {
	return r.db.Model(&model.Api{}).Where("id = ?", api.ID).Updates(api).Error
}

func (r *apiRepository) ApiCreate(ctx context.Context, api *model.Api) error {
	return r.db.Create(api).Error
}

func (r *apiRepository) ApiDelete(ctx context.Context, id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.Api{}).Error
}
