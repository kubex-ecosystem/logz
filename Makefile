# Description: Makefile for building and installing a Go application
# Author: Rafael Mori
# Copyright (c) 2025 Rafael Mori
# License: MIT License

# This Makefile is used to build and install a Go application.
# It provides commands for building the binary, installing it, cleaning up build artifacts,
# and running tests. It also includes a help command to display usage information.
# The Makefile uses color codes for logging messages and provides a consistent interface
# for interacting with the application.

# Define the application name and root directory
define APP_NAME
$(shell echo $(basename $(CURDIR)) | tr '[:upper:]' '[:lower:]')
endef
define ROOT_DIR
$(dir $(abspath $(lastword $(MAKEFILE_LIST))))
endef
define BINARY_NAME
$(ROOT_DIR)$(APP_NAME)
endef

# Define internal dynamic commands variables
private CMD_DIR := $(ROOT_DIR)cmd
INSTALL_SCRIPT=$(call ROOT_DIR)support/scripts/install.sh
private ARGUMENTS := $(MAKECMDGOALS)

CMD_STR := $(strip $(firstword $(ARGUMENTS)))
CMD_STR := $(if $(CMD_STR),$(CMD_STR),$(MAKECMDGOALS))
ARGS := $(filter-out $(strip $(CMD_STR)), $(ARGUMENTS))

# Define the color codes
private COLOR_GREEN := \033[32m
private COLOR_YELLOW := \033[33m
private COLOR_RED := \033[31m
private COLOR_BLUE := \033[34m
private COLOR_RESET := \033[0m

# Logging Functions
log = @printf "%b%s%b %s\n" "$(COLOR_BLUE)" "[LOG]" "$(COLOR_RESET)" "$(1)"
log_info = @printf "%b%s%b %s\n" "$(COLOR_BLUE)" "[INFO]" "$(COLOR_RESET)" "$(1)"
log_success = @printf "%b%s%b %s\n" "$(COLOR_GREEN)" "[SUCCESS]" "$(COLOR_RESET)" "$(1)"
log_warning = @printf "%b%s%b %s\n" "$(COLOR_YELLOW)" "[WARNING]" "$(COLOR_RESET)" "$(1)"
log_break =	 @printf "%b%s%b\n" "$(COLOR_BLUE)" "[INFO]" "$(COLOR_RESET)"
log_error = @printf "%b%s%b %s\n" "$(COLOR_RED)" "[ERROR]" "$(COLOR_RESET)" "$(1)"

# Run dynamic commands with arguments calling the install script.
%:
	@:
	@if [ -z "$(CMD_STR)" ]; then $(call log_error, No command specified. Use make help for usage information.); fi
	@if [[ "$(CMD_STR)" != "build"  &&  "$(CMD_STR)" != "build-dev"  &&  "$(CMD_STR)" != "install"  &&  "$(CMD_STR)" != "clean"  &&  "$(CMD_STR)" != "test" ]]; then $(call log_error, Invalid command: $(CMD_STR). Use make help for usage information.); fi
	@bash $(INSTALL_SCRIPT) $(CMD_STR) $(ARGS)
	@exit $?

# Build the binary using the install script.
build:
	@bash $(INSTALL_SCRIPT) $(CMD_STR) $(ARGS)
	@exit $?

# Install the binary and configure the environment.
install:
	@bash $(INSTALL_SCRIPT) $(CMD_STR) $(ARGS)
	@exit $?

# Clean up build artifacts.
clean, clear:
	@bash $(INSTALL_SCRIPT) $(CMD_STR)
	@exit $?

# Run tests.
test:
	@bash $(INSTALL_SCRIPT) $(CMD_STR)
	@exit $?

# Display help message.
help:
	$(call log, $(APP_NAME) Makefile )
	$(call break, b )
	$(call log, Usage: )
	$(call log,   make [target] [ARGS='--custom-arg value'] )
	$(call break, b )
	$(call log, Available targets: )
	$(call log,   make build      - Build the binary using install script)
	$(call log,   make build-dev  - Build the binary without compressing it)
	$(call log,   make install    - Install the binary and configure environment)
	$(call log,   make clean      - Clean up build artifacts)
	$(call log,   make test       - Run tests)
	$(call log,   make help       - Display this help message)
	$(call break, b )
	$(call log, Usage with arguments: )
	$(call log,   make install ARGS='--custom-arg value' - Pass custom arguments to the install script)
	$(call break, b )
	$(call log, Example: )
	$(call log,   make install ARGS='--prefix /usr/local')
	$(call break, b )
	$(call log, $(APP_NAME) is a tool for managing Kubernetes resources)
	$(call break, b )
	$(call log, For more information, visit: )
	$(call log, 'https://github.com/faelmori/'$(APP_NAME))
	$(call break, b )
	$(call success, End of help message)

.PHONY: build install clean clear test help
.DEFAULT_GOAL := help

# End of Makefile