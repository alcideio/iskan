package trivy

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/alcideio/iskan/pkg/util"
	types "github.com/alcideio/iskan/pkg/vulnprovider/api"
	"k8s.io/klog"
)

const (
	trivyCmd = "trivy"
)

type ImageRef struct {
	Name     string
	Auth     RegistryAuth
	Insecure bool
}

// RegistryAuth wraps registry credentials.
type RegistryAuth interface {
}

type NoAuth struct {
}

type BasicAuth struct {
	Username string
	Password string
}

type BearerAuth struct {
	Token string
}

type Trivy interface {
	Scan(imageRef ImageRef) (ScanReport, error)
	GetVersion() (VersionInfo, error)
}

type wrapper struct {
	config types.TrivyConfig
	runner util.CmdRunner
}

func NewScanner(config types.TrivyConfig, ambassador util.CmdRunner) Trivy {
	return &wrapper{
		config: config,
		runner: ambassador,
	}
}

func (w *wrapper) Scan(imageRef ImageRef) (ScanReport, error) {
	var report ScanReport
	var err error

	reportFile, err := w.runner.TempFile(w.config.ReportsDir, "scan_report_*.json")
	if err != nil {
		return report, err
	}
	klog.V(5).Infof("Saving scan report to tmp file '%v'", reportFile.Name())
	defer func() {
		klog.V(5).Infof("Removing scan report to tmp file '%v'", reportFile.Name())
		err := w.runner.Remove(reportFile.Name())
		if err != nil {
			klog.V(5).Infof("Removing scan report to tmp file '%v'", reportFile.Name())
		}
	}()

	cmd, err := w.prepareScanCmd(imageRef, reportFile.Name())
	if err != nil {
		return report, err
	}

	klog.V(5).Infof("Exec command with args '%v' '%v'", cmd.Path, strings.Join(cmd.Args, ","))

	stdout, err := w.runner.RunCmd(cmd)
	if err != nil {
		klog.V(5).Infof("Running trivy failed '%v' '%v' - \n%v", imageRef.Name, cmd.ProcessState.ExitCode(), string(stdout))
		return report, fmt.Errorf("running trivy: %v: %v", err, string(stdout))
	}

	klog.V(7).Infof("Running trivy failed '%v' '%v' - \n%v", imageRef.Name, cmd.ProcessState.ExitCode(), string(stdout))

	report, err = ScanReportFrom(reportFile)
	return report, err
}

func (w *wrapper) prepareScanCmd(imageRef ImageRef, outputFile string) (*exec.Cmd, error) {
	args := []string{
		"--no-progress",
		"--cache-dir", w.config.CacheDir,
		//"--severity", w.config.Severity,
		//"--vuln-type", w.config.VulnType,
		"--format", "json",
		"--output", outputFile,
		imageRef.Name,
	}

	if w.config.IgnoreUnfixed {
		args = append([]string{"--ignore-unfixed"}, args...)
	}

	if w.config.DebugMode {
		args = append([]string{"--debug"}, args...)
	}

	if w.config.SkipUpdate {
		args = append([]string{"--skip-update"}, args...)
	}

	name, err := w.runner.LookPath(trivyCmd)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(name, args...)

	cmd.Env = w.runner.Environ()

	switch a := imageRef.Auth.(type) {
	case NoAuth:
	case BasicAuth:
		cmd.Env = append(cmd.Env,
			fmt.Sprintf("TRIVY_USERNAME=%s", a.Username),
			fmt.Sprintf("TRIVY_PASSWORD=%s", a.Password))
	case BearerAuth:
		cmd.Env = append(cmd.Env,
			fmt.Sprintf("TRIVY_REGISTRY_TOKEN=%s", a.Token))
	default:
		return nil, fmt.Errorf("invalid type %T", a)
	}

	if imageRef.Insecure {
		cmd.Env = append(cmd.Env, "TRIVY_NON_SSL=true")
	}

	if strings.TrimSpace(w.config.GitHubToken) != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("GITHUB_TOKEN=%s", w.config.GitHubToken))
	}

	if w.config.Insecure {
		cmd.Env = append(cmd.Env, "TRIVY_INSECURE=true")
	}

	return cmd, nil
}

func (w *wrapper) GetVersion() (VersionInfo, error) {
	cmd, err := w.prepareVersionCmd()
	if err != nil {
		return VersionInfo{}, err
	}

	versionOutput, err := w.runner.RunCmd(cmd)
	if err != nil {
		klog.V(5).Infof("Running trivy failed '%v' - \n%v", cmd.ProcessState.ExitCode(), string(versionOutput))
		return VersionInfo{}, fmt.Errorf("running trivy: %v: %v", err, string(versionOutput))
	}

	var vi VersionInfo
	_ = json.Unmarshal(versionOutput, &vi)

	return vi, nil
}

func (w *wrapper) prepareVersionCmd() (*exec.Cmd, error) {
	args := []string{
		"--version",
		"--cache-dir", w.config.CacheDir,
		"--format", "json",
	}

	name, err := w.runner.LookPath(trivyCmd)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(name, args...)
	return cmd, nil
}
