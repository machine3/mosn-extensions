package ipmatcher

import (
	"errors"
	"net"
	"strings"
)

func New(ipStrs []string) (*ipMatcher, error) {
	if len(ipStrs) == 0 {
		return nil, errors.New("missing valid ip argument")
	}
	ips := make([]*net.IP, 0, len(ipStrs))
	ipNets := make([]*net.IPNet, 0, len(ipStrs))
	for i := range ipStrs {
		if strings.Contains(ipStrs[i], "/") {
			_, ipNet, _ := net.ParseCIDR(ipStrs[i])
			if ipNet != nil {
				ipNets = append(ipNets, ipNet)
			}
		} else {
			ip := net.ParseIP(ipStrs[i])
			if ip != nil {
				ips = append(ips, &ip)
			}
		}
	}
	if len(ips) == 0 && len(ipNets) == 0 {
		return nil, errors.New("no valid ips or ipNets")
	}
	return &ipMatcher{
		ips:    ips,
		ipNets: ipNets,
	}, nil
}

type ipMatcher struct {
	ips    []*net.IP
	ipNets []*net.IPNet
}

func (m *ipMatcher) Match(ipStr string) (bool, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, errors.New("ip format error")
	}
	for i := range m.ips {
		if m.ips[i].Equal(ip) {
			return true, nil
		}
	}
	for i := range m.ipNets {
		if m.ipNets[i].Contains(ip) {
			return true, nil
		}
	}
	return false, nil
}
