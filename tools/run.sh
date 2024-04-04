#!/bin/bash

export PGPASSWORD="$(kubectl get secret user-system-my-cluster -o jsonpath={.data.password} | base64 -d)"
go run cmd/main.go