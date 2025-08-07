package common

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/pkg/resp"
)

type CodeService interface {
	GenerateCode(ctx context.Context, param *domain.CodeGenerateReq) (*domain.Code, error)
	CheckCode(ctx context.Context, param *domain.CodeCheckReq) error
}

type CodeHandler struct {
	codeSvc CodeService
}

func NewCodeHandler(g *echo.Group, codeSvc CodeService) {
	handler := &CodeHandler{
		codeSvc: codeSvc,
	}

	g.POST("/code/generate", handler.GenerateCode)
	g.POST("/code/check", handler.CheckCode)
}

func (h *CodeHandler) GenerateCode(c echo.Context) error {
	param := &domain.CodeGenerateReq{}
	if err := c.Bind(param); err != nil {
		return resp.Err(c, err)
	}

	code, err := h.codeSvc.GenerateCode(c.Request().Context(), param)
	if err != nil {
		return resp.Err(c, err)
	}

	return resp.Succ(c, code)
}

func (h *CodeHandler) CheckCode(c echo.Context) error {
	param := &domain.CodeCheckReq{}
	if err := c.Bind(param); err != nil {
		return resp.Err(c, err)
	}

	if err := h.codeSvc.CheckCode(c.Request().Context(), param); err != nil {
		return resp.Err(c, err)
	}

	return resp.Succ(c, nil)
}
