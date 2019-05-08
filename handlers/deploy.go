// Copyright (c) Edward Wilde 2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ewilde/faas-federation/routing"
	"github.com/gorilla/mux"
	"github.com/openfaas/faas/gateway/requests"
	log "github.com/sirupsen/logrus"
)

var functions = map[string]*requests.Function{}

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

func proxyDeployment(proxy http.HandlerFunc, function *requests.CreateFunctionRequest, w http.ResponseWriter, r *http.Request) {
	pathVars := mux.Vars(r)
	if pathVars == nil {
		r = mux.SetURLVars(r, map[string]string{})
		pathVars = mux.Vars(r)
	}

	pathVars["name"] = function.Service
	pathVars["params"] = r.URL.Path
	proxy.ServeHTTP(w, r)
}

func addToFunctionToCache(r *http.Request, providerLookup routing.ProviderLookup) (*requests.CreateFunctionRequest, error) {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	request := &requests.CreateFunctionRequest{}
	if err := json.Unmarshal(body, &request); err != nil {
		return nil, fmt.Errorf("error during unmarshal of create function request. %v", err)
	}

	providerLookup.AddFunction(request)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return request, nil
}
