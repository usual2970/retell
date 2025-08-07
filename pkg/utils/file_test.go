package utils

import (
	"testing"
)

func TestParseUrl(t *testing.T) {
	type args struct {
		urlStr string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				urlStr: "https://www.baidu.com/uploads/1234/ss.jpg",
			},
			want: "/uploads/1234/ss.jpg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseUrl(tt.args.urlStr); got != tt.want {
				t.Errorf("ParseUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFullUrl(t *testing.T) {
	type args struct {
		urlStr string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				urlStr: "https://www.baidu.com/uploads/1234/ss.jpg",
			},
			want: "/uploads/1234/ss.jpg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FullUrl(tt.args.urlStr); got != tt.want {
				t.Errorf("FullUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
