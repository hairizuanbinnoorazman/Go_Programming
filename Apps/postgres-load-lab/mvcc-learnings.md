# PostgreSQL MVCC, autovacuum, and bloat lab

This document is the single guide for provisioning, running, observing, and
removing the Google Compute Engine lab. The lab is deliberately disposable. It
creates sustained write churn so PostgreSQL's MVCC behavior, autovacuum work,
and relation growth become visible over time.

## The mental model

PostgreSQL uses multi-version concurrency control (MVCC). An `UPDATE` normally
creates a new physical tuple version rather than overwriting the tuple that an
older snapshot might still need. A `DELETE` makes a tuple obsolete but does not
immediately remove its storage. This design lets readers and writers overlap,
but it leaves obsolete tuple versions behind.

Once no active snapshot can see an obsolete version, `VACUUM` can mark its
space reusable. Autovacuum decides when to run `VACUUM` and `ANALYZE` in the
background. Ordinary vacuum usually does not shrink a relation's operating-
system file; it makes space inside that relation available for reuse. If vacuum
cannot keep up, cannot remove versions because an old snapshot still needs
them, or space is reused poorly, table and index bloat can grow.

The causal chain is:

```text
UPDATE or DELETE
  -> MVCC leaves an obsolete tuple version
  -> the version becomes dead when no snapshot can see it
  -> VACUUM makes its space reusable
  -> insufficient cleanup or reuse leads to growing relations and extra work
```

Dead tuples and bloat are related but not identical. `n_dead_tup` is an
estimate of obsolete tuples. Bloat describes inefficiently occupied allocated
pages. A table can have few currently dead tuples while retaining a large file
whose free space is ready for future reuse.

## What the experiment creates

The VM runs PostgreSQL, `pgbench`, Prometheus, node exporter, and a PostgreSQL
statistics sampler. PostgreSQL and Prometheus listen only on VM loopback; use
SSH or an SSH tunnel to reach them.

Two 200,000-row tables receive the same updates in each transaction:

| Table | Maintenance policy | Expected behavior |
|---|---|---|
| `mvcc_healthy` | Aggressive table-level autovacuum settings | Dead tuples rise and fall; autovacuum count increases; size should approach a reusable steady state |
| `mvcc_unvacuumed` | Autovacuum deliberately disabled | Dead tuples and allocated heap/index bytes tend to grow until the recovery phase |

Both tables have an index on `updated_at`, and every transaction changes that
column. This intentionally prevents HOT updates so index-version maintenance is
visible too. It is an educational stress case, not a recommended schema.

The workload defaults to 500 transactions per second, eight clients, and two
`pgbench` threads. Each transaction updates one random row in both tables. The
service restarts each hour and continues until paused or the VM is deleted.

## Cost and safety boundary

The default VM is an `e2-standard-2` with a 30 GB balanced persistent disk in
`asia-southeast1-b`. It incurs Google Cloud compute and storage charges until
deleted. Pricing varies by account, discounts, and region; check the current
Google Cloud pricing page before creation.

The provisioning script:

- Requires an explicit `PROJECT_ID`.
- Uses a named zone and instance on every command.
- Gives the VM no Google Cloud service account or API scopes.
- Does not open PostgreSQL, Prometheus, or node-exporter firewall ports.
- Labels the VM `purpose=postgres-mvcc-learning`.
- Requires `--yes` before deleting the VM.

## Prerequisites

Install the Google Cloud CLI, authenticate, and choose an existing project with
billing enabled:

```bash
gcloud auth login
gcloud auth list
gcloud projects describe YOUR_PROJECT_ID
```

You need permission to enable the Compute Engine API, create/delete instances,
and connect using Compute Engine SSH/OS Login. The project must have a usable
network that permits SSH access. The script defaults to the project's default
network.

## Create the lab

From this repository:

```bash
PROJECT_ID=YOUR_PROJECT_ID \
ZONE=asia-southeast1-b \
scripts/gce-mvcc-lab.sh create
```

