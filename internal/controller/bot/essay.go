package bot

import (
	"github.com/usual2970/retell/internal/domain"

	"github.com/pocketbase/pocketbase/core"
)

type essayController struct {
	uc domain.IessayUsecase
}

func (c *essayController) Notify(ctx *core.RequestEvent) error {
	req := &domain.TtsAsyncResp{}
	if err := ctx.BindBody(req); err != nil {
		return err
	}
	return c.uc.Notify(ctx.Request.Context(), req)
}
