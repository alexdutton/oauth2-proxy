package configloader

import (
	"context"
	"fmt"

	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/apis/options"
	"github.com/oauth2-proxy/oauth2-proxy/v7/providers"
)

type Loader struct {
	providersConf     options.Providers             // providers configuration that has been loaded from file at path loader.conf.ProvidersFile
	providersByID     map[string]providers.Provider // providers map, key is provider id
	providersByIssuer map[string]providers.Provider // providers map, key is provider issuer
}

type ProviderWithIssuerURL interface {
	GetIssuerURL() string
}

func New(conf options.Providers) (*Loader, error) {
	loader := &Loader{
		providersConf: conf,
	}
	loader.providersByID = make(map[string]providers.Provider)
	loader.providersByIssuer = make(map[string]providers.Provider)

	for _, providerConf := range loader.providersConf {
		provider, err := providers.NewProvider(providerConf)
		if providerConf.ID == "" {
			return nil, fmt.Errorf("provider ID is not provided")
		}
		if err != nil {
			return nil, fmt.Errorf("invalid provider config(id=%s): %s", providerConf.ID, err.Error())
		}
		loader.providersByID[providerConf.ID] = provider
		providerWithIssuerURL, ok := provider.(ProviderWithIssuerURL)
		if ok {
			loader.providersByIssuer[providerWithIssuerURL.GetIssuerURL()] = provider
		}

	}

	return loader, nil
}

func (l *Loader) Load(_ context.Context, id string) (providers.Provider, error) {
	if provider, ok := l.providersByID[id]; ok {
		return provider, nil
	}
	return nil, fmt.Errorf("no provider found with id='%s'", id)
}

func (l *Loader) LoadByIssuer(_ context.Context, issuer string) (providers.Provider, error) {
	if provider, ok := l.providersByIssuer[issuer]; ok {
		return provider, nil
	}
	return nil, fmt.Errorf("no provider found with issuer='%s'", issuer)
}

func (l *Loader) List(_ context.Context) []providers.Provider {
	providers := make([]providers.Provider, len(l.providersConf))
	for i, providerConf := range l.providersConf {
		providers[i] = l.providersByID[providerConf.ID]
	}
	return providers
}
