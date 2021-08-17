package startup

import (
	"github.com/shirou/gopsutil/v3/host"
)

func GetOSDetails(info *host.InfoStat) string {
	return info.Platform
}
