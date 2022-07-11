package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/windows/registry"

	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func TestWindowsServiceInstall(t *testing.T) {
	t.Run("Test install + uninstall", func(t *testing.T) {
		tempDir := t.TempDir()
		testProductName := "Test Product"

		serviceJSON := filepath.Join(tempDir, "windows-service.json")
		testServiceProgram := filepath.Join(tempDir, "windows-service.exe")
		serviceGoFile, err := filepath.Abs(filepath.Join("testdata", "test-windows-service.go"))
		require.NoError(t, err)

		writeServiceFile(t, serviceJSON, filepath.Join("testdata", "windows-service.json"), serviceGoFile)
		compileProgram(t, serviceGoFile, testServiceProgram)

		defer uninstallService(t)
		createInstallDirRegistryKey(t, testProductName, tempDir)
		defer deleteInstallDirRegistryKey(t, testProductName)

		w := &windowsService{
			newServiceFilePath: serviceJSON,
			serviceName:        "windows-service",
			productName:        testProductName,
		}

		err = w.Install()
		require.NoError(t, err)

		//We want to check that the service was actually loaded
		requireServiceLoadedStatus(t, true)

		err = w.Uninstall()
		require.NoError(t, err)

		//Make sure the service is no longer listed
		requireServiceLoadedStatus(t, false)
	})

	t.Run("Test stop + start", func(t *testing.T) {
		tempDir := t.TempDir()
		testProductName := "Test Product"

		serviceJSON := filepath.Join(tempDir, "windows-service.json")
		testServiceProgram := filepath.Join(tempDir, "windows-service.exe")
		serviceGoFile, err := filepath.Abs(filepath.Join("testdata", "test-windows-service.go"))
		require.NoError(t, err)

		writeServiceFile(t, serviceJSON, filepath.Join("testdata", "windows-service.json"), serviceGoFile)
		compileProgram(t, serviceGoFile, testServiceProgram)

		defer uninstallService(t)
		createInstallDirRegistryKey(t, testProductName, tempDir)
		defer deleteInstallDirRegistryKey(t, testProductName)

		w := &windowsService{
			newServiceFilePath: serviceJSON,
			serviceName:        "windows-service",
			productName:        testProductName,
		}

		err = w.Install()
		require.NoError(t, err)

		// We want to check that the service was actually loaded
		requireServiceLoadedStatus(t, true)

		err = w.Start()
		require.NoError(t, err)

		requireServiceRunningStatus(t, true)

		err = w.Stop()
		require.NoError(t, err)

		requireServiceRunningStatus(t, false)

		err = w.Uninstall()
		require.NoError(t, err)

		// Make sure the service is no longer listed
		requireServiceLoadedStatus(t, false)
	})

	t.Run("Test invalid path for input file", func(t *testing.T) {
		tempDir := t.TempDir()
		testProductName := "Test Product"

		serviceJSON := filepath.Join(tempDir, "windows-service.json")
		testServiceProgram := filepath.Join(tempDir, "windows-service.exe")
		serviceGoFile, err := filepath.Abs(filepath.Join("testdata", "test-windows-service.go"))
		require.NoError(t, err)

		writeServiceFile(t, serviceJSON, filepath.Join("testdata", "windows-service.json"), serviceGoFile)
		compileProgram(t, serviceGoFile, testServiceProgram)

		defer uninstallService(t)
		createInstallDirRegistryKey(t, testProductName, tempDir)
		defer deleteInstallDirRegistryKey(t, testProductName)

		w := &windowsService{
			newServiceFilePath: filepath.Join(tempDir, "not-a-valid-service.json"),
			serviceName:        "windows-service",
			productName:        testProductName,
		}

		err = w.Install()
		require.ErrorContains(t, err, "The system cannot find the file specified.")
		requireServiceLoadedStatus(t, false)
	})

	t.Run("Uninstall fails if not installed", func(t *testing.T) {
		tempDir := t.TempDir()
		testProductName := "Test Product"

		serviceJSON := filepath.Join(tempDir, "windows-service.json")

		w := &windowsService{
			newServiceFilePath: serviceJSON,
			serviceName:        "windows-service",
			productName:        testProductName,
		}

		err := w.Uninstall()
		require.ErrorContains(t, err, "failed to open service")
		requireServiceLoadedStatus(t, false)
	})

	t.Run("Start fails if service not found", func(t *testing.T) {
		tempDir := t.TempDir()
		testProductName := "Test Product"

		serviceJSON := filepath.Join(tempDir, "windows-service.json")

		w := &windowsService{
			newServiceFilePath: serviceJSON,
			serviceName:        "windows-service",
			productName:        testProductName,
		}

		err := w.Start()
		require.ErrorContains(t, err, "failed to open service")
	})

	t.Run("Stop fails if service not found", func(t *testing.T) {
		tempDir := t.TempDir()
		testProductName := "Test Product"

		serviceJSON := filepath.Join(tempDir, "windows-service.json")

		w := &windowsService{
			newServiceFilePath: serviceJSON,
			serviceName:        "windows-service",
			productName:        testProductName,
		}

		err := w.Stop()
		require.ErrorContains(t, err, "failed to open service")
	})
}

