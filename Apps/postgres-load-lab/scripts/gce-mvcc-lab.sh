#!/usr/bin/env bash
set -euo pipefail

project_id=${PROJECT_ID:-}
zone=${ZONE:-asia-southeast1-b}
instance=${INSTANCE_NAME:-postgres-mvcc-lab}
machine_type=${MACHINE_TYPE:-e2-standard-2}
disk_size=${DISK_SIZE_GB:-30}
repo_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)

usage() {
  cat <<'EOF'
Usage:
  PROJECT_ID=your-project scripts/gce-mvcc-lab.sh create
  PROJECT_ID=your-project scripts/gce-mvcc-lab.sh status
  PROJECT_ID=your-project scripts/gce-mvcc-lab.sh watch
  PROJECT_ID=your-project scripts/gce-mvcc-lab.sh pause|resume
  PROJECT_ID=your-project scripts/gce-mvcc-lab.sh vacuum
  PROJECT_ID=your-project scripts/gce-mvcc-lab.sh collect [local-file]
  PROJECT_ID=your-project scripts/gce-mvcc-lab.sh metrics
  PROJECT_ID=your-project scripts/gce-mvcc-lab.sh ssh
  PROJECT_ID=your-project scripts/gce-mvcc-lab.sh destroy --yes

Optional environment variables:
  ZONE=asia-southeast1-b INSTANCE_NAME=postgres-mvcc-lab
  MACHINE_TYPE=e2-standard-2 DISK_SIZE_GB=30

The create command provisions billable Google Cloud resources. PostgreSQL and
monitoring ports stay on VM loopback; only SSH access is used.
EOF
}

die() {
  echo "Error: $*" >&2
  exit 1
}

[[ -n $project_id ]] || die "set PROJECT_ID explicitly"
command -v gcloud >/dev/null || die "gcloud is not installed"

gcloud_args=(--project "$project_id" --zone "$zone")

instance_exists() {
  gcloud compute instances describe "$instance" "${gcloud_args[@]}" >/dev/null 2>&1
}

remote() {
  gcloud compute ssh "$instance" "${gcloud_args[@]}" -- "$@"
}

command=${1:-}
case "$command" in
  create)
    active_account=$(gcloud auth list --filter=status:ACTIVE --format='value(account)' | head -n 1)
    [[ -n $active_account ]] || die "no active gcloud account; run gcloud auth login"
    instance_exists && die "instance $instance already exists in $project_id/$zone"
    echo "Creating $instance ($machine_type, ${disk_size}GB) in $project_id/$zone as $active_account"
    gcloud services enable compute.googleapis.com --project "$project_id"
    gcloud compute instances create "$instance" "${gcloud_args[@]}" \
      --machine-type "$machine_type" \
      --boot-disk-size "${disk_size}GB" \
      --boot-disk-type pd-balanced \
      --image-family debian-12 \
      --image-project debian-cloud \
      --no-service-account \
      --no-scopes \
      --labels purpose=postgres-mvcc-learning
    echo "Waiting for SSH..."
    for attempt in $(seq 1 30); do
      if remote true >/dev/null 2>&1; then
        break
      fi
      [[ $attempt -lt 30 ]] || die "SSH did not become ready; inspect the VM in Google Cloud"
      sleep 5
    done
    archive=$(mktemp --suffix=.tar.gz /tmp/mvcc-lab.XXXXXX)
    trap 'rm -f "$archive"' EXIT
    tar -C "$repo_dir/deploy" -czf "$archive" mvcc-lab
    gcloud compute scp "$archive" "$instance:/tmp/mvcc-lab.tar.gz" "${gcloud_args[@]}"
    remote "rm -rf /tmp/mvcc-lab-deploy && mkdir /tmp/mvcc-lab-deploy && tar -xzf /tmp/mvcc-lab.tar.gz -C /tmp/mvcc-lab-deploy && sudo /tmp/mvcc-lab-deploy/mvcc-lab/install.sh"
    echo "Ready. Run: PROJECT_ID=$project_id ZONE=$zone $0 watch"
    ;;
  status)
    instance_exists || die "instance does not exist"
    gcloud compute instances describe "$instance" "${gcloud_args[@]}" \
      --format='table(name,status,machineType.basename(),disks[0].diskSizeGb,networkInterfaces[0].accessConfigs[0].natIP)'
    remote "systemctl --no-pager --full status mvcc-lab-workload mvcc-lab-sampler | head -n 30"
    ;;
  watch)
    instance_exists || die "instance does not exist"
    gcloud compute ssh "$instance" "${gcloud_args[@]}" -- -t "sudo /opt/mvcc-lab/watch.sh"
    ;;
  pause)
    remote "sudo systemctl stop mvcc-lab-workload"
    ;;
  resume)
    remote "sudo systemctl start mvcc-lab-workload"
    ;;
  vacuum)
    remote "sudo -u postgres psql -d mvcc_lab -v ON_ERROR_STOP=1 -c 'ALTER TABLE mvcc_unvacuumed SET (autovacuum_enabled=true, autovacuum_vacuum_threshold=100, autovacuum_vacuum_scale_factor=0.01)' -c 'VACUUM (VERBOSE, ANALYZE) mvcc_unvacuumed'"
    ;;
  collect)
    destination=${2:-results/mvcc-timeline.csv}
    mkdir -p "$(dirname "$destination")"
    gcloud compute scp "$instance:/var/lib/mvcc-lab/timeline.csv" "$destination" "${gcloud_args[@]}"
    echo "Saved $destination"
    ;;
  metrics)
    echo "Prometheus will be available at http://127.0.0.1:9090; press Ctrl-C to close the tunnel."
    gcloud compute ssh "$instance" "${gcloud_args[@]}" -- -N -L 9090:127.0.0.1:9090
    ;;
  ssh)
    gcloud compute ssh "$instance" "${gcloud_args[@]}"
    ;;
  destroy)
    [[ ${2:-} == --yes ]] || die "destroy requires --yes because it permanently deletes the VM and its lab data"
    gcloud compute instances delete "$instance" "${gcloud_args[@]}" --quiet
    ;;
  help|-h|--help|'')
    usage
    ;;
  *)
    usage >&2
    die "unknown command: $command"
    ;;
esac
