//go:build windows

package debug

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

func addSysFileInfo(fi os.FileInfo, sb *strings.Builder) {
	if w32Attrs, ok := fi.Sys().(*syscall.Win32FileAttributeData); ok {
		sb.WriteString(fmt.Sprintf("0x%08x", w32Attrs.FileAttributes))
	}
}
