module github.com/observiq/observiq-otel-collector/cmd/plugindocgen

go 1.18

require (
	github.com/observiq/observiq-otel-collector/receiver/pluginreceiver v1.9.1
	github.com/spf13/pflag v1.0.5
)

replace github.com/observiq/observiq-otel-collector/receiver/pluginreceiver => ../../receiver/pluginreceiver
