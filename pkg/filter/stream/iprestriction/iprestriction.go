package iprestriction

import (
	"context"
	"github.com/machine3/mosn-extensions/pkg/filter/stream/iprestriction/ipmatcher"
	"mosn.io/api"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
	"net/http"
)

const (
	// filterName is the ip restriction stream filter name.
	filterName = "ip_restriction"

	clientIpHeaderName = "X-Real-Ip"

	IPNotAllowed api.ResponseFlag = 0x4000
)

// streamFilter represents the ip restriction stream filter.
type streamFilter struct {
	config  *config
	handler api.StreamReceiverFilterHandler
}

// newStreamFilter creates ip restriction filter.
func newStreamFilter(conf *config) *streamFilter {
	return &streamFilter{config: conf}
}

func (f *streamFilter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	f.handler = handler
}

// OnReceive creates resource and judges whether current request should be blocked.
func (f *streamFilter) OnReceive(ctx context.Context, headers api.HeaderMap,
	buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	conf := f.config
	clientIp, ok := headers.Get(clientIpHeaderName)
	if !ok {
		log.DefaultContextLogger.Warnf(ctx, "Cannot get the client IP from the header %s", clientIpHeaderName)
		return api.StreamFilterContinue
	}

	if len(conf.Deny) > 0 {
		blocked, err := matchClientIp(conf.Deny, clientIp)
		if err != nil {
			log.DefaultContextLogger.Errorf(ctx, "match client ip error:%v", err)
			return api.StreamFilterContinue
		}
		if blocked {
			f.handler.RequestInfo().SetResponseFlag(IPNotAllowed)
			f.handler.SendHijackReplyWithBody(http.StatusForbidden, headers, "Your IP address is not allowed")
			return api.StreamFilterStop
		}
	}

	if len(conf.Allow) > 0 {
		allowed, err := matchClientIp(conf.Allow, clientIp)
		if err != nil {
			log.DefaultContextLogger.Errorf(ctx, "match client ip error:%v", err)
			return api.StreamFilterContinue
		}
		if !allowed {
			f.handler.RequestInfo().SetResponseFlag(IPNotAllowed)
			f.handler.SendHijackReplyWithBody(http.StatusForbidden, headers, "Your IP address is not allowed")
			return api.StreamFilterStop
		}
	}
	return api.StreamFilterContinue
}

// OnDestroy does some exit tasks.
func (f *streamFilter) OnDestroy() {
}

func matchClientIp(list []string, ip string) (bool, error) {
	if matcher, err := ipmatcher.New(list); err == nil {
		return matcher.Match(ip)
	} else {
		return false, err
	}
}
