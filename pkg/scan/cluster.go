package scan

import (
	"fmt"
	"github.com/alcideio/iskan/pkg/advisor"
	"github.com/alcideio/iskan/pkg/kube"
	"github.com/alcideio/iskan/pkg/types"
	"github.com/alcideio/iskan/pkg/util"
	"github.com/kylelemons/godebug/pretty"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/klog"
)

type ClusterScanner struct {
	Policy          *types.Policy
	ProvidersConfig types.VulnProvidersConfig

	client *kube.KubeClient

	advisorReport *advisor.AdvisorClusterReport

	flowControl flowcontrol.RateLimiter
}

func NewClusterScanner(clusterContext string, policy *types.Policy, providersConfig *types.VulnProvidersConfig) (*ClusterScanner, error) {
	client, err := kube.NewClient(clusterContext)
	if err != nil {
		return nil, fmt.Errorf("Failed to create kubernetes client - %v", err)
	}

	var ratelimiter flowcontrol.RateLimiter
	// Conditionally configure rate limits
	if policy.RateLimit.ApiQPS > 0 {
		ratelimiter = flowcontrol.NewTokenBucketRateLimiter(policy.RateLimit.ApiQPS, int(policy.RateLimit.ApiBurst))
	} else {
		// if rate limits are configured off, c.operationPollRateLimiter.Accept() is a no-op
		ratelimiter = flowcontrol.NewFakeAlwaysRateLimiter()
	}

	return &ClusterScanner{
		Policy:          policy,
		ProvidersConfig: *providersConfig,
		client:          client,
		flowControl:     ratelimiter,
	}, nil
}

func (cs *ClusterScanner) generateSummary(res *types.ScanTaskResult) (*types.ClusterScanReportSummary, error) {
	summary := types.NewClusterScanReportSummary()

	summary.AnalyzedPodCount = uint32(len(res.ScannedPods))
	summary.ExcludedPodCount = uint32(len(res.SkippedPods))

	for _, f := range res.Findings {
		summary.ClusterSeverity.Add(f.Summary)
	}

	failedImages := sets.NewString()

	for _, pod := range res.ScannedPods {
		var podSummary, fixableSummary types.SeveritySummary

		podSummary = types.NewSeveritySummary()
		fixableSummary = types.NewSeveritySummary()

		podContainers := [][]v1.ContainerStatus{
			pod.Status.InitContainerStatuses,
			pod.Status.EphemeralContainerStatuses,
			pod.Status.ContainerStatuses,
		}

		var skipped bool
		var failureCount uint32

		skipped = true
		failureCount = 0
		for _, l := range podContainers {
			for _, c := range l {
				image := util.GetImageId(c.Image, c.ImageID)

				imageFindings, exist := res.Findings[image]
				if !exist {
					continue
				}

				if !imageFindings.CompletedOK {
					failedImages.Insert(imageFindings.Image)
					failureCount++
					continue
				}

				skipped = false
				podSummary.Add(imageFindings.Summary)
				fixableSummary.Add(imageFindings.Fixable)
			}
		}

		summary.FailedOrSkippedImages = failedImages.List()

		if !skipped {
			podSpecInfo := types.PodSpecSummary{
				Name:         pod.Name,
				Namespace:    pod.Namespace,
				Spec:         &pod.Spec,
				Severity:     podSummary,
				Fixable:      fixableSummary,
				ScanFailures: failureCount,
			}

			podKey := fmt.Sprintf("%v/%v", pod.Namespace, pod.Name)
			summary.PodSummary[podKey] = podSpecInfo
			summary.PodSeverity[podKey] = podSummary
			summary.PodFixableSeverity[podKey] = fixableSummary

			nsSummary, exist := summary.NamespaceSeverity[pod.Namespace]
			if !exist {
				nsSummary = types.NewSeveritySummary()
			}
			nsSummary.Add(podSummary)
			summary.NamespaceSeverity[pod.Namespace] = nsSummary
		} else {
			summary.FailedOrSkippedPods = append(summary.FailedOrSkippedPods, fmt.Sprintf("%v/%v", pod.Namespace, pod.Name))
		}
	}

	return summary, nil
}

func (cs *ClusterScanner) GetAdvisorReport() *advisor.AdvisorClusterReport {
	return cs.advisorReport
}

func (cs *ClusterScanner) Scan() (*types.ClusterScanReport, error) {
	pods, err := cs.client.ListPods(v1.NamespaceAll)
	if err != nil {
		return nil, fmt.Errorf("Failed to list pods - %v", err)
	}

	providersConfig := map[string]*types.VulnProviderConfig{}
	for i, r := range cs.ProvidersConfig.Providers {
		providersConfig[r.Repository] = &cs.ProvidersConfig.Providers[i]
	}

	klog.V(10).Infof("ClusterScanner\n%v", pretty.Sprint(cs))
	report := types.NewClusterScanReport()
	report.Policy = *cs.Policy

	errs := []error{}
	res, err := ScanTask(pods, cs.Policy, providersConfig, cs.flowControl)
	if err != nil {
		errs = append(errs, err)
	} else {
		report.Findings = res.Findings
	}

	if clusterUID, err := cs.client.GetClusterUID(); err != nil {
		errs = append(errs, err)
	} else {
		report.ClusterId = clusterUID
	}

	if summary, err := cs.generateSummary(res); err != nil {
		errs = append(errs, err)
	} else {
		report.Summary = *summary
	}

	//FIXME: REMOVE OUTSIDE
	cs.advisorReport, _ = advisor.GenerateAdvisorReport(res)

	return report, errors.NewAggregate(errs)
}
