package routes

import (
	botUC "github.com/usual2970/retell/internal/usecase/bot"

	"github.com/pocketbase/pocketbase/core"
)

func OnessayUpdate(e *core.RecordRequestEvent) error {
	if err := e.Next(); err != nil {
		return err
	}

	uc := botUC.NewessayUsecase()

	return uc.CreateTelegraph(e.Request.Context(), e.Record.Id)
}

func OnessayCreate(e *core.RecordRequestEvent) error {
	if err := e.Next(); err != nil {
		return err
	}

	uc := botUC.NewessayUsecase()

	return uc.CreateTelegraph(e.Request.Context(), e.Record.Id)
}
