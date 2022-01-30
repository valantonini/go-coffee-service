package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

const productServiceUrlEnvKey = "PRODUCT_SERVICE_URL"

var urlStem = os.Getenv(productServiceUrlEnvKey)

type RequestContext struct {
	t          *testing.T
	url        string
	httpMethod string
	body       interface{}
}

var client = http.Client{
	Timeout: time.Second * 5,
}

func DoRequest(requestCtx RequestContext) []byte {
	requestCtx.t.Helper()

	if urlStem == "" {
		requestCtx.t.Errorf("env %v not supplied", productServiceUrlEnvKey)
	}
	url := fmt.Sprintf("%v%v", urlStem, requestCtx.url)

	requestBody, err := json.Marshal(requestCtx.body)
	if err != nil {
		requestCtx.t.Error(err.Error())
	}
	fmt.Printf("requesting %v with\n %v\n", url, string(requestBody))

	req, err := http.NewRequest(requestCtx.httpMethod, url, bytes.NewBuffer(requestBody))
	assert.NoError(requestCtx.t, err)

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	res, err := client.Do(req)
	assert.NoError(requestCtx.t, err)

	if res.Body != nil {
		defer res.Body.Close()
	}

	responseJson, err := ioutil.ReadAll(res.Body)
	assert.NoError(requestCtx.t, err)

	return responseJson
}
