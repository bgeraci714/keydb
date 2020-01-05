package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bgeraci714/keydb/db"
)

// Entry represents the expected shape of the request body for put calls
type Entry struct {
	Key   string
	Value string
}

// Key represents the expected shape of the reuqest body for get calls
type Key struct {
	Key string
}

// little test function to verify server is working correctly
func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, world.\n")
}

// PutHandler takes a pointer to a red black tree and returns
// a callback that adds the key value pair to the tree
func PutHandler(db *db.Db) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var e Entry

		// try to decode the body, could probably be refactored
		err := json.NewDecoder(req.Body).Decode(&e)
		if err != nil || len(e.Key) == 0 {
			log.Printf("Error decoding body: %v", err)
			http.Error(w, "Error occurred when trying to read body.", http.StatusBadRequest)
			return
		}
		err = db.Put([]byte(e.Key), []byte(e.Value))
		if err != nil {
			log.Printf("Error writing to db: %v", err)
			http.Error(w, "Error occurred when trying to write to db.", http.StatusBadRequest)
			return
		}
		log.Printf("Wrote: [%v:%v]", e.Key, e.Value)
		fmt.Fprintf(w, "{successful: true, key: %v}", e.Key)
	}
}

// GetHandler takes a pointer to a red black tree and returns
// a callback retrieves a key if found in the tree
func GetHandler(db *db.Db) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var k Key

		// try to decode the body
		err := json.NewDecoder(req.Body).Decode(&k)
		if err != nil || k.Key == "" {
			log.Printf("Error decoding body: %v", err)
			http.Error(w, "can't decode body", http.StatusBadRequest)
			return
		}

		if val, found := db.Get([]byte(k.Key)); found {
			fmt.Fprintf(w, "{found: true, value: %s}", val)
			return
		}

		fmt.Fprintf(w, "{found: false}")
	}
}
