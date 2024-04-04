# Crate DB Demo

This is a demonstration of a simple application that uses CrateDB in a near-real-world scenario.

It includes necessary bootstrap for running CrateDB on local k8s to perform various experiments with locally started application.

## Prerequisites

- [k3d](https://k3d.io)
- [kubectl](https://kubernetes.io/docs/reference/kubectl/)
- [helm](https://helm.sh/)
- [Golang](https://go.dev/)
- free ports 8080, 4200, 5432

## Usage

Create new k3d cluster, login kubectl into it, install CrateDB with 4 nodes and configure it to serve its UI on 4200 and Postgres wire protocol on 5432:
`tools/setup_on_k3d.sh`

Extract database password and initialize some sample data:
```bash
source tools/get_password.sh
PGUSER=system PGPASSWORD="${PGPASSWORD}" psql -h localhost -f tools/data.sql
```

Run golang application:
`tools/run.sh`

Send one measurement STATE_OF_CHARGE=5 for vehicle with id 1:
`tools/add_measurement.sh 1 5 STATE_OF_CHARGE`

Show state of charge for last 20 minutes aggregated in minute blocks for vehicle 1:
`tools/get_state_of_charge.sh 1`

Running search:
```bash
curl localhost:8080/vehicles/search?q=Route
```

