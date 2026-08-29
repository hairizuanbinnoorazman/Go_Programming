package config

import "testing"

func TestLoad(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("DB_MAX_CONNS", "12")
	t.Setenv("DB_MIN_CONNS", "3")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.MaxConns != 12 || c.MinConns != 3 {
		t.Fatalf("unexpected pool settings: %+v", c)
	}
}

func TestLoadRejectsInvalidPool(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("DB_MAX_CONNS", "2")
	t.Setenv("DB_MIN_CONNS", "3")
	if _, err := Load(); err == nil {
		t.Fatal("expected an error")
	}
}
