package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hairizuan/go-programming/apps/postgres-load-lab/internal/store"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type API struct {
	store    *store.Store
	log      *slog.Logger
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
}

func New(s *store.Store, log *slog.Logger, reg *prometheus.Registry) http.Handler {
	a := &API{store: s, log: log,
		requests: prometheus.NewCounterVec(prometheus.CounterOpts{Name: "loadlab_http_requests_total", Help: "HTTP requests by route, method, and status."}, []string{"route", "method", "status"}),
		duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "loadlab_http_request_duration_seconds", Help: "HTTP request duration by route and method.", Buckets: []float64{.001, .0025, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5}}, []string{"route", "method"}),
	}
	reg.MustRegister(a.requests, a.duration)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /readyz", a.ready)
	mux.Handle("GET /metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	mux.HandleFunc("POST /v1/companies", a.createCompany)
	mux.HandleFunc("GET /v1/companies", a.listCompanies)
	mux.HandleFunc("GET /v1/companies/{id}", a.getCompany)
	mux.HandleFunc("PUT /v1/companies/{id}", a.updateCompany)
	mux.HandleFunc("DELETE /v1/companies/{id}", a.deleteCompany)
	mux.HandleFunc("POST /v1/users", a.createUser)
	mux.HandleFunc("GET /v1/users", a.listUsers)
	mux.HandleFunc("GET /v1/users/{id}", a.getUser)
	mux.HandleFunc("PUT /v1/users/{id}", a.updateUser)
	mux.HandleFunc("DELETE /v1/users/{id}", a.deleteUser)
	mux.HandleFunc("POST /v1/inventory", a.createInventory)
	mux.HandleFunc("GET /v1/inventory", a.listInventory)
	mux.HandleFunc("GET /v1/inventory/{id}", a.getInventory)
	mux.HandleFunc("PUT /v1/inventory/{id}", a.updateInventory)
	mux.HandleFunc("DELETE /v1/inventory/{id}", a.deleteInventory)
	mux.HandleFunc("GET /v1/uuid-study/{version}", a.listUUIDStudy)
	return a.instrument(mux)
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(code int) { w.status = code; w.ResponseWriter.WriteHeader(code) }
func (a *API) instrument(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		route := r.Pattern
		if route == "" {
			route = "unmatched"
		}
		a.requests.WithLabelValues(route, r.Method, strconv.Itoa(rw.status)).Inc()
		a.duration.WithLabelValues(route, r.Method).Observe(time.Since(start).Seconds())
	})
}

func (a *API) ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := a.storeContext(r)
	defer cancel()
	if err := a.store.Pool.Ping(ctx); err != nil {
		writeError(w, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}
func (a *API) storeContext(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), 2*time.Second)
}

type companyInput struct {
	Name string `json:"name"`
}

