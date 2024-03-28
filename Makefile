.PHONY: build run lint clean

build:
	go build -o rates-api .

run:
	go run main.go

lint:
	golangci-lint run

clean:
	rm -f rates-api