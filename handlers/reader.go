// Copyright (c) Edward Wilde 2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ewilde/faas-federation/routing"

	"github.com/openfaas/faas/gateway/requests"
	log "github.com/sirupsen/logrus"
)

// MakeFunctionReader handler for reading functions deployed in the cluster as deployments.
func MakeFunctionReader(providers []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Info("read request")
		functions, err := readServices(providers)
		if err != nil {
			log.Printf("Error getting service list: %s\n", err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		functionBytes, _ := json.Marshal(functions)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(functionBytes)
	}
}

func readServices(providers []string) ([]*requests.Function, error) {
	var list []*requests.Function
	var urls []*url.URL

	for _, v := range providers {
		u, _ := url.Parse(v)
		u.Path = "/system/functions"
		urls = append(urls, u)
	}

	results := routing.Get(urls, len(providers))
	for _, v := range results {
		if v.Err != nil {
			log.Errorf("error fetching function list for %s. %v", providers[v.Index], v.Err)
			break
		}

		if v.Response.StatusCode > 399 {
			log.Errorf("unexpected error code %d while fetching function list for %s. %v", v.Response.StatusCode, providers[v.Index], v.Err)
			break
		}

		var function []*requests.Function
		functionBytes, err := ioutil.ReadAll(v.Response.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response for %s. %v", providers[v.Index], err)
		}

		_ = v.Response.Body.Close()
		err = json.Unmarshal(functionBytes, &function)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling response for %s. %v", providers[v.Index], err)
		}

		list = append(list, function...)
	}

	return list, nil
}

func createToRequest(request *requests.CreateFunctionRequest) *requests.Function {
	return &requests.Function{
		Name:              request.Service,
		Annotations:       request.Annotations,
		EnvProcess:        request.EnvProcess,
		Image:             request.Image,
		Labels:            request.Labels,
		AvailableReplicas: 1,
		Replicas:          1,
	}
}
