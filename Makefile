BINARY_NAME := opa-llm-planner
BIN_DIR := bin
INSTALL_DIR := /usr/local/bin

.PHONY: build clean install uninstall

build:
	go build -o $(BIN_DIR)/$(BINARY_NAME) .

clean:
	rm -rf $(BIN_DIR)

install: build
	install -m 755 $(BIN_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)

uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)
