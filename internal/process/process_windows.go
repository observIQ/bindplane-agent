// +build windows

package process

import (
	"os"
)

// MatchesParent returns a boolean indicating if the parent process matches the supplied ppid.
func MatchesParent(ppid int) bool {
	process, err := os.FindProcess(ppid)
	defer process.Release()
	return err == nil && process != nil
}
