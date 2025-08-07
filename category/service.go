package category

import (
	"context"

	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/domain/constant"
)

type CategoryRepository interface {
	GetByURI(ctx context.Context, uri string) (*domain.Category, error)
	List(ctx context.Context, where map[string]any, pagination *constant.Pagination, order string) ([]domain.Category, error)
}

type Service struct {
	repository CategoryRepository
}

func NewService(repository CategoryRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetCategories(c context.Context, param *domain.CategoriesReq) (*domain.CategoriesResp, error) {
	if param.ParentUri == "" {
		return nil, constant.NewXError(100, "Parent URI is required")
	}

	category, err := s.repository.GetByURI(c, param.ParentUri)
	if err != nil {
		return nil, constant.NewXError(101, "Category not found")
	}

	categories, err := s.repository.List(c, map[string]any{"parent_id": category.ID}, &constant.Pagination{
		Page:     1,
		PageSize: 100,
	}, "id resc")

	if err != nil {
		return nil, constant.NewXError(102, "Failed to list categories:"+err.Error())
	}

	resp := &domain.CategoriesResp{
		Categories: make([]domain.CategoryResp, 0, len(categories)),
	}

	for _, cat := range categories {
		resp.Categories = append(resp.Categories, *cat.Trans2Resp())
	}

	return resp, nil
}
