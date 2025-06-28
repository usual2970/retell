package routes

import (
	"context"

	"github.com/usual2970/retell/internal/domain"
	botUC "github.com/usual2970/retell/internal/usecase/bot"
)

var botUc domain.IBotUsecase

func Register() error {
	var err error
	botUc, err = botUC.New()

	botUc.Start(context.Background())

	return err
}

func UnRegister() {
	if botUc == nil {
		return
	}
	botUc.Stop(context.Background())
}
