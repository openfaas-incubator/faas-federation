package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/openfaas/faas-federation/routing"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"
)

const urlScheme = "http"

// MakeProxy creates a proxy for HTTP web requests which can be routed to a function.
func MakeProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		v, okay := functions[name]
		if !okay {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{ \"status\" : \"Not found\"}"))
		}

		v.InvocationCount = v.InvocationCount + 1
		responseBody := "{ \"status\" : \"Okay\"}"
		w.Write([]byte(responseBody))
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

	providerHostName, err := l.providerLookup.Resolve(name)
	if err != nil {
		return url.URL{}, err
	}

	return l.ResolveContext(context.Background(), providerHostName)
}

// ResolveContext provides an implementation of openfaas-provider proxy.BaseURLResolver with
// context support. See `Resolve`
func (l *FunctionLookup) ResolveContext(ctx context.Context, name string) (u url.URL, err error) {

	u.Host, err = l.byDNSRoundRobin(ctx, name)

	if err != nil {
		return u, err
	}

	u.Scheme = l.scheme
	return u, nil
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