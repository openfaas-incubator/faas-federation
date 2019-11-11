// Copyright (c) OpenFaaS Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"net/http"

	"github.com/openfaas-incubator/faas-federation/routing"
	log "github.com/sirupsen/logrus"
)

// MakeLogHandler to read logs from an endpoint
func MakeLogHandler(proxy http.HandlerFunc, providerLookup routing.ProviderLookup) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("log handler")

		proxy.ServeHTTP(w, r)
	}
}
