package routes

import (
	botUC "ikit-api/internal/usecase/bot"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"

	"ikit-api/internal/controller/bot"
)

func Route(router *router.Router[*core.RequestEvent]) {

	uc, err := botUC.New()
	if err != nil {
		panic(err)
	}

	essayUc := botUC.NewessayUsecase()
	bot.Register(router, uc, essayUc)

}
