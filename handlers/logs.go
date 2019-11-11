// Copyright (c) OpenFaaS Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/openfaas-incubator/faas-federation/routing"
	log "github.com/sirupsen/logrus"
)

// MakeLogHandler to read logs from an endpoint
func MakeLogHandler(proxy http.HandlerFunc, providerLookup routing.ProviderLookup) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("log handler")

		query := r.URL.Query()
		name := query.Get("name")
		uri, err := providerLookup.Resolve(name)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		uriPath := uri.String() + "/system/logs?name=" + name

		log.Infof("URI forwarding logs to: %s", uriPath)

		req, _ := http.NewRequest(http.MethodGet, uriPath, nil)
		res, resErr := http.DefaultClient.Do(req)

		if resErr != nil {
			http.Error(w, resErr.Error(), http.StatusInternalServerError)
			return
		}

		if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusContinue {
			if res.Body != nil {
				defer res.Body.Close()
			}
			http.Error(w, fmt.Sprintf("Incorrect HTTP status code: %d", res.StatusCode), http.StatusInternalServerError)
			return
		}

		io.Copy(w, ioutil.NopCloser(res.Body))
	}
}
