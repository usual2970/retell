package sse

import (
	"errors"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/usual2970/pushlet"
	"github.com/usual2970/retell/pkg/config"
)

type EventType string

const MessageEvent EventType = "message"

var p *pushlet.Pushlet
var pOnce sync.Once

func InitPushlet(e *echo.Echo, path string) {
	pOnce.Do(func() {
		p = pushlet.New()

		conf := config.GetConfig().Redises
		if len(conf) == 0 {
			panic("pushlet redis config is empty, please check your configuration")
		}

		redisConf, ok := conf["paas"]
		if !ok {
			panic("pushlet redis config for 'paas' not found, please check your configuration")
		}

		p.EnableDistributedMode(redisConf.Host+":"+redisConf.Port, redisConf.Password, 2)

		p.Start()

		e.GET(path, func(c echo.Context) error {
			p.HandleSSE(c.Response(), c.Request())
			return nil
		})

	})
}

func Publish(topic, data string) error {
	if p == nil {
		return errors.New("pushlet not initialized, call InitPushlet first")
	}
	p.Publish(topic, string(MessageEvent), data)
	return nil
}