The command enables the Compute Engine API, creates the billable VM, copies the
lab assets, installs packages, initializes both tables, and starts the workload
and sampler. Override the defaults if necessary:

```bash
PROJECT_ID=YOUR_PROJECT_ID \
ZONE=asia-southeast1-c \
INSTANCE_NAME=postgres-mvcc-lab \
MACHINE_TYPE=e2-standard-2 \
DISK_SIZE_GB=30 \
scripts/gce-mvcc-lab.sh create
```

### Completed experiment

The lab ran on 2026-08-29 with this configuration:

| Setting | Deployed value |
|---|---|
| Project | `new-demo-project-462517` |
| Zone | `asia-southeast1-b` |
| Instance | `postgres-mvcc-lab` |
| Machine | `e2-standard-2` |
| Boot disk | 30 GB `pd-balanced` |
| Database | PostgreSQL 15.19 from Debian 12 |
| Workload | 500 transactions/second target, eight clients |
| Active workload period | 11:38:47–13:20:31 UTC (1 hour, 41 minutes, 44 seconds) |
| Sampling period | 11:38:47–13:21:57 UTC, every 10 seconds |

The Compute Engine instance and its auto-delete boot disk were deleted at
approximately 13:24 UTC after collection. Direct resource lookups returned
`not found`, confirming that this completed experiment is no longer accruing VM
or boot-disk charges.

The workload completed approximately 3.05 million transactions, with every
transaction updating one row in each comparison table. The complete 1,223-line
timeline is preserved at `results/mvcc-timeline.csv`. Its SHA-256 checksum is
`2e521d79be56579151650a2b787c97a51d4c2c1d8ad300a59f5cfecf1b01166b`.

## Measured results

### Behavior during sustained writes

Both tables began with 200,000 live rows, a 41,032 KiB heap, and 8,816 KiB of
indexes. The workload deliberately changed the indexed `updated_at` column, so
all 3.05 million updates were non-HOT.

| Measurement | Healthy autovacuum | Autovacuum disabled |
|---|---:|---:|
| Updates | 3,049,253 | 3,049,252 |
| Autovacuums during load | 102 | 0 |
| Maximum estimated dead tuples | 30,347 | 3,049,252 |
| Final estimated dead tuples before recovery | 9,199 | 3,049,252 |
| Final heap size | 40 MiB | 93 MiB |
| Final index size | 33 MiB | 78 MiB |
| Final total size | 73 MiB | 171 MiB |

The healthy table showed the expected sawtooth pattern: dead tuples rose until
the threshold was crossed, then fell when autovacuum ran. Its heap remained at
40 MiB because reclaimed heap space could be reused. Its indexes still grew
from about 9 MiB to 33 MiB, demonstrating that successful autovacuum does not
guarantee minimum relation size or eliminate all index growth.

The unvacuumed table's estimated dead tuples climbed with update count. Its
total allocation reached about 2.34 times the healthy table's size under the
same logical data and write workload. Because there were no old transactions
in the sampled period, this experiment isolated disabled cleanup rather than a
long-running snapshot preventing cleanup.

### Ordinary vacuum recovery

After writes were paused, autovacuum was re-enabled on `mvcc_unvacuumed` and
`VACUUM (VERBOSE, ANALYZE)` was run. PostgreSQL reported:

| Vacuum measurement | Result |
|---|---:|
| Elapsed database time | 2.06 seconds |
| Heap pages scanned | 11,862 |
| Dead index item identifiers removed | 3,049,252 |
| Heap tuples physically removed by vacuum | 4,495 |
| WAL generated | 53,966,922 bytes |
| Dirty buffers | 18,048 |
| Dead tuples afterward | 0 |
| Total relation size afterward | 171 MiB |

