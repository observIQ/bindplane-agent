//go:build windows && integration

package path

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"golang.org/x/sys/windows/registry"
)

func TestInstallDirFromRegistry(t *testing.T) {
	t.Run("Successfully grabs install dir from registry", func(t *testing.T) {
		productName := "default-product-name"
		installDir, err := filepath.Abs("C:/temp")
		require.NoError(t, err)

		defer deleteInstallDirRegistryKey(t, productName)
		createInstallDirRegistryKey(t, productName, installDir)

		dir, err := installDirFromRegistry(zaptest.NewLogger(t), productName)
		require.NoError(t, err)
		require.Equal(t, installDir+`\`, dir)
	})

	t.Run("Registry key does not exist", func(t *testing.T) {
		productName := "default-product-name"

		_, err := installDirFromRegistry(zaptest.NewLogger(t), productName)
		require.ErrorContains(t, err, "failed to open registry key")
	})
}

func deleteInstallDirRegistryKey(t *testing.T, productName string) {
	t.Helper()

	keyPath := fmt.Sprintf(`Software\Microsoft\Windows\CurrentVersion\Uninstall\%s`, productName)
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.WRITE)
	if err != nil {
		// Key may not exist, assume that's why we couldn't open it
		return
	}
	defer key.Close()

	err = registry.DeleteKey(key, "")
	require.NoError(t, err)
}

func createInstallDirRegistryKey(t *testing.T, productName, installDir string) {
	t.Helper()

	installDir, err := filepath.Abs(installDir)
	require.NoError(t, err)
	installDir += `\`

	keyPath := fmt.Sprintf(`Software\Microsoft\Windows\CurrentVersion\Uninstall\%s`, productName)
	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, keyPath, registry.WRITE)
	require.NoError(t, err)
	defer key.Close()

	err = key.SetStringValue("InstallLocation", installDir)
	require.NoError(t, err)
}
