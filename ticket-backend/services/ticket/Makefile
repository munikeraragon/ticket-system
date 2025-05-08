PROTO_DIR=api
PROTO_FILE=$(PROTO_DIR)/ticket.proto
GOOGLEAPIS=../../googleapis
GEN_DIR=gen/go

proto:
	mkdir -p $(GEN_DIR)
	protoc -I $(PROTO_DIR) -I $(GOOGLEAPIS) \
		--go_out=$(GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(GEN_DIR) --grpc-gateway_opt=paths=source_relative \
		$(PROTO_FILE)

