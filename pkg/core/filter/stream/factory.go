package stream

import (
	"context"
	"encoding/json"
	"mosn.io/api"
	"mosn.io/mosn/pkg/cel/extract"
	"mosn.io/mosn/pkg/log"
	"mosn.io/pkg/buffer"
	"runtime/debug"
)

var plugins map[string]Plugin

func init() {
	for pluginName, plugin := range plugins {
		api.RegisterStream(pluginName, plugin.createStreamFilterFactory())
	}
}

// streamFilterFactory represents the stream filter factory.
type streamFilterFactory struct {
	config *Config
}

// CreateFilterChain add the ip restriction stream filter to filter chain.
func (f *streamFilterFactory) CreateFilterChain(ctx context.Context,
	callbacks api.StreamFilterChainFactoryCallbacks) {
	filter := newStreamFilter(ctx, f.config)
	//callbacks.AddStreamReceiverFilter(filter, api.BeforeRoute)
	callbacks.AddStreamReceiverFilter(filter, api.AfterRoute)
	//callbacks.AddStreamReceiverFilter(filter, api.AfterChooseHost)
	callbacks.AddStreamSenderFilter(filter, api.BeforeSend)
	callbacks.AddStreamAccessLog(filter)
}

// newStreamFilter creates ip restriction filter.
func newStreamFilter(config *Config) *streamFilter {
	filter := streamFilter{config: config}
	return &filter
}

type PluginHandler interface {
	Init()
	OnRequest() api.StreamFilterStatus
	OnResponse() api.StreamFilterStatus
	OnDestroy()
	Log()
}

// streamFilter represents the ip restriction stream filter.
type streamFilter struct {
	context context.Context
	plugin  PluginHandler

	receiver api.StreamReceiverFilterHandler
	sender   api.StreamSenderFilterHandler
}

func (f *streamFilter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	return f.plugin.OnRequest()
}

func (f *streamFilter) Append(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	return f.plugin.OnResponse()
}

func (f *streamFilter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	f.receiver = handler
}

func (f *streamFilter) SetSenderFilterHandler(handler api.StreamSenderFilterHandler) {
	f.sender = handler
}

func (f *streamFilter) OnDestroy() {
	f.plugin.OnDestroy()
}

func (f *streamFilter) Log(ctx context.Context, reqHeaders api.HeaderMap, respHeaders api.HeaderMap, requestInfo api.RequestInfo) {
	f.plugin.Log()
}
