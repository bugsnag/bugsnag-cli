#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset
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

have_sudo_access() {
  if [[ ! -x "/usr/bin/sudo" ]]; then
    return 1
  fi

  local -a SUDO=("/usr/bin/sudo")
  if [[ -n "${SUDO_ASKPASS-}" ]]; then
    SUDO+=("-A")
  elif [[ -n "${NONINTERACTIVE-}" ]]; then
    SUDO+=("-n")
  fi

  if [[ -z "${HAVE_SUDO_ACCESS-}" ]]; then
    if [[ -n "${NONINTERACTIVE-}" ]]; then
      "${SUDO[@]}" -l mkdir &>/dev/null
    else
      "${SUDO[@]}" -v && "${SUDO[@]}" -l mkdir &>/dev/null
    fi
    HAVE_SUDO_ACCESS="$?"
  fi

  if [[ -n "${BUGSNAG_CLI_ON_MACOS-}" ]] && [[ "${HAVE_SUDO_ACCESS}" -ne 0 ]]; then
    abort "Need sudo access on macOS (e.g. the user ${USER} needs to be an Administrator)!"
  fi

  return "${HAVE_SUDO_ACCESS}"
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

execute_sudo() {
  local -a args=("$@")
  if have_sudo_access; then
    if [[ -n "${SUDO_ASKPASS-}" ]]; then
      args=("-A" "${args[@]}")
    fi
    ohai "/usr/bin/sudo" "${args[@]}"
    execute "/usr/bin/sudo" "${args[@]}"
  else
    ohai "${args[@]}"
    execute "${args[@]}"
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
BUGSNAG_CLI_GIT_REMOTE="https://api.github.com/repos/bugsnag/bugsnag-cli/releases/latest"

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
elif [[ "${OS}" == "Darwin" ]]; then
  BUGSNAG_CLI_ON_MACOS=1
  OS_NAME="macos"
else
  abort "This install script only works on Linux or Macos"
fi

# Get the OS and ARCH we're using
if [[ -n "${BUGSNAG_CLI_ON_MACOS}" ]]; then
  UNAME_MACHINE="$(/usr/bin/uname -m)"

  if [[ "${UNAME_MACHINE}" == "arm64" ]]; then
    # On ARM macOS, this script installs to /opt/bugsnag only
    BUGSNAG_CLI_PREFIX="/opt/bugsnag"
  else
    # On Intel macOS, this script installs to /usr/local/bugsnag only
    BUGSNAG_CLI_PREFIX="/usr/local/bugsnag"
  fi
  GROUP="admin"
  INSTALL=("/usr/bin/install" -d -o "root" -g "wheel" -m "0755")
else
  UNAME_MACHINE="$(uname -m)"

  # On Linux, this script installs to /usr/local/bugsnag only
  BUGSNAG_CLI_PREFIX="/usr/local/bugsnag"

  GROUP="$(id -gn)"
  INSTALL=("/usr/bin/install" -d -o "${USER}" -g "${GROUP}" -m "0755")
fi

unset HAVE_SUDO_ACCESS # unset this from the environment

# Invalidate sudo timestamp before exiting (if it wasn't active before).
if [[ -x /usr/bin/sudo ]] && ! /usr/bin/sudo -n -v 2>/dev/null; then
  trap '/usr/bin/sudo -k' EXIT
fi

# Things can fail later if `pwd` doesn't exist.
# Also sudo prints a warning message for no good reason
cd "/usr" || exit 1

ohai 'Checking for `sudo` access (which may request your password)...'

if [[ -n "${BUGSNAG_CLI_ON_MACOS-}" ]]; then
  have_sudo_access
elif
  ! [[ -w "${BUGSNAG_CLI_PREFIX}" ]]
  ! have_sudo_access
then
  abort "$(
    cat <<EOABORT
Insufficient permissions to install Bugsnag CLI to \"${BUGSNAG_CLI_PREFIX}\" (the default prefix).
EOABORT
  )"
fi

if [[ -d "${BUGSNAG_CLI_PREFIX}" && ! -x "${BUGSNAG_CLI_PREFIX}" ]]; then
  abort "$(
    cat <<EOABORT
The Homebrew prefix ${tty_underline}${BUGSNAG_CLI_PREFIX}${tty_reset} exists but is not searchable.
If this is not intentional, please restore the default permissions and
try running the installer again:
    sudo chmod 775 ${BUGSNAG_CLI_PREFIX}
EOABORT
  )"
fi

ohai "This script will install:"
echo "${BUGSNAG_CLI_PREFIX}/bin/bugsnag-cli"

