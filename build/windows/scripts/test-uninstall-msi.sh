#!/bin/bash
set -e

vagrant winrm -c "C:/vagrant/cinc-auditor.msi"
vagrant winrm -c "msiexec.exe /x C:\vagrant\observiq-collector.msi /q"
sleep 10
vagrant winrm -c "cinc-auditor exec C:\vagrant\test\uninstall.rb"
