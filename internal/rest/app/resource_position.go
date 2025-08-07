package app

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/pkg/resp"
)

type ResourcePositionService interface {
	GetResource(c context.Context, param *domain.ResourcePositionGetReq) (*domain.ResourcePositionResp, error)
}

type ResourcePositionHandler struct {
	service ResourcePositionService
}

func NewResourceHandler(g *echo.Group, service ResourcePositionService) {
	handler := &ResourcePositionHandler{
		service: service,
	}
	g.GET("/resource-position/:uri", handler.GetResource)

}

func (h *ResourcePositionHandler) GetResource(c echo.Context) error {
	param := &domain.ResourcePositionGetReq{}
	if err := c.Bind(param); err != nil {
		return resp.Err(c, err)
	}
	rs, err := h.service.GetResource(c.Request().Context(), param)
	if err != nil {
		return resp.Err(c, err)
	}
	return resp.Succ(c, rs)
}
