package install

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	"github.com/kballard/go-shellquote"
)

func NewService(latestPath string) Service {
	return &windowsService{
		newServiceFilePath: filepath.Join(latestPath, "install", "windows_service.json"),
		serviceName:        "observiq-otel-collector",
		productName:        "observIQ Distro for OpenTelemetry Collector",
	}
}

type windowsService struct {
	// newServiceFilePath is the file path to the new unit file
	newServiceFilePath string
	// serviceName is the name of the service
	serviceName string
	// productName is the name of the installed product
	productName string
}

// Start the service
func (w windowsService) Start() error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(w.serviceName)
	if err != nil {
		return fmt.Errorf("failed to open service: %w", err)
	}
	defer s.Close()

	if err := s.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	return nil
}

// Stop the service
func (w windowsService) Stop() error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(w.serviceName)
	if err != nil {
		return fmt.Errorf("failed to open service: %w", err)
	}
	defer s.Close()

	if _, err := s.Control(svc.Stop); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	return nil
}

// Installs the service
func (w windowsService) Install() error {
	wsc, err := readWindowsServiceConfig(w.newServiceFilePath)
	if err != nil {
		return fmt.Errorf("failed to read service config: %w", err)
	}

	iDir, err := installDir(w.productName)
	if err != nil {
		return fmt.Errorf("failed to get install dir: %w", err)
	}

	if err = expandArguments(wsc, w.productName, iDir); err != nil {
		return fmt.Errorf("failed to expand arguments in service config: %w", err)
	}

	splitArgs, err := shellquote.Split(wsc.Service.Arguments)
	if err != nil {
		return fmt.Errorf("failed to parse arguments in service config: %w", err)
	}

	startType, delayed, err := startType(wsc)
	if err != nil {
		return fmt.Errorf("failed to parse start type in service config: %w", err)
	}

	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.CreateService(w.serviceName,
		filepath.Join(iDir, wsc.Path),
		mgr.Config{
			Description:      wsc.Service.Description,
			DisplayName:      wsc.Service.DisplayName,
			StartType:        startType,
			DelayedAutoStart: delayed,
		},
		splitArgs...,
	)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	defer s.Close()

	return nil
}

// Uninstalls the service
func (w windowsService) Uninstall() error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(w.serviceName)
	if err != nil {
		return fmt.Errorf("failed to open service: %w", err)
	}
	defer s.Close()

	// TODO: Should this wait for the service to no longer exist?
	// This only marks the service for deletion, no guarantee is made about when it's actually deleted
	if err = s.Delete(); err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	return nil
}

type windowsServiceConfig struct {
	Path string `json:"path"`
	// Note: Name is a part of the on disk config, but we keep the service name hardcoded; We do not want to use a different service name.
	Service struct {
		// Start gives the start type of the service.
		// See: https://wixtoolset.org/documentation/manual/v3/xsd/wix/serviceinstall.html
		Start string `json:"start"`
		// DisplayName is the human-readable name of the service.
		DisplayName string `json:"display-name"`
		// Description is a human-readable description of the service.
		Description string `json:"description"`
		// Arguments is a list of space-separated
		Arguments string `json:"arguments"`
	} `json:"service"`
}

func readWindowsServiceConfig(path string) (*windowsServiceConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var wsc windowsServiceConfig
	err = json.Unmarshal(b, &wsc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return &wsc, nil
}

// expandArguments expands [INSTALLDIR] to the actual install directory and
// expands '&quote;' to the literal '"'
func expandArguments(wsc *windowsServiceConfig, productName, installDir string) error {
	wsc.Service.Arguments = string(replaceInstallDir([]byte(wsc.Service.Arguments), installDir))
	wsc.Service.Arguments = strings.ReplaceAll(wsc.Service.Arguments, "&quot;", `"`)
	return nil
}

func installDir(productName string) (string, error) {
	keyPath := fmt.Sprintf(`Software\Microsoft\Windows\CurrentVersion\Uninstall\%s`, productName)
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.READ)
	if err != nil {
		return "", fmt.Errorf("failed to open registry key: %w", err)
	}
	defer func() {
		err := key.Close()
		if err != nil {
			log.Default().Printf("getInstallDir: failed to close registry key")
		}
	}()

	val, _, err := key.GetStringValue("InstallLocation")
	if err != nil {
		return "", fmt.Errorf("failed to read install dir: %w", err)
	}

	return val, nil
}

// startType converts the start type from the windowsServiceConfig to a start type recongizable by
// kardianos/service.
func startType(wsc *windowsServiceConfig) (startType uint32, delayed bool, err error) {
	switch wsc.Service.Start {
	case "auto":
		startType = mgr.StartAutomatic
	case "demand":
		startType = mgr.StartManual
	case "disabled":
		startType = mgr.StartDisabled
	case "delayed":
		startType = mgr.StartAutomatic
		delayed = true
	default:
		err = fmt.Errorf("invalid start type in service config: %s", wsc.Service.Start)
	}
	return
}
