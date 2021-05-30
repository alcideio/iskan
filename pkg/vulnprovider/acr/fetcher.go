package acr

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	arg "github.com/Azure/azure-sdk-for-go/services/resourcegraph/mgmt/2019-04-01/resourcegraph"
	"github.com/alcideio/iskan/pkg/util"
	types "github.com/alcideio/iskan/pkg/vulnprovider/api"
	"github.com/kylelemons/godebug/pretty"
	"google.golang.org/genproto/googleapis/grafeas/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog"
)

type imageVulnerabilitiesFinder struct {
	client arg.BaseClient

	azureSubscription string
}

func NewImageVulnerabilitiesFinder(cred *types.VulnProviderAPICreds) (types.ImageVulnerabilitiesFinder, error) {
	//klog.V(10).Info("loading supplied credentials", pretty.Sprint(cred))

	if cred == nil || cred.ACR == nil {
		return nil, fmt.Errorf("Missing creds")
	}
	// Create and authorize a ResourceGraph client
	argClient := arg.New()
	authorizer, err := NewAuthorizer(cred.ACR)

	if err != nil {
		return nil, err
	}
	argClient.Authorizer = authorizer

	return &imageVulnerabilitiesFinder{
		client:            argClient,
		azureSubscription: cred.ACR.SubscriptionId,
	}, nil
}

//Which Registry Platform it supports
func (i *imageVulnerabilitiesFinder) Type() string {
	return "acr"
}

