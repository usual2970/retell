package common

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/internal/rest/middleware"
	"github.com/usual2970/retell/pkg/resp"
)

type SettingService interface {
	GetSetting(ctx context.Context, key string) (*domain.SettingResp, error)
	GetUserSetting(ctx context.Context, key string) (*domain.SettingResp, error)
}

type SettingHandler struct {
	service SettingService
}

func NewSettingHandler(g *echo.Group, service SettingService) {
	handler := &SettingHandler{
		service: service,
	}
	g.GET("/setting/:key", handler.GetSetting)
	g.GET("/setting/user/:key", handler.GetUserSetting, middleware.ShouldLogin)
}

func (h *SettingHandler) GetUserSetting(c echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return echo.NewHTTPError(400, "key is required")
	}

	rs, err := h.service.GetUserSetting(c.Request().Context(), key)
	if err != nil {
		return resp.Err(c, err)
	}

	return resp.Succ(c, rs)
}

func (h *SettingHandler) GetSetting(c echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return echo.NewHTTPError(400, "key is required")
	}

	rs, err := h.service.GetSetting(c.Request().Context(), key)
	if err != nil {
		return resp.Err(c, err)
	}

	return resp.Succ(c, rs)
}
