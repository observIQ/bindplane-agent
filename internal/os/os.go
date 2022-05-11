// Package os handles grabbing OS info from the system
package os

import (
	"fmt"
	"net"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v3/host"
)

const unknownValue = "unknown"

// Hostname retrieves the hostname of the machine
func Hostname() (string, error) {
	info, err := host.Info()
	if err != nil {
		return unknownValue, err
	}

	return info.Hostname, nil
}

// Name returns the name of the os.
func Name() (string, error) {
	info, err := host.Info()
	if err != nil {
		return unknownValue, err
	}
	return parseName(info, runtime.GOOS), nil
}

// parseName parses the os name from the supplied info.
func parseName(info *host.InfoStat, os string) string {
	if info == nil {
		return unknownValue
	}

	switch os {
	case "darwin":
		return fmt.Sprintf("macOS %s", parseDarwinVersion(info.PlatformVersion))
	case "linux":
		return fmt.Sprintf("%s %s", formatLinuxName(info.Platform), info.PlatformVersion)
	default:
		return info.Platform
	}
}

// parseDarwinVersion parses a darwin version.
func parseDarwinVersion(platformVersion string) string {
	parts := strings.Split(platformVersion, ".")
	if len(parts) < 2 {
		return "unknown version"
	}
	return fmt.Sprintf("%s.%s", parts[0], parts[1])
}

// formatLinuxName formats a linux distribution name.
func formatLinuxName(platform string) string {
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

// MACAddress returns the MAC address for the host.
func MACAddress() string {
	return findMACAddress(net.Interfaces)
}

// findMACAddress does its best to find the MAC address for the local network interface.
func findMACAddress(interfaces func() ([]net.Interface, error)) string {
	iFaces, err := interfaces()
	if err != nil {
		return unknownValue
	}

	for _, iFace := range iFaces {
		if iFace.HardwareAddr.String() != "" {
			addrs, _ := iFace.Addrs()
			for _, addr := range addrs {
				address := addr.String()
				if strings.Contains(address, "/") {
					address = address[:strings.Index(address, "/")]
				}
				ipAddress := net.ParseIP(address)
				if isValidV4Address(ipAddress) {
					return iFace.HardwareAddr.String()
				}
			}
		}
	}

	return unknownValue
}

// isValidV4Address checks that the IP address is not nil, loopback, or unspecified. It also checks that it is IPv4
func isValidV4Address(address net.IP) bool {
	return address != nil && !address.IsLoopback() && !address.IsUnspecified() && address.To4() != nil
}
