package trivy

import (
	"context"
	"fmt"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"

	"github.com/alcideio/iskan/pkg/types"
	"github.com/alcideio/iskan/pkg/util"
	"google.golang.org/genproto/googleapis/grafeas/v1"
)

type imageVulnerabilitiesFinder struct {
	client Trivy
}

func NewImageVulnerabilitiesFinder(cred *types.VulnProviderAPICreds) (types.ImageVulnerabilitiesFinder, error) {
	if cred.Trivy == nil {
		return nil, fmt.Errorf("Missing Trivy config")
	}

	if cred.Trivy.CacheDir == "" {
		cred.Trivy.CacheDir = filepath.Join(homedir.HomeDir(), ".iskan/trivy/.cache")
		_ = os.MkdirAll(cred.Trivy.CacheDir, os.FileMode(0600))
	}

	if cred.Trivy.ReportsDir == "" {
		cred.Trivy.ReportsDir = filepath.Join(homedir.HomeDir(), ".iskan/trivy/.cache/reports")
		_ = os.MkdirAll(cred.Trivy.ReportsDir, os.FileMode(0600))
	}

	trivy := NewScanner(*cred.Trivy, util.DefaultCmdRunner)

	return &imageVulnerabilitiesFinder{
		client: trivy,
	}, nil
}

//Which Registry Platform it supports
func (i *imageVulnerabilitiesFinder) Type() string {
	return "trivy"
}

func (i *imageVulnerabilitiesFinder) ListOccurrences(_ context.Context, containerImage string) (*types.ImageScanResult, error) {
	findings, err := getImageScanFindings(i.client, containerImage)
	if err != nil {
		return &types.ImageScanResult{
			Findings:    findings,
			CompletedOK: false,
			Reason:      err.Error(),
		}, nil
	}

	return &types.ImageScanResult{
		Findings:    findings,
		CompletedOK: true,
		Reason:      "",
	}, nil
}

func getImageScanFindings(scanner Trivy, containerImage string) ([]*grafeas.Occurrence, error) {
	repo, tag, _, err := util.ParseImageName(containerImage)
	if err != nil {
		return nil, err
	}

	report, err := scanner.Scan(ImageRef{
		Name:     fmt.Sprintf("%v:%v", repo, tag),
		Auth:     NoAuth{},
		Insecure: false,
	})

	if err != nil {
		return nil, err
	}

	vulnOccurences, err := getFindings(report)
	if err != nil {
		return nil, err
	}

	occurrences := make([]*grafeas.Occurrence, 0)

	for _, v := range vulnOccurences {
		o := newImageScanOccurrence(containerImage)
		o.Details = v
		o.Kind = grafeas.NoteKind_VULNERABILITY
		occurrences = append(occurrences, o)
	}

	return occurrences, nil
}

func getFindings(findings ScanReport) ([]*grafeas.Occurrence_Vulnerability, error) {
	vulnerabilityDetails := make([]*grafeas.Occurrence_Vulnerability, 0)

	for _, f := range findings.Vulnerabilities {

		packageSeverity := getVulnerabilitySeverity(f.Severity)

		v := &grafeas.Occurrence_Vulnerability{
			Vulnerability: &grafeas.VulnerabilityOccurrence{
				Severity:          packageSeverity,
				CvssScore:         getCvssScore(&f),
				RelatedUrls:       getRelatedUrls(f.References),
				FixAvailable:      (f.FixedVersion != ""),
				EffectiveSeverity: packageSeverity,
				ShortDescription:  f.Title,
				LongDescription:   f.Description,
				PackageIssue: []*grafeas.VulnerabilityOccurrence_PackageIssue{
					{
						//AffectedCpeUri:  packageURI,
						AffectedPackage: f.PkgName,
						AffectedVersion: &grafeas.Version{
							Kind: grafeas.Version_NORMAL,
							Name: f.InstalledVersion,
						},
						FixedCpeUri:  "",
						FixedPackage: f.PkgName,
						FixedVersion: &grafeas.Version{
							Kind: grafeas.Version_NORMAL,
							Name: f.FixedVersion,
						},
						FixAvailable: (f.FixedVersion != ""),
					},
				},
			},
		}
		vulnerabilityDetails = append(vulnerabilityDetails, v)
	}

	return vulnerabilityDetails, nil
}

func getCvssScore(v *Vulnerability) float32 {
	var score float32 = 0

	for _, cvssInfo := range v.CVSS {
		if cvssInfo.V3Score > score {
			score = cvssInfo.V3Score
		}
	}

	return score
}

func getRelatedUrls(l []string) []*grafeas.RelatedUrl {
	urls := make([]*grafeas.RelatedUrl, len(l))

	for i, e := range l {
		urls[i] = &grafeas.RelatedUrl{
			Url:   e,
			Label: "",
		}
	}

	return urls
}

type ECRImageScanSeverity string

const (
	//UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL
	TrivySeverityCritical = "CRITICAL"
	TrivySeverityHigh     = "HIGH"
	TrivySeverityMedium   = "MEDIUM"
	TrivySeverityLow      = "LOW"
	TrivySeverityUnknown  = "UNKNOWN"
)

func getVulnerabilitySeverity(v string) grafeas.Severity {
	switch v {
	case TrivySeverityCritical:
		return grafeas.Severity_CRITICAL
	case TrivySeverityHigh:
		return grafeas.Severity_HIGH
	case TrivySeverityMedium:
		return grafeas.Severity_MEDIUM
	case TrivySeverityLow:
		return grafeas.Severity_LOW
	case TrivySeverityUnknown:
		fallthrough
	default:
		return grafeas.Severity_SEVERITY_UNSPECIFIED
	}
}

func newImageScanOccurrence(containerImage string) *grafeas.Occurrence {
	o := &grafeas.Occurrence{
		ResourceUri: fmt.Sprintf("%s", containerImage),
		NoteName:    fmt.Sprintf("projects/%s/notes/%s", "alcide", "trivy"),
	}

	return o
}
