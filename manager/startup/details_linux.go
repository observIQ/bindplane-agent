package startup

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/host"
)

func GetOSDetails(info *host.InfoStat) string {
	return fmt.Sprintf("%s %s", getNicePlatform(info.Platform), info.PlatformVersion)
}

func getNicePlatform(platform string) string {
	switch strings.ToLower(platform) {
	case "centos":
		return "CentOS"
	case "rhel":
		return "RedHat Enterprise Linux"
	case "ubuntu":
		return "Ubuntu"
	case "suse", "sles":
		return "SLES Enterprise Linux"
	case "coreos":
		return "CoreOS"
	case "linuxmint":
		return "Linux Mint"
	default:
		return strings.Title(platform)
	}
}
