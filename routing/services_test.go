// Copyright (c) Edward Wilde 2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package routing

import (
	"testing"

	acc "github.com/ewilde/faas-federation/testing"
)

// Test_readServices requires `make up` and `cd examples && faas-cli up`
func Test_readServices(t *testing.T) {
	acc.PreCheckAcc(t)
	type args struct {
		providers []string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "list", args: args{providers: []string{
			"http://faas-provider-a:8082",
			"http://faas-provider-b:8083",
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadServices(tt.args.providers)
			if err != nil {
				t.Error(err)
			}

			if len(got.Providers["faas-provider-a"]) == 0 {
				t.Errorf("Want atleast one function")
			}
		})
	}
}
