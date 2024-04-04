#!/bin/bash


k3d cluster create crate-demo -p 5432:31428@server:0 -p 4200:30262@server:0 --servers 4
helm repo add crate-operator https://crate.github.io/crate-operator
kubectl create namespace crate-operator
helm install crate-operator crate-operator/crate-operator \
    --namespace crate-operator \
    --set env.CRATEDB_OPERATOR_DEBUG_VOLUME_STORAGE_CLASS=local-path


cat << EOF | kubectl apply -f -
apiVersion: cloud.crate.io/v1
kind: CrateDB
metadata:
  name: my-cluster
  namespace: default
spec:
  cluster:
    imageRegistry: crate
    name: crate-dev
    version: 5.0.1
  nodes:
    data:
    - name: my-cluster
      replicas: 4
      resources:
        requests:
          cpu: 1
          memory: 1Gi
        limits:
          cpu: 1
          memory: 1Gi
        disk:
          count: 1
          size: 128GiB
          storageClass: local-path
        heapRatio: 0.25
EOF

sleep 5 # wait for the operator to create the first pod, so we can wait on condition then

# downloading crate image would take some time...
echo "Waiting for CrateDB to be ready... This may take a while for the first time on a fresh cluster."
kubectl wait --for=condition=Ready --timeout=600s pod crate-data-my-cluster-my-cluster-0

kubectl patch service crate-my-cluster --type=merge \
-p '
{ "spec": 
  { "ports" : 
    [ 
      { "name": "http", "nodePort": 30262, "port": 4200, "protocol": "TCP", "targetPort": 4200 }, 
      { "name": "psql", "nodePort": 31428, "port": 5432, "protocol": "TCP", "targetPort": 5432 } 
    ] 
  } 
}'

# this is a workaround for the case when dev machine has big disk with mostly used space,
# don't do this for production, this only makes sense for development machines
psql -h localhost -p 5432 -U system -c "
  SET GLOBAL PERSISTENT cluster.routing.allocation.disk.watermark.flood_stage = '99%';
  SET GLOBAL PERSISTENT cluster.routing.allocation.disk.watermark.high = '97%';
  SET GLOBAL PERSISTENT cluster.routing.allocation.disk.watermark.low = '95%';
"

