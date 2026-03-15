.PHONY: dev stop build fmt lint lint-frontend test test-frontend typecheck e2e quick-check check verify test-diff doctor migrate-up migrate-down migrate-create

# ---------- DEV ----------

dev:
	docker compose up -d

stop:
	docker compose down

# ---------- BUILD ----------

build:
	docker compose build

# ---------- FORMAT ----------

fmt:
	docker compose exec backend gofmt -w .

# ---------- LINT ----------

lint:
	docker compose exec backend go vet ./...

lint-frontend:
	docker compose exec frontend npm run lint

# ---------- TYPE ----------

typecheck:
	docker compose exec frontend npm run typecheck

# ---------- TEST ----------

test:
	docker compose exec backend go test ./...

test-frontend:
	docker compose exec frontend npm run test

test-diff:
	bash scripts/test-diff.sh

# ---------- E2E ----------

e2e:
	docker compose run --rm e2e

# ---------- AI CHECK ----------

quick-check: lint lint-frontend typecheck

check: quick-check test test-frontend

verify: check e2e

# ---------- MIGRATE ----------

migrate-up:
	docker compose exec backend migrate -path /app/migrations -database "mysql://root:root@tcp(db:3306)/expense_tracker" up

migrate-down:
	docker compose exec backend migrate -path /app/migrations -database "mysql://root:root@tcp(db:3306)/expense_tracker" down 1

migrate-create:
	@read -p "Migration name: " name; \
	docker compose exec backend migrate create -ext sql -dir /app/migrations -seq $$name

# ---------- DOCTOR ----------

doctor:
	bash scripts/doctor.sh
