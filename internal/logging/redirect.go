package logging

import (
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

/*
	Redirect console output to a rotating log file
*/
func RedirectConsoleOutput(logPath string) error {
	l := lumberjack.Logger{
		Filename:   logPath,
		MaxBackups: 3,
		MaxSize:    1,
		MaxAge:     7,
	}

	r, w, err := os.Pipe()
	if err != nil {
		return err
	}

	origStdout := os.Stdout
	origStderr := os.Stderr

	os.Stdout = w
	os.Stderr = w

	buf := make([]byte, 1024)
	go func() {
		for {
			n, err := r.Read(buf)
			if err != nil {
				os.Stdout = origStdout
				os.Stderr = origStderr
				return
			}
			_, err = l.Write(buf[:n])
			if err != nil {
				os.Stdout = origStdout
				os.Stderr = origStderr
				return
			}
		}
	}()

	return nil
}
