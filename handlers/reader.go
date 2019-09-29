// Copyright (c) OpenFaaS Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/openfaas-incubator/faas-federation/routing"

	types "github.com/openfaas/faas-provider/types"
	log "github.com/sirupsen/logrus"
)

// MakeFunctionReader handler for reading functions deployed in the cluster as deployments.
func MakeFunctionReader(providers []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Info("read request")

		functions, err := routing.ReadServices(providers)
		if err != nil {
			log.Printf("Error getting service list: %s\n", err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		var result []*types.FunctionStatus
		for _, v := range functions.Providers {
			result = append(result, v...)
		}

		functionBytes, _ := json.Marshal(result)
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
		w.Write(functionBytes)
	}
}
