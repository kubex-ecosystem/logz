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
_REPO_ROOT="$(dirname "$(dirname "$(dirname "$(realpath "$0")")")")"

# The _APP_NAME variable is set to the name of the repository. It is used to identify the application.
_APP_NAME="$(basename "$_REPO_ROOT")"

# The _PROJECT_NAME variable is set to the name of the project. It is used for display purposes.
_PROJECT_NAME="$_APP_NAME"

# The _VERSION variable is set to the version of the project. It is used for display purposes.
_VERSION=$(cat "$_REPO_ROOT/version/CLI_VERSION" 2>/dev/null || echo "v0.0.0")

# The _VERSION_GO variable is set to the version of the Go required by the project.
_VERSION_GO=$(grep '^go ' go.mod | awk '{print $2}')

# The _VERSION variable is set to the version of the project. It is used for display purposes.
_LICENSE="MIT"

# The _ABOUT variable contains information about the script and its usage.
_ABOUT="################################################################################
  This Script is used to install ${_PROJECT_NAME} project, version ${_VERSION}.
  Supported OS: Linux, MacOS, Windows
  Supported Architecture: amd64, arm64, 386
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
    - The script will download the binary from the release URL
    - The script will clean up build artifacts if the clean option is provided.
    - The script will check if the required dependencies are installed.
    - The script will validate the Go version before building the binary.
    - The script will check if the installation directory is in the PATH.
################################################################################"

_BANNER="################################################################################

               ██   ██ ██     ██ ██████   ████████ ██     ██
              ░██  ██ ░██    ░██░█░░░░██ ░██░░░░░ ░░██   ██
              ░██ ██  ░██    ░██░█   ░██ ░██       ░░██ ██
              ░████   ░██    ░██░██████  ░███████   ░░███
              ░██░██  ░██    ░██░█░░░░ ██░██░░░░     ██░██
              ░██░░██ ░██    ░██░█    ░██░██        ██ ░░██
              ░██ ░░██░░███████ ░███████ ░████████ ██   ░░██
              ░░   ░░  ░░░░░░░  ░░░░░░░  ░░░░░░░░ ░░     ░░"

# Variable to store the current running shell
_CURRENT_SHELL=""

# The _CMD_PATH variable is set to the path of the cmd directory. It is used to
# identify the location of the main application code.
_CMD_PATH="${_REPO_ROOT}/cmd"

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

# For internal use only
__PLATFORMS=( "windows" "darwin" "linux" )
__ARCHs=( "amd64" "386" "arm64" )

# The _PLATFORM variable is set to the platform name. It is used to identify the
# platform on which the script is running.
_PLATFORM_WITH_ARCH=""
_PLATFORM=""
_ARCH=""

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
      # shellcheck disable=SC2124
      declare -a _FULL_SCRIPT_ARGS=$@

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
  local debug=${3:-${DEBUG:-false}}

  # With colors
  case $type in
    info|_INFO|-i|-I)
      if test "$debug" = true; then
        printf '%b[_INFO]%b ℹ️  %s\n' "$_INFO" "$_NC" "$message"
      fi
      ;;
    warn|_WARN|-w|-W)
      if test "$debug" = true; then
        printf '%b[_WARN]%b ⚠️  %s\n' "$_WARN" "$_NC" "$message"
      fi
      ;;
    error|_ERROR|-e|-E)
      printf '%b[_ERROR]%b ❌  %s\n' "$_ERROR" "$_NC" "$message"
      ;;
    success|_SUCCESS|-s|-S)
      printf '%b[_SUCCESS]%b ✅  %s\n' "$_SUCCESS" "$_NC" "$message"
      ;;
    *)
      if test "$debug" = true; then
        log "info" "$message"
      fi
      ;;
  esac
}

