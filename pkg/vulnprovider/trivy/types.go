package trivy

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"k8s.io/klog"
)

type ScanReport struct {
	Target          string          `json:"Target"`
	Vulnerabilities []Vulnerability `json:"Vulnerabilities"`
}

type Metadata struct {
	NextUpdate time.Time `json:"NextUpdate"`
	UpdatedAt  time.Time `json:"UpdatedAt"`
}

type VersionInfo struct {
	Version         string    `json:"Version,omitempty"`
	VulnerabilityDB *Metadata `json:"VulnerabilityDB"`
}

type Layer struct {
	Digest string `json:"Digest"`
	DiffID string `json:"DiffID"`
}

type CVSSInfo struct {
	V2Vector string  `json:"V2Vector,omitempty"`
	V3Vector string  `json:"V3Vector,omitempty"`
	V2Score  float32 `json:"V2Score,omitempty"`
	V3Score  float32 `json:"V3Score,omitempty"`
}

type Vulnerability struct {
	VulnerabilityID  string              `json:"VulnerabilityID"`
	PkgName          string              `json:"PkgName"`
	InstalledVersion string              `json:"InstalledVersion"`
	FixedVersion     string              `json:"FixedVersion"`
	Title            string              `json:"Title"`
	Description      string              `json:"Description"`
	Severity         string              `json:"Severity"`
	References       []string            `json:"References"`
	Layer            *Layer              `json:"Layer"`
	CVSS             map[string]CVSSInfo `json:"CVSS"`
}

func ScanReportFrom(f io.Reader) (ScanReport, error) {
	var report ScanReport
	var err error
	var scanReports []ScanReport

	err = json.NewDecoder(f).Decode(&scanReports)
	if err != nil {
		return report, fmt.Errorf("decoding scan report from file: %w", err)
	}

	if len(scanReports) == 0 {
		return report, err
	}

	// Collapse all vulnerabilities to single report
	report.Target = scanReports[0].Target
	report.Vulnerabilities = make([]Vulnerability, 0, len(scanReports[0].Vulnerabilities))
	for _, scanReport := range scanReports {
		klog.V(5).Infof("target=%v Parsing vulnerabilities", scanReport.Target)
		report.Vulnerabilities = append(report.Vulnerabilities, scanReport.Vulnerabilities...)
	}

	return report, err
}
