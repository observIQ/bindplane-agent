package debug

import (
	"log"
	"os"

	"github.com/bmatcuk/doublestar/v4"
)

func RunGlob() {
	if len(os.Args) != 3 {
		printGlobUsageAndExit()
	}

	glob := os.Args[2]

	matches, err := doublestar.FilepathGlob(glob, doublestar.WithFilesOnly(), doublestar.WithFailOnIOErrors())
	if err != nil {
		log.Fatalf("Encountered error while globbing: %s\n", err)
	}

	for _, match := range matches {
		log.Default().Println(match)
	}
}

func printGlobUsageAndExit() {
	log.Fatalf("Usage: %s glob <glob>\n", os.Args[0])
}
