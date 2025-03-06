ROOT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
BINARY_NAME := $(ROOT_DIR)logz
CMD_DIR := $(ROOT_DIR)cmd
INSTALL_SCRIPT=$(ROOT_DIR)scripts/install.sh
ARGS :=

# Alvo para build
build:
	go build -ldflags "-s -w -X main.version=$(git describe --tags) -X main.commit=$(git rev-parse HEAD) -X main.date=$(date +%Y-%m-%d)" -trimpath -o $(BINARY_NAME) ${CMD_DIR} &&\
    upx $(BINARY_NAME) --force-overwrite --lzma --no-progress --no-color

install:
	sh $(INSTALL_SCRIPT) install $(ARGS)

# Limpar o binário gerado
clean:
	rm -f $(BINARY_NAME)

# Alvo de ajuda
help:
	@echo "Available targets:"
	@echo "  make build      - Build the binary using install script"
	@echo "  make install    - Install the binary and configure environment"
	@echo "  make clean      - Clean up build artifacts"
	@echo "  make help       - Display this help message"
	@echo ""
	@echo "Usage with arguments:"
	@echo "  make install ARGS='--custom-arg value' - Pass custom arguments to the install script"
