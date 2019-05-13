package testing

import (
	"os"
	"testing"
)

// PreCheckAcc only allows acceptance tests to run in explicitly enabled via OF_ACC
// environment variable
func PreCheckAcc(t *testing.T) {
	if os.Getenv("OF_ACC") == "" {
		t.Skip("To enable acceptance tests please set environment variable OF_ACC=1")
	}
}
