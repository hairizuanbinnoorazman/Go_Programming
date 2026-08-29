#!/usr/bin/env bash
set -euo pipefail

base_url=${BASE_URL:-http://127.0.0.1:8080}
company=$(curl --fail --silent --show-error -H 'Content-Type: application/json' -d '{"name":"Smoke Company"}' "$base_url/v1/companies")
company_id=$(printf '%s' "$company" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
test -n "$company_id"
curl --fail --silent --show-error -H 'Content-Type: application/json' -d "{\"company_id\":\"$company_id\",\"name\":\"Smoke User\",\"email\":\"smoke-$company_id@example.test\"}" "$base_url/v1/users" >/dev/null
curl --fail --silent --show-error -H 'Content-Type: application/json' -d "{\"company_id\":\"$company_id\",\"sku\":\"SMOKE-$company_id\",\"name\":\"Smoke Item\",\"quantity\":10,\"price_cents\":1250}" "$base_url/v1/inventory" >/dev/null
curl --fail --silent --show-error "$base_url/v1/companies/$company_id" >/dev/null
curl --fail --silent --show-error "$base_url/v1/uuid-study/v4?limit=1" >/dev/null
curl --fail --silent --show-error "$base_url/v1/uuid-study/v7?strategy=offset&limit=1&offset=0" >/dev/null
echo "CRUD smoke test passed for company $company_id"
