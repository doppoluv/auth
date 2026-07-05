PROTO_DIR := proto
GEN_DIR := gen/go
PROTO_FILE := $(PROTO_DIR)/auth/auth.proto

.PHONY: proto
proto:
	mkdir -p $(GEN_DIR)
	protoc -I $(PROTO_DIR) $(PROTO_FILE) \
		--go_out=$(GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative
