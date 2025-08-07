package utils

import (
	"fmt"
	"net/url"

	"github.com/usual2970/retell/pkg/config"
)

// 解析url，返回url中的path
func ParseUrl(urlStr string) string {
	if urlStr == "" {
		return ""
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	return u.Path

}

func FullUrl(urlStr string) string {
	if urlStr == "" {
		return ""
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	conf := config.GetConfig().OSS

	return fmt.Sprintf("https://%s%s", conf.Domain, u.Path)
}

func FullInternalUrl(urlStr string) string {
	if urlStr == "" {
		return ""
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	conf := config.GetConfig().OSS

	return fmt.Sprintf("https://%s%s", conf.InternalDomain, u.Path)
}
