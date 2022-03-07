#!/bin/sh

set -e

username="observiq-collector"

if id "$username" > /dev/null 2>&1; then
    # Skip all user config if already exists
    echo "User ${username} already exists"
    exit 0
fi
  
useradd --shell /sbin/nologin --system "$username"

