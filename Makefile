SRC_DIR=./types
LOGIC_DIR=./examples/counter/logic

protoc:
	@protoc -I=$(SRC_DIR) --go_out=$(SRC_DIR) $(SRC_DIR)/types.proto

example:
	@protoc -I=$(LOGIC_DIR) --go_out=$(LOGIC_DIR) $(LOGIC_DIR)/msg.proto
