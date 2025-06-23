package resp

import (
	"net/http"

	"github.com/usual2970/retell/internal/domain/constant"

	"github.com/pocketbase/pocketbase/core"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Succ(e *core.RequestEvent, data interface{}) error {
	rs := &Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
	return e.JSON(http.StatusOK, rs)
}

func Err(e *core.RequestEvent, err error) error {

	xerr, ok := err.(*constant.XError)
	code := 100
	if ok {
		code = xerr.GetCode()
	}

	rs := &Response{
		Code: code,
		Msg:  err.Error(),
		Data: nil,
	}
	return e.JSON(http.StatusOK, rs)
}

func WecomVerifyUrl(e *core.RequestEvent, msg string) error {

	return e.String(http.StatusOK, msg)
}
