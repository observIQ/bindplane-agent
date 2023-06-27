package debug

import "os"

// If one of the debug commands is specified, do that action and exit.
func TryDebugCommand() {
	if len(os.Args) < 2 {
		// Need second positional arg to tell what to exec
		return
	}

	switch os.Args[1] {
	case "glob":
		RunGlob()
	case "file-info":
		RunPermissions()
	case "read":
		RunRead()
	default:
		return
	}

	os.Exit(0)
}