func (a *API) createCompany(w http.ResponseWriter, r *http.Request) {
	var in companyInput
	if !decode(w, r, &in) || !required(w, in.Name, "name") {
		return
	}
	c, err := a.store.CreateCompany(r.Context(), store.Company{ID: uuid.New(), Name: strings.TrimSpace(in.Name)})
	respond(w, c, err, http.StatusCreated)
}
func (a *API) getCompany(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	c, err := a.store.GetCompany(r.Context(), id)
	respond(w, c, err, http.StatusOK)
}
func (a *API) listCompanies(w http.ResponseWriter, r *http.Request) {
	limit, offset, ok := page(w, r)
	if !ok {
		return
	}
	v, err := a.store.ListCompanies(r.Context(), limit, offset)
	respond(w, v, err, http.StatusOK)
}
func (a *API) updateCompany(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var in companyInput
	if !decode(w, r, &in) || !required(w, in.Name, "name") {
		return
	}
	v, err := a.store.UpdateCompany(r.Context(), store.Company{ID: id, Name: strings.TrimSpace(in.Name)})
	respond(w, v, err, http.StatusOK)
}
func (a *API) deleteCompany(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := a.store.DeleteCompany(r.Context(), id); err != nil {
		respond(w, nil, err, 0)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type userInput struct {
	CompanyID *uuid.UUID `json:"company_id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
}

func (a *API) createUser(w http.ResponseWriter, r *http.Request) {
	var in userInput
	if !decode(w, r, &in) || !validUser(w, in) {
		return
	}
	v, err := a.store.CreateUser(r.Context(), store.User{ID: uuid.New(), CompanyID: in.CompanyID, Name: strings.TrimSpace(in.Name), Email: strings.ToLower(strings.TrimSpace(in.Email))})
	respond(w, v, err, http.StatusCreated)
}
func (a *API) getUser(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	v, err := a.store.GetUser(r.Context(), id)
	respond(w, v, err, http.StatusOK)
}
func (a *API) listUsers(w http.ResponseWriter, r *http.Request) {
	limit, offset, ok := page(w, r)
	if !ok {
		return
	}
	v, err := a.store.ListUsers(r.Context(), limit, offset)
	respond(w, v, err, http.StatusOK)
}
func (a *API) updateUser(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var in userInput
	if !decode(w, r, &in) || !validUser(w, in) {
		return
	}
	v, err := a.store.UpdateUser(r.Context(), store.User{ID: id, CompanyID: in.CompanyID, Name: strings.TrimSpace(in.Name), Email: strings.ToLower(strings.TrimSpace(in.Email))})
	respond(w, v, err, http.StatusOK)
}
func (a *API) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := a.store.DeleteUser(r.Context(), id); err != nil {
		respond(w, nil, err, 0)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func validUser(w http.ResponseWriter, in userInput) bool {
	if !required(w, in.Name, "name") || !required(w, in.Email, "email") {
		return false
	}
	if !strings.Contains(in.Email, "@") {
		writeError(w, http.StatusBadRequest, "email must contain @")
		return false
	}
	return true
}

type inventoryInput struct {
	CompanyID  uuid.UUID `json:"company_id"`
	SKU        string    `json:"sku"`
	Name       string    `json:"name"`
	Quantity   int       `json:"quantity"`
	PriceCents int64     `json:"price_cents"`
}

func (a *API) createInventory(w http.ResponseWriter, r *http.Request) {
	var in inventoryInput
	if !decode(w, r, &in) || !validInventory(w, in) {
		return
	}
	v, err := a.store.CreateInventory(r.Context(), store.Inventory{ID: uuid.New(), CompanyID: in.CompanyID, SKU: strings.TrimSpace(in.SKU), Name: strings.TrimSpace(in.Name), Quantity: in.Quantity, PriceCents: in.PriceCents})
	respond(w, v, err, http.StatusCreated)
}
func (a *API) getInventory(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	v, err := a.store.GetInventory(r.Context(), id)
	respond(w, v, err, http.StatusOK)
}
func (a *API) listInventory(w http.ResponseWriter, r *http.Request) {
	limit, offset, ok := page(w, r)
	if !ok {
		return
	}
	var companyID *uuid.UUID
	if raw := r.URL.Query().Get("company_id"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid company_id")
			return
		}
		companyID = &id
	}
	v, err := a.store.ListInventory(r.Context(), limit, offset, companyID)
	respond(w, v, err, http.StatusOK)
}
func (a *API) updateInventory(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var in inventoryInput
	if !decode(w, r, &in) || !validInventory(w, in) {
		return
	}
	v, err := a.store.UpdateInventory(r.Context(), store.Inventory{ID: id, CompanyID: in.CompanyID, SKU: strings.TrimSpace(in.SKU), Name: strings.TrimSpace(in.Name), Quantity: in.Quantity, PriceCents: in.PriceCents})
	respond(w, v, err, http.StatusOK)
}
func (a *API) deleteInventory(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := a.store.DeleteInventory(r.Context(), id); err != nil {
		respond(w, nil, err, 0)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func validInventory(w http.ResponseWriter, in inventoryInput) bool {
	if in.CompanyID == uuid.Nil {
		writeError(w, http.StatusBadRequest, "company_id is required")
		return false
	}
	if in.Quantity < 0 || in.PriceCents < 0 {
		writeError(w, http.StatusBadRequest, "quantity and price_cents must be non-negative")
		return false
	}
	return required(w, in.SKU, "sku") && required(w, in.Name, "name")
}

type uuidStudyPage struct {
	Items      []store.UUIDStudyItem `json:"items"`
	Strategy   string                `json:"strategy"`
	NextCursor string                `json:"next_cursor,omitempty"`
	NextOffset *int                  `json:"next_offset,omitempty"`
}

func (a *API) listUUIDStudy(w http.ResponseWriter, r *http.Request) {
	version := store.UUIDStudyVersion(r.PathValue("version"))
	if version != store.UUIDStudyV4 && version != store.UUIDStudyV7 {
		writeError(w, http.StatusBadRequest, "version must be v4 or v7")
		return
	}
	limit, offset, ok := page(w, r)
	if !ok {
		return
	}
	strategy := r.URL.Query().Get("strategy")
	if strategy == "" {
		strategy = "keyset"
	}
	if strategy != "keyset" && strategy != "offset" {
		writeError(w, http.StatusBadRequest, "strategy must be keyset or offset")
		return
	}
	var after *uuid.UUID
	if raw := r.URL.Query().Get("after"); raw != "" {
		if strategy != "keyset" {
			writeError(w, http.StatusBadRequest, "after is only valid with keyset strategy")
			return
		}
		id, err := uuid.Parse(raw)
		if err != nil {
			writeError(w, http.StatusBadRequest, "after must be a UUID")
			return
		}
		after = &id
	}
	if strategy == "keyset" && offset != 0 {
		writeError(w, http.StatusBadRequest, "offset is only valid with offset strategy")
		return
	}

	// Fetch one extra row so continuation fields are emitted only when another
	// page exists.
	items, err := a.store.ListUUIDStudy(r.Context(), version, limit+1, offset, after)
	if err != nil {
		respond(w, nil, err, 0)
		return
	}
	result := uuidStudyPage{Items: items, Strategy: strategy}
	if len(items) > limit {
		result.Items = items[:limit]
		if strategy == "keyset" {
			result.NextCursor = result.Items[len(result.Items)-1].ID.String()
		} else {
			next := offset + limit
			result.NextOffset = &next
		}
	}
	writeJSON(w, http.StatusOK, result)
}

func decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return false
	}
	return true
}
func required(w http.ResponseWriter, v, name string) bool {
	if strings.TrimSpace(v) == "" {
		writeError(w, http.StatusBadRequest, name+" is required")
		return false
	}
	return true
}
func pathID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid UUID")
		return uuid.Nil, false
	}
	return id, true
}
func page(w http.ResponseWriter, r *http.Request) (int, int, bool) {
	limit, offset := 50, 0
	var err error
	if v := r.URL.Query().Get("limit"); v != "" {
		limit, err = strconv.Atoi(v)
	}
	if err != nil || limit < 1 || limit > 500 {
		writeError(w, http.StatusBadRequest, "limit must be between 1 and 500")
		return 0, 0, false
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		offset, err = strconv.Atoi(v)
	}
	if err != nil || offset < 0 {
		writeError(w, http.StatusBadRequest, "offset must be non-negative")
		return 0, 0, false
	}
	return limit, offset, true
}
func respond(w http.ResponseWriter, v any, err error, status int) {
	if err == nil {
		writeJSON(w, status, v)
		return
	}
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			writeError(w, http.StatusConflict, "resource already exists")
			return
		case "23503", "23514":
			writeError(w, http.StatusBadRequest, "constraint violation")
			return
		}
	}
	writeError(w, http.StatusInternalServerError, "internal error")
}
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
