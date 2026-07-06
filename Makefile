BIN_DIR := bin
APP_BIN := $(BIN_DIR)/auth

PROTO_DIR := proto
GEN_DIR := gen/go
PROTO_FILE := $(PROTO_DIR)/auth/v1/auth.proto

.PHONY: build
build:
	rm -rf $(BIN_DIR)
	mkdir -p $(BIN_DIR)
	go build -o $(APP_BIN) ./cmd/auth

.PHONY: proto
proto:
	rm -rf $(GEN_DIR)
	mkdir -p $(GEN_DIR)
	protoc -I $(PROTO_DIR) $(PROTO_FILE) \
		--go_out=$(GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative
