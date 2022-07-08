//nolint:typecheck
package utils

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

// TestRUD (Read, Update, Delete) struct to factorize code inside a single feature test.
type TestRUD struct {
	CreateRequest  Request
	CreateResponse ICreateResponse
	SubTests       []SubTest
}

// Runs a series of tests for a Get type endpoint.
func (test TestRUD) Run(t *testing.T) {
	disableLogs()

	w, err := executeRequest(test.CreateRequest.Method, test.CreateRequest.Url, test.CreateRequest.Body)
	if err != nil {
		t.Error(err)
	}

	key, err := test.CreateResponse.KeyCreated(w.Body.Bytes())
	if err != nil {
		t.Error(err)
	}

	// Then run the tests.
	for _, subtest := range test.SubTests {
		t.Run(subtest.Name, func(t *testing.T) {
			url := subtest.Request.Url + key

			w, err := executeRequest(subtest.Request.Method, url, subtest.Request.Body)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, w.Code, subtest.ResponseCode)
			assert.MatchRegex(t, w.Body.String(), subtest.ResponseBodyRegex)
		})
	}
}