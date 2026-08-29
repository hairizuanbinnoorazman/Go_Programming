.PHONY: build test vet compose-up compose-down smoke load

build:
	GOWORK=off go build ./cmd/server ./cmd/loadgen

test:
	GOWORK=off go test ./...

vet:
	GOWORK=off go vet ./...

compose-up:
	docker compose up --build -d --wait

compose-down:
	docker compose down

smoke:
	./scripts/smoke.sh

load:
	GOWORK=off go run ./cmd/loadgen --duration 30s --workers 20 --rate 100
