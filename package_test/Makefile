# Run all tests; Requires cinc, docker, and vagrant to be installed
.PHONY: test-linux-packages
test-linux-packages: test-linux-package-install test-linux-package-service

# Generate the packages for testing.
.PHONY: build-packages
build-packages:
	cd ..; $(MAKE) release-test

# Tests that the `observiq-collector` service is properly installed and 
# enabled on debian and almalinux (centos). These tests relies on both cinc and vagrant being installed.
.PHONY: test-linux-package-service
test-linux-package-service: build-packages
	./service/run.sh

# Tests that the `observiq-collector` package is installed properly using docker images 
# (debian, ubuntu, fedora, and centos). This test relies on both cinc and docker being installed.
.PHONY: test-linux-package-install
test-linux-package-install: build-packages
	./install/run.sh

