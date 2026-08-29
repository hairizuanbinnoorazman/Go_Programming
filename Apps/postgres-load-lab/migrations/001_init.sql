CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

CREATE TABLE IF NOT EXISTS companies (
    id uuid PRIMARY KEY,
    name text NOT NULL CHECK (length(name) BETWEEN 1 AND 200),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    company_id uuid REFERENCES companies(id) ON DELETE SET NULL,
    name text NOT NULL CHECK (length(name) BETWEEN 1 AND 200),
    email text NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS inventory (
    id uuid PRIMARY KEY,
    company_id uuid NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    sku text NOT NULL,
    name text NOT NULL CHECK (length(name) BETWEEN 1 AND 200),
    quantity integer NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    price_cents bigint NOT NULL DEFAULT 0 CHECK (price_cents >= 0),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (company_id, sku)
);

CREATE INDEX IF NOT EXISTS users_company_id_idx ON users (company_id);
CREATE INDEX IF NOT EXISTS inventory_company_id_idx ON inventory (company_id);
CREATE INDEX IF NOT EXISTS inventory_company_updated_idx ON inventory (company_id, updated_at DESC);
