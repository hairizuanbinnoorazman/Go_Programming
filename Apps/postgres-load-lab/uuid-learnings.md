# UUIDv4 vs UUIDv7 pagination study

## Question

Does UUIDv7 make pagination faster than UUIDv4?

The short answer is: UUIDv7 helps data locality and gives an approximately chronological key, but **keyset pagination is the main fix for deep-page latency**. Changing UUID versions does not make `OFFSET` stop walking past all preceding rows.

This study keeps two tables with the same schema and payload. Their only intentional difference is how the primary key is generated:

- `uuid_study_v4`: random UUIDv4 keys
- `uuid_study_v7`: time-ordered UUIDv7 keys

The endpoint orders directly by the UUID primary key so that the effect is visible. UUIDv7 order approximately follows creation order. UUIDv4 order is stable but random and should not be presented to users as chronological order.

## Run it locally

Start PostgreSQL and the API, then load 250,000 rows into each study table and run 50 repetitions of every query:

```bash
docker compose up --build -d --wait
GOWORK=off go run ./cmd/uuidstudy -rows 250000 -runs 50
```

The runner truncates only `uuid_study_v4` and `uuid_study_v7` by default. It bulk-loads equal payloads, runs `VACUUM (ANALYZE)`, measures a 100-row page at 90% depth, reports relation sizes and physical correlation, and includes `EXPLAIN (ANALYZE, BUFFERS)` plans. Useful flags are:

```text
-rows 250000       rows loaded per table
-deep-offset -1    -1 means 90% of the row count
-page-size 100
-runs 40
-reset true        truncate the two study tables first
```

The API supports both pagination strategies:

```bash
# First UUIDv7 keyset page (keyset is the default strategy)
curl -sS 'http://127.0.0.1:8080/v1/uuid-study/v7?limit=100'

# Continue using next_cursor from the response
curl -sS 'http://127.0.0.1:8080/v1/uuid-study/v7?limit=100&after=UUID_FROM_NEXT_CURSOR'

# Deliberately request a deep OFFSET page for comparison
curl -sS 'http://127.0.0.1:8080/v1/uuid-study/v7?strategy=offset&limit=100&offset=225000'

# Replace v7 with v4 to query the random-key table
```

Responses contain `items` plus `next_cursor` for keyset pagination or `next_offset` for offset pagination. A continuation field is omitted on the last page.

## Local result

Measured on 2026-08-29 with the repository's local PostgreSQL 17 Compose container. Each table held 250,000 rows, the page size was 100, the tested offset was 225,000, and each latency result used 50 warm-cache runs. These are local lab numbers, not portable production claims.

| Measurement | UUIDv4 | UUIDv7 | Observed difference |
|---|---:|---:|---:|
| Bulk COPY throughput | 502,176 rows/s | 586,935 rows/s | v7 was 16.9% higher |
| Heap size | 44,040,192 B | 44,040,192 B | equal |
| Primary-key index size | 9,863,168 B | 7,905,280 B | v7 was 19.8% smaller |
| ID/heap physical correlation | -0.006 | 1.000 | v4 random; v7 sequential |
| Deep OFFSET p50 | 60.084 ms | 14.338 ms | v7 was 4.2x faster |
| Deep OFFSET p95 | 62.211 ms | 14.640 ms | v7 was 4.2x faster |
| Keyset p50 | 0.496 ms | 0.397 ms | both stayed sub-millisecond |
| Keyset p95 | 0.611 ms | 0.539 ms | both stayed sub-millisecond |

The plans explain the difference more clearly than latency alone:

| Query | Rows visited | Shared buffers | Plan execution time |
|---|---:|---:|---:|
| v4 OFFSET 225,000 | 225,100 | 226,130 | 70.492 ms |
| v7 OFFSET 225,000 | 225,100 | 5,655 | 23.362 ms |
| v4 keyset after cursor | 100 | 104 | 0.118 ms |
| v7 keyset after cursor | 100 | 6 | 0.045 ms |

Both OFFSET plans still visited 225,100 index entries. UUIDv7 was faster because sequential UUID order matched heap order, so fetching table rows required much less scattered access. The keyset predicate (`WHERE id > $cursor ORDER BY id LIMIT $limit`) let PostgreSQL seek into the primary-key index and visit only the requested page. At p50, keyset was about 121x faster than OFFSET for v4 and 36x faster for v7 in this run.

## Ordering and live-insert simulation

The same 250,000 rows were also checked by comparing each row's generation number (encoded in the test payload) with its position after `ORDER BY id`:

