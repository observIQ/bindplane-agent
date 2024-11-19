module github.com/observiq/bindplane-agent/updater

go 1.22.7

require (
	github.com/google/uuid v1.6.0
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/observiq/bindplane-agent/packagestate v1.64.0
	github.com/oklog/ulid/v2 v2.0.2
	github.com/open-telemetry/opamp-go v0.9.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.9.0
	go.uber.org/zap v1.27.0
	golang.org/x/sys v0.27.0
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/google/go-cmp v0.6.0 // indirect

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

replace github.com/observiq/bindplane-agent/packagestate => ../packagestate
