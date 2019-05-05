package routing

import (
	"testing"

	"github.com/openfaas/faas/gateway/requests"
)

func Test_defaultProviderRouting_Resolve(t *testing.T) {
	type fields struct {
		cache     map[string]*requests.CreateFunctionRequest
		providers map[string]bool
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
				providers: map[string]bool{
					"provider_a" : true,
					"provider_b" : true,
				},
				defaultProvider: "provider_a",
			}, args: args{functionName: "echo"}, wantProviderHostName: "provider_a", wantErr: false,
		},
		{
			name: "provider b is resolved",
			fields: fields{
				cache: map[string]*requests.CreateFunctionRequest{
					"echo": {Service: "echo", Annotations: &map[string]string{"federation.provider_name": "provider_a"}},
					"cat":  {Service: "cat", Annotations: &map[string]string{"federation.provider_name": "provider_b"}},
				},
				providers: map[string]bool{
					"provider_a" : true,
					"provider_b" : true,
				},
				defaultProvider: "provider_a",
			}, args: args{functionName: "cat"}, wantProviderHostName: "provider_b", wantErr: false,
		},
		{
			name: "default provider is resolved, when constraint is missing",
			fields: fields{
				cache: map[string]*requests.CreateFunctionRequest{
					"echo": {Service: "echo", Annotations: &map[string]string{"federation.provider_name": "provider_a"}},
					"cat":  {Service: "cat", Annotations: &map[string]string{}},
				},
				providers: map[string]bool{
					"provider_a" : true,
					"provider_b" : true,
				},
				defaultProvider: "provider_a",
			}, args: args{functionName: "cat"}, wantProviderHostName: "provider_a", wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultProviderRouting{
				cache:     tt.fields.cache,
				providers: tt.fields.providers,
				defaultProvider: tt.fields.defaultProvider,
			}
			gotProviderHostName, err := d.Resolve(tt.args.functionName)
			if (err != nil) != tt.wantErr {
				t.Errorf("defaultProviderRouting.Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotProviderHostName != tt.wantProviderHostName {
				t.Errorf("defaultProviderRouting.Resolve() = got %v, want %v", gotProviderHostName, tt.wantProviderHostName)
			}
		})
	}
}
