#!/bin/bash

curl -XPUT "http://localhost:8080/vehicles/$1/location" -H 'Content-Type: application/json' -v -d "{ \"lon\": ${2}, \"lat\": ${3} }"
