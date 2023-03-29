test:
	go test -count=1 ./...

lint:
	golangci-lint run

PROTO_DIR := ./internal/server/transport/proto

proto:
	protoc --go_out=$(PROTO_DIR) --go-grpc_out=require_unimplemented_servers=false:$(PROTO_DIR) ./api/*.proto
	go mod tidy
