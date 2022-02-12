//go:build integration
// +build integration

package product_service_tests

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

	statusCode, body := DoRequest(req)
	var bd string
	err := json.Unmarshal(body, &bd)

	Is.NoErr(err)
	Is.Equal(statusCode, http.StatusOK)
	Is.Equal(string(body), "\"ok\"\n")
}
