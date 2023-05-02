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
        OTEL_MODULES=$(go list -m -f '{{if not (or .Indirect .Main)}}{{.Path}}{{end}}' all |
            grep -E -e '(?:^github.com/open-telemetry/opentelemetry-collector-contrib)|(?:^go.opentelemetry.io/collector)')

        for mod in $OTEL_MODULES
        do
            if case $mod in go.opentelemetry.io/collector/pdata*) ;; *) false;; esac; then
                # pdata package is versioned separately
                echo "$local_mod: $mod@$PDATA_TARGET_VERSION"
                go mod edit -require "$mod@$PDATA_TARGET_VERSION"
            elif case $mod in github.com/open-telemetry/opentelemetry-collector-contrib*) ;; *) false;; esac; then
                echo "$local_mod: $mod@$CONTRIB_TARGET_VERSION"
                go mod edit -require "$mod@$CONTRIB_TARGET_VERSION"
            else
                echo "$local_mod: $mod@$TARGET_VERSION"
                go mod edit -require "$mod@$TARGET_VERSION"
            fi;
        done
    )
done
