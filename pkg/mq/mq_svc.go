package mq

import (
	"context"
	"sync"

	"github.com/hibiken/asynq"
	"github.com/usual2970/retell/pkg/config"
)

var svc *asynq.Server

var svcOnce sync.Once

var mux *asynq.ServeMux

var muxOnce sync.Once

func getMux() *asynq.ServeMux {
	muxOnce.Do(func() {
		mux = asynq.NewServeMux()
	})
	return mux
}

func GetSvc() *asynq.Server {

	svcOnce.Do(func() {
		conf := config.GetConfig().Redises["paas"]
		address := conf.Host + ":" + conf.Port
		password := conf.Password
		svc = asynq.NewServer(
			asynq.RedisClientOpt{
				Addr:     address,
				Password: password,
				DB:       1,
			},
			asynq.Config{Concurrency: 10},
		)
	})
	return svc
}

type Handler func(ctx context.Context, payload []byte) error

func Handler2Asynq(handler Handler) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		payload := t.Payload()
		return handler(ctx, payload)
	}
}

func AddHandler(topic string, handler Handler) {
	mux := getMux()
	mux.HandleFunc(topic, Handler2Asynq(handler))
}

func Run() error {
	svc := GetSvc()

	mux := getMux()
	return svc.Run(mux)
}

func Shutdown() {
	svc := GetSvc()

	svc.Shutdown()
}
