MAKEFILE_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
CONFIG_PATH ?= $(HOME)/.config/multienv/forward.launchd.yaml
RENDER := $(MAKEFILE_DIR)/render.sh

include $(MAKEFILE_DIR)/../common/launchd.mk
