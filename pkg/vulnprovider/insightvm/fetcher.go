package insightvm

import (
	"context"
	"fmt"
	"github.com/alcideio/iskan/pkg/version"
	"github.com/kylelemons/godebug/pretty"
	"google.golang.org/genproto/googleapis/grafeas/v1"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog"
	"net/url"
	"strings"

	"github.com/alcideio/iskan/pkg/util"
	types "github.com/alcideio/iskan/pkg/vulnprovider/api"

	insightvm "github.com/alcideio/iskan/pkg/vulnprovider/insightvm/client"
)

type imageVulnerabilitiesFinder struct {
	client *insightvm.APIClient

	config *types.InsightVM
}

func NewImageVulnerabilitiesFinder(cred *types.VulnProviderAPICreds) (types.ImageVulnerabilitiesFinder, error) {

	if cred == nil || cred.InsightVM == nil || cred.InsightVM.ApiKey == "" {
		klog.V(10).Info("loading supplied credentials", pretty.Sprint(cred))
		return nil, fmt.Errorf("Failed to create InsightVM client - missing configuration")
	}

	cfg := insightvm.NewConfiguration()
	cfg.UserAgent = fmt.Sprintf("iskan/%v/go", version.Version)

	client := insightvm.NewAPIClient(cfg)
	if client == nil {
		return nil, fmt.Errorf("Failed to create InsightVM client")
	}

	return &imageVulnerabilitiesFinder{
		client: client,
		config: cred.InsightVM,
	}, nil
}

//Which Registry Platform it supports
func (i *imageVulnerabilitiesFinder) Type() string {
	return "insightvm"
}

