#!/bin/sh
# Copyright  observIQ, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


set -e

manage_systemd_service() {
  # Ensure sysv script isn't present, and if it is remove it
  if [ -f /etc/init.d/observiq-otel-collector ]; then
    rm -f /etc/init.d/observiq-otel-collector
  fi
  # Ensure aix env vars file isn't present, and if it is remove it
  if [ -f /opt/observiq-otel-collector/observiq-otel-collector.aix.env ]; then
    rm -f /opt/observiq-otel-collector/observiq-otel-collector.aix.env
  fi

  systemctl daemon-reload

  echo "configured systemd service"

  cat << EOF

The "observiq-otel-collector" service has been configured!

The collector's config file can be found here: 
  /opt/observiq-otel-collector/config.yaml

To view logs from the collector, run:
  sudo tail -F /opt/observiq-otel-collector/log/collector.log

For more information on configuring the collector, see the docs:
  https://github.com/observIQ/bindplane-agent/tree/main#observiq-opentelemetry-collector

To stop the observiq-otel-collector service, run:
  sudo systemctl stop observiq-otel-collector

To start the observiq-otel-collector service, run:
  sudo systemctl start observiq-otel-collector

To restart the observiq-otel-collector service, run:
  sudo systemctl restart observiq-otel-collector

To enable the service on startup, run:
  sudo systemctl enable observiq-otel-collector

If you have any other questions please contact us at support@observiq.com
EOF
}

manage_sysv_service() {
  # Ensure systemd script isn't present, and if it is remove it
  if [ -f /usr/lib/systemd/system/observiq-otel-collector.service ]; then
    rm -f /usr/lib/systemd/system/observiq-otel-collector.service
  fi
  # Ensure aix env vars file isn't present, and if it is remove it
  if [ -f /opt/observiq-otel-collector/observiq-otel-collector.aix.env ]; then
    rm -f /opt/observiq-otel-collector/observiq-otel-collector.aix.env
  fi

  chmod 755 /etc/init.d/observiq-otel-collector
  chmod 644 /etc/sysconfig/observiq-otel-collector
  echo "configured sysv service"

  cat << EOF

The "observiq-otel-collector" service has been configured!

The collector's config file can be found here: 
  /opt/observiq-otel-collector/config.yaml

To view logs from the collector, run:
  sudo tail -F /opt/observiq-otel-collector/log/collector.log

For more information on configuring the collector, see the docs:
  https://github.com/observIQ/bindplane-agent/tree/main#observiq-opentelemetry-collector

To stop the observiq-otel-collector service, run:
  sudo service observiq-otel-collector stop

To start the observiq-otel-collector service, run:
  sudo service observiq-otel-collector start

To restart the observiq-otel-collector service, run:
  sudo service observiq-otel-collector restart

To enable the service on startup, run:
  sudo service observiq-otel-collector enable

If you have any other questions please contact us at support@observiq.com
EOF
}

manage_aix_service() {
  # Ensure sysv script isn't present, and if it is remove it
  if [ -f /etc/init.d/observiq-otel-collector ]; then
    rm -f /etc/init.d/observiq-otel-collector
  fi
  # Ensure systemd script isn't present, and if it is remove it
  if [ -f /usr/lib/systemd/system/observiq-otel-collector.service ]; then
    rm -f /usr/lib/systemd/system/observiq-otel-collector.service
  fi

  # Add the service, removing it if it already exists in order
  # to make sure we have the most recent version
  if lssrc -s observiq-otel-collector > /dev/null 2>&1; then
    rmssys -s observiq-otel-collector
  else
    mkssys -s observiq-otel-collector -p /opt/observiq-otel-collector/observiq-otel-collector -u $(id -u observiq-otel-collector) -S -n15 -f9 -a '--config config.yaml --manager manager.yaml --logging logging.yaml'
  fi

  # Install the service to start on boot
  # Removing it if it exists, in order to have the most recent version
  if lsitab oiqcollector > /dev/null 2>&1; then
    rmitab oiqcollector
  else
    mkitab 'oiqcollector:23456789:respawn:startsrc -s observiq-otel-collector -a start -e "$(cat /opt/observiq-otel-collector/observiq-otel-collector.aix.env)"'
  fi

  echo "configured aix mkssys service"

  cat << EOF

The "observiq-otel-collector" service has been configured!

The collector's config file can be found here: 
  /opt/observiq-otel-collector/config.yaml

To view logs from the collector, run:
  sudo tail -F /opt/observiq-otel-collector/log/collector.log

For more information on configuring the collector, see the docs:
  https://github.com/observIQ/bindplane-agent/tree/main#observiq-opentelemetry-collector

To stop the observiq-otel-collector service, run:
  sudo stopsrc -s observiq-otel-collector

To start the observiq-otel-collector service, run:
  sudo startsrc -s observiq-otel-collector -a start -e "$(cat /opt/observiq-otel-collector/observiq-otel-collector.aix.env)"

To restart the observiq-otel-collector service, run:
  sudo stopsrc -s observiq-otel-collector && sudo startsrc -s observiq-otel-collector -a start -e "$(cat /opt/observiq-otel-collector/observiq-otel-collector.aix.env)"

To enable the service on startup, run:
  sudo mkitab 'oiqcollector:23456789:respawn:startsrc -s observiq-otel-collector -a start -e "$(cat /opt/observiq-otel-collector/observiq-otel-collector.aix.env)"'

If you have any other questions please contact us at support@observiq.com
EOF
}

init_type() {
  # Determine if we need service or systemctl for prereqs
  if command -v systemctl > /dev/null 2>&1; then
    command printf "systemd"
    return
  elif command -v service > /dev/null 2>&1; then
    command printf "service"
    return
  elif command -v mkssys > /dev/null 2>&1; then
    command printf "mkssys"
    return
  fi

  command printf "unknown"
  return
}

manage_service() {
  service_type="$(init_type)"
  case "$service_type" in
    systemd)
      manage_systemd_service
      ;;
    service)
      manage_sysv_service
      ;;
    mkssys)
      manage_aix_service
      ;;
    *)
      echo "could not detect init system, skipping service configuration"
  esac
}

finish_permissions() {
  # Goreleaser does not set plugin file permissions, so do them here
  # We also change the owner of the binary to observiq-otel-collector
  chown -R observiq-otel-collector:observiq-otel-collector /opt/observiq-otel-collector/observiq-otel-collector /opt/observiq-otel-collector/plugins/*
  chmod 0640 /opt/observiq-otel-collector/plugins/*

  # Initialize the log file to ensure it is owned by observiq-otel-collector.
  # This prevents the service (running as root) from assigning ownership to
  # the root user. By doing so, we allow the user to switch to observiq-otel-collector
  # user for 'non root' installs.
  touch /opt/observiq-otel-collector/log/collector.log
  chown observiq-otel-collector:observiq-otel-collector /opt/observiq-otel-collector/log/collector.log
}


finish_permissions
manage_service
