MAKEFILE_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
COMPOSE_FILE := $(MAKEFILE_DIR)/compose.yaml
PASSWORD_FILE := $(MAKEFILE_DIR)/password
ACL_FILE := $(MAKEFILE_DIR)/users.acl
COMPOSE := docker-compose -f "$(COMPOSE_FILE)"

.PHONY: credentials up down status logs

$(PASSWORD_FILE):
	@umask 077; openssl rand -hex 32 > "$@"

$(ACL_FILE): $(PASSWORD_FILE)
	@umask 077; \
		secret=$$(cat "$(PASSWORD_FILE)"); \
		hash=$$(printf '%s' "$$secret" | shasum -a 256 | awk '{print $$1}'); \
		printf 'user default off\nuser app on #%s ~* &* +@all\n' "$$hash" > "$@"

credentials: $(ACL_FILE)

up: $(ACL_FILE)
	@$(COMPOSE) up -d

down:
	@$(COMPOSE) down

status:
	@$(COMPOSE) ps

logs:
	@$(COMPOSE) logs -f
