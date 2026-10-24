BUILD_FLAGS=-ldflags="-s -w" -trimpath
INSTALL_DIR=~/.bin
MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
TUI_BIN := $(MAKEFILE_DIR)/../tui/bin

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

link-%:
	@ln -sf \
		$(TUI_BIN)/$* $(INSTALL_DIR)/$* \
		&& echo "install tui/$* $(shell hostname) OK"

sync:
	@rsync -av --delete --exclude=.git --filter=':- .gitignore' \
		../ dev-vm:/home/dzhiblavi/tools/multienv >/dev/null

local:
	$(MAKE) install-agent # install ~/.bin/multienv-agent
	$(MAKE) install-server # install ~/.bin/multienv-server
	$(MAKE) install-tui # install ~/.bin/multienv-tui

remote:
	$(MAKE) install-agent # install ~/.bin/multienv-agent
	$(MAKE) install-tui # install ~/.bin/multienv-tui

remote-ssh: sync
	ssh dev-vm make -C /home/dzhiblavi/tools/multienv/deploy remote
