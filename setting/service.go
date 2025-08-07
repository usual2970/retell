package setting

import (
	"context"

	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/pkg/jwt"
)

type SettingRepository interface {
	GetSetting(ctx context.Context, key string, userId ...int64) (*domain.Setting, error)
}

type Service struct {
	repo SettingRepository
}

func NewService(repo SettingRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetSetting(ctx context.Context, key string) (*domain.SettingResp, error) {
	rs, err := s.repo.GetSetting(ctx, key)
	if err != nil {
		return nil, err
	}

	return rs.Trans2Resp(), nil
}

func (s *Service) GetUserSetting(ctx context.Context, key string) (*domain.SettingResp, error) {
	userId, err := jwt.GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	rs, err := s.repo.GetSetting(ctx, key, userId)
	if err != nil {
		return nil, err
	}

	return rs.Trans2Resp(), nil
}
