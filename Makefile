.PHONY: dev stop build fmt fmt-frontend fmt-check fmt-check-frontend lint lint-frontend test test-frontend typecheck e2e quick-check check verify test-diff doctor migrate-up migrate-down migrate-create ci-up

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

fmt-frontend:
	docker compose exec frontend npm run format

fmt-check:
	docker compose exec backend gofmt -l . | grep -q . && { echo "gofmt: unformatted files found"; docker compose exec backend gofmt -l .; exit 1; } || true

fmt-check-frontend:
	docker compose exec frontend npm run format:check

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

# ---------- CI ----------

ci-up:
	docker compose up -d --wait

# ---------- DOCTOR ----------

doctor:
	bash scripts/doctor.sh
