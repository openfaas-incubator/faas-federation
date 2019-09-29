// Copyright (c) OpenFaaS Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/openfaas-incubator/faas-federation/routing"
	types "github.com/openfaas/faas-provider/types"
	log "github.com/sirupsen/logrus"
)

// MakeReplicaUpdater updates desired count of replicas
func MakeReplicaUpdater() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("update replicas, nothing to do here")

	}
}

// MakeReplicaReader reads the amount of replicas for a deployment
func MakeReplicaReader(provider routing.ProviderLookup) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("read replicas")

		vars := mux.Vars(r)
		functionName := vars["name"]

		res, ok := provider.GetFunction(functionName)

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		found := &types.FunctionStatus{}
		found.Name = functionName
		found.AvailableReplicas = 1
		found.Annotations = res.Annotations
		found.Labels = res.Labels

		functionBytes, _ := json.Marshal(found)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(functionBytes)
	}
}
