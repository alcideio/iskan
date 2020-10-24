package ecr

import (
	"context"
	"fmt"
	"github.com/kylelemons/godebug/pretty"
	"os"
	"strings"

	"github.com/alcideio/iskan/pkg/types"
	"github.com/alcideio/iskan/pkg/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"google.golang.org/genproto/googleapis/grafeas/v1"

	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog"
)

type imageVulnerabilitiesFinder struct {
	client ecrClient
}

func NewImageVulnerabilitiesFinder(cred *types.VulnProviderAPICreds) (types.ImageVulnerabilitiesFinder, error) {
	var sess *session.Session
	// AWS Session
	if cred == nil || cred.ECR == nil || cred.ECR.AccessKeyId == "" || cred.ECR.SecretAccessKey == "" {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			Config:            *aws.NewConfig(),
			SharedConfigState: session.SharedConfigEnable,
		}))
	} else {
		klog.V(10).Info("loading supplied credentials", pretty.Sprint(cred))

		os.Setenv("AWS_ACCESS_KEY_ID", cred.ECR.AccessKeyId)
		os.Setenv("AWS_SECRET_ACCESS_KEY", cred.ECR.SecretAccessKey)

		creds := credentials.NewChainCredentials(
			[]credentials.Provider{
				&credentials.EnvProvider{},
				&credentials.SharedCredentialsProvider{},
			})

		config := *aws.NewConfig().
			WithCredentials(creds).
			WithRegion(cred.ECR.Region).
			WithCredentialsChainVerboseErrors(true)

		sess = session.Must(session.NewSessionWithOptions(session.Options{
			Config: config,
			//SharedConfigState: session.SharedConfigEnable,
		}))
	}

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

type ecrClient interface {
	DescribeImageScanFindings(input *ecr.DescribeImageScanFindingsInput) (*ecr.DescribeImageScanFindingsOutput, error)
}

func getImageScanFindings(ecrclient ecrClient, containerImage string) ([]*grafeas.Occurrence, error) {
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
	awsRepoName := strings.SplitAfterN(repo, "/", 2)

	input := &ecr.DescribeImageScanFindingsInput{
		RegistryId:     aws.String(ecrAccount),
		RepositoryName: aws.String(awsRepoName[1]),
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
			klog.Errorf("[image=%v][repo=%v],[tag=%v],[digest=%v] - %v - %v", image, repo, tag, digest, input, err)
			errs = append(errs, err)
			break
		}

		klog.V(7).Infof("[image=%v][repo=%v],[tag=%v],[digest=%v] - %v - %v", image, repo, tag, digest, input, pretty.Sprint(resp))

		if resp.ImageScanStatus != nil && resp.ImageScanStatus.Status != nil && aws.StringValue(resp.ImageScanStatus.Status) != ecr.ScanStatusComplete {
			errs = append(errs, fmt.Errorf("Image Scan findings are unavailable ('%v') - %v", aws.StringValue(resp.ImageScanStatus.Status), aws.StringValue(resp.ImageScanStatus.Description)))
		}

		if resp.ImageScanFindings != nil && len(resp.ImageScanFindings.Findings) > 0 {
			findings = append(findings, resp.ImageScanFindings.Findings...)
		}

		if resp.NextToken == nil {
			break
		}

		input.NextToken = resp.NextToken
	}

	if len(errs) > 0 {
		return nil, errors.NewAggregate(errs)
	}

	vulnOccurences, err := getFindings(findings, ecrAccount, ecrRegion, image, tag, digest)
	if err != nil {
		return nil, err
	}

	return vulnOccurences, nil
}

func getFindings(findings []*ecr.ImageScanFinding, ecrAccount string, ecrRegion string, image string, tag string, digest string) ([]*grafeas.Occurrence, error) {
	vulnerabilityDetails := make([]*grafeas.Occurrence, 0)

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
		o := newImageScanOccurrence(ecrAccount, ecrRegion, image, tag, digest, aws.StringValue(p.Name))
		o.Details = v
		o.Kind = grafeas.NoteKind_VULNERABILITY

		vulnerabilityDetails = append(vulnerabilityDetails, o)
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

func newImageScanOccurrence(accountId string, region string, repo string, tag string, digest string, cveId string) *grafeas.Occurrence {
	o := &grafeas.Occurrence{
		ResourceUri: ecrOccurrenceResourceURI(accountId, region, repo, tag, digest),
		NoteName:    cveId,
	}

	return o
}

func ecrOccurrenceResourceURI(account, region, repository, tag, digest string) string {
	return fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/%s:%s@%s", account, region, repository, tag, digest)
}
