#!/bin/sh

# This script will start a container with the root of the repo mounted to
# /app. This allows the built deb / rpm package to be installed. Inspec
# tests are used to verify that the install is correct.

set -e

clean () {
    docker rm -f centos-rpm || true
    docker rm -f fedora-rpm || true
    docker rm -f debian-deb || true
    docker rm -f ubuntu-deb || true
    rm -rf package_test/install/tmp || true
}

# Run a docker container with the working directory mounted to /app
# trigger inspec tests after installing the package.
test_build() {
    name="$1"
    image="$2"

    docker run -d \
        --name "$name" \
        -v "$(pwd):/app" \
        "$image" sleep 3000
    docker exec "$name" bash /app/package_test/install/tmp/install.sh
    cinc-auditor exec package_test/install/integration.rb -t "docker://${name}"
}

test_rpm() {
    echo "yum install -y /app/dist/observiq-collector_*_linux_amd64.rpm" > package_test/install/tmp/install.sh
    test_build "centos-rpm" "gcr.io/gcp-runtimes/centos7@sha256:0f2ee375a95d9eccda1d18506f8d5acd41c7b60901462cfe66f2a72f6d883626"
    test_build "fedora-rpm" "fedora:latest"
}

test_deb() {
    echo "apt-get install -y -f /app/dist/observiq-collector_*_linux_amd64.deb" > package_test/install/tmp/install.sh
    test_build "debian-deb" "gcr.io/google-appengine/debian10@sha256:e58eb64abddb851a5534006fce66aa8a143b69856d041f8e2acdae07c480e9bb"
    test_build "ubuntu-deb" "ubuntu:latest"
}

clean
mkdir package_test/install/tmp
test_rpm
test_deb
clean
