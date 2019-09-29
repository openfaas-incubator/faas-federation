// Copyright (c) OpenFaaS Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	types "github.com/openfaas/faas-provider/types"

	"github.com/gorilla/mux"
	"github.com/openfaas-incubator/faas-federation/routing"
	log "github.com/sirupsen/logrus"
)

// MakeDeployHandler creates a handler to create new functions in the cluster
func MakeDeployHandler(proxy http.HandlerFunc, providerLookup routing.ProviderLookup) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Info("deployment request")

		function, err := addToFunctionToCache(r, providerLookup)
		if err != nil {
			log.Errorln("error during unmarshal of create function request. ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		proxyDeployment(proxy, function, w, r)

		log.Infof("deployment request for function %s path %s", function.Service, r.URL.String())
	}
}

func proxyDeployment(proxy http.HandlerFunc, function *types.FunctionDeployment, w http.ResponseWriter, r *http.Request) {
	pathVars := mux.Vars(r)
	if pathVars == nil {
		r = mux.SetURLVars(r, map[string]string{})
		pathVars = mux.Vars(r)
	}

	pathVars["name"] = function.Service
	pathVars["params"] = r.URL.Path
	proxy.ServeHTTP(w, r)
}

func addToFunctionToCache(r *http.Request, providerLookup routing.ProviderLookup) (*types.FunctionDeployment, error) {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	request := &types.FunctionDeployment{}
	if err := json.Unmarshal(body, &request); err != nil {
		return nil, fmt.Errorf("error during unmarshal of create function request. %v", err)
	}

	providerLookup.AddFunction(request)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return request, nil
}
