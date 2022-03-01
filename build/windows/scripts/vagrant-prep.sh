#!/bin/sh
set -e

vagrant up --provider virtualbox
vagrant winrm -c "setx PATH \"%PATH%;C:/vagrant/wix\;C:/vagrant\""
