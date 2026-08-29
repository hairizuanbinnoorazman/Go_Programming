# PostgreSQL Load Lab

PostgreSQL Load Lab is a deliberately small Go CRUD service for learning where a single-server application reaches its limits. It manages companies, users, and inventory with UUIDv4 identifiers, generates a repeatable mixed workload, and exposes enough telemetry to distinguish application, connection-pool, query, CPU, memory, and storage bottlenecks.

This is an experiment harness, not a production starter. Local Docker Compose is for functional verification only: container resource sharing and desktop virtualization make its performance numbers unsuitable as a server baseline.

See the standalone [before-and-after optimization report](before-after-report.html) for a visual explanation of the measured query-plan improvement.

## What is included

- JSON CRUD API built with `net/http` and `pgxpool`
- PostgreSQL 17 local environment and automatic idempotent schema migration
- Prometheus HTTP latency, database-operation latency, process, and pool metrics
- Mixed Go load generator with concurrency/rate controls and p50/p95/p99 output
- systemd units for the app, load generator, Prometheus, and PostgreSQL exporter
- an installation script for a fresh Debian/Ubuntu-style lab server
- baseline/saturation experiment and PostgreSQL diagnostic scripts

## Run locally

Requirements: Docker with Compose, `curl`, and optionally Go 1.24+.

```bash
cd Apps/postgres-load-lab
docker compose up --build -d --wait
./scripts/smoke.sh
docker compose logs app
```

The API is at `http://127.0.0.1:8080`; metrics are at `/metrics`. Stop it with `docker compose down`. Add `--volumes` only when you intentionally want to erase the lab database.

Example CRUD calls:

```bash
company=$(curl -sS -H 'Content-Type: application/json' \
  -d '{"name":"Acme"}' http://127.0.0.1:8080/v1/companies)
company_id=$(printf '%s' "$company" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')

curl -sS -H 'Content-Type: application/json' \
  -d "{\"company_id\":\"$company_id\",\"name\":\"Ada\",\"email\":\"ada@example.test\"}" \
  http://127.0.0.1:8080/v1/users

curl -sS -H 'Content-Type: application/json' \
  -d "{\"company_id\":\"$company_id\",\"sku\":\"WIDGET-1\",\"name\":\"Widget\",\"quantity\":20,\"price_cents\":1999}" \
  http://127.0.0.1:8080/v1/inventory

curl -sS "http://127.0.0.1:8080/v1/inventory?company_id=$company_id"
curl -sS -X PUT -H 'Content-Type: application/json' -d '{"name":"Acme Ltd"}' \
  "http://127.0.0.1:8080/v1/companies/$company_id"
curl -i -X DELETE "http://127.0.0.1:8080/v1/companies/$company_id"
```

All list endpoints accept `limit` (1–500) and `offset`. Inventory also accepts `company_id`. Deleting a company sets its users' `company_id` to null and cascades deletion of its inventory.

## Load generator

Build it or run it directly:

```bash
GOWORK=off go build -o bin/loadlab-loadgen ./cmd/loadgen
./bin/loadlab-loadgen --target http://127.0.0.1:8080 \
  --duration 2m --workers 25 --rate 200 --output results/capped.json
```

`--rate 0` removes the rate cap and drives maximum throughput. The mixed workload creates, reads, updates, deletes, and lists all three resource types. It prints request count, failure count, requests/second, status distribution, and p50/p95/p99/max latency. Use a capped rate for a stable baseline and increase it in steps; jumping directly to unlimited load only shows the fully saturated state.

Use `--workload inventory-inserts` to isolate the heaviest insert path while still exercising the HTTP API, foreign key, uniqueness constraint, WAL, and inventory indexes. It creates seed companies before the timed run; truncate the lab tables between comparable rounds so dataset growth does not bias later results.

## Install on a small server

Use a disposable Debian/Ubuntu-style Linux host. The script installs PostgreSQL, Prometheus, node exporter, and PostgreSQL exporter; creates least-privilege service users; and manages each process with systemd. Review the script before running it.

Build for the server architecture, then copy this directory:

```bash
./scripts/build-linux.sh                 # set GOARCH=arm64 when needed
rsync -av --exclude results ./ user@server:/tmp/postgres-load-lab/
ssh user@server
cd /tmp/postgres-load-lab
sudo APP_DB_PASSWORD='choose_16_chars_min' \
  EXPORTER_DB_PASSWORD='choose_another_16' ./scripts/install-server.sh
```

The services bind only to loopback. Use SSH forwarding for the API and Prometheus:

```bash
ssh -L 8080:127.0.0.1:8080 -L 9090:127.0.0.1:9090 user@server
```

Useful checks on the server:

```bash
systemctl status loadlab postgresql loadlab-postgres-exporter loadlab-prometheus prometheus-node-exporter
journalctl -u loadlab -f
curl -sS http://127.0.0.1:8080/readyz
curl -sS http://127.0.0.1:9090/-/ready
```

The sample PostgreSQL settings in `deploy/postgresql.conf.example` are measurement aids and conservative starting points, not automatic tuning. Only `pg_stat_statements` and I/O timing are enabled by the installer.

