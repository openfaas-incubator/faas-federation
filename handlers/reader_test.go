// Copyright (c) Edward Wilde 2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"testing"

	acc "github.com/ewilde/faas-federation/testing"
)

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
			"http://localhost:8082",
			"http://localhost:8083",
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readServices(tt.args.providers)
			if err != nil {
				t.Error(err)
			}

			if len(got) == 0 {
				t.Errorf("Want atleast one function")
			}
		})
	}
}