The mismatch between 3.05 million estimated dead tuples and 4,495 heap tuples
physically removed is instructive. Normal update-time page pruning had already
removed many obsolete heap versions, even with autovacuum disabled, while
millions of dead index references still needed index vacuuming. Therefore
`n_dead_tup` is useful as a maintenance signal but is not an exact physical
bloat measurement.

Ordinary vacuum reset the dead-tuple estimate and made space reusable, but did
not reduce the allocated 171 MiB. This directly demonstrates why a successful
vacuum should not be judged by whether the operating-system file shrinks.

### `VACUUM FULL` compaction

With writes still paused, `VACUUM (FULL, VERBOSE, ANALYZE)` rewrote the table:

| Measurement | Before `VACUUM FULL` | After `VACUUM FULL` |
|---|---:|---:|
| Heap | 93 MiB | 20 MiB |
| Indexes | 78 MiB | 8.8 MiB |
| Total | 171 MiB | 28 MiB |
| Live rows | 200,000 | 200,000 |
| Dead rows | 0 | 0 |

This returned most unused allocation to the operating system, but required a
complete rewrite and an `ACCESS EXCLUSIVE` lock. It proves the distinction
between reclamation for internal reuse and physical compaction; it does not
make `VACUUM FULL` appropriate as routine maintenance.

### Conclusions from this run

1. MVCC churn affects more than disk capacity. Without cleanup, the same live
   dataset occupied over twice the space and required much larger heap and
   index working sets.
2. Autovacuum kept dead tuples bounded and allowed heap reuse under a continuous
   500 TPS update workload.
3. Autovacuum is maintenance, not compaction. Even the healthy table's indexes
   grew, and ordinary vacuum did not shrink allocated files.
4. Index maintenance dominated this deliberately non-HOT workload. Schema and
   index choices directly affect how expensive updates and vacuum become.
5. `n_dead_tup` is an estimate and cannot by itself quantify physical bloat.
   Combine it with relation sizes, vacuum history, workload counters, and, when
   needed, more precise inspection tools.
6. `VACUUM FULL` can reclaim operating-system space, but its rewrite and
   exclusive lock make prevention and ordinary vacuum preferable.
7. Long-running-snapshot interference was not present in this run: sampled
   oldest-transaction age stayed at zero seconds. That should be studied as a
   separate experiment rather than inferred from these results.

Check provisioning and service status:

```bash
PROJECT_ID=YOUR_PROJECT_ID scripts/gce-mvcc-lab.sh status
```

If installation fails, connect and inspect both the package installation and
services:

```bash
PROJECT_ID=YOUR_PROJECT_ID scripts/gce-mvcc-lab.sh ssh
sudo systemctl status postgresql mvcc-lab-workload mvcc-lab-sampler
sudo journalctl -u mvcc-lab-workload -u mvcc-lab-sampler --since '30 minutes ago'
```

The Debian distribution's supported PostgreSQL package is installed; record
the exact version shown by the installer rather than assuming a major version.

## Follow the live database state

Open the terminal dashboard:

```bash
PROJECT_ID=YOUR_PROJECT_ID scripts/gce-mvcc-lab.sh watch
```

Every five seconds it shows:

- Estimated live and dead tuples and dead-tuple percentage.
- Total updates and HOT updates.
- Autovacuum count and most recent autovacuum time.
- Allocated heap, index, and total relation sizes.
- Active vacuum progress.
- The oldest active transactions and their wait events.

Expected healthy-table pattern:

```text
dead tuples:       /| /| /|       repeated rise and cleanup
autovacuum count:  _/--/--/       increases over time
allocated bytes:   /-----          grows, then tends toward reuse
```

Expected unvacuumed-table pattern:

```text
dead tuples:       /^^^^^^^^       generally accumulates
autovacuum count:  __________      stays at zero
allocated bytes:   /^^^^^^^^       continues growing for longer
```

These are tendencies, not exact graphs. The statistics collector reports
estimates, physical page reuse is workload-dependent, and autovacuum timing is
asynchronous.

