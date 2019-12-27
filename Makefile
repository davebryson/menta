TYPES_SRC_DIR=./types
STORE_SRC_DIR=./store
COUNTER_SRC_DIR=./examples/services/counter

installproto:
	@go get -u github.com/golang/protobuf/protoc-gen-go

protoc:
	@protoc -I=$(TYPES_SRC_DIR) --go_out=$(TYPES_SRC_DIR) $(TYPES_SRC_DIR)/types.proto
	@protoc -I=$(STORE_SRC_DIR) --go_out=$(STORE_SRC_DIR) $(STORE_SRC_DIR)/data.proto
	@protoc -I=$(COUNTER_SRC_DIR) --go_out=$(COUNTER_SRC_DIR) $(COUNTER_SRC_DIR)/types.proto