# Detect the platform
what_platform() {
  local _platform=""
  _platform="$(uname -o 2>/dev/null || echo "")"

  local _os=""
  _os="$(uname -s)"

  local _arch=""
  _arch="$(uname -m)"

  # Detect the platform and architecture
  case "${_os}" in
  *inux|*nix)
    _os="linux"
    case "${_arch}" in
    "x86_64")
      _arch="amd64"
      ;;
    "armv6")
      _arch="armv6l"
      ;;
    "armv8" | "aarch64")
      _arch="arm64"
      ;;
    .*386.*)
      _arch="386"
      ;;
    esac
    _platform="linux-${_arch}"
    ;;
  *arwin*)
    _os="darwin"
    case "${_arch}" in
    "x86_64")
      _arch="amd64"
      ;;
    "arm64")
      _arch="arm64"
      ;;
    esac
    _platform="darwin-${_arch}"
    ;;
  MINGW|MSYS|CYGWIN|Win*)
    _os="windows"
    case "${_arch}" in
    "x86_64")
      _arch="amd64"
      ;;
    "arm64")
      _arch="arm64"
      ;;
    esac
    _platform="windows-${_arch}"
    ;;
  *)
    _os=""
    _arch=""
    _platform=""
    ;;
  esac

  if [ -z "${_platform}" ]; then
    log "error" "Unsupported platform: ${_os} ${_arch}"
    log "error" "Please report this issue to the project maintainers."
    return 1
  fi

  # Normalize the platform string
  _PLATFORM_WITH_ARCH="${_platform//\-/\_}"
  _PLATFORM="${_os}"
  _ARCH="${_arch}"

  return 0
}

# Get the platform and architecture variables
_get_platform_arch_vars() {
  what_platform || return 1

  local _platform="${1:-"${_PLATFORM:-"${_PLATFORM_WITH_ARCH%_*}"}"}"
  local _arch="${2:-"${_ARCH:-"${_PLATFORM_WITH_ARCH#*\_}"}"}"

  if [ "${_platform}" != "all" ]; then
      _platforms=( "${_platform}" )
      if [[ "${_platform}" != "darwin" ]]; then
        _archS=( "amd64" "386" )
        if [[ "${_arch}" != "all" && "${_arch}" != "arm64" ]]; then
          _archS=( "${_arch}" )
        fi
      else
        _archS=( "amd64" "arm64" )
        if [[ "${_arch}" != "all" && "${_arch}" != "386" ]]; then
          _archS=( "${_arch}" )
        fi
      fi
  else
      _platforms=( "${__PLATFORMS[@]}" )
      _archS=( "${__ARCHs[@]}" )
  fi
  # Print the platform and architecture variables
  echo "("
  for _platform_pos in "${_platforms[@]}"; do
    echo "${_platform_pos} "
  done
  echo ")"
  echo ","
  echo "("
  for _arch_pos in "${_archS[@]}"; do
    echo "${_arch_pos} "
  done
  echo ")"
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
    if [ ! -f "$shell_rc_file" ]; then
        log "error" "Shell configuration file not found: $shell_rc_file"
        return 1
    fi
    echo "$shell_rc_file"
    return 0
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
    local _platforms=( "windows" "darwin" "linux" )
    local _archS=( "amd64" "386" "arm64" )
    for _platform in "${_platforms[@]}"; do
        for _arch in "${_archS[@]}"; do
            local _OUTPUT_NAME="${_BINARY}_${_platform}_${_arch}"
            if [ "${_platform}" != "windows" ]; then
                _COMPRESS_NAME="${_OUTPUT_NAME}.tar.gz"
            else
                _OUTPUT_NAME+=".exe"
                _COMPRESS_NAME="${_BINARY}_${_platform}_${_arch}.zip"
            fi
            rm -f "${_OUTPUT_NAME}" || true
            rm -f "${_COMPRESS_NAME}" || true
            if [ -f "${_OUTPUT_NAME}" ]; then
                if sudo -v; then
                    sudo rm -f "${_OUTPUT_NAME}" || true
                else
                    log "error" "Failed to remove build artifact: ${_OUTPUT_NAME}"
                    log "error" "Please remove it manually with 'sudo rm -f \"${_OUTPUT_NAME}\"'"
                fi
            fi
            if [ -f "${_COMPRESS_NAME}" ]; then
                if sudo -v; then
                    sudo rm -f "${_COMPRESS_NAME}" || true
                else
                    log "error" "Failed to remove build artifact: ${_COMPRESS_NAME}"
                    log "error" "Please remove it manually with 'sudo rm -f \"${_COMPRESS_NAME}\"'"
                fi
            fi
        done
    done
    log "success" "Cleaned up build artifacts."
    return 0
}

