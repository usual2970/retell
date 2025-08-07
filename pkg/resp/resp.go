package resp

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/usual2970/retell/domain/constant"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Succ(e echo.Context, data any) error {
	rs := &Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
	return e.JSON(http.StatusOK, rs)
}

func Err(e echo.Context, err error) error {
	xerr, ok := err.(*constant.XError)
	code := 100
	if ok {
		code = xerr.GetCode()
	}

	rs := &Response{
		Code:    code,
		Message: err.Error(),
		Data:    nil,
	}
	return e.JSON(http.StatusOK, rs)
}
