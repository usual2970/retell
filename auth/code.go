package auth

import (
	"context"

	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/domain/constant"
	"github.com/usual2970/retell/pkg/encode"
)

type AccountRepository interface {
	GetOneByOpenid(ctx context.Context, openid string) (*domain.UserAccount, error)
	Save(ctx context.Context, account *domain.UserAccount,
		profile *domain.UserProfile, privateInfo *domain.UserPrivateInfo,
		beforeSave func(account *domain.UserAccount) bool) error
	GetByUserID(ctx context.Context, userID int64) (*domain.UserAccount, error)
	DeleteFromCache(ctx context.Context, openid string) error
	DelByUserIdFromCache(ctx context.Context, userID int64) error
}

type CodeLoginer struct {
	accountRepo AccountRepository
	codeSvc     CodeService
	platfromId  int8
}

func NewCodeLoginer(codeSvc CodeService, accountRepo AccountRepository) *CodeLoginer {
	return &CodeLoginer{
		accountRepo: accountRepo,
		platfromId:  domain.UserAccountPlatformEmail,
		codeSvc:     codeSvc,
	}
}

func (c *CodeLoginer) Login(ctx context.Context, req *domain.AuthLoginReq) (*domain.UserAccount, error) {
	// check code
	if err := c.codeSvc.CheckCode(ctx, &domain.CodeCheckReq{
		Receiver: req.Identifier,
		Code:     req.Data,
		Purpose:  domain.CodePurposeLogin,
	}); err != nil {
		return nil, err
	}

	// 检测邮箱是否存在
	openid := encode.Openid(domain.UserAccountPlatformEmail, req.Identifier)

	account, err := c.accountRepo.GetOneByOpenid(ctx, openid)
	if err != nil {
		return nil, constant.ErrAccountNotExists
	}

	c.accountRepo.DeleteFromCache(ctx, openid)

	return c.accountRepo.GetByUserID(ctx, account.UserID)

}
