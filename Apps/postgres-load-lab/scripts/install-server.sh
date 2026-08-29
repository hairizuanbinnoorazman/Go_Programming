#!/usr/bin/env bash
set -euo pipefail

if [[ ${EUID} -ne 0 ]]; then
  echo "Run this installer as root." >&2
  exit 1
fi
: "${APP_DB_PASSWORD:?Set APP_DB_PASSWORD to a strong alphanumeric password}"
: "${EXPORTER_DB_PASSWORD:?Set EXPORTER_DB_PASSWORD to a strong alphanumeric password}"
if [[ ! $APP_DB_PASSWORD =~ ^[A-Za-z0-9_-]{16,128}$ ]] || [[ ! $EXPORTER_DB_PASSWORD =~ ^[A-Za-z0-9_-]{16,128}$ ]]; then
  echo "Passwords must be 16-128 characters using letters, digits, underscore, or hyphen." >&2
  exit 1
fi

project_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
for binary in loadlab-server loadlab-loadgen; do
  test -x "$project_dir/bin/$binary" || { echo "Missing bin/$binary; run scripts/build-linux.sh first." >&2; exit 1; }
done

export DEBIAN_FRONTEND=noninteractive
apt-get update
apt-get install -y postgresql postgresql-contrib prometheus prometheus-node-exporter prometheus-postgres-exporter

id loadlab >/dev/null 2>&1 || useradd --system --home /var/lib/loadlab --create-home --shell /usr/sbin/nologin loadlab
install -d -o root -g root -m 0755 /opt/loadlab/bin /opt/loadlab/migrations /etc/loadlab
install -d -o loadlab -g loadlab -m 0750 /var/lib/loadlab
install -d -o prometheus -g prometheus -m 0755 /var/lib/prometheus/loadlab
install -m 0755 "$project_dir/bin/loadlab-server" "$project_dir/bin/loadlab-loadgen" /opt/loadlab/bin/
install -m 0644 "$project_dir"/migrations/*.sql /opt/loadlab/migrations/

if ! runuser -u postgres -- psql -tAc "SELECT 1 FROM pg_roles WHERE rolname='loadlab'" | grep -q 1; then
  runuser -u postgres -- psql -v ON_ERROR_STOP=1 -c "CREATE ROLE loadlab LOGIN PASSWORD '$APP_DB_PASSWORD'"
else
  runuser -u postgres -- psql -v ON_ERROR_STOP=1 -c "ALTER ROLE loadlab PASSWORD '$APP_DB_PASSWORD'"
fi
if ! runuser -u postgres -- psql -tAc "SELECT 1 FROM pg_database WHERE datname='loadlab'" | grep -q 1; then
  runuser -u postgres -- createdb --owner=loadlab loadlab
fi
if ! runuser -u postgres -- psql -tAc "SELECT 1 FROM pg_roles WHERE rolname='loadlab_exporter'" | grep -q 1; then
  runuser -u postgres -- psql -v ON_ERROR_STOP=1 -c "CREATE ROLE loadlab_exporter LOGIN PASSWORD '$EXPORTER_DB_PASSWORD'"
fi
runuser -u postgres -- psql -v ON_ERROR_STOP=1 -c "GRANT pg_monitor TO loadlab_exporter"
runuser -u postgres -- psql -v ON_ERROR_STOP=1 -c "ALTER SYSTEM SET shared_preload_libraries = 'pg_stat_statements'"
runuser -u postgres -- psql -v ON_ERROR_STOP=1 -c "ALTER SYSTEM SET track_io_timing = 'on'"
systemctl restart postgresql
runuser -u postgres -- psql -d loadlab -v ON_ERROR_STOP=1 -c "CREATE EXTENSION IF NOT EXISTS pg_stat_statements"

install -m 0644 "$project_dir/deploy/prometheus.yml" /etc/loadlab/prometheus.yml
install -m 0644 "$project_dir/deploy/loadgen.env" /etc/loadlab/loadgen.env
sed "s#REPLACE_ME#$APP_DB_PASSWORD#" "$project_dir/deploy/app.env.example" >/etc/loadlab/app.env
printf 'DATA_SOURCE_NAME=postgresql://loadlab_exporter:%s@127.0.0.1:5432/loadlab?sslmode=disable\n' "$EXPORTER_DB_PASSWORD" >/etc/loadlab/postgres-exporter.env
chmod 0640 /etc/loadlab/app.env /etc/loadlab/postgres-exporter.env
chown root:loadlab /etc/loadlab/app.env
chown root:postgres /etc/loadlab/postgres-exporter.env

install -m 0644 "$project_dir/deploy/systemd/loadlab.service" /etc/systemd/system/loadlab.service
install -m 0644 "$project_dir/deploy/systemd/loadlab-loadgen.service" /etc/systemd/system/loadlab-loadgen.service
install -m 0644 "$project_dir/deploy/systemd/loadlab-postgres-exporter.service" /etc/systemd/system/loadlab-postgres-exporter.service
install -m 0644 "$project_dir/deploy/systemd/loadlab-prometheus.service" /etc/systemd/system/loadlab-prometheus.service
systemctl daemon-reload
# The Debian package enables its stock Prometheus unit on the same port. The lab
# unit uses the dedicated scrape configuration above, so avoid the port conflict.
systemctl disable --now prometheus.service
systemctl enable --now postgresql prometheus-node-exporter loadlab loadlab-postgres-exporter loadlab-prometheus
echo "Load Lab installed. API: 127.0.0.1:8080; Prometheus: 127.0.0.1:9090"
