package testing

import (
	"os"
	"testing"
)

func PreCheckAcc(t *testing.T) {
	if os.Getenv("OF_ACC") == "" {
		t.Skip("To enable acceptance tests please set environment variable OF_ACC=1")
	}
}
