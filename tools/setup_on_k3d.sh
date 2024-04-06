#!/bin/bash

# create cluster with 4 k8s nodes and CrateDB operator
k3d cluster create crate-demo -p 5432:31428@server:0 -p 4200:30262@server:0 --servers 4
helm repo add crate-operator https://crate.github.io/crate-operator
kubectl create namespace crate-operator
helm install crate-operator crate-operator/crate-operator \
    --namespace crate-operator \
    --set env.CRATEDB_OPERATOR_DEBUG_VOLUME_STORAGE_CLASS=local-path

# create CrateDB cluster with 4 nodes on that cluster
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

sleep 5 # wait for the k8s to create operator pod
operator_pod_name=$(kubectl get pod -l 'app.kubernetes.io/instance=crate-operator' -n crate-operator -o template='{{ (index .items 0).metadata.name }}')
if [ -z "${operator_pod_name}" ]; then
  echo "Operator pod not found, try to run the script again later."
  exit 1
fi

kubectl wait --for=condition=Ready --timeout=600s -n crate-operator pod "${operator_pod_name}"

sleep 5 # wait for the operator to create pods, so we can wait on condition then
echo "Waiting for CrateDB to be ready... This may take a while (>5 minutes) for the first time on a fresh cluster to download the docker image."
kubectl wait --for=condition=Ready --timeout=600s pod crate-data-my-cluster-my-cluster-0
kubectl wait --for=condition=Ready --timeout=600s pod crate-data-my-cluster-my-cluster-1
kubectl wait --for=condition=Ready --timeout=600s pod crate-data-my-cluster-my-cluster-2
kubectl wait --for=condition=Ready --timeout=600s pod crate-data-my-cluster-my-cluster-3

# patch service to pin specific node ports
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
