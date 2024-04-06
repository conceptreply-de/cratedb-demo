#!/bin/bash

export PGPASSWORD="$(kubectl get secret user-system-my-cluster -o jsonpath={.data.password} | base64 -d)"
PGUSER=system PGPASSWORD="${PGPASSWORD}" psql -h localhost
