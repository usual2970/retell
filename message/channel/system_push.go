package channel

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/pkg/sse"
)

type SystemPush struct {
	topic         string
	channelConfig *ChannelConfig
}

func NewSystemPush(
	topic string,
	channelConfig *ChannelConfig,
) *SystemPush {
	return &SystemPush{
		topic:         topic,
		channelConfig: channelConfig,
	}
}

type SystemPushData struct {
	Level   domain.MessageLevel `json:"level"`
	Title   string              `json:"title"`
	Content string              `json:"content"`
}

func (s *SystemPushData) String() string {
	bts, _ := json.Marshal(s)

	return string(bts)
}

func (s *SystemPush) Send(ctx context.Context, title, content string) error {
	data := &SystemPushData{
		Level:   s.channelConfig.Task.Level, // 默认级别为1
		Title:   title,
		Content: content,
	}

	if err := sse.Publish(s.topic, data.String()); err != nil {
		return errors.New("failed to publish system push message" + err.Error())
	}
	return nil
}
