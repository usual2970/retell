package file

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/usual2970/retell/domain"
	"github.com/usual2970/retell/internal/repository"
)

func TestService_Structure(t *testing.T) {
	type fields struct {
		fileRepo FileRepository
	}
	type args struct {
		ctx context.Context
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.FileStructureResp
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				fileRepo: repository.NewFileRepository(),
			},
			args: args{
				ctx: context.Background(),
				url: "https://ikit-blog.oss-cn-hangzhou.aliyuncs.com/uploads/2025/03/12/2503121320ckUROx.zip",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				fileRepo: tt.fields.fileRepo,
			}
			got, err := s.Structure(tt.args.ctx, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Structure() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			rs, _ := json.Marshal(got)
			t.Log(string(rs))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Structure() = %v, want %v", got, tt.want)
			}
		})
	}
}
