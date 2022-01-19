package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
	"valantonini/go-coffee-service/coffee-service/data/entities"
)

var urlStem = os.Getenv("URL_STEM")

func TestCoffees(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") == "" {
		t.Skip()
	}

	t.Run("should get all coffees", func(t *testing.T) {
		req := requestContext{
			t:          t,
			url:        "/list",
			httpMethod: http.MethodGet,
			body:       nil,
		}

		body := doRequest(req)
		bd := entities.Coffees{}
		err := json.Unmarshal(body, &bd)

		assert.NoError(t, err)
		assert.NotZero(t, len(bd))
	})

	t.Run("should add coffee", func(t *testing.T) {
		coffeeName := uuid.New()
		var jsonData = []byte(fmt.Sprintf(`{
			"name": "%v"
		}`, coffeeName))

		req := requestContext{
			t:          t,
			url:        "/add",
			httpMethod: http.MethodPost,
			body:       bytes.NewBuffer(jsonData),
		}

		body := doRequest(req)

		coffee := entities.Coffee{}
		err := json.Unmarshal(body, &coffee)

		assert.NoError(t, err)
		assert.NotZero(t, coffee.ID)
		assert.Equal(t, coffeeName.String(), coffee.Name)
	})
}

type requestContext struct {
	t          *testing.T
	url        string
	httpMethod string
	body       io.Reader
}

var client = http.Client{
	Timeout: time.Second * 2,
}

func doRequest(ctx requestContext) []byte {
	ctx.t.Helper()

	url := fmt.Sprintf("%v%v", urlStem, ctx.url)
	req, err := http.NewRequest(ctx.httpMethod, url, ctx.body)
	assert.NoError(ctx.t, err)

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	res, err := client.Do(req)
	assert.NoError(ctx.t, err)

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(ctx.t, err)

	return body
}
