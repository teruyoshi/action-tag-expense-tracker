.PHONY: dev dev-frontend stop lint test e2e

dev:
	docker compose up -d
	cd frontend && npm run dev

dev-frontend:
	cd frontend && npm run dev

stop:
	docker compose down

lint:
	docker compose exec backend go vet ./...

test:
	docker compose exec backend go test ./...

e2e:
	cd e2e && npx playwright test
