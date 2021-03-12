package ipmatcher

import (
	"net"
	"reflect"
	"testing"
)

var (
	ipStr1       = "127.0.0.1"
	ipStr2       = "2001:DB8:2de:0:0:0:0:e13"
	ipNetStr1    = "192.168.0.0/16"
	ipNetStr2    = "15ba:db5::/64"
	ip1          = net.ParseIP(ipStr1)
	ip2          = net.ParseIP(ipStr2)
	_, ipNet1, _ = net.ParseCIDR(ipNetStr1)
	_, ipNet2, _ = net.ParseCIDR(ipNetStr2)
)

func TestNew(t *testing.T) {
	type args struct {
		ipStrs []string
	}
	tests := []struct {
		name    string
		args    args
		want    *ipMatcher
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no valid addr",
			args: args{
				[]string{"1111", "1.1"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "normal",
			args: args{
				[]string{ipStr1, ipStr2, ipNetStr1, ipNetStr2},
			},
			want: &ipMatcher{
				ips:    []*net.IP{&ip1, &ip2},
				ipNets: []*net.IPNet{ipNet1, ipNet2},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.ipStrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ipMatcher_Match(t *testing.T) {
	type fields struct {
		ips    []*net.IP
		ipNets []*net.IPNet
	}
	type args struct {
		ipStr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "match ip",
			fields: fields{
				ips:    []*net.IP{&ip1, &ip2},
				ipNets: []*net.IPNet{ipNet1, ipNet2},
			},
			args: args{
				ipStr: ipStr2,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "match cidr",
			fields: fields{
				ips:    []*net.IP{&ip1, &ip2},
				ipNets: []*net.IPNet{ipNet1, ipNet2},
			},
			args: args{
				ipStr: "192.168.0.5",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "not match",
			fields: fields{
				ips:    []*net.IP{&ip1, &ip2},
				ipNets: []*net.IPNet{ipNet1, ipNet2},
			},
			args: args{
				ipStr: "192.167.0.5",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ipMatcher{
				ips:    tt.fields.ips,
				ipNets: tt.fields.ipNets,
			}
			got, err := m.Match(tt.args.ipStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Match() got = %v, want %v", got, tt.want)
			}
		})
	}
}