func (i *imageVulnerabilitiesFinder) ListOccurrences(ctx context.Context, containerImage string) (*types.ImageScanResult, error) {
	findings, err := i.getImageScanFindings(ctx, containerImage)
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

func (i *imageVulnerabilitiesFinder) getImageId(containerImage string) (string, error) {
	repo, tag, digest, err := util.ParseImageName(containerImage)
	if err != nil {
		return "", err
	}

	repoUrl, err := url.Parse("https://" + repo)
	if err != nil {
		return "", err
	}

	registryId := repoUrl.Host
	switch registryId {
	case "docker.io":
		registryId = "_DockerHub"
	case "quay.io":
		registryId = "_Quay"
	}

	repoId := strings.TrimPrefix(repoUrl.RequestURI(), "/")

	errs := []error{}
	apiKeys := map[string]insightvm.APIKey{
		"api-key": insightvm.APIKey{
			Key:    i.config.ApiKey,
			Prefix: "",
		},
	}
	apiCtx := context.WithValue(context.Background(), insightvm.ContextAPIKeys, apiKeys)

	nImages := 0
	var page int32 = 0
	nResources := 0 //number of resources

	for {
		var req insightvm.ApiGetRegistryRepositoryImagesRequest

		req = i.client.ContainerApi.GetRegistryRepositoryImages(apiCtx, insightvm.Region(i.config.Region), registryId, repoId).Page(page).Size(100)

		//klog.V(8).Infof("API REQ: \nResources%+v\nResources", req)

		imgs, httpresp, err := req.Execute()
		if err != nil {
			bodyStr := ""
			var code int
			if httpresp != nil {
				body, _ := ioutil.ReadAll(httpresp.Body)
				bodyStr = string(body)
				code = httpresp.StatusCode
			}
			klog.V(3).Infof("[http=%v][registryId=%v][image=%v][repo=%v][tag=%v][digest=%v] - %v - \n%v\n", code, registryId, containerImage, repoId, tag, digest, err, bodyStr)
			errs = append(errs, err)
			return "", errors.NewAggregate(errs)
		}

		metadata := imgs.GetMetadata()
		if nResources == 0 { //Once
			nImages = (int)(*metadata.TotalData)
		}

		for _, img := range imgs.GetData() {
			nResources++
			if len(digest) > 0 {
				for _, imgDigest := range img.GetDigests() {
					if digest == imgDigest.GetDigest() {
						klog.V(5).Infof("[image=%v]-->%v", containerImage, img.GetId())
						return img.GetId(), nil
					}
				}
			} else {
				for _, imgTags := range img.GetTags() {
					if tag == imgTags.GetName() {
						klog.V(5).Infof("[image=%v]-->%v", containerImage, img.GetId())
						return img.GetId(), nil
					}
				}
			}
		}
		page++

		klog.V(7).Infof("[image=%v][registryId=%v][repoId=%v][page=%v][nResources=%v][nImages=%v][index=%v,size=%v]", containerImage, registryId, repoId, page, nResources, nImages, *metadata.Index, *metadata.Size)
		if nResources >= nImages {
			break
		}
	}

	return "", fmt.Errorf("Failed to find %v in InsightVM", containerImage)
}

func (i *imageVulnerabilitiesFinder) getImageScanFindings(ctx context.Context, containerImage string) ([]*grafeas.Occurrence, error) {
	id, err := i.getImageId(containerImage)
	if err != nil {
		return nil, err
	}

	errs := []error{}

	apiKeys := map[string]insightvm.APIKey{
		"api-key": insightvm.APIKey{
			Key:    i.config.ApiKey,
			Prefix: "",
		},
	}
	apiCtx := context.WithValue(ctx, insightvm.ContextAPIKeys, apiKeys)

	imgFindings, httpresp, err := i.client.ContainerApi.GetImage(apiCtx, insightvm.Region(i.config.Region), id).Execute()
	if err != nil {
		bodyStr := ""
		if httpresp != nil {
			body, _ := ioutil.ReadAll(httpresp.Body)
			bodyStr = string(body)
		}
		klog.Errorf("[http=%v][image=%v][id=%v] - %v - \n%v\n", httpresp.StatusCode, containerImage, id, err, bodyStr)
		errs = append(errs, err)
		return nil, errors.NewAggregate(errs)
	}

	klog.V(7).Infof("[image=%v][id=%v] - %v", containerImage, id, pretty.Sprint(imgFindings))

	vulnOccurences, err := getFindings(&imgFindings)
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

func getRelatedUrls(l []insightvm.PackageVulnerabilityReference) []*grafeas.RelatedUrl {
	urls := make([]*grafeas.RelatedUrl, len(l))

	for i, e := range l {
		urls[i] = &grafeas.RelatedUrl{
			Url:   e.GetUrl(),
			Label: "",
		}
	}

	return urls
}

func getPkgInfo(pkgId int64, pkgs []insightvm.Package) *insightvm.Package {
	for i, p := range pkgs {
		if p.GetId() == pkgId {
			return &pkgs[i]
		}
	}

	return nil
}

func getPackageIssues(finding insightvm.PackageVulnerabilityEvaluation, pkgs []insightvm.Package) []*grafeas.VulnerabilityOccurrence_PackageIssue {
	results := finding.GetResults()

	if len(results) == 0 {
		return nil
	}

	result := results[0]
	pkgInfo := getPkgInfo(result.GetPackageId(), pkgs)

	if pkgInfo == nil {
		return nil
	}

	pkgIssues := []*grafeas.VulnerabilityOccurrence_PackageIssue{
		{
			AffectedCpeUri:  "",
			AffectedPackage: pkgInfo.GetName(),
			AffectedVersion: &grafeas.Version{
				Kind:     grafeas.Version_NORMAL,
				Name:     pkgInfo.GetVersion(),
				FullName: fmt.Sprintf("%v:%v", pkgInfo.GetName(), pkgInfo.GetVersion()),
			},
			//FixedCpeUri:          "",
			//FixedPackage:         "",
			//FixedVersion:         &grafeas.Version{
			//				Kind: grafeas.Version_NORMAL,
			//				Name: pkgInfo.GetVersion(),
			//			},
			//FixAvailable:         false,
		},
	}

	return pkgIssues
}

func getCvssScore(vul *insightvm.PackageVulnerability) float32 {
	if vul.HasCvssV3() {
		cvss := vul.GetCvssV3()
		return float32(cvss.GetScore())
	}

	if vul.HasCvssV2() {
		cvss := vul.GetCvssV2()
		return float32(cvss.GetScore())
	}

	return 0
}

func getFindings(report *insightvm.Image) ([]*grafeas.Occurrence_Vulnerability, error) {
	vulnerabilityDetails := make([]*grafeas.Occurrence_Vulnerability, 0)

	assesment := report.GetAssessment()

	for _, f := range assesment.GetFindings() {
		vul := f.GetVulnerability()

		packageSeverity := getVulnerabilitySeverity(vul.GetSeverity())
		vulDescrption := vul.GetDescription()

		v := &grafeas.Occurrence_Vulnerability{
			Vulnerability: &grafeas.VulnerabilityOccurrence{
				Severity:    packageSeverity,
				CvssScore:   getCvssScore(&vul),
				RelatedUrls: getRelatedUrls(vul.GetReferences()),
				//FixAvailable:      (f.FixVersion != ""),
				EffectiveSeverity: packageSeverity,
				ShortDescription:  vul.GetId(),
				LongDescription:   vulDescrption.GetText(),
				PackageIssue:      getPackageIssues(f, report.GetPackages()),
			},
		}

		vulnerabilityDetails = append(vulnerabilityDetails, v)
	}

	return vulnerabilityDetails, nil
}

func getVulnerabilitySeverity(v string) grafeas.Severity {
	// See the enum in openapi.json
	switch v {
	case "critical":
		return grafeas.Severity_CRITICAL
	case "severe":
		return grafeas.Severity_HIGH
	case "moderate":
		return grafeas.Severity_MEDIUM
	case "low":
		return grafeas.Severity_LOW
	case "informational":
		return grafeas.Severity_MINIMAL
	case "none":
		fallthrough
	default:
		return grafeas.Severity_SEVERITY_UNSPECIFIED
	}
}
