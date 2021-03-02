package scan

import (
	"context"
	"fmt"
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

type ImageScanner struct {
	Policy *types.Policy

	ProvidersConfig types.VulnProvidersConfig

	flowControl flowcontrol.RateLimiter
}

func NewImageScanner(policy *types.Policy, providersConfig *types.VulnProvidersConfig) (*ImageScanner, error) {
	var ratelimiter flowcontrol.RateLimiter

	if policy == nil || providersConfig == nil {
		return nil, fmt.Errorf("Invalid call")
	}

	//klog.V(10).Info("providersConfig", pretty.Sprint(providersConfig))

	// Conditionally configure rate limits
	if policy.RateLimit.ApiQPS > 0 {
		ratelimiter = flowcontrol.NewTokenBucketRateLimiter(policy.RateLimit.ApiQPS, int(policy.RateLimit.ApiBurst))
	} else {
		// if rate limits are configured off, c.operationPollRateLimiter.Accept() is a no-op
		ratelimiter = flowcontrol.NewFakeAlwaysRateLimiter()
	}

	return &ImageScanner{
		Policy:          policy,
		ProvidersConfig: *providersConfig,

		flowControl: ratelimiter,
	}, nil
}

func (is *ImageScanner) Scan(image string) (*types.ImageScanResult, error) {

	regsConfig := map[string]*types.VulnProviderConfig{}
	for i, r := range is.ProvidersConfig.Providers {
		regsConfig[r.Repository] = &is.ProvidersConfig.Providers[i]
	}

	registryConfig := RegistryConfigForImage(image, regsConfig)

	return ScanImage(image, is.Policy, registryConfig, is.flowControl)
}
