// Command uuidstudy loads equal UUIDv4 and UUIDv7 datasets, then compares deep
// OFFSET pagination with primary-key keyset pagination. It is intentionally a
// lab utility, not part of the production server image.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type insertResult struct {
	Rows          int     `json:"rows"`
	CopySeconds   float64 `json:"copy_seconds"`
	RowsPerSecond float64 `json:"rows_per_second"`
}

type queryResult struct {
	Runs      int     `json:"runs"`
	P50MS     float64 `json:"p50_ms"`
	P95MS     float64 `json:"p95_ms"`
	MaximumMS float64 `json:"maximum_ms"`
}

type tableResult struct {
	HeapBytes         int64   `json:"heap_bytes"`
	PrimaryIndexBytes int64   `json:"primary_index_bytes"`
	IDCorrelation     float64 `json:"id_physical_correlation"`
}

type report struct {
	RowsPerTable int                     `json:"rows_per_table"`
	PageSize     int                     `json:"page_size"`
	DeepOffset   int                     `json:"deep_offset"`
	Insert       map[string]insertResult `json:"insert"`
	Table        map[string]tableResult  `json:"table"`
	Pagination   map[string]queryResult  `json:"pagination"`
	Plans        map[string][]string     `json:"plans"`
}

type sampleSet struct{ values []time.Duration }

func (s *sampleSet) add(value time.Duration) { s.values = append(s.values, value) }
func (s *sampleSet) result() queryResult {
	sort.Slice(s.values, func(i, j int) bool { return s.values[i] < s.values[j] })
	ms := func(q float64) float64 {
		index := int(float64(len(s.values)-1) * q)
		return float64(s.values[index].Microseconds()) / 1000
	}
	return queryResult{Runs: len(s.values), P50MS: ms(.5), P95MS: ms(.95), MaximumMS: ms(1)}
}

func main() {
	var databaseURL string
	var rows, batchSize, pageSize, deepOffset, runs int
	var reset bool
	flag.StringVar(&databaseURL, "database-url", "postgres://loadlab:loadlab@127.0.0.1:5432/loadlab?sslmode=disable", "PostgreSQL connection URL")
	flag.IntVar(&rows, "rows", 250000, "rows to load into each table")
	flag.IntVar(&batchSize, "batch-size", 10000, "COPY batch size")
	flag.IntVar(&pageSize, "page-size", 100, "rows returned by each pagination query")
	flag.IntVar(&deepOffset, "deep-offset", -1, "offset to test; default is 90% of rows")
	flag.IntVar(&runs, "runs", 40, "timed runs per pagination case")
	flag.BoolVar(&reset, "reset", true, "truncate both study tables before loading")
	flag.Parse()
	if rows < 1 || batchSize < 1 || pageSize < 1 || runs < 1 {
		fatalf("rows, batch-size, page-size, and runs must be positive")
	}
	if deepOffset < 0 {
		deepOffset = rows * 9 / 10
	}
	if deepOffset < 1 || deepOffset+pageSize > rows {
		fatalf("deep-offset must be positive and leave room for one page")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		fatalf("connect: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		fatalf("ping database: %v", err)
	}
	if reset {
		if _, err := pool.Exec(ctx, `TRUNCATE uuid_study_v4, uuid_study_v7`); err != nil {
			fatalf("truncate study tables (has migration 003 run?): %v", err)
		}
	}

	insertDurations := map[string]time.Duration{"v4": 0, "v7": 0}
	insertCounts := map[string]int{"v4": 0, "v7": 0}
	for start := 0; start < rows; start += batchSize {
		count := min(batchSize, rows-start)
		v4Rows, err := makeRows("v4", start, count)
		if err != nil {
			fatalf("generate v4 rows: %v", err)
		}
		v7Rows, err := makeRows("v7", start, count)
		if err != nil {
			fatalf("generate v7 rows: %v", err)
		}
		order := []struct {
			version string
			rows    [][]any
		}{{"v4", v4Rows}, {"v7", v7Rows}}
		if (start/batchSize)%2 == 1 {
			order[0], order[1] = order[1], order[0]
		}
		for _, load := range order {
			started := time.Now()
			copied, err := pool.CopyFrom(ctx, pgx.Identifier{"uuid_study_" + load.version}, []string{"id", "payload"}, pgx.CopyFromRows(load.rows))
			insertDurations[load.version] += time.Since(started)
			if err != nil {
				fatalf("copy %s: %v", load.version, err)
			}
			insertCounts[load.version] += int(copied)
		}
		fmt.Fprintf(os.Stderr, "loaded %d/%d rows per table\r", start+count, rows)
	}
	fmt.Fprintln(os.Stderr)
	for _, table := range []string{"uuid_study_v4", "uuid_study_v7"} {
		if _, err := pool.Exec(ctx, `VACUUM (ANALYZE) `+table); err != nil {
			fatalf("analyze %s: %v", table, err)
		}
	}

	cursors := make(map[string]uuid.UUID)
	for _, version := range []string{"v4", "v7"} {
		query := `SELECT id FROM uuid_study_` + version + ` ORDER BY id LIMIT 1 OFFSET $1`
		var cursor uuid.UUID
		if err := pool.QueryRow(ctx, query, deepOffset-1).Scan(&cursor); err != nil {
			fatalf("find %s cursor: %v", version, err)
		}
		cursors[version] = cursor
	}

	tests := []string{"v4_offset", "v7_offset", "v4_keyset", "v7_keyset"}
	samples := map[string]*sampleSet{}
	for _, name := range tests {
		samples[name] = &sampleSet{}
	}
	// Rotate test order so cache warming and transient host activity are spread
	// across the cases rather than favoring whichever query always runs last.
	for round := 0; round < runs; round++ {
		for index := 0; index < len(tests); index++ {
			name := tests[(round+index)%len(tests)]
			started := time.Now()
			if err := runPage(ctx, pool, name, pageSize, deepOffset, cursors); err != nil {
				fatalf("run %s: %v", name, err)
			}
			samples[name].add(time.Since(started))
		}
	}

	result := report{
		RowsPerTable: rows,
		PageSize:     pageSize,
		DeepOffset:   deepOffset,
		Insert:       make(map[string]insertResult),
		Table:        make(map[string]tableResult),
		Pagination:   make(map[string]queryResult),
		Plans:        make(map[string][]string),
	}
	for _, version := range []string{"v4", "v7"} {
		duration := insertDurations[version]
		result.Insert[version] = insertResult{Rows: insertCounts[version], CopySeconds: duration.Seconds(), RowsPerSecond: float64(insertCounts[version]) / duration.Seconds()}
		stats, err := tableStats(ctx, pool, version)
		if err != nil {
			fatalf("table stats %s: %v", version, err)
		}
		result.Table[version] = stats
	}
	for _, name := range tests {
		result.Pagination[name] = samples[name].result()
		plan, err := explainPage(ctx, pool, name, pageSize, deepOffset, cursors)
		if err != nil {
			fatalf("explain %s: %v", name, err)
		}
		result.Plans[name] = plan
	}
	encoded, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(encoded))
}