| Measurement | UUIDv4 | UUIDv7 |
|---|---:|---:|
| Generation-order/UUID-order correlation | -0.000217 | 1.000000 |
| Rows in their exact generation position | 0 / 250,000 | 250,000 / 250,000 |
| Adjacent order reversals | 125,117 | 0 |

For a live-insert simulation, rows 0–124,999 were treated as the original dataset. A client cursor was placed halfway through that original dataset in UUID order, and rows 125,000–249,999 were treated as records inserted afterward:

| Newly inserted rows relative to existing cursor | UUIDv4 | UUIDv7 |
|---|---:|---:|
| Behind cursor (missed by continuing forward) | 62,674 (50.14%) | 0 (0%) |
| Ahead of cursor | 62,326 (49.86%) | 125,000 (100%) |

This confirms that UUIDv4 has a stable database order but no relationship to generation time. It also confirms exact monotonic UUIDv7 behavior for this run's single Go process and `google/uuid` generator—even though as many as 1,202 rows were inserted within one recorded millisecond. It does **not** prove exact global ordering across processes or machines: independent generators, clock skew, and same-millisecond generation can weaken chronological ordering. For an ascending traversal, later v7 IDs normally land ahead of the cursor; for a newest-first feed, new records normally belong before the current page and are seen when the client refreshes from the beginning.

## What this means

1. **Use keyset pagination to solve deep pagination.** Its work is proportional to the page size rather than the page depth. This matters more than choosing v4 or v7.
2. **UUIDv7 improves write and read locality.** Its smaller index is consistent with fewer disruptive primary-key page splits, and heap order aligned with ID order in this insert-only run.
3. **UUIDv7 makes the ID useful as an approximate time cursor.** UUIDv4 can be keyset-paginated efficiently, but its order is random. With concurrent inserts, a new v4 value can sort before a cursor already passed and be missed by that traversal.
4. **UUIDv7 is not a complete ordering guarantee.** IDs generated in the same millisecond or on different hosts need not reflect exact event order. Use `(created_at, id)` or another explicit immutable sort key when exact application semantics require it.
5. **Do not compare this endpoint directly with the existing list endpoints.** Companies, users, and inventory order by timestamp indexes; changing their ID generator alone would not change those query plans.

## How useful is UUIDv7 here?

UUIDv7 is the better fit when results should normally be presented in approximate chronological order. It combines identity, time-oriented ordering, and a cursor in one value:

```sql
SELECT id, payload, created_at
FROM uuid_study_v7
WHERE id < $cursor
ORDER BY id DESC
LIMIT 100;
```

However, this convenience should not be confused with the pagination optimization itself. PostgreSQL can compare and index UUIDv4 values too, so UUIDv4 supports the same keyset shape:

```sql
SELECT id, payload, created_at
FROM uuid_study_v4
WHERE id < $cursor
ORDER BY id DESC
LIMIT 100;
```

That UUIDv4 query is fast and stable, but its result order is arbitrary rather than chronological. In this run, v4 keyset pagination still had a sub-millisecond p50 (0.496 ms), close to v7's 0.397 ms. Therefore, UUIDv7 is **not especially useful merely to enable fast keyset pagination**—both UUID versions enable it. Its additional value is meaningful approximate time order, better index/heap locality, a smaller measured primary-key index, and more predictable placement of newly generated IDs.

If UUIDv4 records must be displayed chronologically, keyset pagination remains straightforward; the cursor and index just need both the timestamp and UUID tie-breaker:

```sql
SELECT id, payload, created_at
FROM items
WHERE (created_at, id) < ($cursor_created_at, $cursor_id)
ORDER BY created_at DESC, id DESC
LIMIT 100;

CREATE INDEX items_created_id_idx ON items (created_at DESC, id DESC);
```

The practical decision is:

| Requirement | Suitable choice |
|---|---|
| Fast continuation in any stable order | UUIDv4 or UUIDv7 keyset pagination |
| Approximate chronological order with a single-field cursor | UUIDv7 |
| Strict chronological business order | Explicit `(created_at, id)` cursor, regardless of UUID version |
| Direct arbitrary page-number access | `OFFSET`, accepting cost that grows with page depth |

## Study limitations

- The run used warm cache, one local client, an insert-only dataset, and no concurrent writers.
- The runner measures database query time directly so HTTP encoding and network latency do not hide the plan difference.
- Random v4 access suffers more when the heap is larger than memory; results also change after updates, deletes, vacuuming, clustering, and cache eviction.
- `OFFSET` can jump directly to an arbitrary page number. Keyset pagination intentionally trades that feature for stable, bounded continuation latency.
- API cursors should be treated as opaque in a real product, typically by encoding and signing the sort fields rather than exposing a raw UUID contract.
