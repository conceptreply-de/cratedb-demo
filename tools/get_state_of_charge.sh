#!/bin/bash

# usage example:
# tools/get_state_of_charge.sh 1

curl localhost:8080/vehicles/$1/state_of_charge_chart | jq 