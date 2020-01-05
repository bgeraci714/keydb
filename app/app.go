package main

import (
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/bgeraci714/keydb/db"
	"github.com/bgeraci714/keydb/handler"
)

type App struct {
	EntryHandler *EntryHandler
	Database     db.Db
}

func (h *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	if head == "entry" {
		h.EntryHandler.ServeHTTP(res, req)
		return
	}
	log.Printf("head: %v, req: %v", head, req.URL.Path)
	http.Error(res, "Not Found", http.StatusNotFound)
}

type EntryHandler struct {
	db *db.Db
}

func (h *EntryHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		handler.GetHandler(h.db)(res, req)
	case "POST":
		handler.PutHandler(h.db)(res, req)
	}

}

// func TestGet(t *testing.T) {
// 	c := rbserver.RBClient{Port: "8080"}
// 	c.Init()

// 	http.DefaultClient()
// }

func main() {

	// store := rbserver.NewTree()
	// flushCount := 3
	d := db.NewDb(3)
	app := &App{
		EntryHandler: &EntryHandler{&d},
		Database:     d,
	}

	err := http.ListenAndServe(":8080", app)
	if err != nil {
		log.Fatal(err)
	}
}

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
// Used from: https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
