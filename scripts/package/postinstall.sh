#!/bin/sh

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

    systemd_service_file="${systemd_service_dir}/stanza.service"

    cat <<EOF > ${systemd_service_file}
[Unit]
Description=observIQ's distribution of the OpenTelemetry collector
After=network.target
StartLimitIntervalSec=120
StartLimitBurst=5
[Service]
Type=simple
User=observiq-collector
Group=observiq-collector
Environment=PATH=/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin
WorkingDirectory=/opt/observiq-collector
ExecStart=/usr/bin/observiq-collector --config config.yaml
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
    chown -R stanza:stanza /opt/observiq-collector/plugins/*
    chmod 0644 /opt/observiq-collector/plugins/*
}


finish_permissions
install_service
