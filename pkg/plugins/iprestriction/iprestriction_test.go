package iprestriction

import (
	"context"
	"mosn.io/api"
	"mosn.io/pkg/buffer"
	"reflect"
	"testing"
)

func Test_matchClientIp(t *testing.T) {
	type args struct {
		list []string
		ip   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchClientIp(tt.args.list, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("matchClientIp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("matchClientIp() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newStreamFilter(t *testing.T) {
	type args struct {
		conf *config
	}
	tests := []struct {
		name string
		args args
		want *streamFilter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newStreamFilter(tt.args.conf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newStreamFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_streamFilter_OnDestroy(t *testing.T) {
}

func Test_streamFilter_OnReceive(t *testing.T) {
	type fields struct {
		config  *config
		handler api.StreamReceiverFilterHandler
	}
	type args struct {
		ctx      context.Context
		headers  api.HeaderMap
		buf      buffer.IoBuffer
		trailers api.HeaderMap
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   api.StreamFilterStatus
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &streamFilter{
				config:  tt.fields.config,
				handler: tt.fields.handler,
			}
			if got := f.OnReceive(tt.args.ctx, tt.args.headers, tt.args.buf, tt.args.trailers); got != tt.want {
				t.Errorf("OnReceive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_streamFilter_SetReceiveFilterHandler(t *testing.T) {
}
