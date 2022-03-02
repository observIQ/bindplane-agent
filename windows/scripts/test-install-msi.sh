#!/bin/sh
set -e

vagrant winrm -c "cd C:/vagrant; msiexec.exe /i cinc-auditor.msi /qn"
vagrant winrm -c "cd C:/vagrant; msiexec.exe /i observiq-collector.msi /qn"
sleep 10
vagrant winrm -c "cinc-auditor exec C:\vagrant\test\install.rb"
