// +build darwin

package startup

import (
	"testing"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/stretchr/testify/require"
)

func TestMajorMinor(t *testing.T) {
	testCases := []struct {
		name          string
		semver        string
		expectedValue string
	}{
		{
			name:          "Catalina",
			semver:        "10.15.7",
			expectedValue: "macOS 10.15 Catalina",
		},
		{
			name:          "Mojave",
			semver:        "10.14.7",
			expectedValue: "macOS 10.14 Mojave",
		},
		{
			name:          "Catalina",
			semver:        "10.13.9",
			expectedValue: "macOS 10.13 High Sierra",
		},
		{
			name:          "High Sierra",
			semver:        "10.12.712",
			expectedValue: "macOS 10.12 Sierra",
		},
		{
			name:          "Unknown Name",
			semver:        "10.16.0",
			expectedValue: "macOS 10.16",
		},
		{
			name:          "Only Major",
			semver:        "10",
			expectedValue: "macOS",
		},
	}
	for _, tc := range testCases {
		detail := GetOSDetails(&host.InfoStat{
			PlatformVersion: tc.semver,
		})
		require.Equal(t, tc.expectedValue, detail)
	}
}
