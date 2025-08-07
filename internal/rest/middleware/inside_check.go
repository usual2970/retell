package middleware

import (
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/usual2970/retell/pkg/config"
	"github.com/usual2970/retell/pkg/resp"
)

func InsideCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auths := strings.SplitN(c.Request().Header.Get(authorizationKey), " ", 2)
		if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
			return resp.Err(c, jwt.ErrTokenMalformed)
		}
		token := auths[1]

		if token != config.GetConfig().InsideToken {
			return resp.Err(c, jwt.ErrTokenMalformed)
		}
		return next(c)
	}
}
