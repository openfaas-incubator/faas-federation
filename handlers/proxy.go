package handlers

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/openfaas-incubator/faas-federation/routing"
	log "github.com/sirupsen/logrus"
)

const urlScheme = "http"

// MakeProxyHandler creates a handler to invoke functions downstream
func MakeProxyHandler(proxy http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Info("proxy request")

		pathVars := mux.Vars(r)
		if pathVars == nil {
			r = mux.SetURLVars(r, map[string]string{})
			pathVars = mux.Vars(r)
		}

		functionName := strings.Split(r.URL.Path, "/")[2]
		pathVars["name"] = functionName
		pathVars["params"] = r.URL.Path
		log.Infof("proxy request to: %s %s", functionName, r.URL.String())
		proxy.ServeHTTP(w, r)
	}
}

// FunctionLookup is a openfaas-provider proxy.BaseURLResolver that allows the
// caller to verify that a function is resolvable.
type FunctionLookup struct {
	// scheme is the http scheme (http/https) used to proxy the request
	scheme string
	// dnsrrLookup method used to resolve the function IP address, defaults to the internal lookupIP
	// method, which is an implementation of net.LookupIP
	dnsrrLookup    func(context.Context, string) ([]net.IP, error)
	providerLookup routing.ProviderLookup
}

// NewFunctionLookup creates a new FunctionLookup resolver
func NewFunctionLookup(providerLookup routing.ProviderLookup) *FunctionLookup {
	return &FunctionLookup{
		scheme:         urlScheme,
		dnsrrLookup:    lookupIP,
		providerLookup: providerLookup,
	}
}

// Resolve implements the openfaas-provider proxy.BaseURLResolver interface.
func (l *FunctionLookup) Resolve(name string) (u url.URL, err error) {
	log.Infof("resolving function %s", name)
	providerURL, err := l.providerLookup.Resolve(name)
	if err != nil {
		return url.URL{}, err
	}

	log.Infof("using provider %s to for function %s", providerURL.String(), name)

	return *providerURL, nil
}

// resolve the function by checking the available docker DNSRR resolution
func (l *FunctionLookup) byDNSRoundRobin(ctx context.Context, name string) (string, error) {
	entries, lookupErr := l.dnsrrLookup(ctx, fmt.Sprintf("tasks.%s", name))

	if lookupErr != nil {
		return "", lookupErr
	}

	if len(entries) > 0 {
		index := randomInt(0, len(entries))
		return entries[index].String(), nil
	}

	return "", fmt.Errorf("could not resolve '%s' using dnsrr", name)
}

func randomInt(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

// lookupIP implements the net.LookupIP method with context support. It returns a slice of that\
// host's IPv4 and IPv6 addresses.
func lookupIP(ctx context.Context, host string) ([]net.IP, error) {
	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}
	ips := make([]net.IP, len(addrs))
	for i, ia := range addrs {
		ips[i] = ia.IP
	}
	return ips, nil
}
