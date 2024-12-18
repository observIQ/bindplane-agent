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

# This is the list of stable (v1.0.0+) modules in the core repository.
read -d '' STABLE_MODULES << EOF
go.opentelemetry.io/collector/pdata
go.opentelemetry.io/collector/config/configtls
go.opentelemetry.io/collector/config/configretry
go.opentelemetry.io/collector/confmap
go.opentelemetry.io/collector/config/configopaque
go.opentelemetry.io/collector/confmap/provider/envprovider
go.opentelemetry.io/collector/confmap/provider/fileprovider
go.opentelemetry.io/collector/confmap/provider/httpsprovider
go.opentelemetry.io/collector/confmap/provider/yamlprovider
go.opentelemetry.io/collector/config/confignet
go.opentelemetry.io/collector/consumer
EOF

is_stable_module() {
    for stable_mod in $STABLE_MODULES; do
        if [ "$stable_mod" = "$1" ]; then
            return 0
        fi
    done
    return 1
}

is_contrib_module() {
    case $1 in
    github.com/open-telemetry/opentelemetry-collector-contrib*)
        return 0;;
    esac
    return 1
}

TARGET_VERSION=$1
if [ -z "$TARGET_VERSION" ]; then
    echo "Must specify a target version"
    exit 1
fi

CONTRIB_TARGET_VERSION=$2
if [ -z "$CONTRIB_TARGET_VERSION" ]; then
    echo "Must specify a target contrib version"
    exit 1
fi

PDATA_TARGET_VERSION=$3

if [ -z "$PDATA_TARGET_VERSION" ]; then
    echo "Must specify a target pdata version"
    exit 1
fi

LOCAL_MODULES=$(find . -type f -name "go.mod" -exec dirname {} \; | sort)
for local_mod in $LOCAL_MODULES
do
    # Run in a subshell so that the CD doesn't change this shell's current directory
    (
        echo "Updating deps in $local_mod"
        cd "$local_mod" || exit 1
        # go list will not work if module is not tidy, so we tidy first
        go mod tidy -compat=1.22

        OTEL_MODULES=$(go list -m -f '{{if not (or .Indirect .Main)}}{{.Path}}{{end}}' all |
            grep -E -e '(?:^github.com/open-telemetry/opentelemetry-collector-contrib)|(?:^go.opentelemetry.io/collector)')

        for mod in $OTEL_MODULES; do
            if is_stable_module "$mod"; then
                echo "$local_mod: $mod@$PDATA_TARGET_VERSION"
                go mod edit -require "$mod@$PDATA_TARGET_VERSION"
            elif is_contrib_module "$mode"; then
                echo "$local_mod: $mod@$CONTRIB_TARGET_VERSION"
                go mod edit -require "$mod@$CONTRIB_TARGET_VERSION"
            else
                echo "$local_mod: $mod@$TARGET_VERSION"
                go mod edit -require "$mod@$TARGET_VERSION"
            fi
        done
    )
done
