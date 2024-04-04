#!/bin/bash

# Adding measurement to a vehicle
# usage example:
# tools/add_measurement.sh 1 5 STATE_OF_CHARGE

curl localhost:8080/vehicles/$1/measurements \
    -H 'Content-Type: application/json'\
    -XPOST -d '{"value": '"$2"', "type": "'"$3"'"}'