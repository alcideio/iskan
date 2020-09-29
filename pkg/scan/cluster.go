package scan

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/alcideio/iskan/api"
	"github.com/alcideio/iskan/pkg/advisor"
	"github.com/alcideio/iskan/pkg/kube"
	"github.com/alcideio/iskan/pkg/version"
	"github.com/kylelemons/godebug/pretty"
	"google.golang.org/genproto/googleapis/grafeas/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog"
)

type ClusterScanner struct {
	Policy           *api.Policy
	RegistriesConfig api.RegistriesConfig

	client *kube.KubeClient

	advisorReport *advisor.AdvisorClusterReport
}

func NewClusterScanner(clusterContext string, policy *api.Policy, registriesConfig *api.RegistriesConfig) (*ClusterScanner, error) {
	client, err := kube.NewClient(clusterContext)
	if err != nil {
		return nil, fmt.Errorf("Failed to create kubernetes client - %v", err)
	}

	return &ClusterScanner{
		Policy:           policy,
		RegistriesConfig: *registriesConfig,
		client:           client,
	}, nil
}

func (cs *ClusterScanner) generateSummary(res *ScanTaskResult) (*api.ClusterScanReportSummary, error) {
	summary := api.NewClusterScanReportSummary()

	summary.AnalyzedPodCount = uint32(len(res.ScannedPods))
	summary.ExcludedPodCount = uint32(len(res.SkippedPods))

	for _, f := range res.Findings {
		summary.ClusterSeverity.Add(f.Summary)
	}

	failedImages := sets.NewString()

	for _, pod := range res.ScannedPods {
		podSummary := api.NewSeveritySummary()

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
				image := getImageId(c.Image, c.ImageID)

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
			}
		}

		summary.FailedOrSkippedImages = failedImages.List()

		if !skipped {
			podSpecInfo := api.PodSpecSummary{
				Name:         pod.Name,
				Namespace:    pod.Namespace,
				Spec:         &pod.Spec,
				Severity:     podSummary,
				ScanFailures: failureCount,
			}

			podKey := fmt.Sprintf("%v/%v", pod.Namespace, pod.Name)
			summary.PodSummary[podKey] = podSpecInfo
			summary.PodSeverity[podKey] = podSummary

			nsSummary, exist := summary.NamespaceSeverity[pod.Namespace]
			if !exist {
				nsSummary = api.NewSeveritySummary()
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

func (cs *ClusterScanner) Scan() (*api.ClusterScanReport, error) {
	pods, err := cs.client.ListPods(v1.NamespaceAll)
	if err != nil {
		return nil, fmt.Errorf("Failed to list pods - %v", err)
	}

	regsConfig := map[string]*api.RegistryConfig{}
	for i, r := range cs.RegistriesConfig.Registries {
		regsConfig[r.Repository] = &cs.RegistriesConfig.Registries[i]
	}

	klog.V(7).Infof("ClusterScanner\n%v", pretty.Sprint(cs))
	report := api.NewClusterScanReport()
	report.Policy = *cs.Policy

	errs := []error{}
	res, err := ScanTask(pods, cs.Policy, regsConfig)
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
	cs.advisorReport, _ = cs.generateAdvisorReport(res)

	return report, errors.NewAggregate(errs)
}

func (cs *ClusterScanner) generateAdvisorReport(res *ScanTaskResult) (*advisor.AdvisorClusterReport, error) {
	advisorReport := &advisor.AdvisorClusterReport{
		AdvisorReportHeader: advisor.AdvisorReportHeader{
			CreationTimeStamp: time.Now().Format(time.RFC3339),
			ReportUID:         rand.String(10),
			Info:              "Kubernetes Native Container Image Scanner - By Alcide",
			ScannerVersion:    fmt.Sprintf("%v-%v", version.Version, version.Commit),
			MSTimeStamp:       0,
			ProfileID:         "iskan",
		},
		Reports: map[string]*advisor.AdvisorReportData{},
	}

	reportData := &advisor.AdvisorReportData{
		ResourceName:      "Pod Vulnerability Scan",
		ResourceNamespace: "KubeAdvisor",
		ResourceUID:       "iskan",
		ResourceKind:      "Kubernetes Image Vulnerabilities Scan",
		Results:           []advisor.CheckResult{},
	}

	var missingImageInfo, failedImages sets.String
	var check *advisor.CheckResult

	for _, pod := range res.ScannedPods {
		missingImageInfo = sets.NewString()
		failedImages = sets.NewString()

		podSummary := api.NewSeveritySummary()

		podContainers := [][]v1.ContainerStatus{
			pod.Status.InitContainerStatuses,
			pod.Status.EphemeralContainerStatuses,
			pod.Status.ContainerStatuses,
		}

		findings := &advisor.AdditionalFindings{
			SeverityCount: nil,
			Severity:      "",
			Findings:      []*advisor.AdditionalFinding{},
		}

		for _, l := range podContainers {
			for _, c := range l {
				var finding *advisor.AdditionalFinding

				image := getImageId(c.Image, c.ImageID)

				imageFindings, exist := res.Findings[image]
				if !exist {
					missingImageInfo.Insert(image)
					continue
				}

				if !imageFindings.CompletedOK {
					failedImages.Insert(imageFindings.Image)
					continue
				}

				for _, occurence := range imageFindings.Findings {
					finding = new(advisor.AdditionalFinding)

					vul := occurence.GetVulnerability()

					finding.Id = occurence.NoteName
					finding.References = make([]string, len(vul.RelatedUrls))
					for i, url := range vul.RelatedUrls {
						finding.References[i] = url.Url
					}
					finding.Title = vul.ShortDescription
					finding.Description = vul.LongDescription
					finding.Severity = vul.Severity.String()
					finding.Info = map[string]string{}
					finding.Info["Image"] = imageFindings.Image
					finding.Info["FixAvailable"] = fmt.Sprint(vul.FixAvailable)
					finding.Info["EffectiveSeverity"] = vul.EffectiveSeverity.String()
					finding.Info["CvssScore"] = fmt.Sprint(vul.CvssScore)
					for _, pkg := range vul.PackageIssue {
						finding.Info["AffectedPackage"] = pkg.AffectedPackage
						if pkg.AffectedVersion != nil {
							finding.Info["AffectedVersion"] = pkg.AffectedVersion.FullName
						}
						finding.Info["FixedPackage"] = pkg.FixedPackage
						if pkg.FixedVersion != nil {
							finding.Info["FixedVersion"] = pkg.FixedVersion.FullName
						}
					}

					findings.Findings = append(findings.Findings, finding)
				}

				podSummary.Add(imageFindings.Summary)
			}
		}

		check = new(advisor.CheckResult)
		check.Resource.Name = pod.Name
		check.Resource.Namespace = pod.Namespace
		check.Resource.Kind = "Pod"
		check.Resource.Group = ""
		check.Resource.Version = "v1"
		check.ResultUID = strings.ToLower(fmt.Sprintf("iskan.pod@%v@%v@%v", pod.Namespace, pod.Name, pod.UID))
		check.Check.ModuleId = "iskan.1"
		check.Check.ModuleTitle = "Kubernetes Image Vulnerabilities Scan"
		check.Check.GroupId = "1"
		check.Check.GroupTitle = "Pod Vulnerabilities Scan"
		check.Check.CheckId = "1"
		check.Check.CheckTitle = "Container Image Vulnerability Scan"
		check.CheckId = fmt.Sprint("iskan.1.1.1.", pod.UID)

		msgBuf := bytes.NewBufferString("")
		fmt.Fprintf(msgBuf, "Total number of identified vulnerabilties in '%v/%v' by severity: %v\n", pod.Namespace, pod.Name, podSummary.String())

		if missingImageInfo.Len() > 0 {
			fmt.Fprintf(msgBuf,
				"Note: The following images did not have vulnerability information: '%v'\n",
				strings.Join(missingImageInfo.List(), ","))
		}

		if failedImages.Len() > 0 {
			fmt.Fprintf(msgBuf,
				"Note: Extracting vulnerability information for the following images were not successful: '%v'\n",
				strings.Join(failedImages.List(), ","))
		}

		check.Action = advisor.AdmissionAction_Alert.String()
		check.Message = msgBuf.String()
		check.Recommendation = `"Distroless" images contain only your application and its runtime dependencies. 
They do not contain package managers, shells or any other programs you would expect to find in a standard Linux distribution.
Restricting what's in your runtime container to precisely what's necessary for your app is the best practice to employ.
It improves the signal to noise of CVE scanners and reduces the burden of establishing provenance to just what you need.
`
		severity, _ := podSummary.Max()
		switch severity {
		case grafeas.Severity_CRITICAL.String():
			check.Severity = advisor.CheckSeverity_Critical.String()
		case grafeas.Severity_HIGH.String():
			check.Severity = advisor.CheckSeverity_High.String()
		case grafeas.Severity_MEDIUM.String():
			check.Severity = advisor.CheckSeverity_Medium.String()
		case grafeas.Severity_LOW.String():
			check.Severity = advisor.CheckSeverity_Low.String()
		case grafeas.Severity_MINIMAL.String():
			check.Severity = advisor.CheckSeverity_Low.String()
		case "":
			check.Severity = advisor.CheckSeverity_Pass.String()
		}

		check.References = []string{
			"https://github.com/GoogleContainerTools/distroless",
		}

		check.AdditionalFindings = map[string]*advisor.AdditionalFindings{}
		check.AdditionalFindings["iskan"] = findings

		reportData.Results = append(reportData.Results, *check)
	}

	advisorReport.Reports["Kubernetes Image Scan"] = reportData

	return advisorReport, nil
}
