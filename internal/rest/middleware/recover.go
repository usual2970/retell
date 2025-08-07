package middleware

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/labstack/echo/v4"
	"github.com/usual2970/retell/pkg/logger"
	"github.com/usual2970/retell/pkg/resp"
)

// 如果出现了panic,处理成正常的错误返回
func Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				// 日志记录panic的信息, 并附带stack
				logger.WithFields(map[string]interface{}{
					"stack": string(debug.Stack()),
					"error": fmt.Sprintf("%v", r),
				}).Error("panic")
				resp.Err(c, errors.New("internal server error"))
			}
		}()
		return next(c)
	}
}
