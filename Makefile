SRC_DIR=./types

protoc:
	@protoc -I=$(SRC_DIR) --go_out=$(SRC_DIR) $(SRC_DIR)/types.proto

ensure_deps:
	@rm -rf vendor/
	@echo "--> Running dep"
	@dep ensure
