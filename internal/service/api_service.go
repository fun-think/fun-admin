package service

import (
	"context"
	v1 "fun-admin/api/v1"
	"fun-admin/internal/model"
	"fun-admin/internal/repository"
)

type ApiService interface {
	GetApis(ctx context.Context, req *v1.GetApisRequest) (*v1.GetApisResponseData, error)
	ApiUpdate(ctx context.Context, req *v1.ApiUpdateRequest) error
	ApiCreate(ctx context.Context, req *v1.ApiCreateRequest) error
	ApiDelete(ctx context.Context, id uint) error
}

func NewApiService(
	service *Service,
	apiRepository repository.ApiRepository,
) ApiService {
	return &apiService{
		Service:       service,
		apiRepository: apiRepository,
	}
}

type apiService struct {
	*Service
	apiRepository repository.ApiRepository
}

func (s *apiService) GetApis(ctx context.Context, req *v1.GetApisRequest) (*v1.GetApisResponseData, error) {
	list, total, err := s.apiRepository.GetApis(ctx, req)
	if err != nil {
		return nil, err
	}
	groups, err := s.apiRepository.GetApiGroups(ctx)
	if err != nil {
		return nil, err
	}
	data := &v1.GetApisResponseData{
		List:   make([]v1.ApiDataItem, 0),
		Total:  total,
		Groups: groups,
	}
	for _, api := range list {
		data.List = append(data.List, v1.ApiDataItem{
			CreatedAt: api.CreatedAt.Format("2006-01-02 15:04:05"),
			Group:     api.Group,
			ID:        api.ID,
			Method:    api.Method,
			Name:      api.Name,
			Path:      api.Path,
			UpdatedAt: api.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return data, nil
}

func (s *apiService) ApiUpdate(ctx context.Context, req *v1.ApiUpdateRequest) error {
	return s.apiRepository.ApiUpdate(ctx, &model.Api{
		Group:  req.Group,
		Method: req.Method,
		Name:   req.Name,
		Path:   req.Path,
		ID:     req.ID,
	})
}

func (s *apiService) ApiCreate(ctx context.Context, req *v1.ApiCreateRequest) error {
	return s.apiRepository.ApiCreate(ctx, &model.Api{
		Group:  req.Group,
		Method: req.Method,
		Name:   req.Name,
		Path:   req.Path,
	})
}

func (s *apiService) ApiDelete(ctx context.Context, id uint) error {
	return s.apiRepository.ApiDelete(ctx, id)
}
