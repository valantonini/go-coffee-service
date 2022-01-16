package integration_tests

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
	"valantonini/go-coffee-service/coffee-service/data/entities"
)

func TestGetAllCoffees(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") == "" {
		t.Skip()
	}

	url := "http://api:8080/coffees"

	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	bd := entities.Coffees{}
	err = json.Unmarshal(body, &bd)

	assert.NoError(t, err)
	assert.NotZero(t, len(bd))
}
