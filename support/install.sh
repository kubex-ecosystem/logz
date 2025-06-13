#!/usr/bin/env bash

set -euo pipefail
set -o errtrace
set -o functrace
set -o posix

IFS=$'\n\t'

_DEBUG=${DEBUG:-false}
_HIDE_ABOUT=${HIDE_ABOUT:-false}

# Carrega os arquivos de biblioteca
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
#shellcheck source=/dev/null
source "${SCRIPT_DIR}/config.sh"
#shellcheck source=/dev/null
source "${SCRIPT_DIR}/utils.sh"
#shellcheck source=/dev/null
source "${SCRIPT_DIR}/platform.sh"
#shellcheck source=/dev/null
source "${SCRIPT_DIR}/validate.sh"
#shellcheck source=/dev/null
source "${SCRIPT_DIR}/install_funcs.sh"
#shellcheck source=/dev/null
source "${SCRIPT_DIR}/build.sh"
#shellcheck source=/dev/null
source "${SCRIPT_DIR}/info.sh"

# Inicializa os traps
set_trap "$@"

clear_screen

main() {
  if ! what_platform; then
    log error "Plataforma não suportada: ${_PLATFORM}"
    exit 1
  fi

  if [[ "${_DEBUG}" != true ]]; then
    show_headers
    if [[ -z "${_HIDE_ABOUT}" ]]; then
      show_about
    fi
  else
    log info "Modo debug ativado; banner será ignorado..."
    if [[ -z "${_HIDE_ABOUT}" ]]; then
      show_about
    fi
  fi

  _ARGS=( "$@" )
  local default_label='Auto detect'
  local arrArgs=( "${_ARGS[@]:0:$#}" )
  local PLATFORM_ARG
  PLATFORM_ARG=$(_get_os_from_args "${arrArgs[1]:-${_PLATFORM}}")
  local ARCH_ARG
  ARCH_ARG=$(_get_arch_arr_from_args "${arrArgs[2]:-${_ARCH}}")

  log info "Comando: ${arrArgs[0]:-}" true
  log info "Plataforma: ${PLATFORM_ARG:-$default_label}" true
  log info "Arquitetura: ${ARCH_ARG:-$default_label}" true
  log info "Args: ${_ARGS[*]:-}" true

  case "${arrArgs[0]:-}" in
    build|BUILD|-b|-B)
      # validate_versions
      log info "Executando comando de build..."
      build_binary "${PLATFORM_ARG}" "${ARCH_ARG}" || exit 1
      ;;
    install|INSTALL|-i|-I)
      log info "Executando comando de instalação..."
      read -r -p "Deseja baixar o binário pré-compilado? [y/N] (Caso contrário, fará build local): " choice </dev/tty
      log info "Escolha do usuário: ${choice}"
      if [[ "$choice" == "y" || "$choice" == "Y" ]]; then
          log info "Baixando binário pré-compilado..."
          install_from_release
      else
          log info "Realizando build local..."
          validate_versions
          build_binary "${PLATFORM_ARG}" "${ARCH_ARG}" || exit 1
          install_binary
      fi
      summary
      ;;
    clear|clean|CLEAN|-c|-C)
      log info "Executando comando de limpeza..."
      clean_artifacts
      log success "Clean executado com sucesso."
      ;;
    *)
      log error "Comando inválido: ${arrArgs[0]:-}"
      echo "Uso: $0 {build|install|clean}"
      ;;
  esac
}

# Função para limpar artefatos de build
clean_artifacts() {
    log info "Limpando artefatos de build..."
    local platforms=("windows" "darwin" "linux")
    local archs=("amd64" "386" "arm64")
    for platform in "${platforms[@]}"; do
        for arch in "${archs[@]}"; do
            local output_name
            output_name=$(printf '%s_%s_%s' "${_BINARY}" "${platform}" "${arch}")
            if [[ "${platform}" != "windows" ]]; then
                local compress_name="${output_name}.tar.gz"
            else
                output_name="${output_name}.exe"
                local compress_name="${_BINARY}_${platform}_${arch}.zip"
            fi
            rm -f "${output_name}" || true
            rm -f "${compress_name}" || true
        done
    done
    log success "Artefatos de build removidos."
}

# echo "MAKE ARGS: ${ARGS[*]:-}"
log info "Starting installation script..."
main "$@"

