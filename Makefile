SRC_DIR=./types

installproto:
	@go get -u github.com/golang/protobuf/protoc-gen-go

protoc:
	@protoc -I=$(SRC_DIR) --go_out=$(SRC_DIR) $(SRC_DIR)/types.proto

