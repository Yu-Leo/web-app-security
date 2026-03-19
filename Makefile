PROJECT_NAME = web-app-security
COMPOSE = docker compose -p $(PROJECT_NAME)

.PHONY: up down stop add-mock-data db-up db-stop db-down logs ps build

up:
	$(COMPOSE) up --build

build:
	$(COMPOSE) build

down:
	$(COMPOSE) down -v

stop:
	$(COMPOSE) stop

add-mock-data:
	BASE_URL=http://localhost:8001 bash ./scripts/seed-mock.sh

db-up:
	$(COMPOSE) up -d postgres
	$(COMPOSE) up migrations

db-stop:
	$(COMPOSE) stop postgres

db-down:
	$(COMPOSE) stop postgres
	$(COMPOSE) rm -f -v postgres
	@docker volume rm $(PROJECT_NAME)_postgres_data >/dev/null 2>&1 || true

logs:
	$(COMPOSE) logs -f --tail=200

ps:
	$(COMPOSE) ps
