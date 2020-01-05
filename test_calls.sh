#!/bin/bash

port=8080

# test_get <key>
test_get() {
    curl localhost:${port}/entry -X GET -d "{\"key\": \"$1\"}";
}

# test_get <key> <value>
test_put() {
    curl localhost:${port}/entry -X POST -d "{\"key\": \"$1\", \"value\": \"$2\"}";
}

# test_hello
test_hello() {
    curl localhost:${port}/hello
}