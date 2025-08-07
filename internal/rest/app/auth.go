package app

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/internal/rest/middleware"
	"github.com/usual2970/retell/pkg/resp"
)

type AuthService interface {
	Login(ctx context.Context, param *domain.AuthLoginReq) (*domain.AuthLoginResp, error)
	Logout(ctx context.Context) error
	LoginCode(ctx context.Context, param *domain.AuthCodeReq) error
}

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(g *echo.Group, service AuthService) {
	handler := &AuthHandler{
		service: service,
	}
	g.POST("/auth/login", handler.Login)
	g.POST("/auth/logout", handler.Logout, middleware.ShouldLogin)
	g.POST("/auth/login-code", handler.LoginCode)
}

func (h *AuthHandler) LoginCode(c echo.Context) error {
	param := &domain.AuthCodeReq{}
	if err := c.Bind(param); err != nil {
		return resp.Err(c, err)
	}

	if err := h.service.LoginCode(c.Request().Context(), param); err != nil {
		return resp.Err(c, err)
	}

	return resp.Succ(c, nil)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	if err := h.service.Logout(c.Request().Context()); err != nil {
		return resp.Err(c, err)
	} else {
		return resp.Succ(c, nil)
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	param := &domain.AuthLoginReq{}
	if err := c.Bind(param); err != nil {
		return resp.Err(c, err)
	}

	if err := c.Validate(param); err != nil {
		return resp.Err(c, err)
	}

	rs, err := h.service.Login(c.Request().Context(), param)
	if err != nil {
		return resp.Err(c, err)
	}
	return resp.Succ(c, rs)
}
