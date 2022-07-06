package install

import (
	"encoding/json"
	"fmt"
	"internal/syscall/windows/registry"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kardianos/service"
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
	kServiceConfig := &service.Config{
		Name: w.serviceName,
	}

	kService, err := service.New(nil, kServiceConfig)
	if err != nil {
		return fmt.Errorf("failed to create underlying service manager: %w", err)
	}

	err = kService.Start()
	if err != nil {
		return fmt.Errorf("failed to start using underlying service manager: %w", err)
	}

	return nil
}

// Stop the service
func (w windowsService) Stop() error {
	kServiceConfig := &service.Config{
		Name: w.serviceName,
	}

	kService, err := service.New(nil, kServiceConfig)
	if err != nil {
		return fmt.Errorf("failed to create underlying service manager: %w", err)
	}

	err = kService.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop using underlying service manager: %w", err)
	}

	return nil
}

// Installs the service
func (w windowsService) Install() error {
	wsc, err := readWindowsServiceConfig(w.newServiceFilePath)
	if err != nil {
		return fmt.Errorf("failed to read service config: %w", err)
	}

	if err = expandArguments(wsc, w.productName); err != nil {
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

	kServiceConfig := &service.Config{
		Name:        w.serviceName,
		DisplayName: wsc.Service.DisplayName,
		Description: wsc.Service.Description,
		Executable:  wsc.Path,
		Arguments:   splitArgs,
		Option: service.KeyValue{
			service.StartType:  startType,
			"DelayedAutoStart": delayed,
		},
	}

	kService, err := service.New(nil, kServiceConfig)
	if err != nil {
		return fmt.Errorf("failed to create underlying service manager: %w", err)
	}

	if err = kService.Install(); err != nil {
		return fmt.Errorf("failed to install using underlying service manager: %w", err)
	}

	return nil
}

// Uninstalls the service
func (w windowsService) Uninstall() error {
	kServiceConfig := &service.Config{
		Name: w.serviceName,
	}

	kService, err := service.New(nil, kServiceConfig)
	if err != nil {
		return err
	}

	err = kService.Uninstall()
	if err != nil {
		return err
	}

	return nil
}

type windowsServiceConfig struct {
	Path string `json:"path"`
	// Note: Name is a part of this struct, but we keep the service name hardcoded; We do not want to use a different service name.
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
		fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return &wsc, nil
}

// expandArguments expands [INSTALLDIR] to the actual install directory and
// expands '&quote;' to the literal '"'
func expandArguments(wsc *windowsServiceConfig, productName string) error {
	installDir, err := installDir(productName)
	if err != nil {
		return fmt.Errorf("failed to determine install dir: %w", err)
	}

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
func startType(wsc *windowsServiceConfig) (startType string, delayed bool, err error) {
	switch wsc.Service.Start {
	case "auto":
		startType = service.ServiceStartAutomatic
	case "demand":
		startType = service.ServiceStartManual
	case "disabled":
		startType = service.ServiceStartDisabled
	case "delayed":
		startType = service.ServiceStartAutomatic
		delayed = true
	default:
		err = fmt.Errorf("invalid start type in service config: %s", wsc.Service.Start)
	}
	return
}
