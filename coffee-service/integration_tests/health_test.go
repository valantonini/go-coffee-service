package integration_tests

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestHealth(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") == "" {
		t.Skip()
	}
	req := RequestContext{
		t:          t,
		url:        "/health",
		httpMethod: http.MethodGet,
		body:       nil,
	}

	body := DoRequest(req)
	var bd string
	err := json.Unmarshal(body, &bd)

	assert.NoError(t, err)
	assert.Equal(t, "\"ok\"", string(body))
}
