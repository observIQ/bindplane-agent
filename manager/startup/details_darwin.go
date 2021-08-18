package startup

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/host"
)

func GetOSDetails(info *host.InfoStat) string {
	majorMinor := getMajorMinor(info.PlatformVersion)
	return getPrettyVersion(majorMinor)
}

func getMajorMinor(semantic string) string {
	parts := strings.Split(semantic, ".")
	if len(parts) < 2 {
		return ""
	}
	return fmt.Sprintf("%s.%s", parts[0], parts[1])
}

func getPrettyVersion(majorMinor string) string {
	switch majorMinor {
	case "10.15":
		return "macOS 10.15 Catalina"
	case "10.14":
		return "macOS 10.14 Mojave"
	case "10.13":
		return "macOS 10.13 High Sierra"
	// Note Go 1.17 will not support these previous versions of macOS
	case "10.12":
		return "macOS 10.12 Sierra"
	case "10.11":
		return "OS X 10.11 El Capitan"
	case "10.10":
		return "OS X 10.10 Yosemite"
	case "10.9":
		return "OS X 10.9 Mavericks"
	case "10.8":
		return "OS X 10.8 Mountain Lion"
	default:
		if majorMinor == "" {
			return "macOS"
		}
		return fmt.Sprintf("macOS %s", majorMinor)
	}
}
