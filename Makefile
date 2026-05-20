COMPOSE_FILE := dev/docker-compose.yml

.PHONY: stop
stop:
	docker-compose -f $(COMPOSE_FILE) down

.PHONY: start
start:
	docker-compose -f $(COMPOSE_FILE) up -d --build

.PHONY: support
support:
	docker-compose -f $(COMPOSE_FILE) up -d db

.PHONY: build
build:
	docker build -t pack-calculator -f dev/Dockerfile .

.PHONY: test
test:
	go test -v ./... --race

.PHONY: swagger
swagger:
	GOPROXY=https://proxy.golang.org,direct go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g main.go -o ./docs -d ./cmd,./internal/common --parseInternal
