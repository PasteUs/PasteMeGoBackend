package http_client

import (
	"testing"
)

func TestGet(t *testing.T) {
	type args struct {
		IP     string
		rawUrl string
		params map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"100",
			args{
				"127.0.0.1",
				"http://api.pasteme.cn/101,123456",
				map[string]string{},
			},
			"加密测试",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(tt.args.IP, tt.args.rawUrl, tt.args.params); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
