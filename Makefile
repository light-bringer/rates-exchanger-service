.PHONY: build run lint clean

build:
	go build -o rates-api ./cmd/api-service/

run: compose-up
	go run ./cmd/api-service/

lint:
	golangci-lint run

clean: compose-down
	rm -f rates-api

compose-up:
	docker compose -f database-docker-compose.yml up  -d 

compose-down:
	docker compose -f database-docker-compose.yml down