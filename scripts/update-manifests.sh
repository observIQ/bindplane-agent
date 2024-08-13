#!/bin/bash

# Check if the correct number of arguments are provided
if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <new_version_opentelemetry_contrib> <new_version_opentelemetry_collector> <new_version_bindplane_agent>"
    exit 1
fi

# Assign arguments to variables
new_version_opentelemetry_contrib="$1"
new_version_opentelemetry_collector="$2"
new_version_bindplane_agent="$3"

# Create a backup of the original file
cp manifests/observIQ/manifest.yaml manifests/observIQ/manifest.yaml.bak

# Use sed to update the version for matching lines
sed -i -E "s|(github\.com/open-telemetry/opentelemetry-collector-contrib/[^\s]*) v[0-9]+\.[0-9]+\.[0-9]+|\1 $new_version_opentelemetry_contrib|g" manifests/observIQ/manifest.yaml
sed -i -E "s|(go\.opentelemetry\.io/collector/[^\s]*) v[0-9]+\.[0-9]+\.[0-9]+|\1 $new_version_opentelemetry_collector|g" manifests/observIQ/manifest.yaml
sed -i -E "s|(github\.com/observiq/bindplane-agent/[^\s]*) v[0-9]+\.[0-9]+\.[0-9]+|\1 $new_version_bindplane_agent|g" manifests/observIQ/manifest.yaml

echo "Versions updated successfully."
