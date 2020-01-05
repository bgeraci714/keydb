package db

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/bgeraci714/keydb/rbtree"
	"github.com/bgeraci714/keydb/shared"
)

// DbPath is the path for where the segment files will be stored
const DbPath string = "./"

// DBExtension defines the extension used for the db files
const DBExtension string = ".db"
const sname string = "segment"

// Db is the data access object that maintains the store and handles put/get calls
type Db struct {
	store      *rbtree.RBTree
	FlushCount int
}

// NewDb initialiates a new database struct
func NewDb(flushCount int) Db {
	db := Db{FlushCount: flushCount}
	db.store = &rbtree.RBTree{Root: nil}
	return db
}

// Put adds in a new key value pair into the database
func (db *Db) Put(key shared.Key, value shared.Value) error {
	db.store.Insert(key, value)

	// conditionally write the tree to a segment file
	if db.store.Size() >= db.FlushCount {
		if err := saveSegment(db.store); err != nil {
			return fmt.Errorf("Unable to write current store to segment file")
		}

		// initialize new tree
		*db.store = rbtree.RBTree{Root: nil}
	}
	return nil
}

// Get queries the db for a given key
func (db *Db) Get(key shared.Key) (shared.Value, bool) {
	// try to find the value in the tree, indicate if found in response
	if val, found := db.store.Get(key); found {
		return val, true
	} else if val, found := searchSegments(key); found {
		return val, true
	}
	return nil, false
}

// searches through segment files
// TODO: maybe change boolean to string name of file it was found in (could also be number)
func searchSegments(key shared.Key) (shared.Value, bool) {
	// open files using readdir
	files, err := ioutil.ReadDir(DbPath)
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

		segment, err := os.Open(DbPath + f.Name())
		if err != nil {
			log.Fatal(err)
			continue
		}

		// convert segment file to the io.Reader interface then a gob decoder
		r := bufio.NewReader(segment)
		dec := gob.NewDecoder(r)

		// load it into memory as rb tree
		var entries []shared.Entry
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
func search(entries *[]shared.Entry, key shared.Key, low, high int) (shared.Value, bool) {
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
		if (*entries)[i].Key.Compare(key) == 0 {
			return (*entries)[i].Value, true
		}
	}
	return []byte{}, false
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

func generateNextFileName(path string) string {
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
	var entries []shared.Entry
	for k, v := range m {
		entries = append(entries, shared.Entry{Key: []byte(k), Value: v})
	}

	fname := generateNextFileName(DbPath)
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
