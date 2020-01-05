package shared

import "bytes"

// Key is the atomic key for an entry
type Key []byte

// Value is the value of an entry
type Value []byte

// Entry is a pairing of keys and values
type Entry struct {
	Key   Key
	Value Value
}

// Compare returns
// +1 if k > b
// -1 if k < b
//  0 if k == b
func (k Key) Compare(b Key) int {
	return bytes.Compare(k, b)
}

// Compare returns
// +1 if k > b
// -1 if k < b
//  0 if k == b
func (v Value) Compare(b Value) int {
	return bytes.Compare(v, b)
}
