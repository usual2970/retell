package resourceposition

import (
	"context"

	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/domain/constant"
	"github.com/usual2970/retell/pkg/jwt"
)

type ResourcePositionRepository interface {
	Save(ctx context.Context, position *domain.ResourcePosition, resources []domain.Resource, fn func(position *domain.ResourcePosition)) error
	GetByURI(ctx context.Context, uri string) (*domain.ResourcePosition, error)
	GetResourceByURI(ctx context.Context, uri string) (*domain.Resource, error)
	SaveResource(ctx context.Context, resource *domain.Resource) error
	DeleteResource(ctx context.Context, resource *domain.Resource) error
	GetByGameURI(ctx context.Context, uri string) (*domain.ResourcePosition, error)

	Update(ctx context.Context, position *domain.ResourcePosition, toAddResources, toUpdateResources, toDeleteResources []domain.Resource) error
}

type Service struct {
	positionRepo ResourcePositionRepository
}

func NewService(positionRepo ResourcePositionRepository) *Service {
	return &Service{

		positionRepo: positionRepo,
	}
}

func (s *Service) GetResource(ctx context.Context, param *domain.ResourcePositionGetReq) (*domain.ResourcePositionResp, error) {

	position, err := s.positionRepo.GetByURI(ctx, param.URI)
	if err != nil {
		return nil, constant.ErrResourcePositionNotExists
	}

	return position.Trans2Resp(), nil
}

func (s *Service) CreateResource(ctx context.Context, param *domain.ResourcePositionCreateResourceReq) error {
	userID, err := jwt.GetUserID(ctx)
	if err != nil {
		return err
	}

	// 1. 检查参数
	// 2. 检查资源位是否存在
	position, err := s.positionRepo.GetByURI(ctx, param.PositionURI)
	if err != nil {
		return constant.ErrResourcePositionNotExists
	}
	if position.UserID != userID {
		return constant.ErrResourcePositionNotExists
	}
	// 3. 添加资源
	resource := domain.InitResource(param.Item, position)

	if err := s.positionRepo.SaveResource(ctx, &resource); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateResource(ctx context.Context, param *domain.ResourcePositionUpdateResourceReq) error {
	userID, err := jwt.GetUserID(ctx)
	if err != nil {
		return err
	}

	// 1. 检查参数
	// 2. 检查资源位是否存在
	position, err := s.positionRepo.GetByURI(ctx, param.PositionURI)
	if err != nil {
		return constant.ErrResourcePositionNotExists
	}
	if position.UserID != userID {
		return constant.ErrResourcePositionNotExists
	}

	// 3. 检查资源是否存在
	resource, err := s.positionRepo.GetResourceByURI(ctx, param.Item.URI)
	if err != nil {
		return constant.ErrResourceNotExists
	}
	if resource.PositionID != position.ID {
		return constant.ErrResourceNotExists
	}

	// 4. 更新资源
	resource.Update(param.Item)

	if err := s.positionRepo.SaveResource(ctx, resource); err != nil {
		return err
	}
	return nil
}
func (s *Service) DeleteResource(ctx context.Context, param *domain.ResourcePositionDeleteResourceReq) error {
	userID, err := jwt.GetUserID(ctx)
	if err != nil {
		return err
	}

	// 1. 检查参数
	// 2. 检查资源位是否存在
	position, err := s.positionRepo.GetByURI(ctx, param.PositionURI)
	if err != nil {
		return constant.ErrResourcePositionNotExists
	}
	if position.UserID != userID {
		return constant.ErrResourcePositionNotExists
	}

	// 3. 检查资源是否存在
	resource, err := s.positionRepo.GetResourceByURI(ctx, param.Uri)
	if err != nil {
		return constant.ErrResourceNotExists
	}
	if resource.PositionID != position.ID {
		return constant.ErrResourceNotExists
	}

	if err := s.positionRepo.DeleteResource(ctx, resource); err != nil {
		return err
	}
	return nil
}
