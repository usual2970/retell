package auth

import (
	"context"

	"github.com/usual2970/retell/domain"
)

type Loginer interface {
	Login(ctx context.Context, req *domain.AuthLoginReq) (*domain.UserAccount, error)
}

func NewLoginer(t string, accountRepo AccountRepository, codeSvc CodeService) Loginer {
	switch t {
	case domain.AuthLoginTypeCode:
		return NewCodeLoginer(codeSvc, accountRepo)
	case domain.AuthLoginTypePassword:
		return NewPasswordLoginer(accountRepo)
	default:
		return nil
	}
}
