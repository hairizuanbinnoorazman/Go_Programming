COPY (
    SELECT
        clock_timestamp() AS sampled_at,
        s.relname,
        s.n_live_tup,
        s.n_dead_tup,
        s.n_tup_upd,
        s.n_tup_hot_upd,
        s.vacuum_count,
        s.autovacuum_count,
        s.last_vacuum,
        s.last_autovacuum,
        pg_relation_size(s.relid) AS heap_bytes,
        pg_indexes_size(s.relid) AS index_bytes,
        d.xact_commit,
        d.blks_read,
        d.blks_hit,
        d.temp_bytes,
        COALESCE((
            SELECT max(EXTRACT(epoch FROM clock_timestamp() - a.xact_start))::bigint
            FROM pg_stat_activity AS a
            WHERE a.datname = current_database() AND a.xact_start IS NOT NULL
        ), 0) AS oldest_transaction_seconds
    FROM pg_stat_user_tables AS s
    CROSS JOIN pg_stat_database AS d
    WHERE s.relname IN ('mvcc_healthy', 'mvcc_unvacuumed')
      AND d.datname = current_database()
    ORDER BY s.relname
) TO STDOUT WITH (FORMAT csv);
