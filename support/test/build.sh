#!/usr/bin/env bash

# Repository root directory (you should run this script from the root of the repository)
_REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || echo "${PWD}" || echo "${0%/*}")"

# Color codes for logging
_SUCCESS="\033[0;32m"
_ERROR="\033[0;31m"
_WARN="\033[0;33m"
_INFO="\033[0;36m"
_NC="\033[0m"

__var_sourced_name() {
  local _self="${BASH_SOURCE-}"
  _self="${_self//${_REPO_ROOT:-}/}"
  _self="${_self//\.sh/}"
  _self="${_self//\-/_}"
  _self="${_self//\//_}"
  echo "_was_sourced_${_self//__/_}"
  return 0
}
__first(){
  if [ "$EUID" -eq 0 ] || [ "$UID" -eq 0 ]; then
    echo "Its not recommended to run as root. Go dependencies may damage your system. If you are sure, go ahead." 1>&2 > /dev/tty
    exit 2 || kill -9 $$ || true
  elif [ -n "${SUDO_USER:-}" ]; then
    echo "Please do not run as root, but with sudo privileges." 1>&2 > /dev/tty
    exit 2 || kill -9 $$ || true
  else
    local _ws_name="$(__var_sourced_name)"
    if test "${BASH_SOURCE-}" != "${0}"; then
      export "${_ws_name}"="true"
    else
      export "${_ws_name}"="false"
      set -o errexit
      set -o nounset
      set -o pipefail
      set -o errtrace
      set -o functrace
      shopt -s inherit_errexit
    fi
  fi
}
__first "$@" >/dev/tty || exit 11
__main_functions() {
  local _func_name=""
  declare -a _args=()
  if [[ $# -gt 1 ]]; then
    _func_name="${2:-}"
    if [[ $# -gt 2 ]]; then
      _args=("${@:3}")
    fi
  else
    _func_name="${1:-}"
  fi
  if declare -F "$_func_name" >'/dev/null' 2>&1; then
    "$_func_name" "${_args[@]}"
    return $?
  else
    kbx_log error "Function (app '$func_name') not found."
    return 1
  fi
}

__get_type() {
  local file="${1:-}"
  local mime_type

  # Check for specific content patterns
  if grep -qE '^\s*---\s*$' "$file" || grep -qE '^\s*[a-zA-Z0-9_-]+:\s*' "$file"; then
    mime_type="application/x-yaml"
  elif grep -qE '^\s*\{\s*"' "$file"; then
    mime_type="application/json"
  elif grep -qE '^\s*#!/bin/(bash|sh|zsh)' "$file"; then
    mime_type="text/x-shellscript"
  else
    # Use head and strings for additional checks
    local head_output=$(sudo head -n 1 "$file") || return 61
    local strings_output=$(sudo strings "$file" | sudo head -n 10) || return 61

    if echo "$head_output" | grep -qE '^\s*---\s*$' || echo "$strings_output" | grep -qE '^\s*[a-zA-Z0-9_-]+:\s*'; then
      mime_type="application/x-yaml"
    elif echo "$head_output" | grep -qE '^\s*\{\s*"' || echo "$strings_output" | grep -qE '^\s*\{\s*"'; then
      mime_type="application/json"
    elif echo "$head_output" | grep -qE '^\s*#!/bin/(bash|sh|zsh)' || echo "$strings_output" | grep -qE '^\s*#!/bin/(bash|sh|zsh)'; then
      mime_type="text/x-shellscript"
    elif grep -qE '^\s*[a-zA-Z_][a-zA-Z0-9_]*=.*$' "$file"; then
      mime_type="text/x-dotenv"
    else
      # Fallback to file extension check
      case "${file##*.}" in
        yml|yaml)
          mime_type="application/x-yaml"
          ;;
        json)
          mime_type="application/json"
          ;;
        sh)
          mime_type="text/x-shellscript"
          ;;
        *)
          # Use file command as a last resort
          mime_type=$(sudo file -b --mime-type "$file") || return 61
          ;;
      esac
    fi
  fi

  echo "$mime_type" || return 61
  return 0
}
__clear(){
  printf "\033[H\033[2J" > /dev/tty
}
__banner() {
  test -n "${1:-}" && __kbx_clear_scr
  echo -e "$($TEXT_COLOR $CYAN)
\t █████   ████          █████                               █████████
\t░░███   ███░          ░░███                               ███░░░░░███
\t ░███  ███   █████ ████░███████   ██████  █████ █████    ░███    ░███ ████████  ████████   █████
\t ░███████   ░░███ ░███ ░███░░███ ███░░███░░███ ░░███     ░███████████░░███░░███░░███░░███ ███░░
\t ░███░░███   ░███ ░███ ░███ ░███░███████  ░░░█████░      ░███░░░░░███ ░███ ░███ ░███ ░███░░█████
\t ░███ ░░███  ░███ ░███ ░███ ░███░███░░░    ███░░░███     ░███    ░███ ░███ ░███ ░███ ░███ ░░░░███
\t █████ ░░████░░████████████████ ░░██████  █████ █████    █████   █████░███████  ░███████  ██████
\t░░░░░   ░░░░  ░░░░░░░░░░░░░░░░   ░░░░░░  ░░░░░ ░░░░░    ░░░░░   ░░░░░ ░███░░░   ░███░░░  ░░░░░░
\t                                                                      ░███      ░███
\t                                                                      █████     █████
\t                                                                     ░░░░░     ░░░░░
${RESET}"
}
__log() {
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

__list_features() {
  echo ""
}
__select_features() {
  local _menu_items=()

  declare -a _all_files=( $( _list_installers ) )
  declare -a _installers_info=()
  for _file_index in "${!_all_files[@]}"; do
    local _file="${_all_files[$_file_index]}"
    sudo chmod +x "${_file}" || kbx_die 3 "Fail to set execute permission on ${_file}."
    local _tmp_app_info="$(bash "${_file}" get_app_info)"
    declare -A _app_info=()
    eval "_app_info=(${_tmp_app_info})"

    local _counter=$(( _file_index + 1 ))
    local _app_name="${_app_info[name]}"
    local _app_description="${_app_info[description]}"
    local _app_file="${_file}"

    _menu_items+=( "${_counter}" "${_app_name} - ${_app_description}" "off" )
  done

  declare -a _dialog_builder=(
    dialog
    --stdout
    --clear
    --scrollbar
    --extra-button
    --mouse
    --extra-label "Select All"
    --ok-label "Install"
    --cancel-label "Cancel"
    --title "Kubex App Installer"
    --backtitle "Select the applications you want to install"
    --checklist "Select packages to install:" 0 0 0 "${_menu_items[@]}"
    --visit-items
  )

  declare -a _selected_indices=$("${_dialog_builder[@]}")
  declare -a _selected_indexes=( $(echo "${_selected_indices[@]}") )

  for _index in ${!_selected_indexes[@]}; do
    local _selected_index=$(( ${_selected_indexes[$_index]} - 1 ))
    local _file_to_install=""

    _file_to_install="${_all_files[$_selected_index]}"

    kbx_log info "Selected index: ${_selected_index}"
    kbx_log info "Selected file: ${_file_to_install}"

    if [ -z "${_file_to_install}" ]; then
      kbx_die 3 "Fail to find the installer for the selected package."
    fi

    local _tmp_app_info="$(bash "${_file_to_install}" get_app_info)"
    declare -A _app_info_inst=()
    eval "_app_info_inst=(${_tmp_app_info})"

    kbx_banner "Applications - Installer $_index of ${#_selected_indexes[@]}"
    kbx_log info "Loading ${_app_info_inst[name]} installation"
    kbx_wait_with_escape 3

    "${_file_to_install}" install || kbx_die 3 "Fail to install ${_app_info_inst[name]}."

    echo "[\"${_app_info_inst[name]}\"]=\"${_file_to_install}\"" >> "${kubex_cache_dir}/.installed"
  done
  return 0
}
__install_selected_features() {
  declare -a _selected_packages=( $( _select_packages ) )

  return 0
}

install() {
  kbx_banner "Applications - Installer"
  kbx_log info "Welcome to the Kubex App installer!"
  kbx_log info "This package was created to facilitate the installation, management, removal, and configuration of applications."
  kbx_log info "We provide many other features, beyond of what originally comes with the applications."
  kbx_log info "For now, let's list the applications that you can choose to install." && kbx_wait_with_escape 1 '' ' '

  _install_selected_packages

  return 0
}

__main() {
  local _ws_name="$(__var_sourced_name)"
  eval "local _ws_name=\$${_ws_name}" >/dev/null
  if [ $(echo "$_ws_name") != "true" ]; then
    __kbx_banner "."

    __main_functions "$@"
    exit $?
  else
    __kbx_log error "This script is not intended to be run directly. Please use the appropriate command."
    exit 1
  fi
}

__main "$@"
