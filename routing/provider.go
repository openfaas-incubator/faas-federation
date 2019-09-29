package routing

import (
	"fmt"
	"net/url"
	"strings"
	"sync"

	types "github.com/openfaas/faas-provider/types"
	log "github.com/sirupsen/logrus"
)

const federationProviderNameConstraint = "com.openfaas.federation.gateway"

// ProviderLookup allows the federation to determine which provider
// is currently responsible for a given function
type ProviderLookup interface {
	Resolve(functionName string) (providerURI *url.URL, err error)
	AddFunction(f *types.FunctionDeployment)
	GetFunction(name string) (*types.FunctionDeployment, bool)
	GetFunctions() []*types.FunctionDeployment
	ReloadCache() error
}

type defaultProviderRouting struct {
	cache           map[string]*types.FunctionDeployment
	providers       map[string]*url.URL
	defaultProvider *url.URL
	lock            sync.RWMutex
}

// NewDefaultProviderRouting creates a default way to resolve providers currently based
// on name constraint
func NewDefaultProviderRouting(providers []string, defaultProvider string) (ProviderLookup, error) {
	providerMap := map[string]*url.URL{}

	for _, v := range providers {
		pURL, err := url.Parse(v)
		if err != nil {
			return nil, fmt.Errorf("error parsing URL using value %s. %v", v, err)
		}
		providerMap[getHostNameWithoutPorts(pURL)] = pURL
	}

	d, err := url.Parse(defaultProvider)
	if err != nil {
		return nil, fmt.Errorf("error parsing default provider URL using value %s. %v", defaultProvider, err)
	}

	return &defaultProviderRouting{
		cache:           make(map[string]*types.FunctionDeployment),
		providers:       providerMap,
		defaultProvider: d,
	}, nil
}

func (d *defaultProviderRouting) ReloadCache() error {
	log.Info("reloading cache started...")

	var urls []string
	for _, v := range d.providers {
		urls = append(urls, v.String())
	}

	result, err := ReadServices(urls)
	if err != nil {
		return fmt.Errorf("could not reload cache. %v", err)
	}

	for k, v := range result.Providers {
		for _, f := range v {
			cf := requestToCreate(f)
			pURL, _ := url.Parse(k)
			ensureAnnotation(cf, getHostNameWithoutPorts(pURL))
			d.AddFunction(cf)
		}

		log.Infof("   added %d functions for provider %s", len(v), k)
	}

	log.Info("reloading cache completed successfully")
	return nil
}

func (d *defaultProviderRouting) Resolve(functionName string) (providerURI *url.URL, err error) {
	f, ok := d.GetFunction(functionName)

	if !ok {
		log.Warnf("can not find function %s in cache map, will attempt cache reload", functionName)
		if err := d.ReloadCache(); err != nil {
			return nil, fmt.Errorf("can not find function %s in cache map. Attempted to reload cache failed. %v", functionName, err)
		}

		f, ok = d.GetFunction(functionName)
		if !ok {
			return nil, fmt.Errorf("can not find function %s in cache map", functionName)
		}
	}

	log.Infof("Fn: %s, annotations: %v", functionName, f.Annotations)

	c, ok := (*f.Annotations)[federationProviderNameConstraint]
	if !ok {
		log.Infof("%s constraint not found using default provider %s", federationProviderNameConstraint, d.defaultProvider.String())
		return d.defaultProvider, nil
	}

	pURL := d.matchBasedOnName(c)
	if pURL == nil {
		log.Infof("%s constraint value found but does not exist in provider list, using default provider %s", c, d.defaultProvider.String())

		return d.defaultProvider, nil
	}

	return pURL, nil
}

func ensureAnnotation(f *types.FunctionDeployment, defaultValue string) {
	found := false
	if f.Annotations != nil {
		_, found = (*f.Annotations)[federationProviderNameConstraint]
	} else {
		f.Annotations = &map[string]string{}
	}

	if !found {
		(*f.Annotations)[federationProviderNameConstraint] = defaultValue
	}
}

func (d *defaultProviderRouting) matchBasedOnName(v string) *url.URL {
	for _, u := range d.providers {

		if strings.EqualFold(getHostNameWithoutPorts(u), v) {
			return u
		}
	}

	return nil
}

func getHostNameWithoutPorts(v *url.URL) string {
	// return strings.Split(v.Host, ":")[0]
	return v.String()
}

func (d *defaultProviderRouting) AddFunction(f *types.FunctionDeployment) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.cache[f.Service] = f
}

func (d *defaultProviderRouting) GetFunction(name string) (*types.FunctionDeployment, bool) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	v, ok := d.cache[name]

	return v, ok
}

func (d *defaultProviderRouting) GetFunctions() []*types.FunctionDeployment {
	d.lock.RLock()
	defer d.lock.RUnlock()
	var result []*types.FunctionDeployment
	for _, v := range d.cache {
		result = append(result, v)
	}

	return result
}
