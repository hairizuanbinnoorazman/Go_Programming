package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type options struct {
	target                       string
	duration                     time.Duration
	workers, rate, seedCompanies int
	timeout                      time.Duration
	output                       string
	workload                     string
}
type entityPool struct {
	mu  sync.RWMutex
	ids []string
}

func (p *entityPool) add(id string) { p.mu.Lock(); p.ids = append(p.ids, id); p.mu.Unlock() }
func (p *entityPool) pick(r *rand.Rand) string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if len(p.ids) == 0 {
		return ""
	}
	return p.ids[r.Intn(len(p.ids))]
}
func (p *entityPool) take(r *rand.Rand) string {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.ids) == 0 {
		return ""
	}
	n := r.Intn(len(p.ids))
	id := p.ids[n]
	p.ids[n] = p.ids[len(p.ids)-1]
	p.ids = p.ids[:len(p.ids)-1]
	return id
}

type pools struct {
	companies, disposableCompanies, users, disposableUsers, inventory, disposableInventory entityPool
}

type recorder struct {
	mu            sync.Mutex
	latencies     []time.Duration
	total, failed atomic.Uint64
	statuses      sync.Map
	started       time.Time
}

type pacer struct {
	mu       sync.Mutex
	next     time.Time
	interval time.Duration
}

func newPacer(rate int) *pacer {
	return &pacer{next: time.Now(), interval: time.Second / time.Duration(rate)}
}

