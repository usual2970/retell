package routes

import (
	botUC "github.com/usual2970/retell/internal/usecase/bot"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"

	"github.com/usual2970/retell/internal/controller/bot"
)

func Route(router *router.Router[*core.RequestEvent]) {

	uc, err := botUC.New()
	if err != nil {
		panic(err)
	}

	essayUc := botUC.NewessayUsecase()
	bot.Register(router, uc, essayUc)

}
