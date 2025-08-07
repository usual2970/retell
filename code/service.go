package code

import (
	"context"
	"errors"
	"regexp"

	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/domain/constant"
	"github.com/usual2970/retell/pkg/logger"
)

var codeReg = regexp.MustCompile(`^\d{4}$`)

type CodeRepository interface {
	GetByReceiver(ctx context.Context, receiver string, purpose domain.CodePurpose) (*domain.Code, error)
	Save(ctx context.Context, code *domain.Code) error
	Update(ctx context.Context, code *domain.Code) error
	GetByReceiverAndCode(ctx context.Context, receiver, code string, purpose domain.CodePurpose) (*domain.Code, error)
}

type Service struct {
	codeRepo CodeRepository
}

func NewService(codeRepo CodeRepository) *Service {
	return &Service{
		codeRepo: codeRepo,
	}
}

func (s *Service) GenerateCode(ctx context.Context, param *domain.CodeGenerateReq) (*domain.Code, error) {
	l := logger.WithField("module", "code.GenerateCode")

	if c, err := s.codeRepo.GetByReceiver(ctx, param.Receiver, param.Purpose); err != nil && !errors.Is(err, constant.ErrRecordNotFound) {
		return nil, errors.New("send sms code failed:" + err.Error())
	} else if err == nil && !c.IsUsed() {
		return nil, constant.ErrCodeHasSend
	}
	// 没有的话再发送
	code := domain.NewEmailCode(param.Receiver, param.Purpose)
	err := s.codeRepo.Save(ctx, code)
	if err != nil {
		l.WithField("err", err).Error("send code failed")
		return nil, errors.New("send sms code failed:" + err.Error())
	}
	return code, nil

}

func (s *Service) CheckCode(ctx context.Context, req *domain.CodeCheckReq) error {

	code := req.Code
	if !codeReg.MatchString(code) {
		return constant.ErrParamWrongCode
	}

	c, err := s.codeRepo.GetByReceiverAndCode(ctx, req.Receiver, code, req.Purpose)
	if err != nil {
		return constant.ErrCodeExpired
	}
	if c.IsUsed() {
		return constant.ErrParamWrongCode
	}

	c.SetState(domain.CodeStateUsed)

	_ = s.codeRepo.Update(ctx, c)
	return nil
}
