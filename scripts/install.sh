#!/usr/bin/env bash

# This variable are used to customize the script behavior, like repository url and owner
_OWNER="faelmori"

# This function is used to get the release URL for the binary.
# It can be customized to change the URL format or add additional parameters.
# Actually im using the default logic to construct the URL with the release version, the platform and the architecture
# with the format .tar.gz or .zip (for windows). Sweet yourself.
get_release_url() {
    # Default logic for constructing the release URL
    local os="${_PLATFORM%%-*}"
    local arch="${_PLATFORM##*-}"
    # If os is windows, set the format to .zip, otherwise .tar.gz
    local format="${os:zip=tar.gz}"

    echo "https://github.com/${_OWNER}/${_PROJECT_NAME}/releases/download/${_VERSION}/${_PROJECT_NAME}_.${format}"
}

# The _REPO_ROOT variable is set to the root directory of the repository. One above the script directory.
_REPO_ROOT="$(dirname "$(dirname "$(realpath "$0")")")"

# The _APP_NAME variable is set to the name of the repository. It is used to identify the application.
_APP_NAME="$(basename "$_REPO_ROOT")"

# The _PROJECT_NAME variable is set to the name of the project. It is used for display purposes.
_PROJECT_NAME="$_APP_NAME"

# The _VERSION variable is set to the version of the project. It is used for display purposes.
_VERSION=$(cat "$_REPO_ROOT/version/CLI_VERSION" 2>/dev/null || echo "v0.0.0")

# The _VERSION variable is set to the version of the project. It is used for display purposes.
_LICENSE="MIT"

# The _ABOUT variable contains information about the script and its usage.
_ABOUT="'
################################################################################
  This Script is used to install ${_PROJECT_NAME} project, version ${_VERSION}.

  Supported OS: Linux, macOS ---> Windows(not supported)
  Supported Architecture: amd64, arm64
  Source: https://github.com/${_OWNER}/${_PROJECT_NAME}
  Binary Release: https://github.com/${_OWNER}/${_PROJECT_NAME}/releases/latest
  License: ${_LICENSE}
  Notes:
    - [version] is optional; if omitted, the latest version will be used.
    - If the script is run locally, it will try to resolve the version from the
      repo tags if no version is provided.
    - The script will install the binary in the ~/.local/bin directory if the
      user is not root. Otherwise, it will install in /usr/local/bin.
    - The script will add the installation directory to the PATH in the shell
      configuration file.
    - The script will also install UPX if it is not already installed.
    - The script will build the binary if the build option is provided.
    - The script will download the binary from the release URL if the install
      option is provided.
    - The script will clean up build artifacts if the clean option is provided.
    - The script will check if the required dependencies are installed.
    - The script will validate the Go version before building the binary.
    - The script will check if the installation directory is in the PATH.
    - The script will print a summary of the installation.
################################################################################
'"

# Variable to store the current running shell
_CURRENT_SHELL=""

# The _CMD_PATH variable is set to the path of the cmd directory. It is used to
# identify the location of the main application code.
_CMD_PATH="$(dirname "$(realpath "$(dirname "$0")")")/cmd"

# The _BUILD_PATH variable is set to the path of the build directory. It is used
# to identify the location of the build artifacts.
_BUILD_PATH="$(dirname "${_CMD_PATH}")"

# The _BINARY variable is set to the path of the binary file. It is used to
# identify the location of the binary file.
_BINARY="${_BUILD_PATH}/${_APP_NAME}"

# The _LOCAL_BIN variable is set to the path of the local bin directory. It is
# used to identify the location of the local bin directory.
_LOCAL_BIN="${HOME:-"~"}/.local/bin"

# The _GLOBAL_BIN variable is set to the path of the global bin directory. It is
# used to identify the location of the global bin directory.
_GLOBAL_BIN="/usr/local/bin"

# Color codes for logging
_SUCCESS="\033[0;32m"
_WARN="\033[0;33m"
_ERROR="\033[0;31m"
_INFO="\033[0;36m"
_NC="\033[0m"

# The _PLATFORM variable is set to the platform name. It is used to identify the
# platform on which the script is running.
_PLATFORM=""

# Create a temporary directory for script cache
_TEMP_DIR="$(mktemp -d)"

