package kafka

import (
	"context"
	"testing"

	"github.com/usual2970/retell/pkg/config"
)

func TestPushToPartition(t *testing.T) {

	type args struct {
		ctx       context.Context
		topic     string
		partition int
		key       []byte
		value     []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				ctx:       context.Background(),
				topic:     "test-1",
				partition: 0,
				key:       []byte("test-key1"),
				value:     []byte("test-value"),
			},
		},
	}

	SetupDefaultClient(config.Kafka{
		Brokers:  []string{"127.0.0.1:29092"},
		Username: "paas_user",
		Password: "CjGWGK1wv3",
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PushToPartition(tt.args.ctx, tt.args.topic, tt.args.partition, tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("PushToPartition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