func makeRows(version string, start, count int) ([][]any, error) {
	rows := make([][]any, count)
	padding := strings.Repeat("x", 96)
	for index := range rows {
		var id uuid.UUID
		var err error
		if version == "v7" {
			id, err = uuid.NewV7()
		} else {
			id, err = uuid.NewRandom()
		}
		if err != nil {
			return nil, err
		}
		rows[index] = []any{id, fmt.Sprintf("row-%012d-%s", start+index, padding)}
	}
	return rows, nil
}

func runPage(ctx context.Context, pool *pgxpool.Pool, name string, limit, offset int, cursors map[string]uuid.UUID) error {
	parts := strings.Split(name, "_")
	query := `SELECT id,payload,created_at FROM uuid_study_` + parts[0]
	var rows pgx.Rows
	var err error
	if parts[1] == "offset" {
		rows, err = pool.Query(ctx, query+` ORDER BY id LIMIT $1 OFFSET $2`, limit, offset)
	} else {
		rows, err = pool.Query(ctx, query+` WHERE id > $1 ORDER BY id LIMIT $2`, cursors[parts[0]], limit)
	}
	if err != nil {
		return err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var id uuid.UUID
		var payload string
		var createdAt time.Time
		if err := rows.Scan(&id, &payload, &createdAt); err != nil {
			return err
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if count != limit {
		return fmt.Errorf("expected %d rows, got %d", limit, count)
	}
	return nil
}

func tableStats(ctx context.Context, pool *pgxpool.Pool, version string) (tableResult, error) {
	table := "uuid_study_" + version
	var result tableResult
	err := pool.QueryRow(ctx, `SELECT pg_relation_size($1::regclass), pg_relation_size($2::regclass), correlation FROM pg_stats WHERE schemaname=current_schema() AND tablename=$3 AND attname='id'`, table, table+"_pkey", table).Scan(&result.HeapBytes, &result.PrimaryIndexBytes, &result.IDCorrelation)
	return result, err
}

func explainPage(ctx context.Context, pool *pgxpool.Pool, name string, limit, offset int, cursors map[string]uuid.UUID) ([]string, error) {
	parts := strings.Split(name, "_")
	query := `EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT) SELECT id,payload,created_at FROM uuid_study_` + parts[0]
	var rows pgx.Rows
	var err error
	if parts[1] == "offset" {
		rows, err = pool.Query(ctx, query+` ORDER BY id LIMIT $1 OFFSET $2`, limit, offset)
	} else {
		rows, err = pool.Query(ctx, query+` WHERE id > $1 ORDER BY id LIMIT $2`, cursors[parts[0]], limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var lines []string
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
	return lines, rows.Err()
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
