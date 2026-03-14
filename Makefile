.PHONY: dev dev-frontend stop lint lint-frontend test test-frontend typecheck e2e quick-check check verify test-diff doctor

# ---------- DEV ----------

dev:
	docker compose up -d
	cd frontend && npm run dev

dev-frontend:
	cd frontend && npm run dev

stop:
	docker compose down

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

# ---------- DOCTOR ----------

doctor:
	bash scripts/doctor.sh
