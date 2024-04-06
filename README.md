# Crate DB Demo

This is a demonstration of a simple application that uses CrateDB in a near-real-world scenario.

It includes necessary bootstrap for running CrateDB on local k8s to perform various experiments with locally started application.

## Prerequisites

You don't need to install all of those tools, f.e k6 only is needed if you use "load*" scripts, but it is recommended to have them.

- [k3d](https://k3d.io) for a cluster to deploy CrateDB to
- [kubectl](https://kubernetes.io/docs/reference/kubectl/) to do the deployment
- [helm](https://helm.sh/) to do the deployment
- [Golang](https://go.dev/) to compile/run demo application
- [jq](https://jqlang.github.io/jq/) to nicely print some json output
- [bash](https://en.wikipedia.org/wiki/Bash_(Unix_shell)) as a shell for scripts - optionally can be substituded for compatible one (or run commands manually)
- [psql](https://www.postgresql.org/docs/current/app-psql.html) to run queries and execute scripts
- [k6](https://grafana.com/docs/k6/latest/) to load our demo application with lots of queries
- base64 utility - used to decode from base64, look at `tools/get_password.sh`
- free ports: 8080 (application HTTP), 4200 (CrateDB UI), 5432 (CrateDB "psql.port")
- less than 95% taken space on your disk (or cluser won't be healthy)

## Usage

Note: I have only tested this on Linux so far, so some scripts might need adjustments for your system:
 - f.e about base64 utility - you might need to change `tools/get_password.sh` script on some OS, which require `base64 -D`

To create new k3d cluster, login kubectl into it, install CrateDB with 4 nodes and configure it to serve its UI on 4200 and Postgres wire protocol on 5432, run this script:
```bash
tools/setup_on_k3d.sh
```

### Accessing database

Extract database password:

```bash
source tools/get_password.sh
echo system user password is ${PGPASSWORD}
```

That would print out and save password in shell environment, so you could use it later.
Password belongs to `system` user, which is created by default.

Open [localhost:4200](http://localhost:4200) in your browser - it should show CrateDB UI, where you can enter username "system" and password printed before.

Now in case if your dev machine happened to have big and mostly used disk, you might need to adjust some watermarks after cluster is fully up:

```bash
# this is a workaround for the case when dev machine has big disk with mostly used space,
# don't do this for production, this only makes sense for testing purposes
# flood_stage > high > low > space_used_on_host
PGUSER=system psql -h localhost -p 5432 -c "
  SET GLOBAL PERSISTENT cluster.routing.allocation.disk.watermark.flood_stage = '99%';
  SET GLOBAL PERSISTENT cluster.routing.allocation.disk.watermark.high = '97%';
  SET GLOBAL PERSISTENT cluster.routing.allocation.disk.watermark.low = '95%';
"
```

Initialize some sample data:

```bash
# base database schema, indices configs, etc
PGUSER=system psql -h localhost -f ./schema.sql

# sample vehicles and measurements
PGUSER=system psql -h localhost -f tools/data/vehicles.sql

# approximate countries geolocation data to check geo queries
PGUSER=system psql -h localhost -f tools/data/austria.sql
PGUSER=system psql -h localhost -f tools/data/germany.sql
PGUSER=system psql -h localhost -f tools/data/switzerland.sql
```

You can execute any query that database supports that same way:

```bash
# Just show which vehicles we have inserted so far
psql -h localhost -p 5432 -U system -c "
  select * from vehicles;
"

# execute and explain the detailed plan of complex query
# and pretty print json output
psql -h localhost -p 5432 -t -P 'pager=off' -U system -c "
explain analyze WITH avg_measurements AS (
    SELECT vehicle_id,
    DATE_BIN('1 minutes'::INTERVAL, measurement_time, 0) AS period,
    AVG(measurement_value) AS avg_state_of_charge
    FROM vehicle_measurements
    WHERE measurement_type = 'STATE_OF_CHARGE'
    AND measurement_partition > NOW() - INTERVAL '20 minutes'
    GROUP BY 1, 2 
    ORDER BY 1, 2
)
SELECT period period_from,
        period + INTERVAL '1 minutes' period_to,
    am.vehicle_id vehicle_id,
    avg_state_of_charge  
FROM avg_measurements am, vehicles v
WHERE am.vehicle_id = v.id
AND v.id = '1'
LIMIT 100;
" | jq --sort-keys
```

You can also use `tools/db_connect.sh` to connect by psql and run SQL interactively.

### Demo application

Start demo application:

```bash
tools/run.sh
```

Keep that terminal running and open another one.


Send one measurement STATE_OF_CHARGE=5 for vehicle with id 1:

```bash
tools/add_measurement.sh 1 5 STATE_OF_CHARGE
```

Show state of charge for last 20 minutes aggregated in minute blocks for vehicle 1:
```bash
tools/get_state_of_charge.sh 1
```

Running search:
```bash
tools/search_vehicles.sh shuttle
```

This script would set current vehicle location for vehicle "1", to longitude 8.5 and latitude 47:

```bash
tools/set_current_location.sh 1 8.5 47
```

#### Simulating Load

Run this to create artificial load that would attempt to add vehicle measurements at a rate 10/s
(check the script `tools/load_add_measurement_k6.js` for possible configs):

```bash
tools/load_add_measurement.sh
```

Or instead load setting vehicle location:

```bash
tools/load_set_location.sh
```

## What's next?

[Check out what CrateDB is capable of](https://cratedb.com/docs), read code in `cmd`, `internal` and `tools` to get idea how this demo works, add stuff for your other scenarios, try different things, etc.

Some ideas that could be interesting to try here:
- put together setup with Kafka Connect writing directly to database by JDBC driver
- use `skaffold` to quickly bootstrap that Golang app within the cluster, so it would use k8s internal balancing to connect to database and for app to be scaled as well
- use [Vector analysis](https://cratedb.com/solutions/vector-database) to understand how close a particular vehicle measurements profile is to another one

## Cleanup

If deployed by `setup_on_k3d.sh`, simply call this to remove the cluster:

```bash
k3d cluster delete crate-demo
```