## Follow the Prometheus time series

Create the SSH tunnel:

```bash
PROJECT_ID=YOUR_PROJECT_ID scripts/gce-mvcc-lab.sh metrics
```

While it is running, open <http://127.0.0.1:9090>. Useful PromQL queries are:

```promql
mvcc_lab_dead_tuples
mvcc_lab_live_tuples
mvcc_lab_heap_bytes
mvcc_lab_index_bytes
increase(mvcc_lab_autovacuum_total[30m])
rate(node_cpu_seconds_total{mode!="idle"}[5m])
node_memory_MemAvailable_bytes
rate(node_disk_read_bytes_total[5m])
rate(node_disk_written_bytes_total[5m])
```

Graph the two `relation` label values together. The comparison matters more
than any isolated absolute number.

The sampler also appends raw observations every ten seconds. Download them for
analysis or plotting:

```bash
PROJECT_ID=YOUR_PROJECT_ID \
scripts/gce-mvcc-lab.sh collect results/mvcc-timeline.csv
```

The CSV includes table counters and sizes plus database commits, cache reads
and hits, temporary bytes, and oldest-transaction age.

## Suggested experiment sequence

### Phase 1: Establish the contrast

Let the workload run for at least 30–60 minutes. Use both the terminal view and
Prometheus. Record:

1. How frequently `mvcc_healthy` is autovacuumed.
2. The peak and trough of its estimated dead tuples.
3. Whether its heap and indexes reach a rough steady state.
4. How the unvacuumed table's dead tuples and allocated bytes differ.
5. CPU and disk-write behavior while autovacuum runs.

Do not expect autovacuum to keep a table at its minimum possible size. Its goal
is sustainable reuse and required maintenance, not continuous compaction.

### Phase 2: Pause writes and inspect

```bash
PROJECT_ID=YOUR_PROJECT_ID scripts/gce-mvcc-lab.sh pause
PROJECT_ID=YOUR_PROJECT_ID scripts/gce-mvcc-lab.sh watch
```

Pausing separates maintenance behavior from new churn. The sampler remains
active. Resume with:

```bash
PROJECT_ID=YOUR_PROJECT_ID scripts/gce-mvcc-lab.sh resume
```

### Phase 3: Recover the unhealthy table

This command re-enables aggressive autovacuum settings on the unhealthy table
and runs an ordinary verbose vacuum:

```bash
PROJECT_ID=YOUR_PROJECT_ID scripts/gce-mvcc-lab.sh vacuum
```

Observe that dead tuples should fall, but the allocated heap and index files
may not shrink substantially. This is normal: ordinary vacuum makes space
reusable while allowing normal reads and writes to continue.

Resume or leave the workload running and watch whether subsequent updates reuse
the available space rather than growing the relation at its earlier rate.

### Phase 4: Demonstrate physical compaction

Only after recording the ordinary-vacuum results, pause the workload and record
the size:

```bash
PROJECT_ID=YOUR_PROJECT_ID scripts/gce-mvcc-lab.sh pause
PROJECT_ID=YOUR_PROJECT_ID scripts/gce-mvcc-lab.sh ssh
```

On the VM:

```bash
sudo -u postgres psql -d mvcc_lab
```

Then:

```sql
SELECT pg_size_pretty(pg_total_relation_size('mvcc_unvacuumed'));
VACUUM (FULL, ANALYZE) mvcc_unvacuumed;
SELECT pg_size_pretty(pg_total_relation_size('mvcc_unvacuumed'));
```

`VACUUM FULL` rewrites the table and requires an `ACCESS EXCLUSIVE` lock. It is
not routine maintenance and is intentionally not automated by this lab. After
leaving `psql`, resume if desired:

```bash
PROJECT_ID=YOUR_PROJECT_ID scripts/gce-mvcc-lab.sh resume
```

## Optional: prove that an old snapshot limits cleanup

