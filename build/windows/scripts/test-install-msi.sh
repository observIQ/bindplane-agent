#!/bin/sh
set -e

vagrant winrm -c "C:/vagrant/cinc-auditor.msi"
vagrant winrm -c "C:/vagrant/observiq-collector.msi"
sleep 10
vagrant winrm -c "cinc-auditor exec C:\vagrant\test\install.rb"
