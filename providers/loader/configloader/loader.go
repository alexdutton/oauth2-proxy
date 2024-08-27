package configloader

import (
	"context"
	"fmt"

	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/apis/options"
	"github.com/oauth2-proxy/oauth2-proxy/v7/providers"
)

type Loader struct {
	providersConf     options.Providers             // providers configuration that has been loaded from file at path loader.conf.ProvidersFile
	providers         map[string]providers.Provider // providers map, key is provider id
	defaultProviderID string
}

func New(conf options.Providers, defaultProviderID string) (*Loader, error) {
	if defaultProviderID == "" {
		defaultProviderID = conf[0].ID
	}
	loader := &Loader{
		providersConf:     conf,
		defaultProviderID: defaultProviderID,
	}
	loader.providers = make(map[string]providers.Provider)

	for _, providerConf := range loader.providersConf {
		provider, err := providers.NewProvider(providerConf)
		if providerConf.ID == "" {
			return nil, fmt.Errorf("provider ID is not provided")
		}
		if err != nil {
			return nil, fmt.Errorf("invalid provider config(id=%s): %s", providerConf.ID, err.Error())
		}
		loader.providers[providerConf.ID] = provider
	}

	return loader, nil
}

func (l *Loader) Load(_ context.Context, id string) (providers.Provider, error) {
	if provider, ok := l.providers[id]; ok {
		return provider, nil
	}
	return nil, fmt.Errorf("no provider found with id='%s'", id)
}

func (l *Loader) GetDefault(_ context.Context) (providers.Provider, error) {
	if provider, ok := l.providers[l.defaultProviderID]; ok {
		return provider, nil
	}
	return nil, fmt.Errorf("no default provider found with id='%s'", l.defaultProviderID)
}
