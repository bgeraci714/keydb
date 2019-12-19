package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bgeraci714/keydb/rbserver"
)

// func TestGet(t *testing.T) {
// 	c := rbserver.RBClient{Port: "8080"}
// 	c.Init()

// 	http.DefaultClient()
// }

func main() {

	store := rbserver.NewTree()
	flushCount := 3
	port := "8080"

	http.HandleFunc("/put", rbserver.PutHandler(&store, flushCount))
	http.HandleFunc("/get", rbserver.GetHandler(&store))

	fmt.Println("Server will be listening on port: " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
