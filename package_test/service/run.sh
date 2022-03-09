#!/bin/sh
set -e

vagrant destroy -f

vagrant up --provision

# TODO: get port info, run cinc auditor

vagrant destroy -f
