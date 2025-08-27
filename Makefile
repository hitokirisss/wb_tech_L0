APP=order-service

.PHONY: run
run:
	go run ./cmd/$(APP)

.PHONY: build
build:
	go build -o bin/$(APP) ./cmd/$(APP)

.PHONY: tidy
tidy:
	go mod tidy
