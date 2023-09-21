module github.com/observiq/bindplane-agent/updater

go 1.20

require (
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/observiq/bindplane-agent/packagestate v1.35.1
	github.com/open-telemetry/opamp-go v0.2.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.4
	go.uber.org/zap v1.25.0
	golang.org/x/sys v0.12.0
)

require (
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/observiq/bindplane-agent/packagestate => ../packagestate
