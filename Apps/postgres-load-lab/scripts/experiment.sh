#!/usr/bin/env bash
set -euo pipefail

project_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
loadgen=${LOADGEN:-$project_dir/bin/loadlab-loadgen}
target=${TARGET:-http://127.0.0.1:8080}
database_url=${DATABASE_URL:-postgres://loadlab:loadlab@127.0.0.1:5432/loadlab?sslmode=disable}
stamp=$(date -u +%Y%m%dT%H%M%SZ)
result_dir="$project_dir/results/$stamp"
mkdir -p "$result_dir"

curl --fail --silent "$target/metrics" >"$result_dir/metrics-before.txt"
"$loadgen" --target "$target" --duration "${BASELINE_DURATION:-2m}" --workers 2 --rate 10 --output "$result_dir/baseline.json"

"$loadgen" --target "$target" --duration "${LOAD_DURATION:-5m}" --workers "${LOAD_WORKERS:-100}" --rate "${LOAD_RATE:-0}" --seed-companies 100 --output "$result_dir/saturation.json" &
load_pid=$!
sleep 10
"$loadgen" --target "$target" --duration "${PROBE_DURATION:-4m30s}" --workers 2 --rate 10 --output "$result_dir/under-load-probe.json"
wait "$load_pid"

curl --fail --silent "$target/metrics" >"$result_dir/metrics-after.txt"
psql "$database_url" -X -c "SELECT calls,round(total_exec_time::numeric,2) total_ms,round(mean_exec_time::numeric,2) mean_ms,rows,left(query,120) query FROM pg_stat_statements ORDER BY total_exec_time DESC LIMIT 20" >"$result_dir/pg-stat-statements.txt"
echo "Results written to $result_dir"
