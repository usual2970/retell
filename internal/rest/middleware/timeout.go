package middleware

import (
	"context"
	"strings"
	"time"

	echo "github.com/labstack/echo/v4"
)

var skipPaths = []string{"/pushlet", "/developer/v1/file"}

// SetRequestContextWithTimeout will set the request context with timeout for every incoming HTTP Request
func SetRequestContextWithTimeout(d time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			path := c.Request().URL.Path
			for _, skipPath := range skipPaths {
				if strings.HasPrefix(path, skipPath) {
					return next(c) // Skip setting context for paths in skipPaths
				}
			}

			ctx, cancel := context.WithTimeout(c.Request().Context(), d)
			defer cancel()

			newRequest := c.Request().WithContext(ctx)
			c.SetRequest(newRequest)
			return next(c)
		}
	}
}
