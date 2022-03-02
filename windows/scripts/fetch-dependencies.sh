#!/bin/sh
set -e
BASEDIR="$(dirname "$(realpath "$0")")"
PROJECT_BASE="$BASEDIR/../.."

[ -f go-msi.exe ] || curl -f -L -o go-msi.exe https://github.com/observIQ/go-msi/releases/download/v2.0.0/go-msi.exe
[ -f ./cinc-auditor.msi ] || curl -f -L -o cinc-auditor.msi http://downloads.cinc.sh/files/stable/cinc-auditor/4.17.7/windows/2012r2/cinc-auditor-4.17.7-1-x64.msi

[ -f ./wix-binaries.zip ] || curl -f -L -o wix-binaries.zip https://github.com/wixtoolset/wix3/releases/download/wix3112rtm/wix311-binaries.zip

mkdir -p wix
[ -d wix/sdk ] || unzip -o wix-binaries.zip -d wix

$PROJECT_BASE/buildscripts/download-dependencies.sh

cp "$PROJECT_BASE/config/example.yaml" "./config.yaml"
