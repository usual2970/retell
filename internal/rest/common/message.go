package common

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/domain/constant"
	"github.com/usual2970/retell/internal/rest/middleware"
	"github.com/usual2970/retell/pkg/resp"
)

type MessageService interface {
	Send(ctx context.Context, req *domain.MessageSendReq) (*domain.MessageSendResp, error)
	List(ctx context.Context, req *domain.MessageListReq) (*constant.PaginatedResponse[domain.MessageInfoResp], error)
	BatchSetReaded(ctx context.Context) error
	SetRetentionDays(ctx context.Context, req *domain.MessageSetRetentionDaysReq) error
	GetUnreadNum(ctx context.Context) (int64, error)
	Read(ctx context.Context, req *domain.MessageReadReq) error
}

type MessageHandler struct {
	service MessageService
}

func NewMessageHandler(g *echo.Group, service MessageService) {
	handler := &MessageHandler{
		service: service,
	}

	g.POST("/message/send", handler.Send, middleware.InsideCheck)
	g.GET("/messages", handler.List, middleware.ShouldLogin)
	g.POST("/messages/read-all", handler.BatchSetReaded, middleware.ShouldLogin)
	g.POST("/messages/set-retention-days", handler.SetRetentionDays, middleware.ShouldLogin)
	g.GET("/messages/unread-num", handler.GetUnreadNum, middleware.ShouldLogin)
	g.POST("/messages/read", handler.Read, middleware.ShouldLogin)
}

func (h *MessageHandler) Read(c echo.Context) error {
	req := &domain.MessageReadReq{}
	if err := c.Bind(req); err != nil {
		return resp.Err(c, err)
	}

	if err := h.service.Read(c.Request().Context(), req); err != nil {
		return resp.Err(c, err)
	}

	return resp.Succ(c, nil)
}

func (h *MessageHandler) GetUnreadNum(c echo.Context) error {
	num, err := h.service.GetUnreadNum(c.Request().Context())
	if err != nil {
		return resp.Err(c, err)
	}

	return resp.Succ(c, num)
}

func (h *MessageHandler) BatchSetReaded(c echo.Context) error {

	if err := h.service.BatchSetReaded(c.Request().Context()); err != nil {
		return resp.Err(c, err)
	}

	return resp.Succ(c, nil)
}

func (h *MessageHandler) SetRetentionDays(c echo.Context) error {
	req := &domain.MessageSetRetentionDaysReq{}
	if err := c.Bind(req); err != nil {
		return resp.Err(c, err)
	}

	if err := h.service.SetRetentionDays(c.Request().Context(), req); err != nil {
		return resp.Err(c, err)
	}

	return resp.Succ(c, nil)
}

func (h *MessageHandler) List(c echo.Context) error {
	req := &domain.MessageListReq{}
	if err := c.Bind(req); err != nil {
		return resp.Err(c, err)
	}

	rs, err := h.service.List(c.Request().Context(), req)
	if err != nil {
		return resp.Err(c, err)
	}

	return resp.Succ(c, rs)
}

func (h *MessageHandler) Send(c echo.Context) error {
	req := &domain.MessageSendReq{}
	if err := c.Bind(req); err != nil {
		return resp.Err(c, err)
	}

	rs, err := h.service.Send(c.Request().Context(), req)
	if err != nil {
		return resp.Err(c, err)
	}

	return resp.Succ(c, rs)
}