func (p *pacer) wait(ctx context.Context) bool {
	p.mu.Lock()
	now := time.Now()
	if p.next.Before(now) {
		p.next = now
	}
	scheduled := p.next
	p.next = p.next.Add(p.interval)
	p.mu.Unlock()

	timer := time.NewTimer(time.Until(scheduled))
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

const maxLatencySamples = 1_000_000

func (r *recorder) record(d time.Duration, status int, err error) {
	total := r.total.Add(1)
	if err != nil || status < 200 || status >= 400 {
		r.failed.Add(1)
	}
	key := fmt.Sprintf("%d", status)
	if err != nil {
		key = "transport_error"
	}
	v, _ := r.statuses.LoadOrStore(key, &atomic.Uint64{})
	v.(*atomic.Uint64).Add(1)
	r.mu.Lock()
	if len(r.latencies) < maxLatencySamples {
		r.latencies = append(r.latencies, d)
	} else {
		// Deterministic reservoir sampling caps memory during long saturation runs.
		x := total
		x ^= x >> 12
		x ^= x << 25
		x ^= x >> 27
		candidate := (x * 2_685_821_657_736_338_717) % total
		if candidate < maxLatencySamples {
			r.latencies[candidate] = d
		}
	}
	r.mu.Unlock()
}

type report struct {
	Target            string             `json:"target"`
	Workload          string             `json:"workload"`
	DurationSeconds   float64            `json:"duration_seconds"`
	Requests          uint64             `json:"requests"`
	Failed            uint64             `json:"failed"`
	RequestsPerSecond float64            `json:"requests_per_second"`
	LatencySamples    int                `json:"latency_samples"`
	LatencyMS         map[string]float64 `json:"latency_ms"`
	Statuses          map[string]uint64  `json:"statuses"`
}

func (r *recorder) report(target, workload string) report {
	elapsed := time.Since(r.started)
	r.mu.Lock()
	values := append([]time.Duration(nil), r.latencies...)
	r.mu.Unlock()
	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
	pct := func(q float64) float64 {
		if len(values) == 0 {
			return 0
		}
		i := int(float64(len(values)-1) * q)
		return float64(values[i].Microseconds()) / 1000
	}
	statuses := map[string]uint64{}
	r.statuses.Range(func(k, v any) bool { statuses[k.(string)] = v.(*atomic.Uint64).Load(); return true })
	return report{target, workload, elapsed.Seconds(), r.total.Load(), r.failed.Load(), float64(r.total.Load()) / elapsed.Seconds(), len(values), map[string]float64{"p50": pct(.50), "p95": pct(.95), "p99": pct(.99), "max": pct(1)}, statuses}
}

func main() {
	var o options
	flag.StringVar(&o.target, "target", "http://127.0.0.1:8080", "API base URL")
	flag.DurationVar(&o.duration, "duration", 30*time.Second, "test duration")
	flag.IntVar(&o.workers, "workers", 20, "concurrent workers")
	flag.IntVar(&o.rate, "rate", 0, "total requests/second; 0 means unlimited")
	flag.IntVar(&o.seedCompanies, "seed-companies", 20, "companies to create before the timed run")
	flag.DurationVar(&o.timeout, "timeout", 5*time.Second, "per-request timeout")
	flag.StringVar(&o.output, "output", "", "optional JSON result file")
	flag.StringVar(&o.workload, "workload", "mixed", "workload: mixed or inventory-inserts")
	flag.Parse()
	if o.workers < 1 || o.duration <= 0 || o.rate < 0 {
		fmt.Fprintln(os.Stderr, "workers and duration must be positive; rate cannot be negative")
		os.Exit(2)
	}
	if o.workload != "mixed" && o.workload != "inventory-inserts" {
		fmt.Fprintln(os.Stderr, "workload must be mixed or inventory-inserts")
		os.Exit(2)
	}
	o.target = strings.TrimRight(o.target, "/")
	client := &http.Client{Timeout: o.timeout, Transport: &http.Transport{MaxIdleConns: o.workers * 2, MaxIdleConnsPerHost: o.workers, IdleConnTimeout: 90 * time.Second}}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ps := &pools{}
	if err := seed(ctx, client, o, ps); err != nil {
		fmt.Fprintln(os.Stderr, "seed failed:", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(ctx, o.duration)
	defer cancel()
	rec := &recorder{started: time.Now()}
	var limiter *pacer
	if o.rate > 0 {
		limiter = newPacer(o.rate)
	}
	var wg sync.WaitGroup
	for i := 0; i < o.workers; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(worker)))
			for {
				if limiter != nil {
					if !limiter.wait(ctx) {
						return
					}
				} else {
					select {
					case <-ctx.Done():
						return
					default:
					}
				}
				method, path, body, after := operation(r, ps, o.workload)
				if path == "" {
					continue
				}
				started := time.Now()
				status, response, err := request(ctx, client, method, o.target+path, body)
				if err != nil && ctx.Err() != nil {
					return
				}
				rec.record(time.Since(started), status, err)
				if err == nil && status >= 200 && status < 300 && after != nil {
					after(response)
				}
			}
		}(i)
	}
	wg.Wait()
	result := rec.report(o.target, o.workload)
	encoded, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(encoded))
	if o.output != "" {
		if err := os.MkdirAll(filepath.Dir(o.output), 0o755); err != nil {
			fmt.Fprintln(os.Stderr, "create result directory:", err)
			os.Exit(1)
		}
		if err := os.WriteFile(o.output, append(encoded, '\n'), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, "write result:", err)
			os.Exit(1)
		}
	}
}

func seed(ctx context.Context, c *http.Client, o options, p *pools) error {
	for i := 0; i < o.seedCompanies; i++ {
		status, body, err := request(ctx, c, http.MethodPost, o.target+"/v1/companies", map[string]any{"name": fmt.Sprintf("Seed Company %s", uuid.NewString())})
		if err != nil {
			return err
		}
		if status != http.StatusCreated {
			return fmt.Errorf("create company returned %d: %s", status, body)
		}
		var v struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(body, &v); err != nil {
			return err
		}
		p.companies.add(v.ID)
	}
	return nil
}

type callback func([]byte)

