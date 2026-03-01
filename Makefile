BINARY_NAME := opa-llm-planner
BIN_DIR := bin

.PHONY: build clean

build:
	go build -o $(BIN_DIR)/$(BINARY_NAME) .

clean:
	rm -rf $(BIN_DIR)
