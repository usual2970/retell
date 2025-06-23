package routes

import (
	botUC "github.com/usual2970/retell/internal/usecase/bot"

	"github.com/pocketbase/pocketbase/core"
)

func OnessayUpdate(e *core.RecordEvent) error {

	uc := botUC.NewessayUsecase()

	return uc.CreateTelegraph(e.Context, e.Record.Id)
}

func OnessayCreate(e *core.RecordEvent) error {

	uc := botUC.NewessayUsecase()

	return uc.CreateTelegraph(e.Context, e.Record.Id)
}
