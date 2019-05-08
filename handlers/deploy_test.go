package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ewilde/faas-federation/routing"
	acc "github.com/ewilde/faas-federation/testing"
	"github.com/gorilla/mux"
	"github.com/openfaas/faas-provider/proxy"
)

func Test_Deploy(t *testing.T) {
	acc.PreCheckAcc(t)
	req, err := http.NewRequest("POST", "/system/functions", bytes.NewBuffer([]byte(echoDeploy)))
	if err != nil {
		t.Fatal(err)
	}

	mux.NewRouter()
	rr := httptest.NewRecorder()

	providerLookup, err := routing.NewDefaultProviderRouting([]string{"http://provider_a:8082", "http://provider_b:8083"}, "http://provider_a:8082")
	if err != nil {
		t.Fatal(err)
	}

	proxyFunc := proxy.NewHandlerFunc(time.Minute*1, NewFunctionLookup(providerLookup))

	MakeDeployHandler(proxyFunc, providerLookup).ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

const echoDeploy = `{"service":"echo","image":"ewilde/echo:latest","network":"","envProcess":"./handler","envVars":{},"constraints":null,"secrets":[],"labels":{},"annotations":{},"limits":null,"requests":null,"readOnlyRootFilesystem":false}`
