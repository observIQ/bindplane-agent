#!/bin/sh

TARGET_VERSION=$1
if [ -z "$TARGET_VERSION" ]; then
    echo "Must specify a target version"
    exit 1
fi

read -r -d '' DOC_FILES << EOF
docs/processors.md
docs/extensions.md
docs/exporters.md
docs/receivers.md
exporter/googlecloudexporter/README.md
EOF

for doc in $DOC_FILES
do
    echo "$doc"
    # Point contrib links to new version
    sed -i '' -Ee \
        "s|https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/v[^/]*|https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/$TARGET_VERSION|" \
        "$doc"
    # Point core links to new version
    sed -i '' -Ee \
        "s|https://github.com/open-telemetry/opentelemetry-collector/blob/v[^/]*|https://github.com/open-telemetry/opentelemetry-collector/blob/$TARGET_VERSION|" \
        "$doc"
done

