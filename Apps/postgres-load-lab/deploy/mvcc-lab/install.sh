#!/usr/bin/env bash
set -euo pipefail

if [[ ${EUID} -ne 0 ]]; then
  echo "Run this installer as root." >&2
  exit 1
fi

source_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
export DEBIAN_FRONTEND=noninteractive

apt-get update
apt-get install -y postgresql postgresql-contrib prometheus prometheus-node-exporter sysstat

systemctl enable --now postgresql prometheus-node-exporter prometheus

pg_version=$(runuser -u postgres -- psql -Atqc "SHOW server_version" | cut -d. -f1)
conf_file="/etc/postgresql/${pg_version}/main/postgresql.conf"
sed -i "s/^#\?track_io_timing =.*/track_io_timing = on/" "$conf_file"
sed -i "s/^#\?log_autovacuum_min_duration =.*/log_autovacuum_min_duration = 0/" "$conf_file"
systemctl restart postgresql

if ! runuser -u postgres -- psql -Atqc "SELECT 1 FROM pg_database WHERE datname='mvcc_lab'" | grep -qx 1; then
  runuser -u postgres -- createdb mvcc_lab
fi

install -d -o postgres -g postgres -m 0750 /opt/mvcc-lab
install -d -o postgres -g postgres -m 0755 /var/lib/mvcc-lab
install -m 0644 "$source_dir/schema.sql" "$source_dir/churn.pgbench" \
  "$source_dir/sample.sql" "$source_dir/metrics.sql" "$source_dir/watch.sql" /opt/mvcc-lab/
runuser -u postgres -- psql -d mvcc_lab -f /opt/mvcc-lab/schema.sql
# Run initial maintenance in a fresh session, after the loader session has
# flushed its statistics. This avoids a transient double-count in n_live_tup.
runuser -u postgres -- vacuumdb --analyze --table=mvcc_healthy --table=mvcc_unvacuumed mvcc_lab

install -d -o postgres -g postgres -m 0755 /var/lib/mvcc-lab/textfile
cat >/etc/default/prometheus-node-exporter <<'EOF'
ARGS="--web.listen-address=127.0.0.1:9100 --collector.textfile.directory=/var/lib/mvcc-lab/textfile"
EOF
cat >/etc/prometheus/prometheus.yml <<'EOF'
global:
  scrape_interval: 10s
scrape_configs:
  - job_name: prometheus
    static_configs:
      - targets: ['127.0.0.1:9090']
  - job_name: node
    static_configs:
      - targets: ['127.0.0.1:9100']
EOF
cat >/etc/default/prometheus <<'EOF'
ARGS="--web.listen-address=127.0.0.1:9090"
EOF
systemctl restart prometheus-node-exporter prometheus

cat >/opt/mvcc-lab/run-workload.sh <<'SCRIPT'
#!/usr/bin/env bash
set -euo pipefail
while true; do
  /usr/bin/pgbench -n -c "${MVCC_CLIENTS:-8}" -j "${MVCC_THREADS:-2}" \
    -R "${MVCC_RATE:-500}" -T 3600 -P 10 -f /opt/mvcc-lab/churn.pgbench mvcc_lab
done
SCRIPT
chmod 0755 /opt/mvcc-lab/run-workload.sh

cat >/opt/mvcc-lab/sample.sh <<'SCRIPT'
#!/usr/bin/env bash
set -euo pipefail
output=/var/lib/mvcc-lab/timeline.csv
if [[ ! -e $output ]]; then
  echo 'sampled_at,relation,n_live_tup,n_dead_tup,n_tup_upd,n_tup_hot_upd,vacuum_count,autovacuum_count,last_vacuum,last_autovacuum,heap_bytes,index_bytes,xact_commit,blks_read,blks_hit,temp_bytes,oldest_transaction_seconds' >"$output"
fi
while true; do
  psql -X -q -d mvcc_lab -f /opt/mvcc-lab/sample.sql >>"$output"
  metrics_tmp=/var/lib/mvcc-lab/textfile/mvcc.prom.tmp
  psql -X -q -d mvcc_lab -f /opt/mvcc-lab/metrics.sql >"$metrics_tmp"
  mv "$metrics_tmp" /var/lib/mvcc-lab/textfile/mvcc.prom
  sleep 10
done
SCRIPT
chmod 0755 /opt/mvcc-lab/sample.sh

cat >/opt/mvcc-lab/watch.sh <<'SCRIPT'
#!/usr/bin/env bash
set -euo pipefail
watch -n 5 "sudo -u postgres psql -X -d mvcc_lab -f /opt/mvcc-lab/watch.sql"
SCRIPT
chmod 0755 /opt/mvcc-lab/watch.sh

cat >/etc/systemd/system/mvcc-lab-workload.service <<'UNIT'
[Unit]
Description=MVCC lab sustained pgbench workload
After=postgresql.service
Requires=postgresql.service

[Service]
Type=simple
User=postgres
Environment=MVCC_CLIENTS=8
Environment=MVCC_THREADS=2
Environment=MVCC_RATE=500
ExecStart=/opt/mvcc-lab/run-workload.sh
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
UNIT

cat >/etc/systemd/system/mvcc-lab-sampler.service <<'UNIT'
[Unit]
Description=MVCC lab PostgreSQL time-series sampler
After=postgresql.service
Requires=postgresql.service

[Service]
Type=simple
User=postgres
ExecStart=/opt/mvcc-lab/sample.sh
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
UNIT

systemctl daemon-reload
systemctl enable --now mvcc-lab-workload mvcc-lab-sampler

echo "MVCC lab installed with PostgreSQL ${pg_version}."
echo "Workload and sampler are running; data is in /var/lib/mvcc-lab/timeline.csv."
echo "Prometheus is listening on VM loopback port 9090."
