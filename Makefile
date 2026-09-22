.PHONY: up down logs logs-app build E2E-test UNIT-test clean help

DOCKER = docker compose

up:
	$(DOCKER) up -d

down:
	$(DOCKER) down

logs:
	$(DOCKER) logs -f 

logs-app:
	$(DOCKER) logs -f app

build:
	docker build -t order-api .
	
E2E-test:
	go test ./cmd/...

UNIT-test:
	go test -v ./internal/order/...

clean:
	$(DOCKER) down -v

help:
	@echo "up        - поднять всё"
	@echo "down      - остановить"
	@echo "logs      - все логи"
	@echo "logs-app  - логи приложения"
	@echo "build     - собрать образ"
	@echo "E2E-test  - запустить e2e-тесты"
	@echo "UNIT-test - запустить unit-тесты"
	@echo "clean     - полная очистка"