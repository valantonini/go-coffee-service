package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/matryer/is"
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

func DoRequest(requestCtx RequestContext) (statusCode int, responseJson []byte) {
	requestCtx.t.Helper()

	Is := is.New(requestCtx.t)

	if urlStem == "" {
		requestCtx.t.Errorf("env %v not supplied", productServiceUrlEnvKey)
	}
	url := fmt.Sprintf("%v%v", urlStem, requestCtx.url)

	requestBody, err := json.Marshal(requestCtx.body)
	Is.NoErr(err)

	fmt.Printf("requesting %v with\n %v\n", url, string(requestBody))

	req, err := http.NewRequest(requestCtx.httpMethod, url, bytes.NewBuffer(requestBody))
	Is.NoErr(err)

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	res, err := client.Do(req)
	Is.NoErr(err)

	Is.Equal(res.Header.Get("content-type"), "application/json")

	if res.Body != nil {
		defer res.Body.Close()
	}

	responseJson, err = ioutil.ReadAll(res.Body)
	Is.NoErr(err)

	return res.StatusCode, responseJson
}
