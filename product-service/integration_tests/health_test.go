//go:build integration
// +build integration

package integration_tests

import (
	"encoding/json"
	"github.com/matryer/is"
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {
	Is := is.New(t)

	req := RequestContext{
		t:          t,
		url:        "/health",
		httpMethod: http.MethodGet,
		body:       nil,
	}

	body := DoRequest(req)
	var bd string
	err := json.Unmarshal(body, &bd)

	Is.NoErr(err)
	Is.Equal(string(body), "\"ok\"")
}
