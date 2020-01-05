# keydb
Key-value persistent data store written in Go 


## Interfaces 
```golang
type db interface {
    put(key, val) error
    get(key) (val, error)
}
```

## Atomic Data Types 
```golang 
type key []byte
type val []byte
```

## Data Collections
### In Memory
- Red black tree for storing most recent put requests 
- Sparse index of segment files to speed up search 

### Persistent 
- segment files that contain key value pairs


## TODO 
1. ~~Build in memory key-value server using the red black tree implementation for store's structure.~~  
1. Refine Interface to improve testability (build server as object with methods such as put/get to improve testability and consistency)
1. Refactor tree to be use bytes for both sides of storage
1. Add testing to account for this 


## Improvements 
2. Develop routine that saves the current data store's state into a segment file which can then be incorporated into database queries. 
3. Add logging in advance of writing to the store in order to provide crash/fault tolerance. 
4. Develop merge and compression routine for the segment files. 
5. Allow for rollbacks of the database using the log