#!/bin/sh
# This script downloads common dependencies (plugins + the JMX receiver jar)
# Into the directory specified by the first argument. If no argument is specified,
# the dependencies will be downloaded into the current working directory.
set -e

DOWNLOAD_DIR="${1:-.}"
START_DIR=$(pwd)
BASEDIR="$(dirname "$(realpath "$0")")"
PROJECT_BASE="$BASEDIR/.."

curl -fL -o "$DOWNLOAD_DIR/opentelemetry-java-contrib-jmx-metrics.jar" \
    "https://github.com/open-telemetry/opentelemetry-java-contrib/releases/download/$(cat "$PROJECT_BASE/JAVA_CONTRIB_VERSION")/opentelemetry-jmx-metrics.jar"


if [ ! -d "$DOWNLOAD_DIR/stanza-plugins" ]; then
    git clone git@github.com:observIQ/stanza-plugins.git "$DOWNLOAD_DIR/stanza-plugins"
fi

cd "$DOWNLOAD_DIR/stanza-plugins"
git fetch --all
git checkout "$(cat "$PROJECT_BASE/PLUGINS_VERSION")"
cd "$START_DIR"

cp -r "$DOWNLOAD_DIR/stanza-plugins/plugins" "$DOWNLOAD_DIR"
