package handlers

import (
	"net/http"

	"github.com/openfaas-incubator/faas-federation/routing"

	log "github.com/sirupsen/logrus"
)

// MakeUpdateHandler update specified function
func MakeUpdateHandler(proxy http.HandlerFunc, providerLookup routing.ProviderLookup) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("update request")

		function, err := addToFunctionToCache(r, providerLookup)
		if err != nil {
			log.Errorln("error during unmarshal of create function request. ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		proxyDeployment(proxy, function, w, r)

		log.Info("update request successful")
	}
}
