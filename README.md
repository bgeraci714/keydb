# keydb
Key-value persistent data store written in Go 

## TODO 
1. ~~Build in memory key-value server using the red black tree implementation for store's structure.~~  
2. Develop routine that saves the current data store's state into a segment file which can then be incorporated into database queries. 
3. Add logging in advance of writing to the store in order to provide crash/fault tolerance. 
4. Develop merge and compression routine for the segment files. 
5. Allow for rollbacks of the database using the log