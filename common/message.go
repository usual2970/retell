package common

import (
	"context"
	"sync"

	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/internal/repository"
	"github.com/usual2970/retell/message"
)

type MessageService interface {
	Send(ctx context.Context, req *domain.MessageSendReq) (*domain.MessageSendResp, error)
}

var messageSvcOnce sync.Once
var messageService MessageService

func NewMessageService() MessageService {
	messageSvcOnce.Do(func() {

		messageRepo := repository.NewMessageRepository()
		userAccountRepo := repository.NewUserAccountRepository()
		settingRepo := repository.NewSettingRepository()

		messageService = message.NewService(messageRepo, userAccountRepo, settingRepo)
	})
	return messageService
}

func SendMessage(ctx context.Context, req *domain.MessageSendReq) error {
	svc := NewMessageService()
	_, err := svc.Send(ctx, req)
	return err
}
