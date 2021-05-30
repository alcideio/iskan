package harbor

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-openapi/runtime"
	"net/http"
	"strings"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/kylelemons/godebug/pretty"
	"google.golang.org/genproto/googleapis/grafeas/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog"

	"github.com/alcideio/iskan/pkg/util"
	types "github.com/alcideio/iskan/pkg/vulnprovider/api"
	"github.com/alcideio/iskan/pkg/vulnprovider/harbor/client"
	"github.com/alcideio/iskan/pkg/vulnprovider/harbor/client/artifact"
)

type imageVulnerabilitiesFinder struct {
	client *client.HarborAPI

	config *types.HarborConfig
}

func NewImageVulnerabilitiesFinder(cred *types.VulnProviderAPICreds) (types.ImageVulnerabilitiesFinder, error) {

	if cred == nil || cred.Harbor == nil || cred.Harbor.Host == "" {
		klog.V(10).Info("loading supplied credentials", pretty.Sprint(cred))
		return nil, fmt.Errorf("Failed to create Harbor client - missing configuration")
	}

	config := client.DefaultTransportConfig().WithHost(cred.Harbor.Host)
	transport := httptransport.New(config.Host, config.BasePath, config.Schemes)

	transport.Debug = true

	if cred.Harbor.Insecure {
		transport.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := client.New(transport, nil)
	if client == nil {
		return nil, fmt.Errorf("Failed to create Harbor client")
	}

	return &imageVulnerabilitiesFinder{
		client: client,
		config: cred.Harbor,
	}, nil
}

//Which Registry Platform it supports
func (i *imageVulnerabilitiesFinder) Type() string {
	return "ecr"
}

func (i *imageVulnerabilitiesFinder) ListOccurrences(ctx context.Context, containerImage string) (*types.ImageScanResult, error) {
	findings, err := getImageScanFindings(i.client, containerImage, httptransport.BasicAuth(i.config.Username, i.config.Password))
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

func getImageScanFindings(client *client.HarborAPI, containerImage string, authInfo runtime.ClientAuthInfoWriter) ([]*grafeas.Occurrence, error) {
	repo, tag, digest, err := util.ParseImageName(containerImage)
	if err != nil {
		return nil, err
	}

	req := artifact.NewGetAdditionParams()
	parts := strings.Split(repo, "/")
	if len(parts) >= 2 {
		proj := parts[1]
		repoName := strings.Join(parts[2:], "/")

		req.Addition = "vulnerabilities"
		req.RepositoryName = repoName //url.PathEscape(repoName)
		req.ProjectName = proj
	}

	if len(digest) > 0 {
		req.Reference = digest
	} else {
		req.Reference = tag
	}

	errs := []error{}

	resp, err := client.Artifact.GetAddition(req, authInfo)
	if err != nil || resp.Payload == nil {
		klog.Errorf("[image=%v][repo=%v],[tag=%v],[digest=%v] - %v - %v - %v", containerImage, repo, tag, digest, pretty.Sprint(req), pretty.Sprint(client), err)
		errs = append(errs, err)
		return nil, errors.NewAggregate(errs)
	}

	klog.V(7).Infof("[image=%v][repo=%v],[tag=%v],[digest=%v] - %v - %v", containerImage, repo, tag, digest, req, pretty.Sprint(resp))

	//  https://github.com/goharbor/harbor/blob/039733b200cc44ba23829499eb6cc71c63d3b9e6/src/pkg/scan/rest/v1/spec.go#L33
	//	MimeTypeNativeReport = "application/vnd.scanner.adapter.vuln.report.harbor+json; version=1.0"
	//
	reportPayload, exist := resp.Payload.(map[string]interface{})["application/vnd.scanner.adapter.vuln.report.harbor+json; version=1.0"]

	if !exist || reportPayload == "" {
		klog.Errorf("[image=%v][repo=%v],[tag=%v],[digest=%v] - %v - empty payload", containerImage, repo, tag, digest, req)
		errs = append(errs, fmt.Errorf("Empty paylod for application/vnd.scanner.adapter.vuln.report.harbor+json; version=1.0"))
		return nil, errors.NewAggregate(errs)
	}

	//Marshal so we can
	data, err := json.Marshal(reportPayload)
	if err != nil {
		return nil, err
	}

	//Parse Payload
	report := Report{}
	err = json.Unmarshal(data, &report)
	if err != nil {
		klog.Errorf("[image=%v][repo=%v],[tag=%v],[digest=%v] - %v - failed to unmarshal - %v - %v", containerImage, repo, tag, digest, req, err, resp.Payload)
		return nil, err
	}

	vulnOccurences, err := getFindings(&report)
	if err != nil {
		return nil, err
	}

	occurrences := make([]*grafeas.Occurrence, 0)

	for _, v := range vulnOccurences {
		o := newImageScanOccurrence(containerImage, v.Vulnerability.ShortDescription)
		o.Details = v
		o.Kind = grafeas.NoteKind_VULNERABILITY
		occurrences = append(occurrences, o)
	}

	return occurrences, nil
}

func newImageScanOccurrence(containerImage string, cveId string) *grafeas.Occurrence {
	o := &grafeas.Occurrence{
		ResourceUri: fmt.Sprintf("%s", containerImage),
		NoteName:    cveId,
	}

	return o
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

func getFindings(report *Report) ([]*grafeas.Occurrence_Vulnerability, error) {
	vulnerabilityDetails := make([]*grafeas.Occurrence_Vulnerability, 0)

	for _, f := range report.Vulnerabilities {
		packageSeverity := getVulnerabilitySeverity(f.Severity)

		v := &grafeas.Occurrence_Vulnerability{
			Vulnerability: &grafeas.VulnerabilityOccurrence{
				Severity: packageSeverity,
				//CvssScore:         getCvssScore(&f),
				RelatedUrls:       getRelatedUrls(f.Links),
				FixAvailable:      (f.FixVersion != ""),
				EffectiveSeverity: packageSeverity,
				ShortDescription:  f.ID,
				LongDescription:   f.Description,
				PackageIssue: []*grafeas.VulnerabilityOccurrence_PackageIssue{
					{
						//AffectedCpeUri:  packageURI,
						AffectedPackage: f.Package,
						AffectedVersion: &grafeas.Version{
							Kind: grafeas.Version_NORMAL,
							Name: f.Version,
						},
						FixedCpeUri:  "",
						FixedPackage: f.Package,
						FixedVersion: &grafeas.Version{
							Kind: grafeas.Version_NORMAL,
							Name: f.FixVersion,
						},
						FixAvailable: (f.FixVersion != ""),
					},
				},
			},
		}

		vulnerabilityDetails = append(vulnerabilityDetails, v)
	}

	return vulnerabilityDetails, nil
}

func getVulnerabilitySeverity(v string) grafeas.Severity {
	switch v {
	case string(Critical):
		return grafeas.Severity_CRITICAL
	case string(High):
		return grafeas.Severity_HIGH
	case string(Medium):
		return grafeas.Severity_MEDIUM
	case string(Low):
		return grafeas.Severity_LOW
	case string(Unknown):
		return grafeas.Severity_MINIMAL
	default:
		return grafeas.Severity_SEVERITY_UNSPECIFIED
	}
}
