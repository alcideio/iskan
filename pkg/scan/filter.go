package scan

import (
	"github.com/alcideio/iskan/api"
	"google.golang.org/genproto/googleapis/grafeas/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
	"strings"
)

type ScanFilter interface {
	ShouldScan(policy *api.Policy, pod *v1.Pod, image string) bool
}

type ResultFilter interface {
	IncludeResult(policy *api.Policy, occurence *grafeas.Occurrence) bool
}

type scanFilter struct{}

func (s *scanFilter) ShouldScan(policy *api.Policy, pod *v1.Pod, image string) bool {
	include := policy.ScanScope.IsNamespaceIncluded(pod.Namespace)
	if !include {
		return false
	}

	return true
}

type resultFilter struct{}

func (r *resultFilter) IncludeResult(policy *api.Policy, occurence *grafeas.Occurrence) bool {
	if policy == nil || policy.ReportFilter == nil || occurence == nil {
		klog.V(7).Infof("IncludeResult - yes")
		return true
	}

	vul := occurence.GetVulnerability()
	if vul == nil {
		klog.V(7).Infof("IncludeResult - no vul")
		return true
	}

	if vul.GetCvssScore() < policy.ReportFilter.CvssGreaterThan {
		klog.V(7).Infof("Not IncludeResult - %v < %v", vul.GetCvssScore(), policy.ReportFilter.CvssGreaterThan)
		return false
	}

	if policy.ReportFilter.FixableOnly && !vul.FixAvailable {
		klog.V(7).Infof("Not IncludeResult - fixable only '%v'", !vul.FixAvailable)
		return false
	}

	if policy.ReportFilter.Severities != "" && !strings.Contains(policy.ReportFilter.Severities, vul.GetSeverity().String()) {
		klog.V(7).Infof("Not IncludeResult - '%v' not in '%v'", vul.GetSeverity().String(), policy.ReportFilter.Severities)
		return false
	}

	//if policy.ReportFilter.CveNewerThan > 0 && occurence.CreateTime != nil || occurence.UpdateTime != nil {
	//	t := occurence.CreateTime
	//	if occurence.UpdateTime != nil {
	//		t = occurence.UpdateTime
	//	}
	//
	//	if t != nil {
	//		if time.Now().Sub(t.AsTime()) > policy.ReportFilter.CveNewerThan {
	//			klog.V(7).Infof("Not IncludeResult - '%v' is older than '%v'",  vul.GetSeverity().String(), policy.ReportFilter.CveNewerThan.String())
	//			return false
	//		}
	//	}
	//}

	klog.V(7).Infof("IncludeResult - '%v' in '%v'", vul.GetSeverity().String(), policy.ReportFilter.Severities)
	return true
}

var (
	RuntimeScanFilter   ScanFilter   = &scanFilter{}
	RuntimeResultFilter ResultFilter = &resultFilter{}
)