func (i *imageVulnerabilitiesFinder) ListOccurrences(_ context.Context, containerImage string) (*types.ImageScanResult, error) {
	findings, err := i.getImageScanFindings(containerImage)
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

func (i *imageVulnerabilitiesFinder) getImageScanFindings(containerImage string) ([]*grafeas.Occurrence, error) {
	const queryFmt = `securityresources
| where type == "microsoft.security/assessments"
| summarize by assessmentKey=name //the ID of the assessment
| join kind=inner (
    securityresources
     | where type == "microsoft.security/assessments/subassessments"
     | extend assessmentKey = extract(".*assessments/(.+?)/.*",1,  id)
 ) on assessmentKey
  | project assessmentKey, subassessmentKey=name, id, parse_json(properties), resourceGroup, subscriptionId, tenantId
  | where properties.additionalData.repositoryName == "%v"
  | where properties.additionalData.registryHost == "%v"
  | where properties.additionalData.imageDigest == "%v"
  | mvexpand properties.additionalData.cve
  | extend cve = properties_additionalData_cve["title"],
         severity = properties.status.severity,
         impact = properties.impact,
		 cvss3 = properties.additionalData.cvss["3.0"].base,
 		 cvss2 = properties.additionalData.cvss["2.0"].base,
         statusCode = properties.status.code,
         description = properties.description,
         displayName = properties.displayName,
         //resourceId = properties.resourceDetails.id,
         //resourceSource = properties.resourceDetails.source,
         //category = properties.category,
         //additionalData = properties.additionalData,
         timeGenerated = properties.timeGenerated,
         remediation = properties.remediation,
		 fixable = properties.additionalData.patchable,
		 references = properties.additionalData.vendorReferences,
		 link = properties_additionalData_cve["link"],	
		 repositoryName = properties.additionalData.repositoryName,
		 imageDigest = properties.additionalData.imageDigest
   | project-away properties, assessmentKey, subassessmentKey, tenantId, subscriptionId, properties_additionalData_cve, resourceGroup`

	repo, tag, digest, err := util.ParseImageName(containerImage)
	if err != nil {
		return nil, err
	}

	repoUrl, err := url.Parse("https://" + repo)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(queryFmt, strings.TrimPrefix(repoUrl.RequestURI(), "/"), repoUrl.Hostname(), digest)
	// Set options
	RequestOptions := arg.QueryRequestOptions{
		ResultFormat: arg.ResultFormatObjectArray,
	}

	// Create the query request
	Request := arg.QueryRequest{
		Subscriptions: &[]string{i.azureSubscription},
		Query:         &query,
		Options:       &RequestOptions,
	}

	errs := []error{}
	occurrences := make([]*grafeas.Occurrence, 0)

	for {

		results, err := i.client.Resources(context.Background(), Request)
		if err != nil {
			klog.Errorf("[image=%v][repo=%v],[tag=%v],[digest=%v] - %v - %v", containerImage, repo, tag, digest, Request, err)
			return nil, err
		}

		klog.V(7).Infof("[image=%v][repo=%v],[tag=%v],[digest=%v] - \n%v\n", containerImage, repo, tag, digest, pretty.Sprint(results.Data))

		for {

			resData, ok := results.Data.([]interface{})
			if !ok {
				klog.Errorf("Failed to cast result - %v", reflect.TypeOf(results.Data))
				break
			}

			vulnOccurences, err := getFindings(resData)
			if err != nil {
				errs = append(errs, fmt.Errorf("Image Scan findings are unavailable ('%v')", err))
				break
			}

			for _, v := range vulnOccurences {
				o := newImageScanOccurrence(i.azureSubscription, "", containerImage, tag, digest, "XXX")
				o.Details = v
				o.Kind = grafeas.NoteKind_VULNERABILITY
				occurrences = append(occurrences, o)
			}

			break
		}

		Request.Options.SkipToken = results.SkipToken
		if Request.Options.SkipToken == nil {
			break
		}
	}

	return occurrences, errors.NewAggregate(errs)
}

func getFindings(findings []interface{}) ([]*grafeas.Occurrence_Vulnerability, error) {
	vulnerabilityDetails := make([]*grafeas.Occurrence_Vulnerability, 0)

	for _, f := range findings {
		finding, ok := f.(map[string]interface{})
		if !ok {
			klog.Errorf("Failed to cast result - %v", reflect.TypeOf(f))
			continue
		}

		if finding["statusCode"] == "Healthy" {
			continue
		}
		packageSeverity := getVulnerabilitySeverity(finding["severity"].(string))

		v := &grafeas.Occurrence_Vulnerability{
			Vulnerability: &grafeas.VulnerabilityOccurrence{
				CvssScore:         float32(finding["cvss3"].(float64)),
				ShortDescription:  finding["cve"].(string),
				LongDescription:   fmt.Sprintf("%v %v %v", finding["displayName"], finding["description"], finding["remediation"]),
				Severity:          packageSeverity,
				EffectiveSeverity: packageSeverity,
				FixAvailable:      finding["fixable"].(bool),
				RelatedUrls:       getRelatedUrls(finding),
				PackageIssue: []*grafeas.VulnerabilityOccurrence_PackageIssue{
					{
						//AffectedCpeUri:  packageURI,
						//AffectedPackage: packageName,
						//AffectedVersion: &grafeas.Version{
						//	Kind: grafeas.Version_NORMAL,
						//	Name: packageVersion,
						//},
						//FixedCpeUri:  "",
						//FixedPackage: "",
						//FixedVersion: nil,
						//FixAvailable: false,
					},
				},
			},
		}
		vulnerabilityDetails = append(vulnerabilityDetails, v)
	}

	return vulnerabilityDetails, nil
}

func getRelatedUrls(finding map[string]interface{}) []*grafeas.RelatedUrl {
	urls := make([]*grafeas.RelatedUrl, 0)

	urlLink, ok := finding["link"].(string)

	if ok && urlLink != "" {
		urls = append(urls, &grafeas.RelatedUrl{
			Url:   urlLink,
			Label: "",
		})
	}

	//TODO: Add finding["references"][*].link

	return urls
}

type ACRImageScanSeverity string

const (
	ACRSeverityCritical = "Critical"
	ACRSeverityHigh     = "High"
	ACRSeverityMedium   = "Medium"
	ACRSeverityLow      = "Low"
	ACRSeverityUnknown  = "Unknown"
)

func getVulnerabilitySeverity(v string) grafeas.Severity {
	switch v {
	case ACRSeverityCritical:
		return grafeas.Severity_CRITICAL
	case ACRSeverityHigh:
		return grafeas.Severity_HIGH
	case ACRSeverityMedium:
		return grafeas.Severity_MEDIUM
	case ACRSeverityLow:
		return grafeas.Severity_LOW
	case ACRSeverityUnknown:
		return grafeas.Severity_SEVERITY_UNSPECIFIED
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
