package encode

import "testing"

func TestToPinyin(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				text: "测试",
			},
			want: "ceshi",
		},
		{
			name: "test1",
			args: args{
				text: "test",
			},
			want: "test",
		},
		{
			name: "test2",
			args: args{
				text: "te测试st",
			},
			want: "ceshi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToPinyin(tt.args.text); got != tt.want {
				t.Errorf("ToPinyin() = %v, want %v", got, tt.want)
			}
		})
	}
}
