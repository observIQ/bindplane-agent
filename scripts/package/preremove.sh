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

# Determine if we need service or systemctl
if command -v systemctl > /dev/null 2>&1; then
  SVC_PRE=systemctl
elif command -v service > /dev/null 2>&1; then
  SVC_PRE=service
fi

if [ "$SVC_PRE" = "systemctl" ]; then
    info "Stopping service..."
    systemctl stop observiq-otel-collector > /dev/null || { printf 'failed to stop service'; return; }
    succeeded

    info "Disabling service..."
    systemctl disable observiq-otel-collector > /dev/null 2>&1 || { printf 'failed to disable service'; return; }
    succeeded
else
    info "Stopping service..."
    service observiq-otel-collector stop || { printf 'failed to stop service'; return; }
    succeeded

    info "Disabling service..."
    chkconfig observiq-otel-collector on > /dev/null 2>&1 || { printf 'failed to disable service'; return; }
    succeeded
fi
