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

// Test_Delete requires `make up` and `cd examples && faas-cli up`
func Test_Delete(t *testing.T) {
	acc.PreCheckAcc(t)
	req, err := http.NewRequest("DELETE", "/system/functions", bytes.NewBuffer([]byte(echoDelete)))
	if err != nil {
		t.Fatal(err)
	}

	mux.NewRouter()
	rr := httptest.NewRecorder()

	providerLookup, err := routing.NewDefaultProviderRouting([]string{"http://faas-provider-a:8082", "http://faas-provider-b:8083"}, "http://faas-provider-a:8082")
	if err != nil {
		t.Fatal(err)
	}

	err = providerLookup.ReloadCache()
	if err != nil {
		t.Fatalf("error reloading provider cache. %v", err)
	}

	proxyFunc := proxy.NewHandlerFunc(time.Minute*1, NewFunctionLookup(providerLookup))

	MakeDeleteHandler(proxyFunc).ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

const echoDelete = `{"functionName":"echo-b"}`
