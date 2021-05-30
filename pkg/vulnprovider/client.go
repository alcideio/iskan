package vulnprovider

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/alcideio/iskan/pkg/util"
	"github.com/alcideio/iskan/pkg/vulnprovider/acr"
	"github.com/alcideio/iskan/pkg/vulnprovider/api"
	"github.com/alcideio/iskan/pkg/vulnprovider/ecr"
	"github.com/alcideio/iskan/pkg/vulnprovider/gcr"
	"github.com/alcideio/iskan/pkg/vulnprovider/harbor"
	"github.com/alcideio/iskan/pkg/vulnprovider/insightvm"
	"github.com/alcideio/iskan/pkg/vulnprovider/trivy"
)

const (
	ProviderKind_ECR       string = "ecr"
	ProviderKind_GCR       string = "gcr"
	ProviderKind_ACR       string = "acr"
	ProviderKind_HARBOR    string = "harbor"
	ProviderKind_INSIGHTVM string = "insightvm"
	ProviderKind_DTR       string = "dtr"   //Dockdr Hub Enterprise
	ProviderKind_TRIVY     string = "trivy" //Local Processor
	ProviderKind_UNKNOWN   string = "unknown"
)

func NewImageVulnerabilitiesFinder(kind string, cred *api.VulnProviderAPICreds) (api.ImageVulnerabilitiesFinder, error) {
	switch strings.ToLower(kind) {
	case ProviderKind_ECR:
		return ecr.NewImageVulnerabilitiesFinder(cred)
	case ProviderKind_GCR:
		return gcr.NewImageVulnerabilitiesFinder(cred)
	case ProviderKind_ACR:
		return acr.NewImageVulnerabilitiesFinder(cred)
	case ProviderKind_TRIVY:
		return trivy.NewImageVulnerabilitiesFinder(cred)
	case ProviderKind_HARBOR:
		return harbor.NewImageVulnerabilitiesFinder(cred)
	case ProviderKind_INSIGHTVM:
		return insightvm.NewImageVulnerabilitiesFinder(cred)
	case ProviderKind_DTR:
		return nil, fmt.Errorf("registry type '%v' is not supported", kind)
	default:
		return nil, fmt.Errorf("registry type '%v' is not supported", kind)

	}
}

var ecrRegex = regexp.MustCompile(`(^[a-zA-Z0-9][a-zA-Z0-9-_]*)\.dkr\.ecr(\-fips)?\.([a-zA-Z0-9][a-zA-Z0-9-_]*)\.amazonaws\.com(\.cn)?`)
var gcrRegex = regexp.MustCompile("(\\S.)*gcr.io")
var acrRegex = regexp.MustCompile("(\\S.)*azurecr.io")
var dockerHubRegex = regexp.MustCompile("^docker.io/\\S+/")

func DetectProviderKind(image string) (string, error) {
	repo, _, _, err := util.ParseImageName(image)
	if err != nil {
		return ProviderKind_UNKNOWN, fmt.Errorf("failed to detect")
	}

	if gcrRegex.MatchString(repo) {
		return ProviderKind_GCR, nil
	}

	if ecrRegex.MatchString(repo) {
		return ProviderKind_ECR, nil
	}

	if acrRegex.MatchString(repo) {
		return ProviderKind_ACR, nil
	}

	return ProviderKind_UNKNOWN, nil
}
