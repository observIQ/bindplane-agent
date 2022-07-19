package path

// InstallDirFromRegistry gets the installation dir of the given product from the Windows Registry
func InstallDirFromRegistry(productName string) (string, error) {
	// this key is created when installing using the MSI installer
	keyPath := fmt.Sprintf(`Software\Microsoft\Windows\CurrentVersion\Uninstall\%s`, productName)
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.READ)
	if err != nil {
		return "", fmt.Errorf("failed to open registry key: %w", err)
	}
	defer func() {
		err := key.Close()
		if err != nil {
			log.Default().Printf("installDirFromRegistry: failed to close registry key")
		}
	}()

	// This value ("InstallLocation") contains the path to the install folder.
	val, _, err := key.GetStringValue("InstallLocation")
	if err != nil {
		return "", fmt.Errorf("failed to read install dir: %w", err)
	}

	return val, nil
}

// InstallDir returns the filepath to the install directory
func InstallDir() (string, error) {
	return installDirFromRegistry(defaultProductName)
}
