#!/bin/sh
# Copyright  observIQ, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# shellcheck disable=SC2119
# SC2119 -> Use foo "$@" if function's $1 should mean script's $1.

set -e

# Collector Constants
PACKAGE_NAME="observiq-otel-collector_linux"
DOWNLOAD_BASE="https://github.com/observiq/observiq-otel-collector/releases"

# Script Constants
TMP_DIR=${TMPDIR:-"/tmp"} # Allow this to be overriden by cannonical TMPDIR env var
PREREQS="curl printf systemctl sed uname cut"
SCRIPT_NAME="$0"
INDENT_WIDTH='  '
indent=""

# out_file_path is the full path to the downloaded package (e.g. "/tmp/observiq-otel-collector_linux_amd64.deb")
out_file_path="unknown"

# Colors
num_colors=$(tput colors 2>/dev/null)
if test -n "$num_colors" && test "$num_colors" -ge 8; then
  bold="$(tput bold)"
  underline="$(tput smul)"
  # standout can be bold or reversed colors dependent on terminal
  standout="$(tput smso)"
  reset="$(tput sgr0)"
  bg_black="$(tput setab 0)"
  bg_blue="$(tput setab 4)"
  bg_cyan="$(tput setab 6)"
  bg_green="$(tput setab 2)"
  bg_magenta="$(tput setab 5)"
  bg_red="$(tput setab 1)"
  bg_white="$(tput setab 7)"
  bg_yellow="$(tput setab 3)"
  fg_black="$(tput setaf 0)"
  fg_blue="$(tput setaf 4)"
  fg_cyan="$(tput setaf 6)"
  fg_green="$(tput setaf 2)"
  fg_magenta="$(tput setaf 5)"
  fg_red="$(tput setaf 1)"
  fg_white="$(tput setaf 7)"
  fg_yellow="$(tput setaf 3)"
fi

if [ -z "$reset" ]; then
  sed_ignore=''
else
  sed_ignore="/^[$reset]+$/!"
fi

# Helper Functions
printf() {
  if command -v sed >/dev/null; then
    command printf -- "$@" | sed -E "$sed_ignore s/^/$indent/g"  # Ignore sole reset characters if defined
  else
    # Ignore $* suggestion as this breaks the output
    # shellcheck disable=SC2145
    command printf -- "$indent$@"
  fi
}

increase_indent() { indent="$INDENT_WIDTH$indent" ; }
decrease_indent() { indent="${indent#*$INDENT_WIDTH}" ; }

