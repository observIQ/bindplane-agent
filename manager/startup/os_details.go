package startup

import (
	"errors"
	"net"
	"strings"

	"github.com/shirou/gopsutil/v3/host"
)

func GetDetails() string {
	info, err := host.Info()
	if err != nil {
		return "unknown"
	}
	osDetails := GetOSDetails(info)
	return osDetails
}

// FindMACAddressOrUnknown calls FindMACAddress, and returns "unknown" if an error is returned
func FindMACAddressOrUnknown() string {
	macAddress, err := FindMACAddress()
	if err != nil {
		return "unknown"
	}
	return macAddress
}

// FindMACAddress does its best to find the MAC address for the local network interface
func FindMACAddress() (string, error) {
	iFaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iFace := range iFaces {
		if iFace.HardwareAddr.String() != "" {
			addrs, err := iFace.Addrs()
			if err != nil {
				continue
			}

			for _, addr := range addrs {
				address := addr.String()
				if strings.Contains(address, "/") {
					address = address[:strings.Index(address, "/")]
				}
				ipAddress := net.ParseIP(address)
				if isValidV4Address(ipAddress) {
					return iFace.HardwareAddr.String(), nil
				}
			}
		}
	}

	return "", errors.New("Unable to find MAC address")
}

// isValidV4Address checks that the IP address is not nil, loopback, or unspecified. It also checks that it is IPv4
func isValidV4Address(address net.IP) bool {
	return address != nil && !address.IsLoopback() && !address.IsUnspecified() && address.To4() != nil
}
