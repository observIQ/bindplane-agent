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

uninstall_package()
{
    if command -v dpkg > /dev/null 2>&1; then
        dpkg -r "observiq-otel-collector" > /dev/null 2>&1
    elif command -v rpm > /dev/null 2>&1; then
        rpm -e "observiq-otel-collector" > /dev/null 2>&1
    else
        printf "Could not find dpkg or rpm on the system"
    fi
}

main()
{
    # Removes whole install directory
    info "Removing installed artifacts..."
    # This find command gets a list of all artifacts paths except the root directory.
    FILES=$(cd "/opt/observiq-otel-collector"; find "." -not \( -name "." \))
    for f in $FILES
    do
        rm -rf "/opt/observiq-otel-collector/$f" || {printf "Failed to remove artifact /opt/observiq-otel-collector/$f"; return;}
    done
    succeeded

    info "Removing package..."
    uninstall_package()
    succeeded
}