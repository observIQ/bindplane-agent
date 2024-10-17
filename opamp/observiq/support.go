package observiq

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/shirou/gopsutil/v3/host"
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
	AgentID  string
	Version  string
	Goos     string
	GoArch   string
	HostInfo *host.InfoStat
}

func newDiagnosticInfo(agentID, version string) (diagnosticInfo, error) {
	hi, err := host.Info()
	if err != nil {
		return diagnosticInfo{}, fmt.Errorf("stat hostinfo: %w", err)
	}

	return diagnosticInfo{
		AgentID:  agentID,
		Version:  version,
		Goos:     runtime.GOOS,
		GoArch:   runtime.GOARCH,
		HostInfo: hi,
	}, nil
}

func writeSupportPackage(writer io.Writer, di diagnosticInfo) error {
	gzipWriter := gzip.NewWriter(writer)
	defer gzipWriter.Close()

	tw := tar.NewWriter(gzipWriter)

	diYaml, err := yaml.Marshal(di)
	if err != nil {
		return err
	}

	// Write basic agent info
	if err := writeBytesToTar("diagnostic-info.yaml", diYaml, tw); err != nil {
		return fmt.Errorf("write info yaml: %w", err)
	}

	// Write log files
	home := os.Getenv("OIQ_OTEL_COLLECTOR_HOME")
	logsDir := filepath.Join(home, "log")

	logsDirEntries, err := os.ReadDir(logsDir)
	if err != nil {
		return fmt.Errorf("read logs dir entries: %w", err)
	}
	for _, ent := range logsDirEntries {
		if !ent.IsDir() {
			path := filepath.Join(logsDir, ent.Name())
			err := writeFileToTar(path, ent.Name(), tw)
			if err != nil {
				return fmt.Errorf("write log files: %w", err)
			}
		}
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
