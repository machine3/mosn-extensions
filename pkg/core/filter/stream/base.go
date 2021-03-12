package stream

import (
	"mosn.io/api"
)

func loadPluginManager() {

}

type GlobalConfig struct {
}

type PluginManager struct {
	config *GlobalConfig
}

func (p *PluginManager) createStreamFilterFactory() api.StreamFilterFactoryCreator {
	return func(conf map[string]interface{}) (api.StreamFilterChainFactory, error) {
		return &streamFilterFactory{
			config: p.config,
		}, nil
	}
}
