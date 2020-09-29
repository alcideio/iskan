package ecr

import (
	"context"
	"fmt"
	"github.com/alcideio/iskan/pkg/util"
	"github.com/alcideio/iskan/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"google.golang.org/genproto/googleapis/grafeas/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"strings"
)

type imageVulnerabilitiesFinder struct {
	client ecrClient
}

func NewImageVulnerabilitiesFinder(cred *types.RegistryAPICreds) (types.ImageVulnerabilitiesFinder, error) {
	// AWS Session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            *aws.NewConfig(),
		SharedConfigState: session.SharedConfigEnable,
	}))

	// ECR Client
	ecrclient := ecr.New(sess)

	return &imageVulnerabilitiesFinder{
		client: ecrclient,
	}, nil
}

//Which Registry Platform it supports
func (i *imageVulnerabilitiesFinder) Type() string {
	return "ecr"
}

func (i *imageVulnerabilitiesFinder) ListOccurrences(ctx context.Context, containerImage string) (*types.ImageScanResult, error) {
	findings, err := getImageScanFindings(i.client, nil, containerImage)
	if err != nil {
		return nil, err
	}

	return &types.ImageScanResult{Findings: findings}, nil
}

type ecrClient interface {
	DescribeImageScanFindings(input *ecr.DescribeImageScanFindingsInput) (*ecr.DescribeImageScanFindingsOutput, error)
}

func getImageScanFindings(ecrclient ecrClient, policy *types.ScanScope, containerImage string) ([]*grafeas.Occurrence, error) {
	repo, tag, digest, err := util.ParseImageName(containerImage)
	if err != nil {
		return nil, err
	}

	hostPart := strings.Split(repo, ".")
	if len(hostPart) != 6 {
		return nil, fmt.Errorf("Unknown host portion of ECR URL: %v", repo)
	}
	ecrAccount := hostPart[0]
	ecrRegion := hostPart[3]

	image := strings.TrimPrefix(repo, "/")

	input := &ecr.DescribeImageScanFindingsInput{
		RegistryId:     aws.String(ecrAccount),
		RepositoryName: aws.String(repo),
		ImageId:        &ecr.ImageIdentifier{},
	}

	if len(digest) > 0 {
		input.ImageId.ImageDigest = aws.String(digest)
	} else {
		input.ImageId.ImageTag = aws.String(tag)
	}

	errs := []error{}
	var findings []*ecr.ImageScanFinding
	for {
		resp, err := ecrclient.DescribeImageScanFindings(input)
		if err != nil {
			errs = append(errs, err)
			break
		}

		findings = append(findings, resp.ImageScanFindings.Findings...)

		if resp.NextToken == nil {
			break
		}

		input.NextToken = resp.NextToken
	}

	if len(errs) > 0 {
		return nil, errors.NewAggregate(errs)
	}

	vulnOccurences, err := getFindings(findings)
	if err != nil {
		return nil, err
	}

	occurrences := make([]*grafeas.Occurrence, 0)

	for _, v := range vulnOccurences {
		o := newImageScanOccurrence(ecrAccount, ecrRegion, image, tag, digest, "XXX")
		o.Details = v
		o.Kind = grafeas.NoteKind_VULNERABILITY
		occurrences = append(occurrences, o)
	}

	return occurrences, nil
}

func getFindings(findings []*ecr.ImageScanFinding) ([]*grafeas.Occurrence_Vulnerability, error) {
	vulnerabilityDetails := make([]*grafeas.Occurrence_Vulnerability, 0)

	for _, p := range findings {
		var packageURI, packageName, packageVersion string

		packageSeverity := getVulnerabilitySeverity(*p.Severity)

		//FIXME
		for _, k := range p.Attributes {
			if *k.Key == "package_name" {
				packageName = *k.Value
			} else if *k.Key == "package_version" {
				packageVersion = *k.Value
			}
		}

		packageURI = *p.Uri

		v := &grafeas.Occurrence_Vulnerability{
			Vulnerability: &grafeas.VulnerabilityOccurrence{
				Severity: packageSeverity,
				PackageIssue: []*grafeas.VulnerabilityOccurrence_PackageIssue{
					{
						AffectedCpeUri:  packageURI,
						AffectedPackage: packageName,
						AffectedVersion: &grafeas.Version{
							Kind: grafeas.Version_NORMAL,
							Name: packageVersion,
						},
						FixedCpeUri:  "",
						FixedPackage: "",
						FixedVersion: nil,
						FixAvailable: false,
					},
				},
			},
		}
		vulnerabilityDetails = append(vulnerabilityDetails, v)
	}

	return vulnerabilityDetails, nil
}

type ECRImageScanSeverity string

const (
	ECRSeverityCritical      = "CRITICAL"
	ECRSeverityHigh          = "HIGH"
	ECRSeverityMedium        = "MEDIUM"
	ECRSeverityLow           = "LOW"
	ECRSeverityInformational = "INFORMATIONAL"
)

func getVulnerabilitySeverity(v string) grafeas.Severity {
	switch v {
	case ECRSeverityCritical:
		return grafeas.Severity_CRITICAL
	case ECRSeverityHigh:
		return grafeas.Severity_HIGH
	case ECRSeverityMedium:
		return grafeas.Severity_MEDIUM
	case ECRSeverityLow:
		return grafeas.Severity_LOW
	case ECRSeverityInformational:
		return grafeas.Severity_MINIMAL
	default:
		return grafeas.Severity_SEVERITY_UNSPECIFIED
	}
}

func newImageScanOccurrence(accountId string, region string, repo string, tag string, digest string, queueName string) *grafeas.Occurrence {
	o := &grafeas.Occurrence{
		ResourceUri: ecrOccurrenceResourceURI(accountId, region, repo, tag, digest),
		NoteName:    ecrOccurrenceNote(queueName),
	}

	return o
}

func ecrOccurrenceResourceURI(account, region, repository, tag, digest string) string {
	return fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/%s:%s@%s", account, region, repository, tag, digest)
}

func ecrOccurrenceNote(queueName string) string {
	return fmt.Sprintf("projects/%s/notes/%s", "alcide", queueName)
}
