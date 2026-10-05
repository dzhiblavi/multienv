BUILD_FLAGS=-ldflags="-s -w" -trimpath
INSTALL_DIR=~/.bin
MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
TUI_BIN := $(MAKEFILE_DIR)/../tui/bin

install-server:
	@go build $(BUILD_FLAGS) \
		-o $(INSTALL_DIR)/multienv-server \
		$(MAKEFILE_DIR)/../server/cmd \
		&& echo "install server $(shell hostname) OK"

install-%-agent:
	@go build $(BUILD_FLAGS) \
		-o $(INSTALL_DIR)/$*-agent \
		$(MAKEFILE_DIR)/../$*-agent/cmd \
		&& echo "install $*-agent $(shell hostname) OK"

link-%:
	@ln -sf \
		$(TUI_BIN)/$* $(INSTALL_DIR)/$* \
		&& echo "install tui/$* $(shell hostname) OK"

sync:
	@rsync -av --delete --exclude=.git --filter=':- .gitignore' \
		../ dev-vm:/home/dzhiblavi/tools/multienv >/dev/null

local:
	$(MAKE) install-codex-agent # install ~/.bin/codex-agent
	$(MAKE) install-tmux-agent # install ~/.bin/tmux-agent
	$(MAKE) install-server # install ~/.bin/multienv-server
	$(MAKE) link-multienv # install workspace entrypoint ~/.bin/multienv
	$(MAKE) link-multienv-tui # install workspace TUI ~/.bin/multienv-tui

remote:
	$(MAKE) install-codex-agent # install ~/.bin/codex-agent
	$(MAKE) install-tmux-agent # install ~/.bin/tmux-agent
	$(MAKE) link-multienv-tui # install workspace TUI ~/.bin/multienv-tui

remote-ssh: sync
	ssh dev-vm make -C /home/dzhiblavi/tools/multienv/deploy remote
