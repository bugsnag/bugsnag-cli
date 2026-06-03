#!/usr/bin/env bash

# Bugsnag CLI installer
# This script detects the current OS and CPU architecture,
# downloads the appropriate prebuilt bugsnag-cli binary,
# and installs it into a user-local prefix.

# Fail fast on errors and broken pipelines
set -o errexit
set -o pipefail
#set -o nounset
#set -o xtrace

# Print an error message to stderr and exit the script
abort() {
  printf "%s\n" "$@" >&2
  exit 1
}

# Terminal formatting helpers (used for colored and styled output)
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

# Join command arguments for readable logging output
shell_join() {
  local arg
  printf "%s" "$1"
  shift
  for arg in "$@"; do
    printf " "
    printf "%s" "${arg// /\ }"
  done
}

# Remove a trailing newline from a string
chomp() {
  printf "%s" "${1/"$'\n'"/}"
}

# Emit an audible bell when running in an interactive terminal
ring_bell() {
  # Use the shell's audible bell.
  if [[ -t 1 ]]; then
    printf "\a"
  fi
}

# Check whether a path exists but is not readable, writable, and executable
exists_but_not_writable() {
  [[ -e "$1" ]] && ! [[ -r "$1" && -w "$1" && -x "$1" ]]
}

# User-facing logging helpers
ohai() {
  printf "${tty_blue}==>${tty_bold} %s${tty_reset}\n" "$(shell_join "$@")"
}

warn() {
  printf "${tty_red}Warning${tty_reset}: %s\n" "$(chomp "$1")"
}

# Execute a command and abort with context if it fails
execute() {
  if ! "$@"; then
    abort "$(printf "Failed during: %s" "$(shell_join "$@")")"
  fi
}

# File ownership and permission helpers
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

# Print usage and help information
display_help() {
  echo ""
  cat <<EOS
  Usage: ./install.sh

  Flags:
    --help                          Display help and usage information
    --version=<VERSION_NUMBER>      Specify the version to install
EOS
}

VERSION="3.10.2"

# Parse command-line arguments
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

# Ensure USER is set for ownership and permission operations
if [[ -z "${USER-}" ]]; then
  USER="$(chomp "$(id -un)")"
  export USER
fi

# Detect supported operating systems (Linux and macOS only)
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

# Normalize CPU architecture names to match release asset naming
UNAME_MACHINE="$(uname -m)"

case "$UNAME_MACHINE" in
  aarch64|arm64)
    ARCH="arm64"
    ;;
  x86_64|amd64)
    ARCH="x86_64"
    ;;
  i386|i686)
    UNAME_MACHINE="i386"
    ;;
  *)
    abort "Unsupported architecture: $UNAME_MACHINE"
    ;;
esac

# Map OS name to platform identifier used in download artifacts
if [[ "$OS_NAME" == "linux" ]]; then
  PLATFORM="linux"
elif [[ "$OS_NAME" == "macos" ]]; then
  PLATFORM="macos"
else
  abort "Unsupported OS: $OS_NAME"
fi

BIN_NAME="bugsnag-cli"
DOWNLOAD_NAME="${ARCH}-${PLATFORM}-${BIN_NAME}"
BUGSNAG_CLI_PREFIX="${HOME}/.local/bugsnag"

ohai "This script will install:"
echo "${BUGSNAG_CLI_PREFIX}/bin/bugsnag-cli"

# Determine which directories need to be created under the install prefix
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

# Create required directories with correct ownership and permissions
if [[ "${#mkdirs[@]}" -gt 0 ]]; then
  execute "mkdir" "-p" "${mkdirs[@]}"
  execute "chmod" "ug=rwx" "${mkdirs[@]}"
  execute "chown" "${USER}" "${mkdirs[@]}"
  execute "chgrp" "${GROUP}" "${mkdirs[@]}"
fi

# Download the appropriate bugsnag-cli binary and install it
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
      execute "curl" "-#" "-L" "$url/download/v${VERSION}/${DOWNLOAD_NAME}" "-o" "$output_file"
  fi

  execute "chmod" "ug=rwx" "${output_file}"

) ||
  exit 1

# Warn the user if the install location is not on PATH
if [[ ":${PATH}:" != *":${BUGSNAG_CLI_PREFIX}/bin:"* ]]; then
  warn "${BUGSNAG_CLI_PREFIX}/bin is not in your PATH.
  Instructions on how to configure your shell for Bugsnag CLI
  can be found in the 'Next steps' section below."
fi

ohai "Installation successful!"
echo

ring_bell

# Detect the user's shell profile for PATH modification instructions
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

# Print post-install instructions if bugsnag-cli is not yet discoverable
if [[ "$(which bugsnag-cli)" != "${output_file}" ]]; then
  ohai "Next steps:"

  cat <<EOS
- Run these three commands in your terminal to add Bugsnag CLI to your ${tty_bold}PATH${tty_reset}:
    echo "# Bugsnag CLI" >> ${shell_profile}
    echo 'PATH="${BUGSNAG_CLI_PREFIX}/bin:\$PATH"' >> ${shell_profile}
    source "${shell_profile}"
EOS
fi
