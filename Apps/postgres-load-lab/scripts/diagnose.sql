\pset pager off
\timing on

SELECT now(), version();
SELECT datname, numbackends, xact_commit, xact_rollback, blks_read, blks_hit,
       CASE WHEN blks_hit + blks_read = 0 THEN 0 ELSE round(100.0 * blks_hit / (blks_hit + blks_read), 2) END AS cache_hit_pct,
       temp_files, pg_size_pretty(temp_bytes) AS temp_bytes, deadlocks
FROM pg_stat_database WHERE datname = current_database();

SELECT wait_event_type, wait_event, state, count(*)
FROM pg_stat_activity WHERE datname = current_database()
GROUP BY wait_event_type, wait_event, state ORDER BY count(*) DESC;

SELECT relname, n_live_tup, n_dead_tup, seq_scan, idx_scan,
       last_autovacuum, last_autoanalyze
FROM pg_stat_user_tables ORDER BY n_live_tup DESC;

SELECT calls, round(total_exec_time::numeric, 2) AS total_ms,
       round(mean_exec_time::numeric, 2) AS mean_ms,
       shared_blks_read, shared_blks_hit, temp_blks_written,
       left(query, 140) AS query
FROM pg_stat_statements
ORDER BY total_exec_time DESC LIMIT 20;

SELECT pg_size_pretty(pg_database_size(current_database())) AS database_size;
