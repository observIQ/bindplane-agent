package observiq

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
)

const (
	diagnosticsReportV1Capability = "com.bindplane.diagnostics.v1"
	diagnosticsRequestType        = "requestDiagnosticPackage"
)

type diagnosticRequestCustomMessage struct {
	ReportURL string            `yaml:"report_url"`
	Headers   map[string]string `yaml:"headers"`
}

type diagnosticInfo struct {
	AgentID string
	Version string
	Goos    string
	GoArch  string
}

func newDiagnosticInfo(agentID, version string) diagnosticInfo {
	return diagnosticInfo{
		AgentID: agentID,
		Version: version,
		Goos:    runtime.GOOS,
		GoArch:  runtime.GOARCH,
	}
}

func writeSupportPackage(writer io.Writer, di diagnosticInfo) error {
	gzipWriter := gzip.NewWriter(writer)
	defer gzipWriter.Close()

	tw := tar.NewWriter(gzipWriter)

	diYaml, err := yaml.Marshal(di)
	if err != nil {
		return err
	}

	if err := writeBytesToTar("diagnostic-info.yaml", diYaml, tw); err != nil {
		return fmt.Errorf("write info yaml: %w", err)
	}

	return tw.Close()
}

func writeBytesToTar(file string, by []byte, tw *tar.Writer) error {
	err := tw.WriteHeader(&tar.Header{
		Typeflag: tar.TypeReg,
		Name:     file,
		Size:     int64(len(by)),
		Mode:     0666,
	})
	if err != nil {
		return err
	}

	_, err = tw.Write(by)
	if err != nil {
		return err
	}

	return nil
}

func writeFileToTar(filePath, tarFile string, tw *tar.Writer) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("stat file: %w", err)
	}

	err = tw.WriteHeader(&tar.Header{
		Typeflag: tar.TypeReg,
		Name:     tarFile,
		Size:     fi.Size(),
		Mode:     0666,
	})
	if err != nil {
		return fmt.Errorf("write tar header: %w", err)
	}

	if _, err = io.Copy(tw, f); err != nil {
		return fmt.Errorf("copy file to tar: %w", err)
	}

	return nil
}