## Measure baseline and loaded latency

Run the generator on a different machine when testing the server; otherwise the generator competes for the CPU being measured. Reset query statistics, then run a quiet baseline:

```bash
sudo -u postgres psql -d loadlab -c 'SELECT pg_stat_statements_reset()'
./bin/loadlab-loadgen --target http://SERVER:8080 --duration 5m \
  --workers 2 --rate 10 --output results/baseline.json
```

Next run one uncapped saturation workload and a small, fixed-rate probe at the same time. The probe represents user-facing latency while the heavy process creates contention:

```bash
./bin/loadlab-loadgen --target http://SERVER:8080 --duration 10m \
  --workers 100 --rate 0 --output results/saturation.json &
load_pid=$!
sleep 10
./bin/loadlab-loadgen --target http://SERVER:8080 --duration 9m30s \
  --workers 2 --rate 10 --output results/under-load-probe.json
wait "$load_pid"
```

On the server, `sudo systemctl start loadlab-loadgen` runs the settings in `/etc/loadlab/loadgen.env`. For a local all-in-one experiment, `scripts/experiment.sh` captures baseline/load JSON, application metrics, and the top `pg_stat_statements` queries under a timestamped `results/` directory.

## Identify the bottleneck

Compare the baseline and under-load probe p95/p99, then inspect the same time window in Prometheus. Useful queries are:

```promql
histogram_quantile(0.95, sum by (le, route) (rate(loadlab_http_request_duration_seconds_bucket[1m])))
histogram_quantile(0.95, sum by (le, operation) (rate(loadlab_db_query_duration_seconds_bucket[1m])))
loadlab_db_pool_acquired_connections / loadlab_db_pool_max_connections
rate(node_cpu_seconds_total{mode!="idle"}[1m])
node_memory_MemAvailable_bytes
rate(node_disk_io_time_seconds_total[1m])
rate(pg_stat_database_xact_commit{datname="loadlab"}[1m])
rate(pg_stat_database_blks_read{datname="loadlab"}[1m])
```

Interpret them together:

| Evidence | Likely limit | Test before changing anything |
|---|---|---|
| HTTP p95 rises, DB p95 stays flat | App CPU, queueing, network, or load-generator limit | Check app/server CPU and run generator remotely |
| Pool utilization stays near 1 and HTTP exceeds DB latency | Requests wait for connections | Check PostgreSQL CPU/I/O; cautiously test a larger pool |
| DB p95 rises and `pg_stat_statements` shows high mean time/blocks | Query plan or storage | Run `EXPLAIN (ANALYZE, BUFFERS)` for that exact query |
| CPU is saturated while disk wait is low | CPU-bound app/query work | Find the process and hottest statements before scaling |
| Disk I/O time and block reads rise, cache hit falls | Working set or storage limit | Check dataset size, cache, scans, and storage latency |
| Locks/deadlocks or wait events dominate | Contention | Inspect `pg_stat_activity` and shorten conflicting transactions |
| Errors rise but server metrics are quiet | Client/network/timeouts | Inspect status distribution and load-generator host saturation |

Capture a point-in-time database report with:

```bash
psql "$DATABASE_URL" -f scripts/diagnose.sql
```

## Demonstrate an improvement

Use one change per run, identical dataset/load/duration, and compare throughput, error rate, p95/p99, CPU/I/O, and top statements. For the inventory-by-company query, obtain a real company ID and inspect the plan:

```sql
EXPLAIN (ANALYZE, BUFFERS)
SELECT id, company_id, sku, name, quantity, price_cents, created_at, updated_at
FROM inventory
WHERE company_id = 'UUID_HERE'
ORDER BY updated_at DESC
LIMIT 100;
```

The schema includes `inventory_company_updated_idx (company_id, updated_at DESC)` for this access pattern and foreign-key indexes for joins/deletes. It also includes matching `created_at DESC` indexes for the unfiltered list endpoints. Without those list indexes, a local validation dataset of about 27,000 companies produced a sequential scan plus top-N sort: 406 shared buffers and 3.77 ms execution time. After adding `companies_created_at_idx`, the same query used an index scan, touched 51 buffers, and executed in 0.10 ms (the exact result is machine- and cache-dependent). This is the intended proof pattern: matching load and data, an actual plan change, fewer buffers, then improved application p95—not just the existence of an index.

On a sufficiently populated table, the desired evidence is an index scan with fewer buffers and lower execution time. If PostgreSQL still chooses a sequential scan, do not force the index blindly: the table may be small, statistics may be stale, or the query may return enough rows that a scan is cheaper. Run `ANALYZE inventory`, check selectivity, and retest.

Other evidence-led changes include lowering an oversized app pool when PostgreSQL context switching dominates, raising it only when pool waits dominate and PostgreSQL has headroom, batching write-heavy workflows, replacing deep `OFFSET` pagination with keyset pagination, and tuning autovacuum for a high-churn table. A successful result must improve the selected latency/throughput target without merely moving the bottleneck or increasing errors.

## Development checks

```bash
GOWORK=off gofmt -w cmd internal
GOWORK=off go test ./...
GOWORK=off go vet ./...
docker compose config --quiet
```
