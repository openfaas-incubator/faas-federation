package routing

import (
	"fmt"
	"github.com/openfaas/faas/gateway/requests"
)

const federationProviderNameConstraint = "federation.provider_name"
// ProviderLookup allows the federation to determine which provider
// is currently responsible for a given function
type ProviderLookup interface {
	Resolve(functionName string) (providerHostName string, err error)
}

type defaultProviderRouting struct {
	cache map[string] *requests.CreateFunctionRequest
	providers map[string]bool
	defaultProvider string
}

func NewDefaultProviderRouting(providers []string, defaultProvider string) ProviderLookup {
	providerMap := map[string]bool{}

	for _, v := range providers {
		providerMap[v] = true
	}

	return &defaultProviderRouting{
		cache:make(map[string]*requests.CreateFunctionRequest),
		providers: providerMap,
		defaultProvider: defaultProvider,
	}
}


func (d *defaultProviderRouting) Resolve(functionName string) (providerHostName string, err error) {
	v, ok := d.cache[functionName]
	if !ok {
		return "", fmt.Errorf("can not find function %s in cache map", functionName)
	}

	c, ok := (*v.Annotations)[federationProviderNameConstraint]
	if !ok {
		return d.defaultProvider, nil
	}

	_, ok = d.providers[c]
	if !ok {
		return d.defaultProvider, nil
	}

	return c, nil
}