# Install the binary to the appropriate directory
install_binary() {
    local _SUFFIX="${_PLATFORM_WITH_ARCH}"
    local _BINARY_TO_INSTALL="${_BINARY}${_SUFFIX:+_${_SUFFIX}}"
    log "info" "Installing binary: '$_BINARY_TO_INSTALL' like '$_APP_NAME'"

    if [ "$(id -u)" -ne 0 ]; then
        log "info" "You are not root. Installing in $_LOCAL_BIN..."
        mkdir -p "$_LOCAL_BIN"
        cp "$_BINARY_TO_INSTALL" "$_LOCAL_BIN/$_APP_NAME" || exit 1
        add_to_path "$_LOCAL_BIN"
    else
        log "info" "Root detected. Installing in $_GLOBAL_BIN..."
        cp "$_BINARY_TO_INSTALL" "$_GLOBAL_BIN/$_APP_NAME" || exit 1
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
    # shellcheck disable=SC2317
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
# shellcheck disable=SC2207,SC2116,SC2091,SC2155,SC2005
build_binary() {
  # Get the platform and architecture variables
  local __ARGS=( "${1:-}" "${2:-}" )
  local __vars_arrays=$(_get_platform_arch_vars "${__ARGS[@]}")
  # __platforms
  local _platforms=( )
  eval _platforms="$(echo "$(printf '%s\n' "${__vars_arrays%%\,*}")")"
  #echo "${_platforms[@]}"
  local _archS=( "$(echo "$(printf '%s\n' "${__vars_arrays#*,}" | tr -d '(' | tr -d ')' | tr -d '\n')")" )
  eval _archS="( "$(echo "$(printf '%s' "${_archS[@]}")")" )"
  #echo "${_archS[@]}"

  for _platform_pos in "${_platforms[@]}"; do
    if test -z "$_platform_pos"; then
      continue
    fi
    for _arch_pos in "${_archS[@]}"; do
      if test -z "$_arch_pos"; then
        continue
      fi
      if [[ "$_platform_pos" != "darwin" && "$_arch_pos" == "arm64" ]]; then
        continue
      fi
      if [[ "$_platform_pos" != "windows" && "$_arch_pos" == "386" ]]; then
        continue
      fi
      local _OUTPUT_NAME="$(printf '%s_%s_%s' "${_BINARY}" "$_platform_pos" "$_arch_pos")"
      if [[ "$_platform_pos" == "windows" ]]; then
        _OUTPUT_NAME="$(printf '%s.exe' "${_OUTPUT_NAME}")"
      fi

      local _build_env=(
        "export GOOS='$_platform_pos' &&"
        "export GOARCH='$_arch_pos' &&"
      )
      local _build_args=(
        "-ldflags '-s -w -X main.version=$(git describe --tags) -X main.commit=$(git rev-parse HEAD) -X main.date=$(date +%Y-%m-%d)'"
        "-trimpath -o '${_OUTPUT_NAME}' '${_CMD_PATH}'"
      )
      local _build_cmd=(
        "${_build_env[*]}"
        "go build"
        "${_build_args[*]}"
      )
      local _build_cmd_str="${_build_cmd[*]}"
      log "info" "$(printf '%s %s/%s\n' "Building the binary for" "${_platform_pos}" "${_arch_pos}")"
      log "info" "Command: ${_build_cmd_str}"
      # Build the binary using the environment variables and arguments
      if ! bash -c "$_build_cmd_str"; then
        log "error" "Failed to build the binary for ${_platform_pos} ${_arch_pos}"
        log "error" "Command: ${_build_cmd_str}"
        return 1
      else
        # If the build was successful, check if UPX is installed and compress the binary (if not Windows)
        if [[ "$_platform_pos" != "windows" ]]; then
            install_upx
            log "info" "Packing/compressing the binary with UPX..."
            upx "${_OUTPUT_NAME}" --force-overwrite --lzma --no-progress --no-color -qqq || true
            log "success" "Binary packed/compressed successfully: ${_OUTPUT_NAME}"
        fi
        # Check if the binary was created successfully (if not Windows)
        if [[ ! -f "${_OUTPUT_NAME}" ]]; then
          log "error" "Binary not found after build: ${_OUTPUT_NAME}"
          log "error" "Command: ${_build_cmd_str}"
          return 1
        else
          local compress_vars=( "$_platform_pos" "$_arch_pos" )
          compress_binary "${compress_vars[@]}" || return 1
          log "success" "Binary created successfully: ${_OUTPUT_NAME}"
        fi
      fi
    done
  done

  echo ""
  log "success" "All builds completed successfully!"
  echo ""
  return 0
}

# Compress the binary into a single tar.gz/zip file
# shellcheck disable=SC2207,SC2116,SC2091,SC2155,SC2005
compress_binary() {
  # Get the platform and architecture variables
  local __ARGS=( "${1:-}" "${2:-}" )
  local __vars_arrays=$(_get_platform_arch_vars "${__ARGS[@]}")
  # __platforms
  local _platforms=( )
  eval _platforms="$(echo "$(printf '%s\n' "${__vars_arrays%%\,*}")")"
  #echo "${_platforms[@]}"
  local _archS=( "$(echo "$(printf '%s\n' "${__vars_arrays#*,}" | tr -d '(' | tr -d ')' | tr -d '\n')")" )
  eval _archS="( "$(echo "$(printf '%s' "${_archS[@]}")")" )"
  #echo "${_archS[@]}"

  for _platform_pos in "${_platforms[@]}"; do
    if test -z "${_platform_pos}"; then
      continue
    fi
    for _arch_pos in "${_archS[@]}"; do
      if test -z "${_arch_pos}"; then
        continue
      fi
      if [[ "${_platform_pos}" != "darwin" && "${_arch_pos}" == "arm64" ]]; then
        continue
      fi
      if [[ "${_platform_pos}" == "linux" && "${_arch_pos}" == "386" ]]; then
        continue
      fi
      local _BINARY_NAME="${_BINARY}_${_platform_pos}_${_arch_pos}"
      local _OUTPUT_NAME=""
      log "info" "Compressing the binary for ${_platform_pos} ${_arch_pos} into ${_OUTPUT_NAME}..."
      if [ "${_platform_pos}" != "windows" ]; then
        _OUTPUT_NAME="${_BINARY_NAME}.tar.gz"
        tar -czf "${_OUTPUT_NAME}" "${_BINARY_NAME}" || return 1
      else
        _OUTPUT_NAME="${_BINARY_NAME}.zip"
        log "info" "Compressing the binary for ${_platform_pos} ${_arch_pos} into ${_OUTPUT_NAME}..."
        zip -r -9 "${_OUTPUT_NAME}" "${_BINARY_NAME}.exe" || return 1
      fi
      if [ ! -f "${_OUTPUT_NAME}" ]; then
        log "error" "Failed to create tar.gz file: ${_OUTPUT_NAME}"
        return 1
      fi
    done
  done

  log "success" "All binaries compressed successfully!"

  return 0
}

# Validate the Go version
validate_versions() {
    REQUIRED_GO_VERSION="${_VERSION_GO:-1.20.0}"
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    if [ "$(printf '%s\n' "$REQUIRED_GO_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_GO_VERSION" ]; then
        log "error" "Go version must be >= $REQUIRED_GO_VERSION. Detected: $GO_VERSION"
        exit 1
    fi
    log "success" "Go version is valid: $GO_VERSION"
    go mod tidy || return 1
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
    # Check if the platform and architecture are set
    what_platform || return 1

    # Check if the Go version is valid
    validate_versions

    local _PLATFORM_ARG="${1:-${_PLATFORM}}"
    local _ARCH_ARG="${2:-${_ARCH}}"

    local _WHICH_COMPILE_ARG=( )
    case "${_PLATFORM_ARG}" in
        all|ALL|a|A|-a|-A)
            log "info" "Building for all platforms..."
            _WHICH_COMPILE_ARG+=( "all" )
            ;;
        *)
            case "${_PLATFORM_ARG}" in
                win|WIN|windows|WINDOWS|w|W|-w|-W)
                  log "info" "Building for Windows..."
                  _WHICH_COMPILE_ARG+=( "windows" )
                  ;;
                linux|LINUX|l|L|-l|-L)
                  log "info" "Building for Linux..."
                  _WHICH_COMPILE_ARG+=( "linux" )
                  ;;
                darwin|DARWIN|macOS|MACOS|m|M|-m|-M)
                  log "info" "Building for MacOS..."
                  _WHICH_COMPILE_ARG+=( "darwin" )
                  ;;
                *)
                  log "error" "build_and_validate: Unsupported platform: '$(printf '%s' "${_WHICH_COMPILE_ARG[*]}")'."
                  log "error" "Please specify a valid platform (windows, linux, darwin)."
                  return 1
                  ;;
            esac
    esac
    case "${_ARCH_ARG}" in
        all|ALL|a|A|-a|-A)
            log "info" "Building for all architectures..."
            _WHICH_COMPILE_ARG+=( "all" )
            ;;
        *)
            case "${_ARCH_ARG}" in
                amd64|AMD64|x86_64|X86_64|x64|X64)
                  log "info" "Building for AMD64..."
                  _WHICH_COMPILE_ARG+=( "amd64" )
                  ;;
                arm64|ARM64|aarch64|AARCH64)
                  log "info" "Building for ARM64..."
                  _WHICH_COMPILE_ARG+=( "arm64" )
                  ;;
                386|i386|I386)
                  log "info" "Building for 386..."
                  _WHICH_COMPILE_ARG+=( "386" )
                  ;;
                *)
                  log "error" "build_and_validate: Unsupported architecture: '${__ARCH}'. Please specify a valid architecture (amd64, arm64, 386)."
                  return 1
                  ;;
            esac
    esac

    # Call the build_binary function with the platform and architecture arguments
    build_binary "${_WHICH_COMPILE_ARG[@]}" || return 1

    return 0
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

