#!/bin/bash

# usage example:
# source tools/get_password.sh

export PGPASSWORD="$(kubectl get secret user-system-my-cluster -o jsonpath={.data.password} | base64 -d)"