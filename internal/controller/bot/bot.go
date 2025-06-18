package bot

import (
	"ikit-api/internal/domain"
	"ikit-api/internal/util/resp"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

type controller struct {
	uc domain.IBotUsecase
}

func (c *controller) Start(ctx *core.RequestEvent) error {
	c.uc.Start(ctx.Request.Context())
	return resp.Succ(ctx, nil)
}

func (c *controller) Stop(ctx *core.RequestEvent) error {
	c.uc.Stop(ctx.Request.Context())
	return resp.Succ(ctx, nil)
}

func Register(route *router.Router[*core.RequestEvent], uc domain.IBotUsecase, essayUc domain.IessayUsecase) {
	c := &controller{uc: uc}

	group := route.Group("/api/v1/bot")
	group.POST("/start", c.Start)
	group.POST("/stop", c.Stop)

	essayController := &essayController{uc: essayUc}
	group.POST("/notify", essayController.Notify)

}