# Diretório temporário para baixar o arquivo
if test -d "${_TEMP_DIR}"; then
    log "info" "Temporary directory created: ${_TEMP_DIR}"
else
    log "error" "Failed to create temporary directory."
    return 1
fi

# Function to clear the script cache
clear_script_cache() {
  # Disable the trap for cleanup
  trap - EXIT HUP INT QUIT ABRT ALRM TERM

  # Check if the temporary directory exists, if not, return
  if ! test -d "${_TEMP_DIR}"; then
    return 0
  fi

  # Remove the temporary directory
  rm -rf "${_TEMP_DIR}" || true
  if test -d "${_TEMP_DIR}"; then
    sudo rm -rf "${_TEMP_DIR}"
    if test -d "${_TEMP_DIR}"; then
      log "error" "Failed to remove temporary directory: ${_TEMP_DIR}"
      return 1
    else
      log "success" "Temporary directory removed successfully."
    fi
  fi

  return 0
}

# Function to get the current shell
get_current_shell() {
  _CURRENT_SHELL="$(cat /proc/$$/comm)"

  case "${0##*/}" in
    ${_CURRENT_SHELL}*)
      shebang="$(head -1 "${0}")"
      _CURRENT_SHELL="${shebang##*/}"
      ;;
  esac

  echo "${_CURRENT_SHELL}"

  return 0
}

# Set a trap to clean up the temporary directory on exit
set_trap(){
  # Get the current shell
  get_current_shell

  # Set the trap for the current shell and enable error handling, if applicable
  case "${_CURRENT_SHELL}" in
    *ksh|*zsh|*bash)

      # Collect all arguments passed to the script into an array without modifying or running them
      declare -a _FULL_SCRIPT_ARGS="$@"

      # Check if the script is being run in debug mode, if so, enable debug mode on the script output
      if [[ ${_FULL_SCRIPT_ARGS[*]} =~ ^.*-d.*$ ]]; then
          set -x
      fi

      # Set for the current shell error handling and some other options
      if test "${_CURRENT_SHELL}" = "bash"; then
        set -o errexit
        set -o pipefail
        set -o errtrace
        set -o functrace
        shopt -s inherit_errexit
      fi

      # Set the trap to clear the script cache on exit.
      # It will handle the following situations: command line exit, hangup, interrupt, quit, abort, alarm, and termination.
      trap 'clear_script_cache' EXIT HUP INT QUIT ABRT ALRM TERM
      ;;
  esac

  return 0
}

# Call the set_trap function to set up the trap
set_trap "$@"

# Clear the screen. If the script gets here, it means the script passed the
# initial checks and the temporary directory was created successfully.
clear

# Log messages with different levels
# Arguments:
#   $1 - log level (info, warn, error, success)
#   $2 - message to log
log() {
  local type=
  type=${1:-info}
  local message=
  message=${2:-}

  # With colors
  case $type in
    info|_INFO|-i|-I)
      printf '%b[_INFO]%b ℹ️  %s\n' "$_INFO" "$_NC" "$message"
      ;;
    warn|_WARN|-w|-W)
      printf '%b[_WARN]%b ⚠️  %s\n' "$_WARN" "$_NC" "$message"
      ;;
    error|_ERROR|-e|-E)
      printf '%b[_ERROR]%b ❌  %s\n' "$_ERROR" "$_NC" "$message"
      ;;
    success|_SUCCESS|-s|-S)
      printf '%b[_SUCCESS]%b ✅  %s\n' "$_SUCCESS" "$_NC" "$message"
      ;;
    *)
      log "info" "$message"
      ;;
  esac
}

# Detect the platform
what_platform() {
  local _os=""
  _os="$(uname -s)"

  local _arch=""
  _arch="$(uname -m)"

  case "${_os}" in
  "Linux")
    case "${_arch}" in
    "x86_64")
      arch=amd64
      ;;
    "armv6")
      arch=armv6l
      ;;
    "armv8" | "aarch64")
      arch=arm64
      ;;
    .*386.*)
      arch=386
      ;;
    esac
    platform="linux-${arch}"
    ;;
  "Darwin")
    case "${_arch}" in
    "x86_64")
      arch=amd64
      ;;
    "arm64")
      arch=arm64
      ;;
    esac
    platform="darwin-${_arch}"
    ;;
  "MINGW" | "MSYS" | "CYGWIN")
    case "${_arch}" in
    "x86_64")
      arch=amd64
      ;;
    "arm64")
      arch=arm64
      ;;
    esac
    platform="windows-${arch}"
    ;;
  esac

  if [ -z "${platform}" ]; then
    log "error" "Unsupported platform: ${_os} ${_arch}"
    log "error" "Please report this issue to the project maintainers."
    return 1
  fi

  _PLATFORM="${platform}"

  echo "${platform}"
  return 0
}

