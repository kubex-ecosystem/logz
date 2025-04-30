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
private APP_NAME := $(shell echo $(basename $(CURDIR)) | tr '[:upper:]' '[:lower:]')
private ROOT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
private BINARY_NAME := $(ROOT_DIR)$(APP_NAME)
private CMD_DIR := $(ROOT_DIR)cmd

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

ARGUMENTS := $(MAKECMDGOALS)
INSTALL_SCRIPT=$(ROOT_DIR)support/scripts/install.sh
CMD_STR := $(strip $(firstword $(ARGUMENTS)))
ARGS := $(filter-out $(strip $(CMD_STR)), $(ARGUMENTS))

# Build the binary using the install script.
build:
	$(call log_info, Building $(APP_NAME) binary)
	$(call log_info, Args: $(ARGS))
	@#$(INSTALL_SCRIPT) clean $(ARGS) 2>&1 >/dev/null || exit 1
	@$(INSTALL_SCRIPT) build $(ARGS)
    $(shell exit 0)

# Install the binary and configure the environment.
install:
	$(call log_info, Installing $(APP_NAME) binary)
	$(call log_info, Args: $(ARGS))
	@#bash $(INSTALL_SCRIPT) clean $(ARGS) 2>&1 >/dev/null || exit 1
	bash $(INSTALL_SCRIPT) install $(ARGS)
	$(shell exit 0)

# Clean up build artifacts.
clean:
	$(call log_info, Cleaning up build artifacts)
	$(call log_info, Args: $(ARGS))
	@bash $(INSTALL_SCRIPT) clean $(ARGS)
	$(shell exit 0)

# Run tests.
test:
	$(call log_info, Running tests)
	$(call log_info, Args: $(ARGS))
	@bash $(INSTALL_SCRIPT) test $(ARGS)
	$(shell exit 0)

## Run dynamic commands with arguments calling the install script.
%:
	@:
	$(call log_info, Running command: $(CMD_STR))
	$(call log_info, Args: $(ARGS))
	@bash $(INSTALL_SCRIPT) $(CMD_STR) $(ARGS)
	$(shell exit 0)

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
	$(shell exit 0)


# End of Makefile
