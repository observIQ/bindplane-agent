#!/bin/sh
set -e

BASEDIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

# We execute these commands in subshells so that we stay in our CWD,
# in case the script exits unexpectedly
(cd "$BASEDIR"; vagrant destroy -f)

(cd "$BASEDIR"; vagrant up --provision)

# Gets host names + port for ssh; Format is:
# Host <hostname>
# Port <port>
# ...
# Which will repeat, once for every VM spun up.
HOST_PORTS=$(cd "$BASEDIR"; vagrant ssh-config | grep --color=never -E "(Port )|(Host )" | sed -E 's/^ +//g')

HOST=""
# Set separator in for-loop below to only be newlines
IFS='
'
for l in $HOST_PORTS; do
    case "$l" in 
        Port*)
            PORT=$(printf "%s" "$l" | sed -E 's/Port //g')
            echo "Running cinc auditor over ssh on host $HOST, port $PORT"
            (cd "$BASEDIR"; cinc-auditor exec ./integration.rb -t "ssh://vagrant@localhost:$PORT")
            ;;
        Host*)
            HOST=$(printf "%s" "$l" | sed -E 's/Host //g') 
            ;;
    esac
done

(cd "$BASEDIR"; vagrant destroy -f)
