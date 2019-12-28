TYPES_SRC_DIR=./types
COUNTER_SRC_DIR=./examples/services/counter
NS_SRC_DIR=./storage

installproto:
	@go get -u github.com/golang/protobuf/protoc-gen-go

protoc:
	@protoc -I=$(TYPES_SRC_DIR) --go_out=$(TYPES_SRC_DIR) $(TYPES_SRC_DIR)/types.proto
	@protoc -I=$(COUNTER_SRC_DIR) --go_out=$(COUNTER_SRC_DIR) $(COUNTER_SRC_DIR)/types.proto
	@protoc -I=$(NS_SRC_DIR) --go_out=$(NS_SRC_DIR) $(NS_SRC_DIR)/data.proto