# Show about information
show_about() {
    # Print the ABOUT message
    printf '%s\n\n' "${_ABOUT:-}"
}

# Show banner information
show_banner() {
    # Print the ABOUT message
    printf '\n%s\n\n' "${_BANNER:-}"
}

# Show headers information
show_headers() {
    # Print the BANNER message
    show_banner || return 1
    # Print the ABOUT message
    show_about || return 1
}

# Main function to handle command line arguments
main() {
  # Detect the platform if not provided, will be used in the build command
  what_platform || exit 1

  local _WHICH_COMPILE_ARG=( "${_PLATFORM}" "${_ARCH}" )

  # Show the banner information
  if test "$debug" != true; then
    show_headers
  else
    log "info" "Debug mode enabled. Skipping banner..."
    if test -z "${HIDE_ABOUT}"; then
      show_about
    fi
  fi

  # Check if the user has provided a command
  case "${1:-}" in
      build|BUILD|"-b"|"-B")
        _ARGS=( "$@" )
        local _default_label='Auto detect'
        log "info" "$(printf '%s%s' "Executing build command for platform: " "${_ARGS[1]:-${_default_label}}'")"
        # shellcheck disable=SC2124
        local _arrArgs=( "${_ARGS[1]:1:-1}" )
        if [[ $# -gt 2 ]]; then
            # If the user provided a platform argument, use it
            _WHICH_COMPILE_ARG=( "${_arrArgs[1]}" "${_arrArgs[2]}" )
        elif [[ $# -gt 1 ]]; then
            # If the user provided a platform argument, use it
            _WHICH_COMPILE_ARG=( "${_arrArgs[1]}" "${_ARCH}" )
        else
            # Otherwise, use the default platform and architecture
            _WHICH_COMPILE_ARG=( "${_PLATFORM}" "${_ARCH}" )
        fi
        # Call the build function with the detected platform
        build_and_validate "${_WHICH_COMPILE_ARG[@]}" || exit 1
        ;;
      install|INSTALL|"-i"|"-I")
          log "info" "Executing install command..."
          read -r -p "Do you want to download the precompiled binary? [y/N] (No will build locally): " c </dev/tty
          log "info" "User choice: ${c}"

          if [ "$c" = "y" ] || [ "$c" = "Y" ]; then
              log "info" "Downloading precompiled binary..." true
              install_from_release "${_WHICH_COMPILE_ARG[@]}" || exit 1
          else
              log "info" "Building locally..." true
              build_and_validate "${_WHICH_COMPILE_ARG[@]}" || exit 1
              install_binary "${_WHICH_COMPILE_ARG[@]}" || exit 1
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