# Detect the shell configuration file
# Returns:
#   Shell configuration file path
detect_shell_rc() {
    shell_rc_file=""
    user_shell=$(basename "$SHELL")
    case "$user_shell" in
        bash) shell_rc_file="$HOME/.bashrc" ;;
        zsh) shell_rc_file="$HOME/.zshrc" ;;
        sh) shell_rc_file="$HOME/.profile" ;;
        fish) shell_rc_file="$HOME/.config/fish/config.fish" ;;
        *)
            log "warn" "Unsupported shell, modify PATH manually."
            return 1
            ;;
    esac
    log "info" "$shell_rc_file"
}

# Add a directory to the PATH in the shell configuration file
# Arguments:
#   $1 - target path to add to PATH
add_to_path() {
    target_path="$1"
    shell_rc_file=$(detect_shell_rc)
    if [ -z "$shell_rc_file" ]; then
        log "error" "Could not determine shell configuration file."
        return 1
    fi

    if grep -q "export PATH=.*$target_path" "$shell_rc_file" 2>/dev/null; then
        log "success" "$target_path is already in $shell_rc_file."
        return 0
    fi

    echo "export PATH=$target_path:\$PATH" >> "$shell_rc_file"
    log "success" "Added $target_path to PATH in $shell_rc_file."
    log "success" "Run 'source $shell_rc_file' to apply changes."
}

# Clean up build artifacts
clean() {
    log "info" "Cleaning up build artifacts..."
    rm -f "$_BINARY" || true
    log "success" "Cleaned up build artifacts."
}

# Install the binary to the appropriate directory
install_binary() {
    if [ "$(id -u)" -ne 0 ]; then
        log "info" "You are not root. Installing in $_LOCAL_BIN..."
        mkdir -p "$_LOCAL_BIN"
        cp "$_BINARY" "$_LOCAL_BIN/$_APP_NAME" || exit 1
        add_to_path "$_LOCAL_BIN"
    else
        log "info" "Root detected. Installing in $_GLOBAL_BIN..."
        cp "$_BINARY" "$_GLOBAL_BIN/$_APP_NAME" || exit 1
        add_to_path "$_GLOBAL_BIN"
    fi
    clean
}

# Install UPX if it is not already installed
install_upx() {
    if ! command -v upx > /dev/null; then
        log "info" "Installing UPX..."
        if [ "$(uname)" = "Darwin" ]; then
            brew install upx
        elif command -v apt-get > /dev/null; then
            sudo apt-get install -y upx
        else
            log "error" 'Install UPX manually from https://upx.github.io/'
            exit 1
        fi
    else
        log "success" ' UPX is already installed.'
    fi
}

# Check if the required dependencies are installed
# Arguments:
#   $@ - list of dependencies to check
check_dependencies() {
    for dep in "$@"; do
        if ! command -v "$dep" > /dev/null; then
            log "error" "$dep is not installed."
            exit 1
        else
            log "success" "$dep is installed."
        fi
    done
}

# Build the binary
build_binary() {
    log "info" "Building the binary..."
    go build -ldflags "-s -w -X main.version=$(git describe --tags) -X main.commit=$(git rev-parse HEAD) -X main.date=$(date +%Y-%m-%d)" -trimpath -o "$_BINARY" "$_CMD_PATH"
    install_upx
    upx "$_BINARY" --force-overwrite --lzma --no-progress --no-color -qqq
}

