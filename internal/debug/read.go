package debug

import (
	"io"
	"log"
	"os"

	"github.com/bmatcuk/doublestar/v4"
)

func RunRead() {
	if len(os.Args) != 3 {
		printReadUsageAndExit()
	}

	glob := os.Args[2]

	matches, err := doublestar.FilepathGlob(glob, doublestar.WithFilesOnly(), doublestar.WithFailOnIOErrors())
	if err != nil {
		log.Fatalf("Encountered error while globbing: %s\n", err)
	}

	for _, match := range matches {
		f, err := os.Open(match)
		if err != nil {
			log.Default().Printf("Failed to open %s: %s", match, err)
			continue
		}

		log.Default().Printf("Reading %s", match)
		_, err = io.Copy(io.Discard, f)
		if err != nil {
			log.Default().Printf("failed to read file: %s (close err: %v)", err, f.Close())
			continue
		}

		log.Default().Printf("Read %s (close error: %v)", match, f.Close())
	}
}

func printReadUsageAndExit() {
	log.Fatalf("Usage: %s read <glob>\n", os.Args[0])
}
