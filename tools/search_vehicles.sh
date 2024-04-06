#!/bin/bash

curl 'localhost:8080/vehicles/search?q='"${1}" | jq
