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

var urlStem = os.Getenv("PRODUCT_SERVICE_URL")

// var urlStem = "http://localhost:8080"

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

	url := fmt.Sprintf("%v%v", urlStem, requestCtx.url)

	jsonMapAsStringFormat, err := json.Marshal(requestCtx.body)
	if err != nil {
		requestCtx.t.Error(err.Error())
	}
	fmt.Printf("requesting %v with\n %v\n", url, string(jsonMapAsStringFormat))

	req, err := http.NewRequest(requestCtx.httpMethod, url, bytes.NewBuffer(jsonMapAsStringFormat))
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