func operation(r *rand.Rand, p *pools, workload string) (string, string, any, callback) {
	n := r.Intn(100)
	suffix := uuid.NewString()
	company := p.companies.pick(r)
	if workload == "inventory-inserts" {
		return http.MethodPost, "/v1/inventory", map[string]any{
			"company_id":  company,
			"sku":         "LOAD-" + suffix,
			"name":        "Load Item " + suffix,
			"quantity":    r.Intn(1000),
			"price_cents": r.Intn(100000),
		}, nil
	}
	switch {
	case n < 3:
		return http.MethodPost, "/v1/companies", map[string]any{"name": "Disposable Company " + suffix}, func(b []byte) { addID(b, &p.disposableCompanies) }
	case n < 8:
		return http.MethodPost, "/v1/companies", map[string]any{"name": "Company " + suffix}, func(b []byte) { addID(b, &p.companies) }
	case n < 10:
		id := p.disposableCompanies.take(r)
		if id == "" {
			return "", "", nil, nil
		}
		return http.MethodDelete, "/v1/companies/" + id, nil, nil
	case n < 15:
		id := p.companies.pick(r)
		if id == "" {
			return "", "", nil, nil
		}
		return http.MethodPut, "/v1/companies/" + id, map[string]any{"name": "Updated " + suffix}, nil
	case n < 20:
		return http.MethodGet, "/v1/companies?limit=50", nil, nil
	case n < 24:
		if company == "" {
			return "", "", nil, nil
		}
		return http.MethodPost, "/v1/users", map[string]any{"company_id": company, "name": "Disposable User " + suffix, "email": suffix + "@load.test"}, func(b []byte) { addID(b, &p.disposableUsers) }
	case n < 40:
		if company == "" {
			return "", "", nil, nil
		}
		return http.MethodPost, "/v1/users", map[string]any{"company_id": company, "name": "User " + suffix, "email": suffix + "@load.test"}, func(b []byte) { addID(b, &p.users) }
	case n < 47:
		id := p.users.pick(r)
		if id == "" {
			return http.MethodGet, "/v1/users?limit=50", nil, nil
		}
		return http.MethodGet, "/v1/users/" + id, nil, nil
	case n < 53:
		id := p.users.pick(r)
		if id == "" || company == "" {
			return "", "", nil, nil
		}
		return http.MethodPut, "/v1/users/" + id, map[string]any{"company_id": company, "name": "Updated User " + suffix, "email": suffix + "@load.test"}, nil
	case n < 58:
		id := p.disposableUsers.take(r)
		if id == "" {
			return "", "", nil, nil
		}
		return http.MethodDelete, "/v1/users/" + id, nil, nil
	case n < 62:
		if company == "" {
			return "", "", nil, nil
		}
		return http.MethodPost, "/v1/inventory", map[string]any{"company_id": company, "sku": "DISPOSABLE-" + suffix, "name": "Disposable Item " + suffix, "quantity": r.Intn(1000), "price_cents": r.Intn(100000)}, func(b []byte) { addID(b, &p.disposableInventory) }
	case n < 78:
		if company == "" {
			return "", "", nil, nil
		}
		return http.MethodPost, "/v1/inventory", map[string]any{"company_id": company, "sku": "SKU-" + suffix, "name": "Item " + suffix, "quantity": r.Intn(1000), "price_cents": r.Intn(100000)}, func(b []byte) { addID(b, &p.inventory) }
	case n < 85:
		id := p.inventory.pick(r)
		if id == "" {
			return http.MethodGet, "/v1/inventory?limit=50", nil, nil
		}
		return http.MethodGet, "/v1/inventory/" + id, nil, nil
	case n < 91:
		id := p.inventory.pick(r)
		if id == "" || company == "" {
			return "", "", nil, nil
		}
		return http.MethodPut, "/v1/inventory/" + id, map[string]any{"company_id": company, "sku": "SKU-" + suffix, "name": "Updated Item", "quantity": r.Intn(1000), "price_cents": r.Intn(100000)}, nil
	case n < 96:
		id := p.disposableInventory.take(r)
		if id == "" {
			return "", "", nil, nil
		}
		return http.MethodDelete, "/v1/inventory/" + id, nil, nil
	default:
		if company != "" {
			return http.MethodGet, "/v1/inventory?company_id=" + company + "&limit=100", nil, nil
		}
		return http.MethodGet, "/v1/inventory?limit=100", nil, nil
	}
}
func addID(body []byte, p *entityPool) {
	var v struct {
		ID string `json:"id"`
	}
	if json.Unmarshal(body, &v) == nil && v.ID != "" {
		p.add(v.ID)
	}
}
func request(ctx context.Context, c *http.Client, method, url string, payload any) (int, []byte, error) {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return 0, nil, err
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return 0, nil, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	return resp.StatusCode, b, err
}
