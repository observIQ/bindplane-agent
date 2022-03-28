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

set -e

# Constants
DOWNLOAD_URL_BASE="https://github.com/observiq/observiq-otel-collector/releases"
TMP_DIR=${TMPDIR:-"/tmp"} # Allow this to be overriden by cannonical TMPDIR env var

# Global variables

# arch is the architecture of the system.
# May be "arm64", "amd64", or "arm"
arch="unknown"

# package type is the type of package we want to install.
# may be "deb" or "rpm"
package_type="unknown"

# package_file_name is the name of the package file (e.g. "observiq-otel-collector_linux_amd64.deb")
package_file_name="unknown"

# out_file_path is the full path to the downloaded package (e.g. "/tmp/observiq-otel-collector_linux_amd64.deb")
out_file_path="unknown"

set_arch() {
    os_arch=$(uname -m)
    case "$os_arch" in
    # arm64 strings. These are from https://stackoverflow.com/questions/45125516/possible-values-for-uname-m
    aarch64|arm64|aarch64_be|armv8b|armv8l)
        arch="arm64"
        ;;
    x86_64)
        arch="amd64"
        ;;
    # armv6/32bit. These are what raspberry pi can return, which is the main reason we support 32-bit arm
    arm|armv6l|armv7l)
        arch="arm"
        ;;
    *)
        echo "Installation failed: Unsupported os arch: $os_arch"
        exit 1
        ;;
    esac
}

set_package_type() {
    if command -v dpkg > /dev/null 2>&1; then
        package_type="deb"
    elif command -v rpm > /dev/null 2>&1; then
        package_type="rpm"
    else
        echo "Installation failed: Could not find dpkg or rpm on the system"
        exit 1
    fi
}

set_file_name() {
    package_file_name="observiq-otel-collector_linux_${arch}.${package_type}"
    out_file_path="$TMP_DIR/$package_file_name"
}

download_package() {
    download_url="${DOWNLOAD_URL_BASE}/latest/download/${package_file_name}"
    
    echo "Downloading ${download_url} to ${out_file_path}..."
    curl -L "$download_url" -o "$out_file_path" --progress-bar --fail
    echo "${out_file_path} successfully downloaded!"
}

package_version() {
    case "$package_type" in
        deb)
            dpkg -s observiq-otel-collector | grep -i '^Version:' | cut -d' ' -f2
            ;;
        rpm)
            rpm -q --queryformat '%{VERSION}' "observiq-otel-collector"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

install_package() {
    case "$package_type" in
        deb)
            dpkg -i "$out_file_path"
            ;;
        rpm)
            # -U installs OR upgrades
            # just -i will fail if previous versions exists
            rpm -U "$out_file_path"
            ;;
        *)
            echo "Installation failed: Could not determine the package manager to use for installation (package_type: ${package_type})"
            exit 1
            ;;
    esac

    if command -v systemctl > /dev/null 2>&1; then
        systemctl enable --now observiq-otel-collector
    fi
}

uninstall_package() {
    if command -v dpkg > /dev/null 2>&1; then
        if dpkg -l "observiq-otel-collector" > /dev/null 2>&1; then
            desired_state="$(dpkg -l observiq-otel-collector | tail -n1 | cut -b 1)"
            if [ "$desired_state" = "i" ]; then
                echo "Existing DEB installation detected, uninstalling..."
                dpkg -r "observiq-otel-collector"
            else
                # The package is already uninstalled;
                # This case can be hit when the package was uninstalled, but config files are left behind,
                # which dpkg still keeps track of.
                echo "dpkg was detected, but no observiq-otel-collector package is currently installed, skipping..."
            fi
        else
            echo "dpkg was detected, but no observiq-otel-collector package is currently installed, skipping..."
        fi
    fi

    if command -v rpm > /dev/null 2>&1; then
        if rpm -q observiq-otel-collector > /dev/null 2>&1; then
            echo "Existing RPM installation detected, uninstalling..."
            rpm -e "observiq-otel-collector"
        else
            echo "rpm was detected, but no observiq-otel-collector package is currently installed, skipping..."
        fi
    fi
}

main() {
    # Determine system parameters
    set_arch
    set_package_type
    set_file_name

    COMMAND=${1:-"install"}

    case "$COMMAND" in
        install)
            echo "Installing the observIQ OpenTelemetry collector..."
            # Download and install the package
            download_package
            install_package
            echo ""
            echo ""
            echo "Successfully installed version $(package_version) of the observIQ OpenTelemetry collector."
            ;;
        uninstall)
            echo "Uninstalling the observIQ OpenTelemetry collector..."
            uninstall_package
            echo "Successfully uninstalled the observIQ OpenTelemetry collector."
            ;;
        *)
            echo "Unrecognized command: $COMMAND"
            echo "Usage: $0 [install|uninstall]"
            exit 1
            ;;
    esac
}

main "$@"
