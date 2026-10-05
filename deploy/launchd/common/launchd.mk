.DEFAULT_GOAL := render

DOMAIN := gui/$(shell id -u)

LABEL = $(shell CONFIG_PATH="$(CONFIG_PATH)" bash "$(RENDER)" --get label)
STDOUT_PATH = $(shell CONFIG_PATH="$(CONFIG_PATH)" bash "$(RENDER)" --get stdout_path)
STDERR_PATH = $(shell CONFIG_PATH="$(CONFIG_PATH)" bash "$(RENDER)" --get stderr_path)
PLIST_SOURCE = /tmp/$(LABEL).plist
PLIST_INSTALLED = $(HOME)/Library/LaunchAgents/$(LABEL).plist
SERVICE = $(DOMAIN)/$(LABEL)

.PHONY: check-config render enable disable status logs

check-config:
	@CONFIG_PATH="$(CONFIG_PATH)" bash "$(RENDER)" --get label >/dev/null

render: check-config
	@CONFIG_PATH="$(CONFIG_PATH)" bash "$(RENDER)" >"$(PLIST_SOURCE)"
	@plutil -lint "$(PLIST_SOURCE)"

enable: render
	@mkdir -p "$(dir $(PLIST_INSTALLED))" \
		"$$(dirname "$(STDOUT_PATH)")" "$$(dirname "$(STDERR_PATH)")"
	@install -m 644 "$(PLIST_SOURCE)" "$(PLIST_INSTALLED)"
	@if launchctl print "$(SERVICE)" >/dev/null 2>&1; then \
		launchctl bootout "$(SERVICE)"; fi
	@launchctl enable "$(SERVICE)"
	@launchctl bootstrap "$(DOMAIN)" "$(PLIST_INSTALLED)"

disable: check-config
	@if launchctl print "$(SERVICE)" >/dev/null 2>&1; then \
		launchctl bootout "$(SERVICE)"; fi
	@rm -f "$(PLIST_INSTALLED)"

status: check-config
	@launchctl print "$(SERVICE)"

logs: check-config
	@tail -n 100 -f "$(STDOUT_PATH)" "$(STDERR_PATH)"
