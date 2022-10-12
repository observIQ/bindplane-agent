package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/observiq/observiq-otel-collector/receiver/pluginreceiver"
	"github.com/spf13/pflag"
)

func main() {
	pluginDir := pflag.String("plugins", "./plugins", "The directory containing plugins")
	pflag.Parse()

	entries, err := os.ReadDir(*pluginDir)
	if err != nil {
		log.Fatalln("Failed to read plugin directory", err)
	}

	for _, entry := range entries {
		entryName := entry.Name()
		fullFilePath, err := filepath.Abs(filepath.Join(*pluginDir, entryName))
		if err != nil {
			log.Fatalln("Failed to determine path of plugin file", entryName, ":", err)
		}

		plugin, err := pluginreceiver.LoadPlugin(fullFilePath)
		if err != nil {
			log.Fatalln("Failed to load plugin", entryName, ":", err)
		}

	}
}
