package app

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/pkg/resp"
)

type CategoryService interface {
	GetCategories(c context.Context, param *domain.CategoriesReq) (*domain.CategoriesResp, error)
}

type CategoryHandler struct {
	service CategoryService
}

func NewCategoryHandler(g *echo.Group, service CategoryService) *CategoryHandler {

	handler := &CategoryHandler{
		service: service,
	}
	g.GET("/categories", handler.GetCategories)
	return handler
}

func (h *CategoryHandler) GetCategories(c echo.Context) error {
	param := &domain.CategoriesReq{}
	if err := c.Bind(param); err != nil {
		return resp.Err(c, err)
	}

	// Assuming you have a service to handle the logic
	categories, err := h.service.GetCategories(c.Request().Context(), param)
	if err != nil {
		return resp.Err(c, err)
	}

	return resp.Succ(c, categories)
}
