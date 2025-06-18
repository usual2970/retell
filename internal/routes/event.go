package routes

import (
	botUC "ikit-api/internal/usecase/bot"

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
