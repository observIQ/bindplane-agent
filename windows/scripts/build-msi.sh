#!/bin/sh
set -e
BASEDIR="$(dirname "$(realpath "$0")")"
PROJECT_BASE="$BASEDIR/../.."

cp "$PROJECT_BASE/dist/collector_windows_amd64.exe" "observiq-collector.exe"

vagrant winrm -c \
    "cd C:/vagrant; go-msi.exe make -m observiq-collector.msi --version v0.0.1 --arch amd64"
