.PHONY: dev dev-frontend stop lint lint-frontend test test-frontend typecheck e2e check

dev:
	docker compose up -d
	cd frontend && npm run dev

dev-frontend:
	cd frontend && npm run dev

stop:
	docker compose down

lint:
	docker compose exec backend go vet ./...

lint-frontend:
	cd frontend && npm run lint

test:
	docker compose exec backend go test ./...

test-frontend:
	cd frontend && npm run test

typecheck:
	cd frontend && npm run typecheck

e2e:
	cd e2e && npx playwright test

check: lint lint-frontend typecheck test
