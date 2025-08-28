#!/usr/bin/env bash

set -o errexit
set -o pipefail
#set -o nounset
#set -o xtrace

abort() {
  printf "%s\n" "$@" >&2
  exit 1
}

# string formatters
if [[ -t 1 ]]; then
  tty_escape() { printf "\033[%sm" "$1"; }
else
  tty_escape() { :; }
fi
tty_mkbold() { tty_escape "1;$1"; }
tty_underline="$(tty_escape "4;39")"
tty_blue="$(tty_mkbold 34)"
tty_red="$(tty_mkbold 31)"
tty_bold="$(tty_mkbold 39)"
tty_reset="$(tty_escape 0)"

shell_join() {
  local arg
  printf "%s" "$1"
  shift
  for arg in "$@"; do
    printf " "
    printf "%s" "${arg// /\ }"
  done
}

chomp() {
  printf "%s" "${1/"$'\n'"/}"
}

ring_bell() {
  # Use the shell's audible bell.
  if [[ -t 1 ]]; then
    printf "\a"
  fi
}

exists_but_not_writable() {
  [[ -e "$1" ]] && ! [[ -r "$1" && -w "$1" && -x "$1" ]]
}

ohai() {
  printf "${tty_blue}==>${tty_bold} %s${tty_reset}\n" "$(shell_join "$@")"
}

warn() {
  printf "${tty_red}Warning${tty_reset}: %s\n" "$(chomp "$1")"
}

execute() {
  if ! "$@"; then
    abort "$(printf "Failed during: %s" "$(shell_join "$@")")"
  fi
}

STAT_PRINTF=("stat" "-f")

get_owner() {
  "${STAT_PRINTF[@]}" "%u" "$1"
}

file_not_owned() {
  [[ "$(get_owner "$1")" != "$(id -u)" ]]
}

get_group() {
  "${STAT_PRINTF[@]}" "%g" "$1"
}

file_not_grpowned() {
  [[ " $(id -G "${USER}") " != *" $(get_group "$1") "* ]]
}

display_help() {
  echo ""
  cat <<EOS
  Usage: ./install.sh

  Flags:
    --help                          Display help and usage information
    --version=<VERSION_NUMBER>      Specify the version to install
EOS
}

VERSION="3.3.2"

while [[ "$#" -gt 0 ]]; do
  case "$1" in
    --version=*)
      VERSION=${1#*=}
      ;;
    --help)
      display_help
      exit 0
      ;;
    *)
      ohai "Invalid argument: $1"
      display_help
      exit 1
      ;;
  esac
  shift
done

# USER isn't always set so provide a fall back for the installer and subprocesses.
if [[ -z "${USER-}" ]]; then
  USER="$(chomp "$(id -un)")"
  export USER
fi

# First check OS.
OS="$(uname)"
if [[ "${OS}" == "Linux" ]]; then
  BUGSNAG_CLI_ON_LINUX=1
  OS_NAME="linux"
  GROUP="$(id -gn)"
elif [[ "${OS}" == "Darwin" ]]; then
  BUGSNAG_CLI_ON_MACOS=1
  OS_NAME="macos"
  GROUP="admin"
else
  abort "This install script only works on Linux or Macos"
fi

UNAME_MACHINE="$(uname -m)"
BUGSNAG_CLI_PREFIX="${HOME}/.local/bugsnag"

ohai "This script will install:"
echo "${BUGSNAG_CLI_PREFIX}/bin/bugsnag-cli"

directories=(
  bin
)

mkdirs=()

for dir in "${directories[@]}"; do
  if ! [[ -d "${BUGSNAG_CLI_PREFIX}/${dir}" ]]; then
    mkdirs+=("${BUGSNAG_CLI_PREFIX}/${dir}")
  fi
done

if [[ "${#mkdirs[@]}" -gt 0 ]]; then
  ohai "The following new directories will be created:"
  printf "%s\n" "${mkdirs[@]}"
fi

if [[ "${#mkdirs[@]}" -gt 0 ]]; then
  execute "mkdir" "-p" "${mkdirs[@]}"
  execute "chmod" "ug=rwx" "${mkdirs[@]}"
  execute "chown" "${USER}" "${mkdirs[@]}"
  execute "chgrp" "${GROUP}" "${mkdirs[@]}"
fi

ohai "Downloading and installing Bugsnag CLI..."
(
  cd "${BUGSNAG_CLI_PREFIX}" >/dev/null || return

  url="https://github.com/bugsnag/bugsnag-cli/releases"
  output_file="${BUGSNAG_CLI_PREFIX}/bin/bugsnag-cli"

  http_status_code=$(execute "curl" "-s" "-o" "/dev/null" "-w" "%{http_code}" "-#" "-L" "$url/tag/v${VERSION}")

  if [ "${http_status_code}" -eq 404 ]; then
      abort "Unable to download bugsnag-cli v${VERSION}. Please check https://github.com/bugsnag/bugsnag-cli/releases for a list of releases."
  elif [ "${http_status_code}" -ne 200 ]; then
      abort "The URL returned a non-404 error with status code ${http_status_code}."
  else
      execute "curl" "-#" "-L" "$url/download/v${VERSION}/${UNAME_MACHINE}-${OS_NAME}-bugsnag-cli" "-o" "$output_file"
  fi

  execute "chmod" "ug=rwx" "${BUGSNAG_CLI_PREFIX}/bin/bugsnag-cli"

) ||
  exit 1

if [[ ":${PATH}:" != *":${BUGSNAG_CLI_PREFIX}/bin:"* ]]; then
  warn "${BUGSNAG_CLI_PREFIX}/bin is not in your PATH.
  Instructions on how to configure your shell for Bugsnag CLI
  can be found in the 'Next steps' section below."
fi

ohai "Installation successful!"
echo

ring_bell

case "${SHELL}" in
*/bash*)
  if [[ -r "${HOME}/.bash_profile" ]]; then
    shell_profile="${HOME}/.bash_profile"
  else
    shell_profile="${HOME}/.profile"
  fi
  ;;
*/zsh*)
  shell_profile="${HOME}/.zprofile"
  ;;
*)
  shell_profile="${HOME}/.profile"
  ;;
esac

if [[ "$(which bugsnag-cli)" != "${BUGSNAG_CLI_PREFIX}/bin/bugsnag-cli" ]]; then
  ohai "Next steps:"

  cat <<EOS
- Run these three commands in your terminal to add Bugsnag CLI to your ${tty_bold}PATH${tty_reset}:
    echo "# Bugsnag CLI" >> ${shell_profile}
    echo 'PATH="${BUGSNAG_CLI_PREFIX}/bin:\$PATH"' >> ${shell_profile}
    source "${shell_profile}"
EOS
fi
