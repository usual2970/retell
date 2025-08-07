package common

import (
	"context"

	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/internal/repository"
)

func GetUserAccountByUID(ctx context.Context, userID int64) (*domain.UserAccount, error) {
	repo := repository.NewUserAccountRepository()
	return repo.GetByUserID(ctx, userID)
}
