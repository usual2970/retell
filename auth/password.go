package auth

import (
	"context"
	"errors"

	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/domain/constant"
	"github.com/usual2970/retell/pkg/encode"
)

type PasswordLoginer struct {
	accountRepo AccountRepository
}

func NewPasswordLoginer(accountRepo AccountRepository) *PasswordLoginer {
	return &PasswordLoginer{
		accountRepo: accountRepo,
	}
}

func (p *PasswordLoginer) Login(ctx context.Context, req *domain.AuthLoginReq) (*domain.UserAccount, error) {
	openid := encode.Openid(domain.UserAccountPlatformEmail, req.Identifier)

	account, err := p.accountRepo.GetOneByOpenid(ctx, openid)
	if err != nil && !errors.Is(err, constant.ErrRecordNotFound) {
		return nil, err
	}

	// openid不存在
	if err != nil {
		return nil, constant.ErrAccountOrPasswordIncorrect
	}

	// 检测密码
	privateInfo := account.UserPrivateInfo
	if !privateInfo.CheckPassword(req.Data) {
		return nil, constant.ErrAccountOrPasswordIncorrect
	}

	return account, nil
}
