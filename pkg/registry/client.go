package registry

import (
	"fmt"
	"github.com/alcideio/iskan/pkg/util"
	"regexp"
	"strings"

	"github.com/alcideio/iskan/api"
	"github.com/alcideio/iskan/pkg/registry/ecr"
	"github.com/alcideio/iskan/pkg/registry/gcr"
)

const (
	RegistryKind_ECR string = "ecr"
	RegistryKind_GCR string = "gcr"
	RegistryKind_DTR string = "dtr" //Dockdr Hub Enterprise

	RegistryKind_UNKNOWN string = "unknown"
)

func NewImageVulnerabilitiesFinder(kind string, cred *api.RegistryAPICreds) (api.ImageVulnerabilitiesFinder, error) {
	switch strings.ToLower(kind) {
	case RegistryKind_ECR:
		return ecr.NewImageVulnerabilitiesFinder(cred)
	case RegistryKind_GCR:
		return gcr.NewImageVulnerabilitiesFinder(cred)
	case RegistryKind_DTR:
		return nil, fmt.Errorf("registry type '%v' is not supported", kind)
	default:
		return nil, fmt.Errorf("registry type '%v' is not supported", kind)

	}
}

var ecrRegex = regexp.MustCompile(`(^[a-zA-Z0-9][a-zA-Z0-9-_]*)\.dkr\.ecr(\-fips)?\.([a-zA-Z0-9][a-zA-Z0-9-_]*)\.amazonaws\.com(\.cn)?`)
var gcrRegex = regexp.MustCompile("(\\S.)*gcr.io")
var dockerHubRegex = regexp.MustCompile("^docker.io/\\S+/")

func DetectRegistryKind(image string) (string, error) {
	repo, _, _, err := util.ParseImageName(image)
	if err != nil {
		return RegistryKind_UNKNOWN, fmt.Errorf("failed to detect")
	}

	if gcrRegex.MatchString(repo) {
		return RegistryKind_GCR, nil
	}

	if ecrRegex.MatchString(repo) {
		return RegistryKind_ECR, nil
	}

	return RegistryKind_UNKNOWN, nil
}
