package routing

import (
	"fmt"
	"github.com/openfaas/faas/gateway/requests"
	"net/url"
	"strings"
	"sync"
)

const federationProviderNameConstraint = "federation.provider_name"

// ProviderLookup allows the federation to determine which provider
// is currently responsible for a given function
type ProviderLookup interface {
	Resolve(functionName string) (providerUri *url.URL, err error)
	AddFunction(f *requests.CreateFunctionRequest)
	GetFunction(name string) (*requests.CreateFunctionRequest, bool)
}

type defaultProviderRouting struct {
	cache map[string] *requests.CreateFunctionRequest
	providers map[string]*url.URL
	defaultProvider *url.URL
	lock sync.RWMutex
}

func NewDefaultProviderRouting(providers []string, defaultProvider string) (ProviderLookup, error) {
	providerMap := map[string]*url.URL{}

	for _, v := range providers {
		pURL, err := url.Parse(v)
		if err != nil {
			return nil, fmt.Errorf("error parsing URL using value %s. %v", v, err)
		}
		providerMap[getHostName(pURL)] = pURL
	}

	d, err := url.Parse(defaultProvider)
	if err != nil {
		return nil, fmt.Errorf("error parsing default provider URL using value %s. %v", defaultProvider, err)
	}

	return &defaultProviderRouting{
		cache:make(map[string]*requests.CreateFunctionRequest),
		providers: providerMap,
		defaultProvider: d,
	}, nil
}


func (d *defaultProviderRouting) Resolve(functionName string) (providerUri *url.URL, err error) {
	f, ok := d.GetFunction(functionName)
	if !ok {
		return nil, fmt.Errorf("can not find function %s in cache map", functionName)
	}

	c, ok := (*f.Annotations)[federationProviderNameConstraint]
	if !ok {
		return d.defaultProvider, nil
	}

	pURL, ok := d.providers[c]
	if !ok {
		return d.defaultProvider, nil
	}

	return pURL, nil
}

func getHostName(v *url.URL) string {
	return strings.Split(v.Host, ":")[0]
}

func (d *defaultProviderRouting) AddFunction(f *requests.CreateFunctionRequest) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.cache[f.Service] = f
}

func (d *defaultProviderRouting) GetFunction(name string) (*requests.CreateFunctionRequest, bool) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	v, ok := d.cache[name]


	return v, ok
}