directories=(
  bin bin/bugsnag-cli
)
group_chmods=()
for dir in "${directories[@]}"; do
  if exists_but_not_writable "${BUGSNAG_CLI_PREFIX}/${dir}"; then
    group_chmods+=("${BUGSNAG_CLI_PREFIX}/${dir}")
  fi
done

directories=(
  bin
)
mkdirs=()
for dir in "${directories[@]}"; do
  if ! [[ -d "${BUGSNAG_CLI_PREFIX}/${dir}" ]]; then
    mkdirs+=("${BUGSNAG_CLI_PREFIX}/${dir}")
  fi
done

chmods=()
if [[ "${#group_chmods[@]}" -gt 0 ]]; then
  chmods+=("${group_chmods[@]}")
fi

chowns=()
chgrps=()
if [[ "${#chmods[@]}" -gt 0 ]]; then
  for dir in "${chmods[@]}"; do
    if file_not_owned "${dir}"; then
      chowns+=("${dir}")
    fi
    if file_not_grpowned "${dir}"; then
      chgrps+=("${dir}")
    fi
  done
fi

if [[ "${#group_chmods[@]}" -gt 0 ]]; then
  ohai "The following existing directories will be made group writable:"
  printf "%s\n" "${group_chmods[@]}"
fi
if [[ "${#chowns[@]}" -gt 0 ]]; then
  ohai "The following existing directories will have their owner set to ${tty_underline}${USER}${tty_reset}:"
  printf "%s\n" "${chowns[@]}"
fi
if [[ "${#chgrps[@]}" -gt 0 ]]; then
  ohai "The following existing directories will have their group set to ${tty_underline}${GROUP}${tty_reset}:"
  printf "%s\n" "${chgrps[@]}"
fi
if [[ "${#mkdirs[@]}" -gt 0 ]]; then
  ohai "The following new directories will be created:"
  printf "%s\n" "${mkdirs[@]}"
fi

if [[ -d "${BUGSNAG_CLI_PREFIX}" ]]; then
  if [[ "${#chmods[@]}" -gt 0 ]]; then
    execute_sudo "chmod" "u+rwx" "${chmods[@]}"
  fi
  if [[ "${#group_chmods[@]}" -gt 0 ]]; then
    execute_sudo "chmod" "g+rwx" "${group_chmods[@]}"
  fi
  if [[ "${#chowns[@]}" -gt 0 ]]; then
    execute_sudo "chown" "${USER}" "${chowns[@]}"
  fi
  if [[ "${#chgrps[@]}" -gt 0 ]]; then
    execute_sudo "chgrp" "${GROUP}" "${chgrps[@]}"
  fi
else
  execute_sudo "${INSTALL[@]}" "${BUGSNAG_CLI_PREFIX}"
fi

if [[ "${#mkdirs[@]}" -gt 0 ]]; then
  execute_sudo "mkdir" "-p" "${mkdirs[@]}"
  execute_sudo "chmod" "ug=rwx" "${mkdirs[@]}"
  execute_sudo "chown" "${USER}" "${mkdirs[@]}"
  execute_sudo "chgrp" "${GROUP}" "${mkdirs[@]}"
fi

ohai "Downloading and installing Bugsnag CLI..."
(
  cd "${BUGSNAG_CLI_PREFIX}" >/dev/null || return

  DOWNLOAD_URL=$(curl -sS --no-progress-meter ${BUGSNAG_CLI_GIT_REMOTE} |
    grep "${UNAME_MACHINE}-${OS_NAME}-bugsnag-cli*" |
    grep "browser_download_url" |
    cut -d : -f 2,3 |
    tr -d \" |
    xargs)

  execute "curl" "-L" "--no-progress-meter" "${DOWNLOAD_URL}" "-o" "${BUGSNAG_CLI_PREFIX}/bin/bugsnag-cli"

  execute_sudo "chmod" "ug=rwx" "${BUGSNAG_CLI_PREFIX}/bin/bugsnag-cli"

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

ohai "Next steps:"
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

# `which` is a shell function defined above.
# shellcheck disable=SC2230
if [[ "$(which bugsnag-cli)" != "${BUGSNAG_CLI_PREFIX}/bin/bugsnag-cli" ]]; then
  cat <<EOS
- Run these three commands in your terminal to add Bugsnag CLI to your ${tty_bold}PATH${tty_reset}:
    echo "# Bugsnag CLI" >> ${shell_profile}
    echo 'PATH="${BUGSNAG_CLI_PREFIX}/bin:\$PATH"' >> ${shell_profile}
    source "${shell_profile}"
EOS
fi
