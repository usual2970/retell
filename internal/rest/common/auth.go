package common

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/internal/rest/middleware"
	"github.com/usual2970/retell/pkg/resp"
)

type AuthService interface {
	Login(ctx context.Context, param *domain.AuthLoginReq) (*domain.AuthLoginResp, error)
	Register(ctx context.Context, param *domain.AuthRegisterReq) (*domain.AuthRegisterResp, error)
	Logout(ctx context.Context) error
	LoginCode(ctx context.Context, param *domain.AuthCodeReq) error
	RegisterCode(ctx context.Context, param *domain.AuthCodeReq) error

	ForgetPasswordCode(ctx context.Context, param *domain.AuthCodeReq) error
	ForgetPassword(ctx context.Context, param *domain.AuthForgetPasswordReq) error
	ResetPassword(ctx context.Context, param *domain.AuthResetPasswordReq) error
	ResetHeadImg(ctx context.Context, param *domain.AuthResetHeadImgReq) error
}

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(g *echo.Group, service AuthService) {
	handler := &AuthHandler{
		service: service,
	}
	g.POST("/auth/login", handler.Login)
	g.POST("/auth/register", handler.Register)
	g.POST("/auth/logout", handler.Logout, middleware.ShouldLogin)
	g.POST("/auth/login-code", handler.LoginCode)
	g.POST("/auth/register-code", handler.RegisterCode)

	g.POST("/auth/forget-password-code", handler.ForgetPasswordCode)
	g.POST("/auth/forget-password", handler.ForgetPassword)
	g.POST("/auth/reset-password", handler.ResetPassword, middleware.ShouldLogin)
	g.POST("/auth/reset-headimg", handler.ResetHeadImg, middleware.ShouldLogin)

}

func (h *AuthHandler) ResetHeadImg(c echo.Context) error {
	param := &domain.AuthResetHeadImgReq{}
	if err := c.Bind(param); err != nil {
		return resp.Err(c, err)
	}

	if err := h.service.ResetHeadImg(c.Request().Context(), param); err != nil {
		return resp.Err(c, err)
	}

	return resp.Succ(c, nil)
}

func (h *AuthHandler) ForgetPasswordCode(c echo.Context) error {
	param := &domain.AuthCodeReq{}
	if err := c.Bind(param); err != nil {
		return resp.Err(c, err)
	}

	if err := h.service.ForgetPasswordCode(c.Request().Context(), param); err != nil {
		return resp.Err(c, err)
	}

	return resp.Succ(c, nil)
}

func (h *AuthHandler) ForgetPassword(c echo.Context) error {
	param := &domain.AuthForgetPasswordReq{}
	if err := c.Bind(param); err != nil {
		return resp.Err(c, err)
	}

	if err := h.service.ForgetPassword(c.Request().Context(), param); err != nil {
		return resp.Err(c, err)
	}

	return resp.Succ(c, nil)
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	param := &domain.AuthResetPasswordReq{}
	if err := c.Bind(param); err != nil {
		return resp.Err(c, err)
	}

	if err := h.service.ResetPassword(c.Request().Context(), param); err != nil {
		return resp.Err(c, err)
	}

	return resp.Succ(c, nil)
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

func (h *AuthHandler) RegisterCode(c echo.Context) error {
	param := &domain.AuthCodeReq{}
	if err := c.Bind(param); err != nil {
		return resp.Err(c, err)
	}

	if err := h.service.RegisterCode(c.Request().Context(), param); err != nil {
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

func (h *AuthHandler) Register(c echo.Context) error {

	param := &domain.AuthRegisterReq{}
	if err := c.Bind(param); err != nil {
		return resp.Err(c, err)
	}
	rs, err := h.service.Register(c.Request().Context(), param)
	if err != nil {
		return resp.Err(c, err)
	}
	return resp.Succ(c, rs)
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
