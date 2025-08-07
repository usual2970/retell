package rest

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/usual2970/retell/domain/constant"
)

func GetPagination(c echo.Context) *constant.Pagination {
	page := 1
	pageSize := 10
	if p := c.QueryParam("page"); p != "" {
		page, _ = strconv.Atoi(p)
	}
	if p := c.QueryParam("pageSize"); p != "" {
		pageSize, _ = strconv.Atoi(p)
	}
	return &constant.Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}
