-- These indexes match the unfiltered list endpoints' ORDER BY clauses. Without
-- them, PostgreSQL must scan and sort each growing table before applying LIMIT.
CREATE INDEX IF NOT EXISTS companies_created_at_idx ON companies (created_at DESC);
CREATE INDEX IF NOT EXISTS users_created_at_idx ON users (created_at DESC);
CREATE INDEX IF NOT EXISTS inventory_created_at_idx ON inventory (created_at DESC);
