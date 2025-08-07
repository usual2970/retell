package mq

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/usual2970/retell/pkg/config"
	"github.com/usual2970/retell/pkg/logger"
)

var client *asynq.Client

var clientOnce sync.Once

func getClient() *asynq.Client {
	conf := config.GetConfig().Redises["paas"]
	address := conf.Host + ":" + conf.Port
	password := conf.Password
	clientOnce.Do(func() {
		client = asynq.NewClient(asynq.RedisClientOpt{
			Addr:     address,
			Password: password,
			DB:       1,
		})
	})

	return client
}

type ProduceOption struct {
	Delay    time.Duration
	TaskId   string
	MaxRetry int
}

type Option func(*ProduceOption)

func WithDuration(duration time.Duration) Option {
	return func(o *ProduceOption) {
		o.Delay = duration
	}
}

func WithTaskId(taskId string) Option {
	return func(o *ProduceOption) {
		o.TaskId = taskId
	}
}

func WithMaxRetry(maxRetry int) Option {
	return func(o *ProduceOption) {
		o.MaxRetry = maxRetry
	}
}

func Produce(topic string, data interface{}, option ...Option) error {
	cli := getClient()
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	task := asynq.NewTask(topic, payload)

	produceOption := &ProduceOption{}

	for _, op := range option {
		op(produceOption)
	}

	ops := make([]asynq.Option, 0)

	if produceOption.Delay > 0 {
		ops = append(ops, asynq.ProcessIn(produceOption.Delay))
	}

	if produceOption.MaxRetry > 0 {
		ops = append(ops, asynq.MaxRetry(produceOption.MaxRetry))
	}

	if produceOption.TaskId != "" {
		ops = append(ops, asynq.TaskID(produceOption.TaskId))
	}

	info, err := cli.Enqueue(task, ops...)

	logger.WithField("module", "mq_product").
		WithField("task_info", info).
		WithField("param", []interface{}{topic, data}).Info("produce task")

	return err
}
