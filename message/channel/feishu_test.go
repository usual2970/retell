package channel

import (
	"context"
	"testing"
)

func TestFeishu_Send(t *testing.T) {
	type fields struct {
		botToken string
	}
	type args struct {
		ctx     context.Context
		title   string
		content string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test Feishu Send",
			fields: fields{
				botToken: "95e57567-7652-49cb-acfa-90f7d3b286be",
			},
			args: args{
				ctx:     context.Background(),
				title:   "",
				content: `{"msg_type":"post","content":{"post":{"zh_cn":{"title":"项目更新通知","content":[[{"tag":"text","text":"项目有更新: "},{"tag":"a","text":"请查看","href":"http://www.example.com/"}]]}}}}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Feishu{
				botToken: tt.fields.botToken,
			}
			if err := f.Send(tt.args.ctx, tt.args.title, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("Feishu.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