// uninstallService is a helper that uninstalls the service manually for test setup, in case it is somehow leftover.
func uninstallService(t *testing.T) {
	m, err := mgr.Connect()
	require.NoError(t, err)
	defer m.Disconnect()

	s, err := m.OpenService("windows-service")
	if err != nil {
		// Failed to open the service, we assume it doesn't exist
		return
	}
	defer s.Close()

	status, err := s.Control(svc.Stop)
	// If we get an error, the service is likely already in a stopped state.
	if err == nil {
		for status.State != svc.Stopped {
			time.Sleep(100 * time.Millisecond)
			status, err = s.Query()
			require.NoError(t, err)
		}
	}

	err = s.Delete()
	require.NoError(t, err)
}

func requireServiceLoadedStatus(t *testing.T, loaded bool) {
	t.Helper()

	m, err := mgr.Connect()
	require.NoError(t, err, "failed to connect to service manager")
	defer m.Disconnect()

	s, err := m.OpenService("windows-service")
	if err != nil {
		require.False(t, loaded, "Could not connect open service, but service should be loaded")
		return
	}
	defer s.Close()

	require.True(t, loaded, "Connected to open service, but it should not be loaded")

}

func requireServiceRunningStatus(t *testing.T, running bool) {
	t.Helper()

	m, err := mgr.Connect()
	require.NoError(t, err, "failed to connect to service manager")
	defer m.Disconnect()

	s, err := m.OpenService("windows-service")
	require.NoError(t, err, "Failed to open service")
	defer s.Close()

	status, err := s.Query()
	require.NoError(t, err, "Failed to query service state")

	if running {
		require.Contains(t, []svc.State{svc.StartPending, svc.Running}, status.State)
	} else {
		require.Contains(t, []svc.State{svc.StopPending, svc.Stopped}, status.State)
	}
}

func writeServiceFile(t *testing.T, outPath, inPath, serviceGoPath string) {
	t.Helper()

	b, err := os.ReadFile(inPath)
	require.NoError(t, err)

	fileStr := string(b)
	fileStr = os.Expand(fileStr, func(s string) string {
		switch s {
		case "SERVICE_PATH":
			return strings.ReplaceAll(serviceGoPath, `\`, `\\`)
		}
		return ""
	})

	err = os.WriteFile(outPath, []byte(fileStr), 0666)
	require.NoError(t, err)
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

	keyPath := fmt.Sprintf(`Software\Microsoft\Windows\CurrentVersion\Uninstall\%s`, productName)
	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, keyPath, registry.WRITE)
	require.NoError(t, err)
	defer key.Close()

	err = key.SetStringValue("InstallLocation", installDir)
	require.NoError(t, err)
}

func compileProgram(t *testing.T, inPath, outPath string) {
	t.Helper()

	cmd := exec.Command("go.exe", "build", "-o", outPath, inPath)
	err := cmd.Run()
	require.NoError(t, err)
}
