package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/usual2970/retell/pkg/config"
)

// InternalCheck 中间件，检查请求是否来自内部域名
func InternalCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 获取内部域名配置
		internalDomain := config.GetConfig().PaasInternalDomain
		internalU, _ := url.Parse(internalDomain)
		internalDomain = internalU.Host
		if colonIndex := strings.Index(internalDomain, ":"); colonIndex != -1 {
			internalDomain = internalDomain[:colonIndex]
		}

		// 获取当前请求的主机名
		host := c.Request().Host

		// 从主机名中提取域名部分（移除端口号）
		domain := host
		if colonIndex := strings.Index(host, ":"); colonIndex != -1 {
			domain = host[:colonIndex]
		}

		// 检查请求域名是否与内部域名匹配
		if !strings.EqualFold(domain, internalDomain) {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Unauthorized: Access restricted to internal domain",
			})
		}

		// 域名匹配，允许请求通过
		return next(c)
	}
}
