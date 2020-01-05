package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bgeraci714/keydb/db"
	"github.com/bgeraci714/keydb/shared"
)

func TestPutHandler(t *testing.T) {
	k := []byte("apple")
	v := []byte("sauce")

	body := getSamplePutRequestBody(k, v)
	req, err := http.NewRequest("POST", "localhost:8080/entry", body)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// initialization variables
	flushCount := 2
	d := db.NewDb(flushCount)

	rec := httptest.NewRecorder()
	// make the call

	PutHandler(&d)(rec, req)

	// verify the status is OK
	res := rec.Result()
	defer res.Body.Close() // good form to close out the body
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.StatusCode)
	}

	actual, found := d.Get(k)
	if !found {
		t.Fatalf("expected to find key: %s in the store but none were found", k)
	} else if actual.Compare(v) != 0 {
		t.Fatalf("found key: %s in the store expecting %s but instead found %v", k, v, actual)
	}

	// verify the return value
	// b, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	t.Fatalf("could not read response: %v", err)
	// }

}

func getSamplePutRequestBody(key shared.Key, value shared.Value) *bytes.Buffer {
	body, err := json.Marshal(map[string]string{
		"key":   string(key),
		"value": string(value),
	})
	if err != nil {
		log.Fatalf("could not create request body: %v", err)
	}
	return bytes.NewBuffer(body)
}
