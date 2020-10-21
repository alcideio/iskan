package scan

import (
	"context"
	"github.com/kylelemons/godebug/pretty"
	"k8s.io/client-go/util/flowcontrol"

	"github.com/alcideio/iskan/pkg/types"
	"github.com/alcideio/iskan/pkg/vulnprovider"
	"google.golang.org/genproto/googleapis/grafeas/v1"
	"k8s.io/klog"
)

func ScanImage(image string, policy *types.Policy, config *types.VulnProviderConfig, flowControl flowcontrol.RateLimiter) (*types.ImageScanResult, error) {

	//Apply Ratelimits
	flowControl.Accept()

	klog.V(10).Infof("[image=%v][%+v]", image, *policy)
	summary := types.NewSeveritySummary()
	fixable := types.NewSeveritySummary()

	s, err := vulnprovider.NewImageVulnerabilitiesFinder(config.Kind, &config.Creds)
	if err != nil {
		klog.V(5).Infof("[image=%v] failed to create vuln provider client - %v", image, err)
		return &types.ImageScanResult{
			Image:       image,
			CompletedOK: false,
			Reason:      err.Error(),
			Findings:    nil,
			Summary:     summary,
			Fixable:     fixable,
		}, nil
	}

	res, err := s.ListOccurrences(context.Background(), image)
	if err != nil {
		klog.V(5).Infof("[image=%v] failed to list findings - %v", image, err)
		return &types.ImageScanResult{
			Image:       image,
			CompletedOK: false,
			Reason:      err.Error(),
			Findings:    nil,
			Summary:     summary,
			Fixable:     fixable,
		}, nil
	}

	filter := RuntimeResultFilter
	filtered := []*grafeas.Occurrence{}
	result := &types.ImageScanResult{
		Image:       image,
		CompletedOK: res.CompletedOK,
		Reason:      res.Reason,
		Summary:     summary,
		Fixable:     fixable,
	}

	for _, vul := range res.Findings {
		if filter.IncludeResult(policy, vul) {
			filtered = append(filtered, vul)
			result.Summary[vul.GetVulnerability().Severity] = result.Summary[vul.GetVulnerability().Severity] + 1

			if vul.GetVulnerability().FixAvailable {
				result.Fixable[vul.GetVulnerability().Severity] = result.Fixable[vul.GetVulnerability().Severity] + 1
			}
		} else {
			result.ExcludeCount++
		}
	}

	result.Findings = filtered

	klog.V(8).Infof("%v", pretty.Sprint(result))

	return result, nil
}
