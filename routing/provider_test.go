package routing

import (
	"net/url"
	"testing"

	"github.com/openfaas/faas/gateway/requests"
)

func Test_defaultProviderRouting_Resolve(t *testing.T) {
	type fields struct {
		cache     map[string]*requests.CreateFunctionRequest
		providers map[string]*url.URL
		defaultProvider string
	}
	type args struct {
		functionName string
	}
	tests := []struct {
		name                 string
		fields               fields
		args                 args
		wantProviderHostName string
		wantErr              bool
	}{
		{
			name: "provider a is resolved",
			fields: fields{
				cache: map[string]*requests.CreateFunctionRequest{
					"echo": {Service: "echo", Annotations: &map[string]string{"federation.provider_name": "provider_a"}},
					"cat":  {Service: "cat", Annotations: &map[string]string{"federation.provider_name": "provider_b"}},
				},
				providers: map[string]*url.URL{
					"provider_a" : parseURL("http://provider_a:8080"),
					"provider_b" : parseURL("http://provider_b:8080"),
				},
				defaultProvider: "http://provider_a:8080",
			}, args: args{functionName: "echo"}, wantProviderHostName: "provider_a:8080", wantErr: false,
		},
		{
			name: "provider b is resolved",
			fields: fields{
				cache: map[string]*requests.CreateFunctionRequest{
					"echo": {Service: "echo", Annotations: &map[string]string{"federation.provider_name": "provider_a"}},
					"cat":  {Service: "cat", Annotations: &map[string]string{"federation.provider_name": "provider_b"}},
				},
				providers: map[string]*url.URL{
					"provider_a" : parseURL("http://provider_a:8080"),
					"provider_b" : parseURL("http://provider_b:8080"),
				},
				defaultProvider: "http://provider_a:8080",
			}, args: args{functionName: "cat"}, wantProviderHostName: "provider_b:8080", wantErr: false,
		},
		{
			name: "default provider is resolved, when constraint is missing",
			fields: fields{
				cache: map[string]*requests.CreateFunctionRequest{
					"echo": {Service: "echo", Annotations: &map[string]string{"federation.provider_name": "provider_a"}},
					"cat":  {Service: "cat", Annotations: &map[string]string{}},
				},
				providers: map[string]*url.URL{
					"provider_a" : parseURL("http://provider_a:8080"),
					"provider_b" : parseURL("http://provider_b:8080"),
				},
				defaultProvider: "http://provider_a:8080",
			}, args: args{functionName: "cat"}, wantProviderHostName: "provider_a:8080", wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultProviderRouting{
				cache:     tt.fields.cache,
				providers: tt.fields.providers,
				defaultProvider: parseURL(tt.fields.defaultProvider),
			}
			gotProviderHostName, err := d.Resolve(tt.args.functionName)
			if (err != nil) != tt.wantErr {
				t.Errorf("defaultProviderRouting.Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotProviderHostName == nil {
				t.Errorf("defaultProviderRouting.Resolve() = nil")
			}

			if gotProviderHostName.Host != tt.wantProviderHostName {
				t.Errorf("defaultProviderRouting.Resolve() = got %v, want %v", gotProviderHostName.Host, tt.wantProviderHostName)
			}
		})
	}
}


func parseURL(v string) *url.URL {
	u, _ := url.Parse(v)

	return u
}