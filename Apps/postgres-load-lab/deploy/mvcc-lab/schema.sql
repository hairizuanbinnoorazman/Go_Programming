\set ON_ERROR_STOP on

DROP TABLE IF EXISTS mvcc_healthy;
DROP TABLE IF EXISTS mvcc_unvacuumed;

CREATE TABLE mvcc_healthy (
    id bigint PRIMARY KEY,
    version bigint NOT NULL DEFAULT 0,
    payload text NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT clock_timestamp()
) WITH (
    fillfactor = 90,
    autovacuum_enabled = true,
    autovacuum_vacuum_threshold = 100,
    autovacuum_vacuum_scale_factor = 0.01,
    autovacuum_analyze_threshold = 100,
    autovacuum_analyze_scale_factor = 0.01
);

CREATE TABLE mvcc_unvacuumed (
    LIKE mvcc_healthy INCLUDING ALL
) WITH (
    fillfactor = 90,
    autovacuum_enabled = false
);

-- Updating this indexed column makes each workload update non-HOT. This is
-- deliberate: it makes both heap-version churn and index maintenance visible.
CREATE INDEX mvcc_healthy_updated_idx ON mvcc_healthy (updated_at);
CREATE INDEX mvcc_unvacuumed_updated_idx ON mvcc_unvacuumed (updated_at);

INSERT INTO mvcc_healthy (id, payload)
SELECT n, repeat(md5(n::text), 4)
FROM generate_series(1, 200000) AS n;

INSERT INTO mvcc_unvacuumed (id, payload)
SELECT id, payload FROM mvcc_healthy;
