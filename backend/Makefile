.PHONY: run build test tidy clean docker-build

APP_NAME = anr-platform-backend
BIN_DIR = bin

run:
	go run cmd/server/main.go

build:
	mkdir -p $(BIN_DIR)
	go build -ldflags="-w -s" -o $(BIN_DIR)/server cmd/server/main.go

test:
	go test -v ./...

tidy:
	go mod tidy

clean:
	rm -rf $(BIN_DIR)

docker-build:
	docker build -t $(APP_NAME):latest .
