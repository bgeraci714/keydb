package rbserver

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/bgeraci714/keydb/rbtree"
)

const path string = "./"
const sname string = "segment"

// DBExtension defines the extension used for the db files
const DBExtension string = ".db"

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
func PutHandler(tree *rbtree.RBTree, flushSize int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var e Entry

		// try to decode the body, could probably be refactored
		err := json.NewDecoder(req.Body).Decode(&e)
		if err != nil || e.Key == "" {
			log.Printf("Error decoding body: %v", err)
			http.Error(w, "Error occurred when trying to read body.", http.StatusBadRequest)
			return
		}

		tree.Insert(e.Key, e.Value)

		// conditionally write the tree to a segment file
		if tree.Size() >= flushSize {
			if err := saveSegment(tree); err != nil {
				panic(err)
			}

			// initialize new tree
			*tree = NewTree()
		}

		// return results and print state of tree (purely for testing purposes initially)
		fmt.Printf("tree state:\n" + tree.ToString())
	}
}

// GetHandler takes a pointer to a red black tree and returns
// a callback retrieves a key if found in the tree
func GetHandler(tree *rbtree.RBTree) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var k Key

		// try to decode the body
		err := json.NewDecoder(req.Body).Decode(&k)
		if err != nil || k.Key == "" {
			log.Printf("Error decoding body: %v", err)
			http.Error(w, "can't decode body", http.StatusBadRequest)
			return
		}

		// try to find the value in the tree, indicate if found in response
		if val, found := tree.Get(k.Key); found {
			fmt.Fprintf(w, "{found: true, value: %s}", val)
			return
		} else if val, found := searchSegments(k.Key); found {
			fmt.Fprintf(w, "{found: true, value: %s}", val)
			return
		}

		fmt.Fprintf(w, "{found: false}")
	}
}

// searches through segment files
// TODO: maybe change boolean to string name of file it was found in (could also be number)
func searchSegments(key string) (interface{}, bool) {
	// open files using readdir
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	// iterate through the files (sorted alphanumerically) from oldest to newest
	for i := len(files) - 1; i >= 0; i-- {
		f := files[i]

		// skip if not the globally defined db extension
		if filepath.Ext(f.Name()) != DBExtension {
			continue
		}

		segment, err := os.Open(path + f.Name())
		if err != nil {
			log.Fatal(err)
			continue
		}

		// convert segment file to the io.Reader interface then a gob decoder
		r := bufio.NewReader(segment)
		dec := gob.NewDecoder(r)

		// load it into memory as rb tree
		var entries []Entry
		if err := dec.Decode(&entries); err != nil {
			panic(err)
		}

		// search the entries slice
		if val, found := search(&entries, key, 0, len(entries)-1); found {
			return val, found
		}
	}

	return nil, false
}

// linear search for the entries
// entries are coming in unsorted
func search(entries *[]Entry, key string, low, high int) (interface{}, bool) {
	// mid := (high + low) / 2
	// fmt.Println(*entries)
	// // fmt.Println((*entries)[low].Key, "<", (*entries)[mid].Key, "<", (*entries)[(high+len(*entries))/2].Key)
	// if high < low {
	// 	return nil, false
	// }
	// if (*entries)[mid].Key == key {
	// 	return (*entries)[mid].Value, true
	// } else if (*entries)[mid].Key > key {
	// 	return search(entries, key, mid+1, high)
	// } else { // if entries[mid].Key < key { // less than
	// 	return search(entries, key, low, mid-1)
	// }

	// linear search
	for i := range *entries {
		if (*entries)[i].Key == key {
			return (*entries)[i].Value, true
		}
	}
	return nil, false
}

func getNextSegmentID(last string) int {
	// segment-0.db
	re := regexp.MustCompile("[0-9]+.")
	j := re.FindString(last)
	if j == "" {
		panic("ahhhhh")
	}
	i, err := strconv.Atoi(j[:len(j)-1])
	if err != nil {
		panic("nooooo")
	}
	return i + 1
}

func getFileName(path string) string {
	// open files using readdir
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		// create file name
		return sname + "-0" + DBExtension
	}

	var i int
	for i = len(files) - 1; i >= 0; i-- {
		if filepath.Ext(files[i].Name()) != DBExtension {
			continue
		}
		// inputFmt:=input[:len(input)-1]
		id := getNextSegmentID(files[i].Name())
		return sname + "-" + strconv.Itoa(id) + DBExtension
	}

	return sname + "-0" + DBExtension
}

func saveSegment(t *rbtree.RBTree) error {
	m := t.ToMap()
	var entries []Entry
	for k, v := range m {
		entries = append(entries, Entry{k, v.(string)})
	}

	fname := getFileName(path)
	f, err := os.Create(fname)
	if err != nil {
		return err
	}

	// close f on exit and check for its returned error
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	// make a write buffer
	w := bufio.NewWriter(f)
	enc := gob.NewEncoder(w)

	if err := enc.Encode(entries); err != nil {
		log.Fatal("encoding error:", err)
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}

// NewTree is a short function to make a new tree. This is a might be a bad code smell
func NewTree() rbtree.RBTree {
	return rbtree.RBTree{
		Root: nil,
		Compare: func(a, b interface{}) int {
			return strings.Compare(a.(string), b.(string))
		},
	}
}

// RBClient is the main client struct to combine the store, max items in tree, and port
// type RBClient struct {
// 	Port     string
// 	storeMax int
// 	store    *rbtree.RBTree
// }

// // Init initializes client and starts the server
// func (c *RBClient) Init(flushSize int) {
// 	c.initStore(flushSize)
// 	c.initServer(flushSize)
// }

// func (c *RBClient) initServer(flushSize int) {
// 	http.HandleFunc("/hello", hello)
// 	http.HandleFunc("/put", put(c.store, flushSize))
// 	http.HandleFunc("/get", get(c.store))

// 	http.ListenAndServe(":"+c.Port, nil)
// }

// // // InitStore initializes the store for the client
// // func (c *RBClient) initStore(sMax int) {
// // 	c.storeMax = sMax
// // 	*c.store = newTree()
// // }
