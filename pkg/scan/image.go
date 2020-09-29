package scan

import (
	"context"

	"github.com/alcideio/iskan/api"
	"github.com/alcideio/iskan/pkg/registry"
	"google.golang.org/genproto/googleapis/grafeas/v1"
	"k8s.io/klog"
)

func ScanImage(image string, policy *api.Policy, config *api.RegistryConfig) (*api.ImageScanResult, error) {
	klog.V(5).Infof("[image=%v][%+v]", image, *config)

	summary := api.NewSeveritySummary()

	s, err := registry.NewImageVulnerabilitiesFinder(config.Kind, &config.Creds)
	if err != nil {
		return &api.ImageScanResult{
			Image:       image,
			CompletedOK: false,
			Reason:      err.Error(),
			Findings:    nil,
			Summary:     summary,
		}, nil
	}

	res, err := s.ListOccurrences(context.Background(), image)
	if err != nil {
		return &api.ImageScanResult{
			Image:       image,
			CompletedOK: false,
			Reason:      err.Error(),
			Findings:    nil,
			Summary:     summary,
		}, nil
	}

	filter := RuntimeResultFilter
	filtered := []*grafeas.Occurrence{}
	result := &api.ImageScanResult{
		Image:       image,
		CompletedOK: true,
		Reason:      "",
		Summary:     summary,
	}

	for _, vul := range res.Findings {
		if filter.IncludeResult(policy, vul) {
			filtered = append(filtered, vul)
			result.Summary[vul.GetVulnerability().Severity] = result.Summary[vul.GetVulnerability().Severity] + 1
		} else {
			result.ExcludeCount++
		}
	}

	result.Findings = filtered

	//klog.V(8).Infof("%v", pretty.Sprint(result))

	return result, nil
}
