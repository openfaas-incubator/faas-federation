package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/openfaas-incubator/faas-federation/routing"
	acc "github.com/openfaas-incubator/faas-federation/testing"
	"github.com/openfaas/faas-provider/proxy"
)

func Test_Update(t *testing.T) {
	acc.PreCheckAcc(t)
	req, err := http.NewRequest("PUT", "/system/functions", bytes.NewBuffer([]byte(echoUpdate)))
	if err != nil {
		t.Fatal(err)
	}

	mux.NewRouter()
	rr := httptest.NewRecorder()

	providerLookup, err := routing.NewDefaultProviderRouting([]string{"http://faas-provider-a:8082", "http://faas-provider-b:8083"}, "http://faas-provider-a:8082")
	if err != nil {
		t.Fatal(err)
	}

	proxyFunc := proxy.NewHandlerFunc(time.Minute*1, NewFunctionLookup(providerLookup))

	MakeUpdateHandler(proxyFunc, providerLookup).ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

const echoUpdate = `{"service":"echo-a","image":"openfaas/echo:latest","network":"","envProcess":"./handler","envVars":{},"constraints":null,"secrets":[],"labels":{},"annotations":{},"limits":null,"requests":null,"readOnlyRootFilesystem":false}`
