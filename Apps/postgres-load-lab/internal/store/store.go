package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

var ErrNotFound = errors.New("not found")

type Company struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        uuid.UUID  `json:"id"`
	CompanyID *uuid.UUID `json:"company_id,omitempty"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Inventory struct {
	ID         uuid.UUID `json:"id"`
	CompanyID  uuid.UUID `json:"company_id"`
	SKU        string    `json:"sku"`
	Name       string    `json:"name"`
	Quantity   int       `json:"quantity"`
	PriceCents int64     `json:"price_cents"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UUIDStudyItem struct {
	ID        uuid.UUID `json:"id"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

type UUIDStudyVersion string

const (
	UUIDStudyV4 UUIDStudyVersion = "v4"
	UUIDStudyV7 UUIDStudyVersion = "v7"
)

type Store struct {
	Pool          *pgxpool.Pool
	timeout       time.Duration
	queryDuration *prometheus.HistogramVec
}

func New(pool *pgxpool.Pool, timeout time.Duration, reg prometheus.Registerer) *Store {
	h := prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "loadlab_db_query_duration_seconds", Help: "Database operation duration.", Buckets: prometheus.DefBuckets}, []string{"operation"})
	reg.MustRegister(h, newPoolCollector(pool))
	return &Store{Pool: pool, timeout: timeout, queryDuration: h}
}

func (s *Store) observe(operation string) func() {
	start := time.Now()
	return func() { s.queryDuration.WithLabelValues(operation).Observe(time.Since(start).Seconds()) }
}

func (s *Store) ctx(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, s.timeout)
}

func Migrate(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		body, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		if _, err := pool.Exec(ctx, string(body)); err != nil {
			return fmt.Errorf("apply migration %s: %w", name, err)
		}
	}
	return nil
}

func (s *Store) CreateCompany(ctx context.Context, c Company) (Company, error) {
	defer s.observe("company_create")()
	ctx, cancel := s.ctx(ctx)
	defer cancel()
	err := s.Pool.QueryRow(ctx, `INSERT INTO companies(id,name) VALUES($1,$2) RETURNING created_at,updated_at`, c.ID, c.Name).Scan(&c.CreatedAt, &c.UpdatedAt)
	return c, err
}
func (s *Store) GetCompany(ctx context.Context, id uuid.UUID) (Company, error) {
	defer s.observe("company_get")()
	ctx, cancel := s.ctx(ctx)
	defer cancel()
	var c Company
	err := s.Pool.QueryRow(ctx, `SELECT id,name,created_at,updated_at FROM companies WHERE id=$1`, id).Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	return c, mapErr(err)
}
func (s *Store) ListCompanies(ctx context.Context, limit, offset int) ([]Company, error) {
	defer s.observe("company_list")()
	ctx, cancel := s.ctx(ctx)
	defer cancel()
	rows, err := s.Pool.Query(ctx, `SELECT id,name,created_at,updated_at FROM companies ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Company, 0)
	for rows.Next() {
		var c Company
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, rows.Err()
}
func (s *Store) UpdateCompany(ctx context.Context, c Company) (Company, error) {
	defer s.observe("company_update")()
	ctx, cancel := s.ctx(ctx)
	defer cancel()
	err := s.Pool.QueryRow(ctx, `UPDATE companies SET name=$2,updated_at=now() WHERE id=$1 RETURNING created_at,updated_at`, c.ID, c.Name).Scan(&c.CreatedAt, &c.UpdatedAt)
	return c, mapErr(err)
}
func (s *Store) DeleteCompany(ctx context.Context, id uuid.UUID) error {
	defer s.observe("company_delete")()
	ctx, cancel := s.ctx(ctx)
	defer cancel()
	tag, err := s.Pool.Exec(ctx, `DELETE FROM companies WHERE id=$1`, id)
	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

func (s *Store) CreateUser(ctx context.Context, u User) (User, error) {
	defer s.observe("user_create")()
	ctx, cancel := s.ctx(ctx)
	defer cancel()
	err := s.Pool.QueryRow(ctx, `INSERT INTO users(id,company_id,name,email) VALUES($1,$2,$3,$4) RETURNING created_at,updated_at`, u.ID, u.CompanyID, u.Name, u.Email).Scan(&u.CreatedAt, &u.UpdatedAt)
	return u, err
}
func (s *Store) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	defer s.observe("user_get")()
	ctx, cancel := s.ctx(ctx)
	defer cancel()
	var u User
	err := s.Pool.QueryRow(ctx, `SELECT id,company_id,name,email,created_at,updated_at FROM users WHERE id=$1`, id).Scan(&u.ID, &u.CompanyID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	return u, mapErr(err)
}
func (s *Store) ListUsers(ctx context.Context, limit, offset int) ([]User, error) {
	defer s.observe("user_list")()
	ctx, cancel := s.ctx(ctx)
	defer cancel()
	rows, err := s.Pool.Query(ctx, `SELECT id,company_id,name,email,created_at,updated_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]User, 0)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.CompanyID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, u)
	}
	return items, rows.Err()
}
func (s *Store) UpdateUser(ctx context.Context, u User) (User, error) {
	defer s.observe("user_update")()
	ctx, cancel := s.ctx(ctx)
	defer cancel()
	err := s.Pool.QueryRow(ctx, `UPDATE users SET company_id=$2,name=$3,email=$4,updated_at=now() WHERE id=$1 RETURNING created_at,updated_at`, u.ID, u.CompanyID, u.Name, u.Email).Scan(&u.CreatedAt, &u.UpdatedAt)
	return u, mapErr(err)
}
func (s *Store) DeleteUser(ctx context.Context, id uuid.UUID) error {
	defer s.observe("user_delete")()
	ctx, cancel := s.ctx(ctx)
	defer cancel()
	tag, err := s.Pool.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

func (s *Store) CreateInventory(ctx context.Context, i Inventory) (Inventory, error) {
	defer s.observe("inventory_create")()
	ctx, cancel := s.ctx(ctx)
	defer cancel()
	err := s.Pool.QueryRow(ctx, `INSERT INTO inventory(id,company_id,sku,name,quantity,price_cents) VALUES($1,$2,$3,$4,$5,$6) RETURNING created_at,updated_at`, i.ID, i.CompanyID, i.SKU, i.Name, i.Quantity, i.PriceCents).Scan(&i.CreatedAt, &i.UpdatedAt)
	return i, err
}
func (s *Store) GetInventory(ctx context.Context, id uuid.UUID) (Inventory, error) {
	defer s.observe("inventory_get")()
	ctx, cancel := s.ctx(ctx)
	defer cancel()
	var i Inventory
	err := s.Pool.QueryRow(ctx, `SELECT id,company_id,sku,name,quantity,price_cents,created_at,updated_at FROM inventory WHERE id=$1`, id).Scan(&i.ID, &i.CompanyID, &i.SKU, &i.Name, &i.Quantity, &i.PriceCents, &i.CreatedAt, &i.UpdatedAt)
	return i, mapErr(err)
}
func (s *Store) ListInventory(ctx context.Context, limit, offset int, companyID *uuid.UUID) ([]Inventory, error) {
	defer s.observe("inventory_list")()
	ctx, cancel := s.ctx(ctx)
	defer cancel()
	query := `SELECT id,company_id,sku,name,quantity,price_cents,created_at,updated_at FROM inventory ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	args := []any{limit, offset}
	if companyID != nil {
		query = `SELECT id,company_id,sku,name,quantity,price_cents,created_at,updated_at FROM inventory WHERE company_id=$1 ORDER BY updated_at DESC LIMIT $2 OFFSET $3`
		args = []any{*companyID, limit, offset}
	}
	rows, err := s.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Inventory, 0)
	for rows.Next() {
		var i Inventory
		if err := rows.Scan(&i.ID, &i.CompanyID, &i.SKU, &i.Name, &i.Quantity, &i.PriceCents, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}
func (s *Store) UpdateInventory(ctx context.Context, i Inventory) (Inventory, error) {
	defer s.observe("inventory_update")()
	ctx, cancel := s.ctx(ctx)
	defer cancel()
	err := s.Pool.QueryRow(ctx, `UPDATE inventory SET company_id=$2,sku=$3,name=$4,quantity=$5,price_cents=$6,updated_at=now() WHERE id=$1 RETURNING created_at,updated_at`, i.ID, i.CompanyID, i.SKU, i.Name, i.Quantity, i.PriceCents).Scan(&i.CreatedAt, &i.UpdatedAt)
	return i, mapErr(err)
}
func (s *Store) DeleteInventory(ctx context.Context, id uuid.UUID) error {
	defer s.observe("inventory_delete")()
	ctx, cancel := s.ctx(ctx)
	defer cancel()
	tag, err := s.Pool.Exec(ctx, `DELETE FROM inventory WHERE id=$1`, id)
	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

// ListUUIDStudy deliberately orders by the UUID primary key. For v7 this is
// approximately creation order; for v4 it is a stable but random order.
func (s *Store) ListUUIDStudy(ctx context.Context, version UUIDStudyVersion, limit, offset int, after *uuid.UUID) ([]UUIDStudyItem, error) {
	defer s.observe("uuid_study_" + string(version) + "_list")()
	ctx, cancel := s.ctx(ctx)
	defer cancel()

	var table string
	switch version {
	case UUIDStudyV4:
		table = "uuid_study_v4"
	case UUIDStudyV7:
		table = "uuid_study_v7"
	default:
		return nil, fmt.Errorf("unsupported UUID study version %q", version)
	}

	query := `SELECT id,payload,created_at FROM ` + table + ` ORDER BY id LIMIT $1 OFFSET $2`
	args := []any{limit, offset}
	if after != nil {
		query = `SELECT id,payload,created_at FROM ` + table + ` WHERE id > $1 ORDER BY id LIMIT $2`
		args = []any{*after, limit}
	}
	rows, err := s.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]UUIDStudyItem, 0, limit)
	for rows.Next() {
		var item UUIDStudyItem
		if err := rows.Scan(&item.ID, &item.Payload, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func mapErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

type poolCollector struct {
	pool                       *pgxpool.Pool
	acquired, idle, total, max *prometheus.Desc
}

func newPoolCollector(p *pgxpool.Pool) *poolCollector {
	return &poolCollector{p, prometheus.NewDesc("loadlab_db_pool_acquired_connections", "Acquired pool connections.", nil, nil), prometheus.NewDesc("loadlab_db_pool_idle_connections", "Idle pool connections.", nil, nil), prometheus.NewDesc("loadlab_db_pool_total_connections", "Total pool connections.", nil, nil), prometheus.NewDesc("loadlab_db_pool_max_connections", "Maximum pool connections.", nil, nil)}
}
func (c *poolCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.acquired
	ch <- c.idle
	ch <- c.total
	ch <- c.max
}
func (c *poolCollector) Collect(ch chan<- prometheus.Metric) {
	s := c.pool.Stat()
	ch <- prometheus.MustNewConstMetric(c.acquired, prometheus.GaugeValue, float64(s.AcquiredConns()))
	ch <- prometheus.MustNewConstMetric(c.idle, prometheus.GaugeValue, float64(s.IdleConns()))
	ch <- prometheus.MustNewConstMetric(c.total, prometheus.GaugeValue, float64(s.TotalConns()))
	ch <- prometheus.MustNewConstMetric(c.max, prometheus.GaugeValue, float64(s.MaxConns()))
}
