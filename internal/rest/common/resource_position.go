package common

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/internal/rest/middleware"
	"github.com/usual2970/retell/pkg/resp"
)

type ResourcePositionService interface {
	GetResource(c context.Context, param *domain.ResourcePositionGetReq) (*domain.ResourcePositionResp, error)

	CreateResource(c context.Context, param *domain.ResourcePositionCreateResourceReq) error
	UpdateResource(c context.Context, param *domain.ResourcePositionUpdateResourceReq) error
	DeleteResource(c context.Context, param *domain.ResourcePositionDeleteResourceReq) error
}

type ResourcePositionHandler struct {
	service ResourcePositionService
}

func NewResourceHandler(g *echo.Group, service ResourcePositionService) {
	handler := &ResourcePositionHandler{
		service: service,
	}
	g.GET("/resource-position/:uri", handler.GetResource)
	g.POST("/resource-position/resource", handler.CreateResource, middleware.ShouldLogin)
	g.PUT("/resource-position/resource/:uri", handler.UpdateResource, middleware.ShouldLogin)
	g.POST("/resource-position/resource-delete", handler.DeleteResource, middleware.ShouldLogin)

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

func (h *ResourcePositionHandler) CreateResource(c echo.Context) error {
	param := &domain.ResourcePositionCreateResourceReq{}
	if err := c.Bind(param); err != nil {
		return resp.Err(c, err)
	}
	err := h.service.CreateResource(c.Request().Context(), param)
	if err != nil {
		return resp.Err(c, err)
	}
	return resp.Succ(c, nil)
}

func (h *ResourcePositionHandler) UpdateResource(c echo.Context) error {
	param := &domain.ResourcePositionUpdateResourceReq{}
	if err := c.Bind(param); err != nil {
		return resp.Err(c, err)
	}
	err := h.service.UpdateResource(c.Request().Context(), param)
	if err != nil {
		return resp.Err(c, err)
	}
	return resp.Succ(c, nil)
}

func (h *ResourcePositionHandler) DeleteResource(c echo.Context) error {
	param := &domain.ResourcePositionDeleteResourceReq{}
	if err := c.Bind(param); err != nil {
		return resp.Err(c, err)
	}
	err := h.service.DeleteResource(c.Request().Context(), param)
	if err != nil {
		return resp.Err(c, err)
	}
	return resp.Succ(c, nil)
}
