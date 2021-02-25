package iprestriction

import (
	"context"
	"encoding/json"
	"mosn.io/api"
)

func init() {
	api.RegisterStream(filterName, createStreamFilterFactory)
}

func createStreamFilterFactory(conf map[string]interface{}) (api.StreamFilterChainFactory, error) {
	c, err := parseToConfig(conf)
	if err != nil {
		return nil, err
	}
	return &streamFilterFactory{
		config: c,
	}, nil
}

// streamFilterFactory represents the stream filter factory.
type streamFilterFactory struct {
	config *config
}

// CreateFilterChain add the ip restriction stream filter to filter chain.
func (f *streamFilterFactory) CreateFilterChain(context context.Context,
	callbacks api.StreamFilterChainFactoryCallbacks) {
	filter := newStreamFilter(f.config)
	callbacks.AddStreamReceiverFilter(filter, api.AfterRoute)
}

// parseToConfig
func parseToConfig(conf map[string]interface{}) (*config, error) {
	c := &config{}
	b, err := json.Marshal(conf)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, c); err != nil {
		return nil, err
	}
	return c, nil
}
