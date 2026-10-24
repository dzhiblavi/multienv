BUILD_FLAGS=-ldflags="-s -w" -trimpath
INSTALL_DIR=~/.bin
MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
LAUNCHER_BIN := $(MAKEFILE_DIR)/../launcher/bin
REMOTE_HOST ?= dev-vm

install-server:
	@go build $(BUILD_FLAGS) \
		-o $(INSTALL_DIR)/multienv-server \
		$(MAKEFILE_DIR)/../server/cmd \
		&& echo "install server $(shell hostname) OK"

install-agent:
	@go build $(BUILD_FLAGS) \
		-o $(INSTALL_DIR)/multienv-agent \
		$(MAKEFILE_DIR)/../agent/cmd \
		&& echo "install multienv-agent $(shell hostname) OK"

install-tui:
	@go build $(BUILD_FLAGS) \
		-o $(INSTALL_DIR)/multienv-tui \
		$(MAKEFILE_DIR)/../tui/cmd \
		&& echo "install multienv-tui $(shell hostname) OK"

install-launcher:
	@ln -sf \
		$(LAUNCHER_BIN)/multienv $(INSTALL_DIR)/multienv \
		&& echo "install multienv $(shell hostname) OK"

local:
	$(MAKE) install-agent # install ~/.bin/multienv-agent
	$(MAKE) install-server # install ~/.bin/multienv-server
	$(MAKE) install-launcher # install ~/.bin/multienv
	$(MAKE) install-tui # install ~/.bin/multienv-tui

remote:
	$(MAKE) install-agent # install ~/.bin/multienv-agent
	$(MAKE) install-tui # install ~/.bin/multienv-tui