# Color functions reset only when given an argument
bold() { command printf "$bold$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
underline() { command printf "$underline$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
standout() { command printf "$standout$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
# Ignore "parameters are never passed"
# shellcheck disable=SC2120
reset() { command printf "$reset$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
bg_black() { command printf "$bg_black$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
bg_blue() { command printf "$bg_blue$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
bg_cyan() { command printf "$bg_cyan$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
bg_green() { command printf "$bg_green$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
bg_magenta() { command printf "$bg_magenta$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
bg_red() { command printf "$bg_red$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
bg_white() { command printf "$bg_white$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
bg_yellow() { command printf "$bg_yellow$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
fg_black() { command printf "$fg_black$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
fg_blue() { command printf "$fg_blue$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
fg_cyan() { command printf "$fg_cyan$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
fg_green() { command printf "$fg_green$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
fg_magenta() { command printf "$fg_magenta$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
fg_red() { command printf "$fg_red$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
fg_white() { command printf "$fg_white$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }
fg_yellow() { command printf "$fg_yellow$*$(if [ -n "$1" ]; then command printf "$reset"; fi)" ; }

# Intentionally using variables in format string
# shellcheck disable=SC2059
info() { printf "$*\\n" ; }
# Intentionally using variables in format string
# shellcheck disable=SC2059
warn() {
  increase_indent
  printf "$fg_yellow$*$reset\\n"
  decrease_indent
}
# Intentionally using variables in format string
# shellcheck disable=SC2059
error() {
  increase_indent
  printf "$fg_red$*$reset\\n"
  decrease_indent
}
# Intentionally using variables in format string
# shellcheck disable=SC2059
success() { printf "$fg_green$*$reset\\n" ; }
# Ignore 'arguments are never passed'
# shellcheck disable=SC2120
prompt() {
  if [ "$1" = 'n' ]; then
    command printf "y/$(fg_red '[n]'): "
  else
    command printf "$(fg_green '[y]')/n: "
  fi
}

observiq_banner()
{
  fg_cyan "           888                                        8888888 .d88888b.\\n"
  fg_cyan "           888                                          888  d88P\" \"Y88b\\n"
  fg_cyan "           888                                          888  888     888\\n"
  fg_cyan "   .d88b.  88888b.  .d8888b   .d88b.  888d888 888  888  888  888     888\\n"
  fg_cyan "  d88\"\"88b 888 \"88b 88K      d8P  Y8b 888P\"   888  888  888  888     888\\n"
  fg_cyan "  888  888 888  888 \"Y8888b. 88888888 888     Y88  88P  888  888 Y8b 888\\n"
  fg_cyan "  Y88..88P 888 d88P      X88 Y8b.     888      Y8bd8P   888  Y88b.Y8b88P\\n"
  fg_cyan "   \"Y88P\"  88888P\"   88888P'  \"Y8888  888       Y88P  8888888 \"Y888888\"\\n"
  fg_cyan "                                                                   Y8b  \\n"

  reset
}

separator() { printf "===================================================\\n" ; }

banner()
{
  printf "\\n"
  separator
  printf "| %s\\n" "$*" ;
  separator
}

usage()
{
  increase_indent
  USAGE=$(cat <<EOF
Usage:
  $(fg_yellow '-v, --version')
      Defines the version of the observIQ OpenTelemetry Collector.
      If not provided, this will default to the latest version.
      Alternatively the COLLECTOR_VERSION environment variable can be
      set to configure the collector version.
      Example: '-v 1.2.12' will download 1.2.12.

  $(fg_yellow '-l, --url')
      Defines the URL that the components will be downloaded from.
      If not provided, this will default to observIQ OpenTelemetry Collector\'s GitHub releases.
      Example: '-l http://my.domain.org/observiq-otel-collector' will download from there.

  $(fg_yellow '-x, --proxy')
      Defines the proxy server to be used for communication by the install script.
      Example: $(fg_blue -x) $(fg_magenta http\(s\)://server-ip:port/).

  $(fg_yellow '-U, --proxy-user')
      Defines the proxy user to be used for communication by the install script.

  $(fg_yellow '-P, --proxy-password')
      Defines the proxy password to be used for communication by the install script.

EOF
  )
  info "$USAGE"
  decrease_indent
  return 0
}

force_exit()
{
  # Exit regardless of subshell level with no "Terminated" message
  kill -PIPE $$
  # Call exit to handle special circumstances (like running script during docker container build)
  exit 1
}

error_exit()
{
  line_num=$(if [ -n "$1" ]; then command printf ":$1"; fi)
  error "ERROR ($SCRIPT_NAME$line_num): ${2:-Unknown Error}" >&2
  if [ -n "$0" ]; then
    increase_indent
    error "$*"
    decrease_indent
  fi
  force_exit
}

print_prereq_line()
{
  if [ -n "$2" ]; then
    command printf "\\n${indent}  - "
    command printf "[$1]: $2"
  fi
}

check_failure()
{
  if [ "$indent" != '' ]; then increase_indent; fi
  command printf "${indent}${fg_red}ERROR: %s check failed!${reset}" "$1"

  print_prereq_line "Issue" "$2"
  print_prereq_line "Resolution" "$3"
  print_prereq_line "Help Link" "$4"
  print_prereq_line "Rerun" "$5"

  command printf "\\n"
  if [ "$indent" != '' ]; then decrease_indent; fi
  force_exit
}

succeeded()
{
  increase_indent
  success "Succeeded!"
  decrease_indent
}

failed()
{
  error "Failed!"
}

# This will set all installation variables
# at the beginning of the script.
setup_installation()
{
    banner "Configuring Installation Variables"
    increase_indent

    # Installation variables
    set_os_arch
    set_package_type
    set_download_urls
    set_proxy
    set_file_name

    success "Configuration complete!"
    decrease_indent
}

set_file_name() {
    package_file_name="${PACKAGE_NAME}_${arch}.${package_type}"
    out_file_path="$TMP_DIR/$package_file_name"
}

set_proxy()
{
  if [ -n "$proxy" ]; then
    info "Using proxy from arguments: $proxy"
    if [ -n "$proxy_user" ]; then
      while [ -z "$proxy_password" ] ; do
        increase_indent
        command printf "${indent}$(fg_blue "$proxy_user@$proxy")'s password: "
        stty -echo
        read -r proxy_password
        stty echo
        info
        if [ -z "$proxy_password" ]; then
          warn "The password must be provided!"
        fi
        decrease_indent
      done
      protocol="$(echo "$proxy" | cut -d'/' -f1)"
      host="$(echo "$proxy" | cut -d'/' -f3)"
      full_proxy="$protocol//$proxy_user:$proxy_password@$host"
    fi
  fi

  if [ -z "$full_proxy" ]; then
    full_proxy="$proxy"
  fi
}


set_os_arch()
{
  os_arch=$(uname -m)
  case "$os_arch" in 
    # arm64 strings. These are from https://stackoverflow.com/questions/45125516/possible-values-for-uname-m
    aarch64|arm64|aarch64_be|armv8b|armv8l)
      os_arch="arm64"
      ;;
    x86_64)
      os_arch="amd64"
      ;;
    # armv6/32bit. These are what raspberry pi can return, which is the main reason we support 32-bit arm
    arm|armv6l|armv7l)
      os_arch="arm"
      ;;
    *)
      error_exit "$LINENO" "Unsupported os arch: $os_arch"
      ;;
  esac
}

# Set the package type before install
set_package_type()
{
  if command -v dpkg > /dev/null 2>&1; then
      package_type="deb"
  elif command -v rpm > /dev/null 2>&1; then
      package_type="rpm"
  else
      error_exit "$LINENO" "Could not find dpkg or rpm on the system"
  fi
}

# This will set the urls to use when downloading the collector and its plugins.
# These urls are constructed based on the --version flag or COLLECTOR_VERSION env variable.
# If not specified, the version defaults to "latest".
set_download_urls()
{
  if [ -z "$version" ] ; then
    # shellcheck disable=SC2153
    version=$COLLECTOR_VERSION
  fi

  if [ -z "$url" ] ; then
    url=$DOWNLOAD_BASE
  fi

  if [ -z "$version" ] ; then
    collector_download_url="$url/latest/download/${PACKAGE_NAME}_${os_arch}.${package_type}"
  else
    collector_download_url="$url/download/v$version/${PACKAGE_NAME}_${os_arch}.${package_type}"
  fi
}

# This will check all prerequisites before running an installation.
check_prereqs()
{
  banner "Checking Prerequisites"
  increase_indent
  root_check
  os_check
  os_arch_check
  package_type_check
  dependencies_check
  success "Prerequisite check complete!"
  decrease_indent
}

# This
root_check()
{
  system_user_name=$(id -un)
  if [ "${system_user_name}" != 'root' ]
  then
    failed
    error_exit "$LINENO" "Script needs to be run as root or with sudo"
  fi
}

# This will check if the operating system is supported.
os_check()
{
  info "Checking that the operating system is supported..."
  os_type=$(uname -s)
  case "$os_type" in
    Linux)
      succeeded
      ;;
    *)
      failed
      error_exit "$LINENO" "The operating system $(fg_yellow "$os_type") is not supported by this script."
      ;;
  esac
}

# This will check if the system architecture is supported.
os_arch_check()
{
  info "Checking for valid operating system architecture..."
  arch=$(uname -m)
  case "$arch" in 
    x86_64|aarch64|arm64|aarch64_be|armv8b|armv8l|arm|armv6l|armv7l)
      succeeded
      ;;
    *)
      failed
      error_exit "$LINENO" "The operating system architecture $(fg_yellow "$arch") is not supported by this script."
      ;;
  esac
}

# This will check if the current environment has
# all required shell dependencies to run the installation.
dependencies_check()
{
  info "Checking for script dependencies..."
  FAILED_PREREQS=''
  for prerequisite in $PREREQS; do
    if command -v "$prerequisite" >/dev/null; then
      continue
    else
      if [ -z "$FAILED_PREREQS" ]; then
        FAILED_PREREQS="${fg_red}$prerequisite${reset}"
      else
        FAILED_PREREQS="$FAILED_PREREQS, ${fg_red}$prerequisite${reset}"
      fi
    fi
  done

  if [ -n "$FAILED_PREREQS" ]; then
    failed
    error_exit "$LINENO" "The following dependencies are required by this script: [$FAILED_PREREQS]"
  fi
  succeeded
}

# This will check to ensure either dpkg or rpm is installedon the system
package_type_check()
{
  info "Checking for package manager..."
  if command -v dpkg > /dev/null 2>&1; then
      succeeded
  elif command -v rpm > /dev/null 2>&1; then
      succeeded
  else
      failed
      error_exit "$LINENO" "Could not find dpkg or rpm on the system"
  fi
}

# This will install the package by downloading the archived collector,
# extracting the binaries, and then removing the archive.
install_package()
{
  banner "Installing observIQ OpenTelemetry Collector"
  increase_indent

  proxy_args=""
  if [ -n "$proxy" ]; then
    proxy_args="-x $proxy"
    if [ -n "$proxy_user" ]; then
      proxy_args="$proxy_args -U $proxy_user:$proxy_password"
    fi
  fi

  if [ -n "$proxy" ]; then
    info "Downloading package using proxy..."
  fi 

  info "Downloading package..."
  eval curl -L "$proxy_args" "$collector_download_url" -o "$out_file_path" --progress-bar --fail || error_exit "$LINENO" "Failed to download package"
  succeeded

  info "Installing package..."
  unpack_package || error_exit "$LINENO" "Failed to extract package"
  succeeded

  info "Enabling service..."
  systemctl enable --now observiq-otel-collector > /dev/null 2>&1 || error_exit "$LINENO" "Failed to enable service"
  succeeded

  success "observIQ OpenTelemetry Collector installation complete!"
  decrease_indent
}

unpack_package()
{
  case "$package_type" in
    deb)
      dpkg -i "$out_file_path" > /dev/null
      ;;
    rpm)
      rpm -U "$out_file_path" > /dev/null
      ;;
    *)
      error "Unrecognized package type"
      return 1
      ;;
  esac
  return 0
}

# This will display the results of an installation
display_results()
{
    banner 'Information'
    increase_indent
    info "Collector Home:     $(fg_cyan "/opt/observiq-otel-collector")$(reset)"
    info "Collector Config:   $(fg_cyan "/opt/observiq-otel-collector/config.yaml")$(reset)"
    info "Start Command:      $(fg_cyan "sudo systemctl start observiq-otel-collector")$(reset)"
    info "Stop Command:       $(fg_cyan "sudo systemctl stop observiq-otel-collector")$(reset)"
    info "Logs Command:       $(fg_cyan "sudo journalctl -u observiq-otel-collector.service")$(reset)"
    decrease_indent

    banner 'Support'
    increase_indent
    info "For more information on configuring the collector, see the docs:"
    increase_indent
    info "$(fg_cyan "https://github.com/observiq/observiq-otel-collector/tree/main#observiq-opentelemetry-collector")$(reset)"
    decrease_indent
    info "If you have any other questions please contact us at $(fg_cyan support@observiq.com)$(reset)"
    increase_indent
    decrease_indent
    decrease_indent

    banner "$(fg_green Installation Complete!)"
    return 0
}

uninstall_package()
{
  case "$package_type" in
    deb)
      dpkg -r "observiq-otel-collector" > /dev/null 2>&1
      ;;
    rpm)
      rpm -e "observiq-otel-collector" > /dev/null 2>&1
      ;;
    *)
      error "Unrecognized package type"
      return 1
      ;;
  esac
  return 0
}

uninstall()
{
  observiq_banner

  set_package_type
  banner "Uninstalling observIQ OpenTelemetry Collector"
  increase_indent

  info "Checking permissions..."
  root_check
  succeeded

  info "Stopping service..."
  systemctl stop observiq-otel-collector > /dev/null || error_exit "$LINENO" "Failed to stop service"
  succeeded

  info "Disabling service..."
  systemctl disable observiq-otel-collector > /dev/null 2>&1 || error_exit "$LINENO" "Failed to disable service"
  succeeded

  info "Removing package..."
  uninstall_package || error_exit "$LINENO" "Failed to remove package"
  succeeded
  decrease_indent

  banner "$(fg_green Uninstallation Complete!)"
}

main()
{
  if [ $# -ge 1 ]; then
    while [ -n "$1" ]; do
      case "$1" in
        -v|--version)
          version=$2 ; shift 2 ;;
        -l|--url)
          url=$2 ; shift 2 ;;
        -x|--proxy)
          proxy=$2 ; shift 2 ;;
        -U|--proxy-user)
          proxy_user=$2 ; shift 2 ;;
        -P|--proxy-password)
          proxy_password=$2 ; shift 2 ;;
        -r|--uninstall)
          uninstall
          force_exit
          ;;
        -h|--help)
          usage
          force_exit
          ;;
      --)
        shift; break ;;
      *)
        error "Invalid argument: $1"
        usage
        force_exit
        ;;
      esac
    done
  fi

  observiq_banner
  check_prereqs
  setup_installation
  install_package
  display_results
}

main "$@"
