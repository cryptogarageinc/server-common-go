.PHONY: setup install install-tools deps

setup: install install-tools deps
	@echo "setup done"

install:
	go mod download

install-tools: install
	$(shell cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %)

deps:
	go mod tidy

vendor:
	go mod vendor

unit-test:
	gotestsum -- -cover ./...

build-examples:
	go build -tags=examples -o ./bin/grpcserver ./api/examples/grpc/server.go
	go build -tags=examples -o ./bin/restserver ./api/examples/rest/server.go

help:
	@make2help $(MAKEFILE_LIST)
