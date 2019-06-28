SRC_DIR=./types
ACCOUNT_DIR=./x/accounts

installproto:
	@go get -u github.com/golang/protobuf/protoc-gen-go

protoc:
	@protoc -I=$(SRC_DIR) --go_out=$(SRC_DIR) $(SRC_DIR)/types.proto

account:
	@protoc -I=$(ACCOUNT_DIR) --go_out=$(ACCOUNT_DIR) $(ACCOUNT_DIR)/types.proto