# Validate the Go version
validate_versions() {
    REQUIRED_GO_VERSION="1.18"
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    if [ "$(printf '%s\n' "$REQUIRED_GO_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_GO_VERSION" ]; then
        log "error" "Go version must be >= $REQUIRED_GO_VERSION. Detected: $GO_VERSION"
        exit 1
    fi
    log "success" "Go version is valid: $GO_VERSION"
}

# Print a summary of the installation
summary() {
    install_dir="$_BINARY"
    log "success" "Build and installation complete!"
    log "success" "Binary: $_BINARY"
    log "success" "Installed in: $install_dir"
    check_path "$install_dir"
}

# Build the binary and validate the Go version
build_and_validate() {
    validate_versions
    build_binary
}

# Check if the installation directory is in the PATH
# Arguments:
#   $1 - installation directory
check_path() {
    log "info" "Checking if the installation directory is in the PATH..."
    if ! echo "$PATH" | grep -q "$1"; then
        log "warn" "$1 is not in the PATH."
        log "warn" "Add the following to your ~/.bashrc, ~/.zshrc, or equivalent file:"
        log "warn" "export PATH=$1:\$PATH"
    else
        log "success" "$1 is already in the PATH."
    fi
}

# Download the binary from the release URL
download_binary() {
    # Obtem o sistema operacional e a arquitetura
    if ! what_platform > /dev/null; then
        log "error" "Failed to detect platform."
        return 1
    fi

    # Validação: Verificar se o sistema operacional ou a arquitetura são suportados
    if test -z "${_PLATFORM}"; then
        log "error" "Unsupported platform: ${_PLATFORM}"
        return 1
    fi

    # Obter a versão mais recente de forma robusta (fallback para "latest")
    version=$(curl -s "https://api.github.com/repos/${_OWNER}/${_PROJECT_NAME}/releases/latest" | \
        grep "tag_name" | cut -d '"' -f 4 || echo "latest")

    if [ -z "$version" ]; then
        log "error" "Failed to determine the latest version."
        return 1
    fi

    # Construir a URL de download usando a função customizável
    release_url=$(get_release_url)
    log "info" "Downloading ${_APP_NAME} binary for OS=$os, ARCH=$arch, Version=$version..."
    log "info" "Release URL: ${release_url}"

    archive_path="${_TEMP_DIR}/${_APP_NAME}.tar.gz"

    # Realizar o download e validar sucesso
    if ! curl -L -o "${archive_path}" "${release_url}"; then
        log "error" "Failed to download the binary from: ${release_url}"
        return 1
    fi
    log "success" "Binary downloaded successfully."

    # Extração do arquivo para o diretório binário
    log "info" "Extracting binary to: $(dirname "${_BINARY}")"
    if ! tar -xzf "${archive_path}" -C "$(dirname "${_BINARY}")"; then
        log "error" "Failed to extract the binary from: ${archive_path}"
        rm -rf "${_TEMP_DIR}"
        exit 1
    fi

    # Limpar artefatos temporários
    rm -rf "${_TEMP_DIR}"
    log "success" "Binary extracted successfully."

    # Verificar se o binário foi extraído com sucesso
    if [ ! -f "$_BINARY" ]; then
        log "error" "Binary not found after extraction: $_BINARY"
        exit 1
    fi

    log "success" "Download and extraction of ${_APP_NAME} completed!"
}

# Install the binary from the release URL
install_from_release() {
    download_binary
    install_binary
}

# Show banner information
show_banner() {
    # Print the ABOUT message
    printf '\n%s\n\n' "${_ABOUT}"
}

# Main function to handle command line arguments
main() {
  # Show the banner information
  show_banner

  # Check if the user has provided a command
  case "$1" in
      build|BUILD|"-b"|"-B")
          log "info" "Executing build command..."
          build_and_validate || exit 1
          ;;
      install|INSTALL|"-i"|"-I")
          log "info" "Executing install command..."
          read -r -p "Do you want to download the precompiled binary? [y/N] (No will build locally): " c </dev/tty
          log "info" "User choice: ${c}"

          if [ "$c" = "y" ] || [ "$c" = "Y" ]; then
              log "info" "Downloading precompiled binary..."
              install_from_release || exit 1
          else
              log "info" "Building locally..."
              build_and_validate || exit 1
              install_binary || exit 1
          fi
          summary
          ;;
      clean|CLEAN|"-c"|"-C")
          log "info" "Executing clean command..."
          clean || exit 1
          ;;
      *)
          log "error" "Invalid command: $1"
          echo "Usage: $0 {build|install|clean}"
          exit 1
          ;;
  esac
}

# Execute the main function with all script arguments
main "$@"

exit $?