Pause the workload first so the result is easier to interpret. Open two SSH
sessions. In session A:

```sql
BEGIN ISOLATION LEVEL REPEATABLE READ;
SELECT version FROM mvcc_healthy WHERE id = 1;
```

Keep that transaction open. In session B, update the same row many times and
vacuum:

```sql
DO $$
BEGIN
  FOR n IN 1..10000 LOOP
    UPDATE mvcc_healthy
    SET version = version + 1, updated_at = clock_timestamp()
    WHERE id = 1;
  END LOOP;
END $$;
VACUUM (VERBOSE, ANALYZE) mvcc_healthy;
```

Session A still needs the old version visible to its snapshot, so vacuum cannot
remove everything that transaction might see. Commit session A and vacuum
again:

```sql
COMMIT;
VACUUM (VERBOSE, ANALYZE) mvcc_healthy;
```

This is why old transactions, abandoned `idle in transaction` sessions, and
some replication slots can prevent normal cleanup even though autovacuum is
enabled and running.

## How autovacuum decides to run

For update/delete churn, the approximate table threshold is:

```text
autovacuum_vacuum_threshold
  + autovacuum_vacuum_scale_factor * estimated table rows
```

The healthy table uses a base threshold of 100 and a scale factor of 0.01, so a
200,000-row table becomes eligible at roughly 2,100 obsolete tuples. Eligibility
does not mean vacuum starts at that exact instant: the launcher schedule,
available workers, cost throttling, and other eligible tables affect timing.

Inspect actual global and table settings:

```sql
SELECT name, setting, unit
FROM pg_settings
WHERE name LIKE 'autovacuum%'
ORDER BY name;

SELECT relname, reloptions
FROM pg_class
WHERE relname IN ('mvcc_healthy', 'mvcc_unvacuumed');
```

Disabling autovacuum globally is not a legitimate tuning strategy. PostgreSQL
also needs vacuum to maintain visibility information and protect against
transaction-ID wraparound.

## Interpretation pitfalls

- `n_dead_tup` is an estimate, not an exact bloat percentage.
- Ordinary vacuum usually returns space to the table, not the operating system.
- Relation size remaining flat after vacuum can mean successful future reuse,
  not failed vacuuming.
- This workload deliberately updates an indexed column. Updating only
  non-indexed columns may permit HOT updates and produce different index growth.
- Autovacuum consumes CPU and I/O because cleanup is real work. The goal is to
  keep maintenance sustainable, not eliminate its resource use.
- A short local or cloud experiment does not produce universal production
  tuning values.
- The intentionally disabled table is a teaching control. Do not copy that
  setting into production.

## Stop charges and remove the lab

Collect the CSV first if you want to preserve the results. Then permanently
delete the VM and its boot disk:

```bash
PROJECT_ID=YOUR_PROJECT_ID scripts/gce-mvcc-lab.sh destroy --yes
```

Verify that it is gone:

```bash
gcloud compute instances list \
  --project YOUR_PROJECT_ID \
  --filter='labels.purpose=postgres-mvcc-learning'
```

## Primary references

- [PostgreSQL 15 concurrency control and MVCC](https://www.postgresql.org/docs/15/mvcc.html)
- [PostgreSQL 15 routine vacuuming and autovacuum](https://www.postgresql.org/docs/15/routine-vacuuming.html)
- [PostgreSQL 15 `VACUUM`](https://www.postgresql.org/docs/15/sql-vacuum.html)
- [PostgreSQL 15 statistics views](https://www.postgresql.org/docs/15/monitoring-stats.html)
- [PostgreSQL 15 vacuum progress reporting](https://www.postgresql.org/docs/15/progress-reporting.html)
- [Google Compute Engine instance creation](https://cloud.google.com/compute/docs/instances/instance-creation-overview)
- [Google Cloud CLI `compute instances create`](https://cloud.google.com/sdk/gcloud/reference/compute/instances/create)
