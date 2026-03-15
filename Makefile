.PHONY: dev dev-frontend stop build fmt lint lint-frontend test test-frontend typecheck e2e quick-check check verify test-diff doctor migrate-up migrate-down migrate-create

# ---------- DEV ----------

dev:
	docker compose up -d
	cd frontend && npm run dev

dev-frontend:
	cd frontend && npm run dev

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
	cd frontend && npm run lint

# ---------- TYPE ----------

typecheck:
	cd frontend && npm run typecheck

# ---------- TEST ----------

test:
	docker compose exec backend go test ./...

test-frontend:
	cd frontend && npm run test

test-diff:
	bash scripts/test-diff.sh

# ---------- E2E ----------

e2e:
	cd e2e && npx playwright test

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
