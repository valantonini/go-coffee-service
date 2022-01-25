package integration_tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

var urlStem = os.Getenv("COFFEE_SERVICE_URL")

type RequestContext struct {
	t          *testing.T
	url        string
	httpMethod string
	body       io.Reader
}

var client = http.Client{
	Timeout: time.Second * 2,
}

func DoRequest(requestCtx RequestContext) []byte {
	requestCtx.t.Helper()

	url := fmt.Sprintf("%v%v", urlStem, requestCtx.url)

	fmt.Printf("requesting %v\n", url)

	req, err := http.NewRequest(requestCtx.httpMethod, url, requestCtx.body)
	assert.NoError(requestCtx.t, err)

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	res, err := client.Do(req)
	assert.NoError(requestCtx.t, err)

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(requestCtx.t, err)

	return body
}
