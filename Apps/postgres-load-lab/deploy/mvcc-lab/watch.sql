\pset pager off
SELECT
    clock_timestamp()::timestamp(0) AS observed_at,
    relname,
    n_live_tup,
    n_dead_tup,
    round(100.0 * n_dead_tup / GREATEST(n_live_tup + n_dead_tup, 1), 1) AS dead_pct,
    n_tup_upd,
    n_tup_hot_upd,
    autovacuum_count AS auto_vacuums,
    to_char(last_autovacuum, 'HH24:MI:SS') AS last_autovac,
    pg_size_pretty(pg_relation_size(relid)) AS heap,
    pg_size_pretty(pg_indexes_size(relid)) AS indexes,
    pg_size_pretty(pg_total_relation_size(relid)) AS total
FROM pg_stat_user_tables
WHERE relname IN ('mvcc_healthy', 'mvcc_unvacuumed')
ORDER BY relname;

SELECT
    p.pid,
    p.relid::regclass AS relation,
    p.phase,
    p.heap_blks_scanned,
    p.heap_blks_total
FROM pg_stat_progress_vacuum AS p;

SELECT
    pid,
    state,
    clock_timestamp() - xact_start AS transaction_age,
    wait_event_type,
    wait_event,
    left(query, 70) AS query
FROM pg_stat_activity
WHERE datname = current_database()
  AND xact_start IS NOT NULL
ORDER BY xact_start
LIMIT 5;
