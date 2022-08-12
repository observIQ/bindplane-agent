#!/bin/sh

TARGET_VERSION=$1
if [ -z "$TARGET_VERSION" ]; then
    echo "Must specify a target version"
    exit 1
fi

LOCAL_MODULES=$(find . -type f -name "go.mod" -exec dirname {} \; | sort)
for local_mod in $LOCAL_MODULES
do
    # Run in a subshell so that the CD doesn't change this shell's current directory
    (
        echo "Updating deps in $local_mod"
        cd "$local_mod" || exit 1
        go mod tidy -compat=1.18
        OTEL_MODULES=$(go list -m -f '{{if not (or .Indirect .Main)}}{{.Path}}{{end}}' all |
            grep -E -e '(?:^github.com/open-telemetry/opentelemetry-collector-contrib)|(?:^go.opentelemetry.io/collector)')

        for mod in $OTEL_MODULES
        do
            echo "$local_mod: $mod@$TARGET_VERSION"
            go mod edit -require "$mod@$TARGET_VERSION"
        done

        go mod tidy -compat=1.18
    )
done


