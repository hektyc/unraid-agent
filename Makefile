SHELL := /bin/bash
PLUGIN_NAME := unraid-mcp
VERSION := 0.0.1
PLG_FILE := $(PLUGIN_NAME)-$(VERSION).plg
BUILD_DIR := /tmp/$(PLUGIN_NAME)-build
INSTALL_DIR := $(BUILD_DIR)/usr/local/emhttp/plugins/$(PLUGIN_NAME)
CONFIG_DIR := $(BUILD_DIR)/boot/config/plugins/$(PLUGIN_NAME)

.PHONY: all build clean package

all: package

build:
	@echo "Building Go binary..."
	cd go && go build -o bin/unraid-mcp cmd/unraid-mcp/main.go
	@echo "Build complete."

package: build
	@echo "Packaging plugin..."
	rm -rf "$(BUILD_DIR)"
	mkdir -p "$(INSTALL_DIR)"
	mkdir -p "$(CONFIG_DIR)"
	cp -r plugin/* "$(INSTALL_DIR)/"
	cp -r go "$(INSTALL_DIR)/"
	cp -r memory "$(INSTALL_DIR)/"
	mkdir -p "$(INSTALL_DIR)/bin"
	cp go/bin/unraid-mcp "$(INSTALL_DIR)/bin/unraid-mcp" 2>/dev/null || echo "Binary not found. Build first."
	chmod +x "$(INSTALL_DIR)/scripts/"*.sh
	tar -C "$(BUILD_DIR)" -cf - . | xz -z > "$(PLG_FILE)"
	rm -rf "$(BUILD_DIR)"
	@echo "Package created: $(PLG_FILE)"

clean:
	rm -rf "$(BUILD_DIR)"
	rm -f "$(PLUGIN_NAME)-*.plg"
	rm -rf go/bin/
