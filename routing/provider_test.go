package routing

import (
	"net/url"
	"testing"

	acc "github.com/openfaas-incubator/faas-federation/testing"
	types "github.com/openfaas/faas-provider/types"
)

func Test_defaultProviderRouting_Resolve(t *testing.T) {
	type fields struct {
		cache           map[string]*types.FunctionDeployment
		providers       map[string]*url.URL
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
				cache: map[string]*types.FunctionDeployment{
					"echo": {Service: "echo", Annotations: &map[string]string{federationProviderNameConstraint: "faas-provider-a:8080"}},
					"cat":  {Service: "cat", Annotations: &map[string]string{federationProviderNameConstraint: "faas-provider-b:8080"}},
				},
				providers: map[string]*url.URL{
					"faas-provider-a": parseURL("http://faas-provider-a:8080"),
					"faas-provider-b": parseURL("http://faas-provider-b:8080"),
				},
				defaultProvider: "http://faas-provider-a:8080",
			}, args: args{functionName: "echo"}, wantProviderHostName: "faas-provider-a:8080", wantErr: false,
		},
		{
			name: "provider b is resolved",
			fields: fields{
				cache: map[string]*types.FunctionDeployment{
					"echo": {Service: "echo", Annotations: &map[string]string{federationProviderNameConstraint: "http://faas-provider-a:8080"}},
					"cat":  {Service: "cat", Annotations: &map[string]string{federationProviderNameConstraint: "http://faas-provider-b:8080"}},
				},
				providers: map[string]*url.URL{
					"http://faas-provider-a:8080": parseURL("http://faas-provider-a:8080"),
					"http://faas-provider-b:8080": parseURL("http://faas-provider-b:8080"),
				},
				defaultProvider: "http://faas-provider-a:8080",
			}, args: args{functionName: "cat"}, wantProviderHostName: "faas-provider-b:8080", wantErr: false,
		},
		{
			name: "default provider is resolved, when constraint is missing",
			fields: fields{
				cache: map[string]*types.FunctionDeployment{
					"echo": {Service: "echo", Annotations: &map[string]string{federationProviderNameConstraint: "faas-provider-a:8080"}},
					"cat":  {Service: "cat", Annotations: &map[string]string{}},
				},
				providers: map[string]*url.URL{
					"faas-provider-a": parseURL("http://faas-provider-a:8080"),
					"faas-provider-b": parseURL("http://faas-provider-b:8080"),
				},
				defaultProvider: "http://faas-provider-a:8080",
			}, args: args{functionName: "cat"}, wantProviderHostName: "faas-provider-a:8080", wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultProviderRouting{
				cache:           tt.fields.cache,
				providers:       tt.fields.providers,
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

// Test_reloadCache requires `make up` and `cd examples && faas-cli up`
func Test_reloadCache(t *testing.T) {
	acc.PreCheckAcc(t)

	d := &defaultProviderRouting{
		providers: map[string]*url.URL{
			"faas-provider-a": parseURL("http://faas-provider-a:8082"),
			"faas-provider-b": parseURL("http://faas-provider-b:8083"),
		},
		defaultProvider: parseURL("http://faas-provider-a:8083"),
		cache:           map[string]*types.FunctionDeployment{},
	}

	err := d.ReloadCache()
	if err != nil {
		t.Fatal(err)
	}

	if len(d.cache) == 0 {
		t.Error("no items found in cache, check you have deployed examples to localhost:8080")
	}

	echoAConstraint := (*(d.cache["echo-a"].Annotations))[federationProviderNameConstraint]
	if echoAConstraint != "faas-provider-a" {
		t.Errorf("want: faas-provider-a got: %s", echoAConstraint)
	}
}

func parseURL(v string) *url.URL {
	u, _ := url.Parse(v)

	return u
}

func Test_ensureAnnotation(t *testing.T) {
	type args struct {
		f            *types.FunctionDeployment
		defaultValue string
		expected     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "annotation is present",
			args: args{
				defaultValue: "http://provider-c:8080",
				f: &types.FunctionDeployment{
					Annotations: &map[string]string{federationProviderNameConstraint: "http://provider-b:8080"},
				},
				expected: "http://provider-b:8080",
			},
		},
		{
			name: "annotation is missing, have other annotations. expect default value",
			args: args{
				defaultValue: "http://provider-c:8080",
				f: &types.FunctionDeployment{
					Annotations: &map[string]string{"bill": "ben"},
				},
				expected: "http://provider-c:8080",
			},
		},
		{
			name: "annotation is missing, have no annotations. expect default value",
			args: args{
				defaultValue: "http://provider-c:8080",
				f:            &types.FunctionDeployment{},
				expected:     "http://provider-c:8080",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ensureAnnotation(tt.args.f, tt.args.defaultValue)
			if tt.args.f.Annotations == nil {
				t.Error("Annotations should not be nil")
			}

			actual, okay := (*tt.args.f.Annotations)[federationProviderNameConstraint]
			if !okay {
				t.Errorf("%s annotation missing", federationProviderNameConstraint)
			}

			if actual != tt.args.expected {
				t.Errorf("want %s, got %s", tt.args.expected, actual)
			}
		})
	}
}
