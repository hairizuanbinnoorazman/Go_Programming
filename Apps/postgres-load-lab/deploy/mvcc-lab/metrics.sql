\pset tuples_only on
\pset format unaligned
SELECT '# HELP mvcc_lab_live_tuples Estimated live tuples in the lab table.';
SELECT '# TYPE mvcc_lab_live_tuples gauge';
SELECT format('mvcc_lab_live_tuples{relation="%s"} %s', relname, n_live_tup)
FROM pg_stat_user_tables
WHERE relname IN ('mvcc_healthy', 'mvcc_unvacuumed') ORDER BY relname;
SELECT '# HELP mvcc_lab_dead_tuples Estimated dead tuples in the lab table.';
SELECT '# TYPE mvcc_lab_dead_tuples gauge';
SELECT format('mvcc_lab_dead_tuples{relation="%s"} %s', relname, n_dead_tup)
FROM pg_stat_user_tables
WHERE relname IN ('mvcc_healthy', 'mvcc_unvacuumed') ORDER BY relname;
SELECT '# HELP mvcc_lab_heap_bytes Allocated heap bytes.';
SELECT '# TYPE mvcc_lab_heap_bytes gauge';
SELECT format('mvcc_lab_heap_bytes{relation="%s"} %s', relname, pg_relation_size(relid))
FROM pg_stat_user_tables
WHERE relname IN ('mvcc_healthy', 'mvcc_unvacuumed') ORDER BY relname;
SELECT '# HELP mvcc_lab_index_bytes Allocated index bytes.';
SELECT '# TYPE mvcc_lab_index_bytes gauge';
SELECT format('mvcc_lab_index_bytes{relation="%s"} %s', relname, pg_indexes_size(relid))
FROM pg_stat_user_tables
WHERE relname IN ('mvcc_healthy', 'mvcc_unvacuumed') ORDER BY relname;
SELECT '# HELP mvcc_lab_autovacuum_total Autovacuum runs observed for the table.';
SELECT '# TYPE mvcc_lab_autovacuum_total counter';
SELECT format('mvcc_lab_autovacuum_total{relation="%s"} %s', relname, autovacuum_count)
FROM pg_stat_user_tables
WHERE relname IN ('mvcc_healthy', 'mvcc_unvacuumed') ORDER BY relname;
