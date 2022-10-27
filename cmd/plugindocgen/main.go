package main

import (
	"github.com/spf13/pflag"
)

func main() {
	// Parse flags
	pluginDir := pflag.String("plugins", "../../plugins", "The directory containing plugins")
	pluginExtraDetailsDir := pflag.String("extra_details", "./extra_details", "The directory containing extra details for plugin documentation")
	pluginDocsDir := pflag.String("docs", "../../docs/plugins", "The directory containing plugin documentation")
	pflag.Parse()

	generator := newPluginDocGenerator(*pluginDir, *pluginExtraDetailsDir, *pluginDocsDir)

	generator.GeneratePluginDocs()
}
