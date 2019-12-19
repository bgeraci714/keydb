package rbserver

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPutHandler(t *testing.T) {
	k, v := "apple", "sauce"

	body := getSamplePutRequestBody(k, v)
	req, err := http.NewRequest("POST", "localhost:8080/put", body)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// initialization variables
	store := NewTree()
	flushCount := 2
	rec := httptest.NewRecorder()

	// make the call
	PutHandler(&store, flushCount)(rec, req)

	// verify the status is OK
	res := rec.Result()
	defer res.Body.Close() // good form to close out the body
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.StatusCode)
	}

	if actual, found := store.Get(k); !found || actual != v {
		t.Fatalf("expected to find key: %v in the store with value %v", k, v)
	}

	// verify the return value
	// b, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	t.Fatalf("could not read response: %v", err)
	// }

}

func getSamplePutRequestBody(key, value string) *bytes.Buffer {
	body, err := json.Marshal(map[string]string{
		"key":   key,
		"value": value,
	})
	if err != nil {
		log.Fatalf("could not create request body: %v", err)
	}
	return bytes.NewBuffer(body)
}
