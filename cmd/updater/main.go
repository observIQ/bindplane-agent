package main

import (
	"fmt"

	"github.com/observiq/observiq-otel-collector/internal/version"
	"github.com/spf13/pflag"
)

// Unimplemented
func main() {
	var showVersion = pflag.BoolP("version", "v", false, "prints the version of the collector")
	pflag.Parse()

	if *showVersion {
		fmt.Println("observiq-otel-collector updater version", version.Version())
		fmt.Println("commit:", version.GitHash())
		fmt.Println("built at:", version.Date())
		return
	}
}
