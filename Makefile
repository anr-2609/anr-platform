.PHONY: help up down logs ps backend-run backend-build backend-test

help:
	@echo "ANR Platform Makefile Commands:"
	@echo "  make up            - Khởi động Docker Compose stack (backend, postgres, redis)"
	@echo "  make down          - Dừng Docker Compose stack"
	@echo "  make logs          - Xem logs thời gian thực của Docker Compose"
	@echo "  make ps            - Kiểm tra trạng thái containers"
	@echo "  make backend-run   - Chạy Go backend trực tiếp trên máy"
	@echo "  make backend-build - Build binary Go backend"
	@echo "  make backend-test  - Chạy unit tests cho Go backend"

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

ps:
	docker compose ps

backend-run:
	cd backend && go run cmd/server/main.go

backend-build:
	cd backend && go build -o bin/server cmd/server/main.go

backend-test:
	cd backend && go test -v ./...
