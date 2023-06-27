package debug

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

func RunPermissions() {
	if len(os.Args) != 3 {
		printPermissionsUsageAndExit()
	}

	glob := os.Args[2]

	matches, err := doublestar.FilepathGlob(glob, doublestar.WithFilesOnly(), doublestar.WithFailOnIOErrors())
	if err != nil {
		log.Fatalf("Encountered error while globbing: %s\n", err)
	}

	for _, match := range matches {
		fi, err := os.Lstat(match)
		sb := strings.Builder{}

		sb.WriteString(match)
		sb.WriteString(" - ")
		if err != nil {
			sb.WriteString("lstat error: ")
			sb.WriteString(err.Error())
			sb.WriteString("\n")
			continue
		}
		sb.WriteString(fi.Mode().String())
		sb.WriteString(" - ")
		sb.WriteString(fmt.Sprintf("%d", fi.Size()))
		sb.WriteString(" - ")
		sb.WriteString(fi.ModTime().String())
		sb.WriteString(" - ")
		addSysFileInfo(fi, &sb)
		sb.WriteString("\n")

		log.Default().Print(sb.String())
	}
}

func printPermissionsUsageAndExit() {
	log.Fatalf("Usage: %s file-info <glob>\n", os.Args[0])
}
