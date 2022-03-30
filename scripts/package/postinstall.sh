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

create_systemd_service() {
    if [ -d "/usr/lib/systemd/system" ]; then
      systemd_service_dir="/usr/lib/systemd/system"
    elif [ -d "/lib/systemd/system" ]; then
      systemd_service_dir="/lib/systemd/system"
    elif [ -d "/etc/systemd/system" ]; then
      systemd_service_dir="/etc/systemd/system"
    else
      echo "failed to detect systemd service file directory"
      exit 1
    fi

    echo "detected service file directory: ${systemd_service_dir}"

    systemd_service_file="${systemd_service_dir}/observiq-otel-collector.service"

    cat <<EOF > ${systemd_service_file}
[Unit]
Description=observIQ's distribution of the OpenTelemetry collector
After=network.target
StartLimitIntervalSec=120
StartLimitBurst=5
[Service]
Type=simple
User=observiq-otel-collector
Group=observiq-otel-collector
Environment=PATH=/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin
WorkingDirectory=/opt/observiq-otel-collector
ExecStart=/opt/observiq-otel-collector/observiq-otel-collector --config config.yaml
SuccessExitStatus=0
TimeoutSec=120
StandardOutput=journal
Restart=on-failure
RestartSec=5s
[Install]
WantedBy=multi-user.target
EOF

    chmod 0644 "${systemd_service_file}"
    chown root:root "${systemd_service_file}"

    systemctl daemon-reload

    echo "configured systemd service"

    cat << EOF

The "observiq-otel-collector" service has been configured!

The collector's config file can be found here: 
  /opt/observiq-otel-collector/config.yaml

To view logs from the collector, run:
  sudo journalctl -u observiq-otel-collector.service

For more information on configuring the collector, see the docs:
  https://github.com/observiq/observiq-otel-collector/tree/main#observiq-opentelemetry-collector

To stop the observiq-otel-collector service, run:
  sudo systemctl stop observiq-otel-collector

To start the observiq-otel-collector service, run:
  sudo systemctl start observiq-otel-collector

To restart the observiq-otel-collector service, run:
  sudo systemctl restart observiq-otel-collector

If you have any other questions please contact us at support@observiq.com
EOF
}

init_type() {
  systemd_test="$(systemctl 2>/dev/null || : 2>&1)"
  if command printf "$systemd_test" | grep -q '\-.mount'; then
    command printf "systemd"
    return
  fi

  command printf "unknown"
  return
}

install_service() {
  service_type="$(init_type)"
  case "$service_type" in
    systemd)
      create_systemd_service
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
}


finish_permissions
install_service
